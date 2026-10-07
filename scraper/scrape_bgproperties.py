#!/usr/bin/env python3
"""
Scraper for bulgarianproperties.com retail shop rental listings in Sofia.

Three-pass approach:
  1. Paginate search results pages — each contains ~30 listing cards
  2. Visit each listing detail page to extract exact Google Maps coordinates
  3. Output: bgproperties_listings.jsonl, bgproperties_listings.csv,
     and upsert into PostGIS bgproperties_locations table.
"""

import asyncio
import csv
import json
import logging
import math
import os
import random
import re
import sys
from datetime import datetime, timezone
from pathlib import Path

import httpx
from bs4 import BeautifulSoup

# ---------------------------------------------------------------------------
# Configuration
# ---------------------------------------------------------------------------

# Commercial properties for rent in Sofia (stip=16, sady=2 = rent only)
SEARCH_URL = (
    "https://www.bulgarianproperties.com/Search/index.php"
    "?stown=4732&stip=16&sady=2&scntr=1&c=Search"
)

HEADERS = {
    "User-Agent": (
        "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 "
        "(KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
    ),
    "Accept-Language": "en,bg;q=0.9",
}

DATABASE_URL = os.getenv(
    "DATABASE_URL",
    "postgres://geopulse:geopulse@localhost:5432/geopulse?sslmode=disable",
)

CONCURRENCY = 3
MAX_RETRIES = 3
PER_PAGE = 30
OUTPUT_DIR = Path(__file__).parent

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s %(message)s",
    datefmt="%H:%M:%S",
)
log = logging.getLogger("scraper")

# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------


def sleep_random():
    return random.uniform(1.0, 2.5)


async def fetch_with_retry(
    client: httpx.AsyncClient, url: str, *, retries: int = MAX_RETRIES
):
    """GET with exponential backoff. Returns None on permanent failure."""
    for attempt in range(retries):
        try:
            resp = await client.get(url, headers=HEADERS, follow_redirects=True)
            if resp.status_code == 404:
                return None
            resp.raise_for_status()
            return resp
        except httpx.HTTPStatusError as exc:
            if exc.response.status_code == 404:
                return None
            wait = 2**attempt + random.random()
            log.warning("HTTP %s for %s — retry %d in %.1fs", exc.response.status_code, url, attempt + 1, wait)
            await asyncio.sleep(wait)
        except httpx.RequestError as exc:
            wait = 2**attempt + random.random()
            log.warning("Request error for %s: %s — retry %d in %.1fs", url, exc, attempt + 1, wait)
            await asyncio.sleep(wait)
    return None


def parse_price(raw):
    """Extract numeric EUR price from strings like '€ 121 435' or '€ 3 100/month'."""
    if not raw:
        return None
    cleaned = raw.replace("\xa0", " ").replace(",", "").strip()
    cleaned = re.sub(r"[€EUR]", "", cleaned, flags=re.I)
    cleaned = re.sub(r"/\s*month", "", cleaned, flags=re.I)
    cleaned = cleaned.strip()
    m = re.search(r"[\d\s]+", cleaned)
    if m:
        num_str = m.group().replace(" ", "")
        try:
            return float(num_str)
        except ValueError:
            pass
    return None


def parse_area(raw):
    """Extract area in sqm from text like 'Area: 46.00 m2'."""
    if not raw:
        return None, None
    m = re.search(r"([\d,.]+)\s*m2", raw, re.I)
    area = None
    if m:
        try:
            area = float(m.group(1).replace(",", "."))
        except ValueError:
            pass

    floor = None
    text_lower = raw.lower()
    if "groundfloor" in text_lower or "ground floor" in text_lower:
        floor = 0
    elif "basement" in text_lower:
        floor = -1
    else:
        fm = re.search(r"(\d+)(?:st|nd|rd|th)\s*floor", text_lower)
        if fm:
            floor = int(fm.group(1))

    return area, floor


def parse_price_per_sqm(raw):
    """Extract price per sqm from text like '(€/m2 2 639)'."""
    if not raw:
        return None
    m = re.search(r"€/m2\s*([\d\s]+)", raw.replace("\xa0", " "))
    if m:
        num_str = m.group(1).replace(" ", "")
        try:
            return float(num_str)
        except ValueError:
            pass
    return None


def parse_location(raw):
    """Normalize location string whitespace."""
    if not raw:
        return raw
    return re.sub(r"\s+", " ", raw).strip()


def extract_neighborhood(location):
    """Extract the neighborhood part from location."""
    if not location:
        return None
    loc = re.sub(r"^Near\s+", "", location, flags=re.I).strip()
    loc = re.sub(r"^Sofia\s*,\s*", "", loc, flags=re.I).strip()
    loc = re.sub(r"^Quarter\s+", "", loc, flags=re.I).strip()
    return loc if loc else None


# ---------------------------------------------------------------------------
# Phase 1 — Search page pagination
# ---------------------------------------------------------------------------


def detect_total_pages(soup):
    """Parse total results from the button-wrapper to compute total pages."""
    wrapper = soup.select_one(".button-wrapper")
    if wrapper:
        text = wrapper.get_text()
        m = re.search(r"(\d+)\s+from\s+(\d+)", text, re.I)
        if m:
            total = int(m.group(2))
            return math.ceil(total / PER_PAGE)
    page_links = soup.select(".button-wrapper a")
    if page_links:
        max_page = 0
        for link in page_links:
            href = link.get("href", "")
            pm = re.search(r"page=(\d+)", href)
            if pm:
                max_page = max(max_page, int(pm.group(1)))
        if max_page > 0:
            return max_page + 1
    return 1


def extract_listing_from_card(card):
    """Extract all fields from a .component-property-item element."""
    prop_id = card.get("data-preference-prop-id")
    if prop_id:
        try:
            prop_id = int(prop_id)
        except ValueError:
            pass

    status_el = card.select_one(".top-labels .standard-label")
    status = status_el.get_text(strip=True) if status_el else None

    tag_els = card.select(".bottom-labels span")
    tags = [t.get_text(strip=True) for t in tag_els if t.get_text(strip=True)]

    title_el = card.select_one("a.title")
    title = title_el.get_text(strip=True) if title_el else None
    listing_url = title_el.get("href", "") if title_el else ""
    if listing_url and not listing_url.startswith("http"):
        listing_url = f"https://www.bulgarianproperties.com{listing_url}"

    location_el = card.select_one("span.location")
    location = parse_location(location_el.get_text() if location_el else None)

    price_el = card.select_one(".property-prices .regular-price")
    price_raw = price_el.get_text(strip=True) if price_el else None
    price_eur = parse_price(price_raw)

    size_el = card.select_one(".size")
    size_text = size_el.get_text() if size_el else None
    area_sqm, floor = parse_area(size_text)
    price_per_sqm = parse_price_per_sqm(size_text)

    if price_per_sqm is None and price_eur and area_sqm and area_sqm > 0:
        price_per_sqm = round(price_eur / area_sqm, 2)

    type_el = card.select_one(".type")
    property_type = None
    if type_el:
        type_text = type_el.get_text(strip=True)
        property_type = re.sub(r"^Type of property:\s*", "", type_text, flags=re.I)

    subtitle_el = card.select_one("span.list-subtitle")
    subtitle = subtitle_el.get_text(strip=True) if subtitle_el else None

    desc_el = card.select_one("span.list-description")
    description = desc_el.get_text(strip=True) if desc_el else None

    img_el = card.select_one("a.image img")
    image_url = img_el.get("src") if img_el else None

    agent_name_el = card.select_one(".broker-info .name")
    agent_name = agent_name_el.get_text(strip=True) if agent_name_el else None

    agent_role_el = card.select_one(".broker-info .info")
    agent_role = agent_role_el.get_text(strip=True) if agent_role_el else None

    neighborhood = extract_neighborhood(location)

    return {
        "url": listing_url,
        "property_id": prop_id,
        "title": title,
        "status": status,
        "property_type": property_type,
        "neighborhood": neighborhood,
        "location_raw": location,
        "area_sqm": area_sqm,
        "price_eur": price_eur,
        "price_per_sqm": price_per_sqm,
        "floor": floor,
        "tags": tags,
        "subtitle": subtitle,
        "description": description,
        "image_url": image_url,
        "agent_name": agent_name,
        "agent_role": agent_role,
        "lat": None,
        "lng": None,
        "geo_source": None,
        "scraped_at": datetime.now(timezone.utc).isoformat(),
    }


async def scrape_search_pages(client):
    """Paginate through search results, extract listing cards."""
    all_listings = {}

    resp = await fetch_with_retry(client, SEARCH_URL)
    if resp is None:
        log.error("Failed to fetch first search page — aborting")
        return []

    soup = BeautifulSoup(resp.text, "lxml")
    total_pages = detect_total_pages(soup)
    log.info("Detected %d total pages", total_pages)

    cards = soup.select(".component-property-item")
    for card in cards:
        listing = extract_listing_from_card(card)
        if listing["url"]:
            all_listings[listing["url"]] = listing

    log.info("Page 1/%d: %d listings (total so far: %d)", total_pages, len(cards), len(all_listings))

    for page_num in range(1, total_pages):
        await asyncio.sleep(sleep_random())

        sep = "&" if "?" in SEARCH_URL else "?"
        url = f"{SEARCH_URL}{sep}page={page_num}"
        resp = await fetch_with_retry(client, url)
        if resp is None:
            log.error("Failed to fetch page %d — stopping", page_num + 1)
            break

        soup = BeautifulSoup(resp.text, "lxml")
        cards = soup.select(".component-property-item")

        if not cards:
            log.info("Page %d/%d: no listings — reached end", page_num + 1, total_pages)
            break

        for card in cards:
            listing = extract_listing_from_card(card)
            if listing["url"]:
                all_listings[listing["url"]] = listing

        log.info(
            "Page %d/%d: %d listings (total so far: %d)",
            page_num + 1, total_pages, len(cards), len(all_listings),
        )

    return list(all_listings.values())


# ---------------------------------------------------------------------------
# Phase 2 — Detail page coordinate extraction
# ---------------------------------------------------------------------------

# Pattern: google.com/maps?q=42.668531154449,23.327800750264
GMAPS_COORD_RE = re.compile(
    r"google\.com/maps\?q=([0-9]+\.[0-9]+),([0-9]+\.[0-9]+)"
)


async def enrich_listing(
    client: httpx.AsyncClient,
    listing: dict,
    sem: asyncio.Semaphore,
) -> None:
    """Visit individual listing page to extract Google Maps coordinates."""
    async with sem:
        await asyncio.sleep(sleep_random())

        url = listing["url"]
        if not url:
            return

        resp = await fetch_with_retry(client, url)
        if resp is None:
            log.warning("Failed to fetch detail page: %s", url)
            listing["geo_source"] = "none"
            return

        m = GMAPS_COORD_RE.search(resp.text)
        if m:
            listing["lat"] = float(m.group(1))
            listing["lng"] = float(m.group(2))
            listing["geo_source"] = "detail_page"
        else:
            listing["geo_source"] = "none"


# ---------------------------------------------------------------------------
# Phase 3 — Database insert
# ---------------------------------------------------------------------------


def upsert_to_db(listings):
    """Upsert all listings into the bgproperties_locations table."""
    try:
        import psycopg2
    except ImportError:
        log.warning("psycopg2 not installed — skipping database insert")
        return

    conn_str = DATABASE_URL
    if conn_str.startswith("postgres://"):
        conn_str = conn_str.replace("postgres://", "postgresql://", 1)

    try:
        conn = psycopg2.connect(conn_str)
    except Exception as exc:
        log.warning("Cannot connect to database: %s — skipping DB insert", exc)
        return

    # Clear old data and insert fresh
    cur = conn.cursor()
    cur.execute("DELETE FROM bgproperties_locations")
    inserted = 0

    for listing in listings:
        lat = listing.get("lat")
        lng = listing.get("lng")
        has_coords = lat is not None and lng is not None

        try:
            cur.execute(
                """
                INSERT INTO bgproperties_locations
                    (url, property_id, title, status, property_type, neighborhood,
                     area_sqm, price_eur, price_per_sqm, floor, tags,
                     subtitle, description, image_url, agent_name, agent_role,
                     lat, lng, geo_source, scraped_at, location)
                VALUES
                    (%(url)s, %(property_id)s, %(title)s, %(status)s,
                     %(property_type)s, %(neighborhood)s,
                     %(area_sqm)s, %(price_eur)s, %(price_per_sqm)s, %(floor)s,
                     %(tags)s,
                     %(subtitle)s, %(description)s, %(image_url)s,
                     %(agent_name)s, %(agent_role)s,
                     %(lat)s, %(lng)s, %(geo_source)s, %(scraped_at)s,
                     CASE WHEN %(has_coords)s
                          THEN ST_SetSRID(ST_MakePoint(%(lng)s, %(lat)s), 4326)
                          ELSE NULL END)
                ON CONFLICT (url) DO UPDATE SET
                    property_id = EXCLUDED.property_id,
                    title = EXCLUDED.title,
                    status = EXCLUDED.status,
                    property_type = EXCLUDED.property_type,
                    neighborhood = EXCLUDED.neighborhood,
                    area_sqm = EXCLUDED.area_sqm,
                    price_eur = EXCLUDED.price_eur,
                    price_per_sqm = EXCLUDED.price_per_sqm,
                    floor = EXCLUDED.floor,
                    tags = EXCLUDED.tags,
                    subtitle = EXCLUDED.subtitle,
                    description = EXCLUDED.description,
                    image_url = EXCLUDED.image_url,
                    agent_name = EXCLUDED.agent_name,
                    agent_role = EXCLUDED.agent_role,
                    lat = EXCLUDED.lat,
                    lng = EXCLUDED.lng,
                    geo_source = EXCLUDED.geo_source,
                    scraped_at = EXCLUDED.scraped_at,
                    location = EXCLUDED.location
                """,
                {
                    **listing,
                    "tags": listing.get("tags") or [],
                    "has_coords": has_coords,
                },
            )
            inserted += 1
        except Exception as exc:
            log.error("DB error for %s: %s", listing.get("url", "?"), exc)
            conn.rollback()
            continue

    conn.commit()
    cur.close()
    conn.close()
    log.info("Database: %d inserted", inserted)


# ---------------------------------------------------------------------------
# Phase 4 — Output files
# ---------------------------------------------------------------------------

CSV_FIELDS = [
    "url", "property_id", "title", "status", "property_type", "neighborhood",
    "location_raw", "area_sqm", "price_eur", "price_per_sqm", "floor",
    "tags", "subtitle", "description", "image_url",
    "agent_name", "agent_role",
    "lat", "lng", "geo_source", "scraped_at", "geom_wkt",
]


def write_outputs(listings):
    """Write bgproperties_listings.jsonl and bgproperties_listings.csv."""
    jsonl_path = OUTPUT_DIR / "bgproperties_listings.jsonl"
    csv_path = OUTPUT_DIR / "bgproperties_listings.csv"

    with open(jsonl_path, "w", encoding="utf-8") as f:
        for listing in listings:
            f.write(json.dumps(listing, ensure_ascii=False) + "\n")

    with open(csv_path, "w", encoding="utf-8", newline="") as f:
        writer = csv.DictWriter(f, fieldnames=CSV_FIELDS)
        writer.writeheader()
        for listing in listings:
            row = {**listing}
            if row.get("lat") and row.get("lng"):
                row["geom_wkt"] = f"POINT({row['lng']} {row['lat']})"
            else:
                row["geom_wkt"] = ""
            if isinstance(row.get("tags"), list):
                row["tags"] = ", ".join(row["tags"])
            writer.writerow(row)

    log.info("Written %s (%d listings)", jsonl_path, len(listings))
    log.info("Written %s (%d listings)", csv_path, len(listings))


def print_summary(listings):
    """Print scrape summary."""
    total = len(listings)
    with_coords = sum(1 for l in listings if l.get("geo_source") == "detail_page")
    none_coords = sum(1 for l in listings if l.get("geo_source") in ("none", None))

    type_counts = {}
    for l in listings:
        t = l.get("property_type") or "Unknown"
        type_counts[t] = type_counts.get(t, 0) + 1

    print("\n" + "=" * 50)
    print(f"Total listings scraped:        {total}")
    print(f"With map coordinates:          {with_coords}")
    print(f"No coordinates:                {none_coords}")
    print(f"\nProperty type breakdown:")
    for t, c in sorted(type_counts.items(), key=lambda x: -x[1]):
        print(f"  {t}: {c}")
    print("=" * 50 + "\n")


# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------


async def main():
    async with httpx.AsyncClient(http2=True, timeout=30.0) as client:
        # Phase 1: search pages
        log.info("=== Phase 1: Scraping search result pages ===")
        listings = await scrape_search_pages(client)
        log.info("Phase 1 complete: %d listings collected", len(listings))

        if not listings:
            log.error("No listings found — exiting")
            return

        # Phase 2: visit detail pages for exact coordinates
        log.info("=== Phase 2: Extracting coordinates from detail pages ===")
        sem = asyncio.Semaphore(CONCURRENCY)
        tasks = [enrich_listing(client, listing, sem) for listing in listings]
        await asyncio.gather(*tasks)
        log.info("Phase 2 complete")

    # Output
    write_outputs(listings)
    upsert_to_db(listings)
    print_summary(listings)


if __name__ == "__main__":
    asyncio.run(main())
