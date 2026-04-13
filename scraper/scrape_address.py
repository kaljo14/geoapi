#!/usr/bin/env python3
"""
Scraper for address.bg commercial rental listings in Sofia.

Two-pass approach:
  1. Paginate search results pages — each contains ~20 listings with bulk data
  2. Visit individual listing pages for missing coordinates and street addresses
  3. Nominatim geocoding fallback for listings still without coordinates

Output: listings.jsonl, listings.csv, and upsert into PostGIS adres_locations table.
"""

import asyncio
import csv
import json
import logging
import os
import random
import re
import sys
import html
from datetime import datetime, timezone
from pathlib import Path

import httpx
from bs4 import BeautifulSoup

# ---------------------------------------------------------------------------
# Configuration
# ---------------------------------------------------------------------------

SEARCH_URL = (
    "https://address.bg/rent/sofia/l4451"
    "?estateTypes=16,6,1,13,20,28,29,26,22,19,21,25"
)

HEADERS = {
    "User-Agent": (
        "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 "
        "(KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
    ),
    "Accept-Language": "bg,en;q=0.9",
}

NOMINATIM_URL = "https://nominatim.openstreetmap.org/search"
NOMINATIM_HEADERS = {"User-Agent": "mustaci-scraper/1.0"}

DATABASE_URL = os.getenv(
    "DATABASE_URL",
    "postgres://geopulse:geopulse@localhost:5432/geopulse?sslmode=disable",
)

CONCURRENCY = 3
MAX_RETRIES = 3
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


def parse_floor(raw):
    """Normalize floor string to integer. Партер/партер → 0."""
    if raw is None:
        return None
    raw = raw.strip().lower()
    if raw in ("партер", "parter"):
        return 0
    if raw in ("сутерен", "мазе"):
        return -1
    m = re.search(r"(\d+)", raw)
    return int(m.group(1)) if m else None


# ---------------------------------------------------------------------------
# Phase 1 — Search page pagination
# ---------------------------------------------------------------------------


def extract_offers_json(raw_html):
    """Extract the offers JSON from the raw HTML :offers-object attribute.

    The attribute uses &quot; for ALL double quotes (both JSON-structural and
    those inside string values like location_info). Decoding &quot; naively via
    html.unescape() breaks JSON because inner quotes become unescaped.

    Strategy: extract the raw attribute value, decode all HTML entities EXCEPT
    &quot; (using a placeholder), then restore &quot; → " which preserves the
    JSON structure.
    """
    marker = ":offers-object=\""
    idx = raw_html.find(marker)
    if idx < 0:
        return None
    start = idx + len(marker)
    # Attribute value ends at the next bare double-quote (not &quot;)
    end = raw_html.index("\"", start)
    raw_json_html = raw_html[start:end]

    # Protect &quot; from html.unescape, then decode everything else
    decoded = raw_json_html.replace("&quot;", "\x00Q\x00")
    decoded = html.unescape(decoded)
    decoded = decoded.replace("\x00Q\x00", "\"")

    return json.loads(decoded)


async def scrape_search_pages(client):
    """Paginate through search results, extract the embedded JSON payload."""
    all_listings = {}
    page = 1

    while True:
        url = "{}&page={}".format(SEARCH_URL, page)
        resp = await fetch_with_retry(client, url)
        if resp is None:
            log.error("Failed to fetch search page %d — stopping pagination", page)
            break

        try:
            payload = extract_offers_json(resp.text)
        except (ValueError, json.JSONDecodeError):
            log.error("Failed to parse JSON on page %d", page)
            payload = None

        if payload is None:
            log.warning("No offers-object found on page %d — stopping", page)
            break

        data = payload.get("data", [])
        if not data:
            log.info("Page %d: no listings — reached end", page)
            break

        last_page = payload.get("last_page", page)

        for item in data:
            offer_id = item.get("id")
            if offer_id is None:
                continue

            quarter = item.get("quarter") or {}
            translated = quarter.get("translated") or {}
            neighborhood = translated.get("name")

            price = item.get("price")
            square = item.get("square")
            price_per_sqm = None
            if price and square and square > 0:
                price_per_sqm = round(price / square, 2)

            listing = {
                "offer_id": offer_id,
                "url": item.get("url", ""),
                "property_type": item.get("estateTypeLabel"),
                "neighborhood": neighborhood,
                "area_sqm": square,
                "price_eur": price,
                "price_per_sqm": price_per_sqm,
                "floor": None,  # enriched in pass 2
                "address_text": None,  # enriched in pass 2
                "lat": item.get("latitude"),
                "lng": item.get("longitude"),
                "geo_source": "search_api" if item.get("latitude") else None,
                "scraped_at": datetime.now(timezone.utc).isoformat(),
            }
            all_listings[offer_id] = listing

        log.info(
            "Page %d/%d: %d listings (total so far: %d)",
            page, last_page, len(data), len(all_listings),
        )

        if page >= last_page:
            break

        page += 1
        await asyncio.sleep(sleep_random())

    return list(all_listings.values())


# ---------------------------------------------------------------------------
# Phase 2 — Individual listing pages
# ---------------------------------------------------------------------------

# Regex patterns for Bulgarian street addresses
ADDRESS_PATTERNS = [
    re.compile(r'(?:ул\.|улица)\s*["\u201e\u201c]?([^"\u201c\u201d]+)["\u201c\u201d]?\s*(?:\u2116|No\.?|\u043d\u043e\u043c\u0435\u0440)?\s*(\d+[а-яА-Я]?)?', re.I),
    re.compile(r'бул\.\s*["\u201e\u201c]?([^"\u201c\u201d]+)["\u201c\u201d]?', re.I),
    re.compile(r'пл\.\s*["\u201e\u201c]?([^"\u201c\u201d]+)["\u201c\u201d]?', re.I),
]


def extract_address(text):
    """Try to extract a street address from listing description."""
    for pattern in ADDRESS_PATTERNS:
        m = pattern.search(text)
        if m:
            parts = [p for p in m.groups() if p]
            return " ".join(parts).strip()
    return None


async def enrich_listing(
    client: httpx.AsyncClient,
    listing: dict,
    sem: asyncio.Semaphore,
    skipped_log: Path,
) -> None:
    """Visit individual listing page for coordinates and address details."""
    async with sem:
        await asyncio.sleep(sleep_random())

        url = listing["url"]
        if not url:
            return

        # Ensure absolute URL
        if not url.startswith("http"):
            url = f"https://address.bg{url}"

        resp = await fetch_with_retry(client, url)
        if resp is None:
            with open(skipped_log, "a") as f:
                f.write(f"{listing['offer_id']}\t{url}\t404/failed\n")
            return

        text = resp.text
        soup = BeautifulSoup(text, "lxml")

        # Extract coordinates from <map-offer :lat="..." :lng="...">
        map_tag = soup.find(attrs={":lat": True, ":lng": True})
        if map_tag:
            try:
                lat = float(map_tag[":lat"])
                lng = float(map_tag[":lng"])
                if lat and lng:
                    # Only override if we don't have coords yet
                    if listing["lat"] is None:
                        listing["lat"] = lat
                        listing["lng"] = lng
                        listing["geo_source"] = "listing_page"
            except (ValueError, KeyError):
                pass

        # Extract floor from listing details
        floor_tag = soup.find(string=re.compile(r"Етаж|етаж", re.I))
        if floor_tag:
            parent = floor_tag.find_parent()
            if parent:
                next_el = parent.find_next_sibling()
                if next_el:
                    listing["floor"] = parse_floor(next_el.get_text(strip=True))

        # Try to extract floor from a different pattern (key-value pairs)
        if listing["floor"] is None:
            for dt in soup.find_all(["dt", "th", "span"]):
                if re.search(r"Етаж|етаж", dt.get_text()):
                    dd = dt.find_next_sibling(["dd", "td", "span"])
                    if dd:
                        listing["floor"] = parse_floor(dd.get_text(strip=True))
                        break

        # Extract address from description
        desc_tag = soup.find(class_=re.compile(r"description|desc", re.I))
        if desc_tag is None:
            # Fallback: look for translated description in the page
            desc_tag = soup.find(attrs={"v-html": re.compile(r"description", re.I)})
        if desc_tag is None:
            desc_tag = soup.find("div", class_="offer-description")

        if desc_tag:
            desc_text = desc_tag.get_text()
            addr = extract_address(desc_text)
            if addr:
                listing["address_text"] = addr


# ---------------------------------------------------------------------------
# Phase 3 — Nominatim geocoding fallback
# ---------------------------------------------------------------------------


async def geocode_listing(
    client: httpx.AsyncClient,
    listing: dict,
) -> None:
    """Geocode a listing via Nominatim if it has an address but no coordinates."""
    if listing["lat"] is not None:
        return
    if not listing.get("address_text"):
        listing["geo_source"] = "none"
        return

    query = f"{listing['address_text']}, Sofia, Bulgaria"
    try:
        resp = await client.get(
            NOMINATIM_URL,
            params={"q": query, "format": "json", "limit": 1},
            headers=NOMINATIM_HEADERS,
        )
        resp.raise_for_status()
        results = resp.json()
        if results:
            listing["lat"] = float(results[0]["lat"])
            listing["lng"] = float(results[0]["lon"])
            listing["geo_source"] = "geocoded"
            return
    except Exception as exc:
        log.warning("Geocoding failed for offer %s: %s", listing["offer_id"], exc)

    listing["geo_source"] = "none"


# ---------------------------------------------------------------------------
# Phase 4 — Database insert
# ---------------------------------------------------------------------------


def upsert_to_db(listings):
    """Upsert all listings into the adres_locations table."""
    try:
        import psycopg2
    except ImportError:
        log.warning("psycopg2 not installed — skipping database insert")
        return

    # Parse DATABASE_URL for psycopg2
    conn_str = DATABASE_URL
    if conn_str.startswith("postgres://"):
        conn_str = conn_str.replace("postgres://", "postgresql://", 1)

    try:
        conn = psycopg2.connect(conn_str)
    except Exception as exc:
        log.warning("Cannot connect to database: %s — skipping DB insert", exc)
        return

    cur = conn.cursor()
    inserted = 0
    updated = 0

    for listing in listings:
        lat = listing.get("lat")
        lng = listing.get("lng")
        has_coords = lat is not None and lng is not None

        try:
            cur.execute(
                """
                INSERT INTO adres_locations
                    (offer_id, url, property_type, neighborhood, area_sqm,
                     price_eur, price_per_sqm, floor, address_text,
                     lat, lng, geo_source, scraped_at, location)
                VALUES
                    (%(offer_id)s, %(url)s, %(property_type)s, %(neighborhood)s,
                     %(area_sqm)s, %(price_eur)s, %(price_per_sqm)s, %(floor)s,
                     %(address_text)s, %(lat)s, %(lng)s, %(geo_source)s,
                     %(scraped_at)s,
                     CASE WHEN %(has_coords)s
                          THEN ST_SetSRID(ST_MakePoint(%(lng)s, %(lat)s), 4326)
                          ELSE NULL END)
                ON CONFLICT (offer_id) DO UPDATE SET
                    url = EXCLUDED.url,
                    property_type = EXCLUDED.property_type,
                    neighborhood = EXCLUDED.neighborhood,
                    area_sqm = EXCLUDED.area_sqm,
                    price_eur = EXCLUDED.price_eur,
                    price_per_sqm = EXCLUDED.price_per_sqm,
                    floor = EXCLUDED.floor,
                    address_text = EXCLUDED.address_text,
                    lat = EXCLUDED.lat,
                    lng = EXCLUDED.lng,
                    geo_source = EXCLUDED.geo_source,
                    scraped_at = EXCLUDED.scraped_at,
                    location = EXCLUDED.location
                """,
                {
                    **listing,
                    "has_coords": has_coords,
                },
            )
            if cur.rowcount == 1:
                inserted += 1
            else:
                updated += 1
        except Exception as exc:
            log.error("DB error for offer %s: %s", listing["offer_id"], exc)
            conn.rollback()
            continue

    conn.commit()
    cur.close()
    conn.close()
    log.info("Database: %d inserted, %d updated", inserted, updated)


# ---------------------------------------------------------------------------
# Phase 5 — Output files
# ---------------------------------------------------------------------------

CSV_FIELDS = [
    "offer_id", "url", "property_type", "neighborhood", "area_sqm",
    "price_eur", "price_per_sqm", "floor", "address_text",
    "lat", "lng", "geo_source", "scraped_at", "geom_wkt",
]


def write_outputs(listings):
    """Write listings.jsonl and listings.csv."""
    jsonl_path = OUTPUT_DIR / "listings.jsonl"
    csv_path = OUTPUT_DIR / "listings.csv"

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
            writer.writerow(row)

    log.info("Written %s (%d listings)", jsonl_path, len(listings))
    log.info("Written %s (%d listings)", csv_path, len(listings))


def print_summary(listings):
    """Print scrape summary."""
    total = len(listings)
    exact = sum(1 for l in listings if l.get("geo_source") in ("search_api", "listing_page"))
    geocoded = sum(1 for l in listings if l.get("geo_source") == "geocoded")
    none_coords = sum(1 for l in listings if l.get("geo_source") in ("none", None))

    type_counts = {}
    for l in listings:
        t = l.get("property_type") or "Unknown"
        type_counts[t] = type_counts.get(t, 0) + 1

    print("\n" + "=" * 50)
    print(f"Total listings scraped:        {total}")
    print(f"With exact coordinates:        {exact}")
    print(f"Geocoded from address:         {geocoded}")
    print(f"No coordinates (fallback):     {none_coords}")
    print(f"\nProperty type breakdown:")
    for t, c in sorted(type_counts.items(), key=lambda x: -x[1]):
        print(f"  {t}: {c}")
    print("=" * 50 + "\n")


# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------


async def main():
    skipped_log = OUTPUT_DIR / "skipped.log"
    # Clear skipped log
    skipped_log.write_text("")

    async with httpx.AsyncClient(http2=True, timeout=30.0) as client:
        # Pass 1: search pages
        log.info("=== Pass 1: Scraping search result pages ===")
        listings = await scrape_search_pages(client)
        log.info("Pass 1 complete: %d listings collected", len(listings))

        if not listings:
            log.error("No listings found — exiting")
            return

        # Pass 2: individual listing pages (all listings, prioritize those without coords)
        log.info("=== Pass 2: Enriching individual listing pages ===")
        sem = asyncio.Semaphore(CONCURRENCY)

        # Sort: listings without coordinates first
        listings.sort(key=lambda l: (l["lat"] is not None, l["offer_id"]))

        tasks = [
            enrich_listing(client, listing, sem, skipped_log)
            for listing in listings
        ]
        await asyncio.gather(*tasks)
        log.info("Pass 2 complete")

        # Pass 3: Nominatim geocoding for listings still without coordinates
        log.info("=== Pass 3: Geocoding fallback ===")
        no_coords = [l for l in listings if l["lat"] is None and l.get("address_text")]
        log.info("%d listings with address but no coordinates — geocoding", len(no_coords))
        for listing in no_coords:
            await geocode_listing(client, listing)
            await asyncio.sleep(1.1)  # Nominatim rate limit

        # Mark remaining without coords
        for listing in listings:
            if listing["lat"] is None and listing["geo_source"] is None:
                listing["geo_source"] = "none"

    # Output
    write_outputs(listings)
    upsert_to_db(listings)
    print_summary(listings)


if __name__ == "__main__":
    asyncio.run(main())
