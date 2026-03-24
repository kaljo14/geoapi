# GeoPulse Backend

Admin and data-seeding service for the GeoPulse platform. Handles Google Places scraping, place detail enrichment, and bulk CSV import/export. Not exposed to end users.

**Port:** `8081`
**Database:** PostgreSQL + PostGIS (shared with BFF — migrations owned by `bff/migrations/`)

---

## Prerequisites

- Go 1.22+
- Docker (with Compose v2 — `docker compose` not `docker-compose`)
- [`venom`](https://github.com/ovh/venom) (E2E tests) — install:
  ```bash
  mkdir -p ~/.local/bin
  curl -fsSL https://github.com/ovh/venom/releases/download/v1.2.0/venom.darwin-arm64 -o ~/.local/bin/venom && chmod +x ~/.local/bin/venom
  # Intel Mac: use venom.darwin-amd64
  export PATH="$HOME/.local/bin:$PATH"
  ```
- [`oapi-codegen`](https://github.com/oapi-codegen/oapi-codegen) (codegen only)
- [`sqlc`](https://sqlc.dev/) (codegen only)
- [`tsp`](https://typespec.io/) via Node.js (codegen only)
- A Google Places API key (for scrape/enrich endpoints)

> **Database setup:** Migrations live in `../bff/migrations/`. Run `make seed` from the `bff/` directory before starting this service.

---

## Quick Start

```bash
# From bff/ — spin up postgres and apply migrations (if not already done)
cd ../bff && make seed

# From backend/
export DATABASE_URL=postgres://geopulse:geopulse@localhost:5432/geopulse?sslmode=disable
export GOOGLE_API_KEY=your_key_here
make run
```

Verify:

```bash
curl localhost:8081/livez   # {"status":"ok"}
curl localhost:8081/readyz  # {"status":"ok","db":"ok"}
```

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | — | PostgreSQL connection string. If unset, server starts without a database |
| `PORT` | `8081` | HTTP listen port |
| `GOOGLE_API_KEY` | — | Google Places API key. Required for `/scrape` and `/enrich` |
| `SCRAPE_LAT` | `42.6977` | Latitude center for nearby search (default: Sofia, Bulgaria) |
| `SCRAPE_LNG` | `23.3219` | Longitude center for nearby search |
| `SCRAPE_RADIUS` | `5000` | Search radius in meters |
| `SCRAPE_TYPES` | `restaurant` | Google Places type(s) to scrape |

---

## API Endpoints

### Data Seeding

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/scrape` | Start Google Places nearby search scraper (background goroutine). Returns `202` immediately. |
| `POST` | `/enrich` | Start Google Place Details enricher for all un-enriched places (background goroutine). Returns `202` immediately. |

```bash
curl -X POST localhost:8081/scrape
# {"message":"scraper started"}

curl -X POST localhost:8081/enrich
# {"message":"enricher started"}
```

### CSV Import / Export

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/places/export` | Download all places as a full 29-column CSV |
| `GET` | `/api/places/export-simple` | Download all places as a 3-column CSV (`name`, `place_id`, `website`) |
| `POST` | `/api/places/import` | Upload a CSV file (`multipart/form-data`, field `file`, max 10 MB) |

**Import auto-detects format by column count:**

| Columns | Format | Fields |
|---------|--------|--------|
| 3 | Simple | `name`, `place_id`, `website` |
| 8 | Basic | `place_id`, `name`, `address`, `lat`, `lng`, `rating`, `business_status`, `category` |
| 29 | Full | All columns matching export header order |

```bash
# Export
curl localhost:8081/api/places/export -o places.csv
curl localhost:8081/api/places/export-simple -o places-simple.csv

# Import
curl -X POST localhost:8081/api/places/import \
  -F "file=@places.csv"
# {"imported":42}
```

### Probes

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/livez` | Always `200 {"status":"ok"}` |
| `GET` | `/readyz` | `200` if DB is reachable, `503` otherwise |

---

## How the Scraper Works

`POST /scrape` triggers a background goroutine that:
1. Calls Google Places Nearby Search API with the configured `lat/lng/radius/type`
2. Pages through all results (up to 3 pages × 20 results), sleeping 2s between pages (required by Google)
3. Upserts each place into the `places` table with `scraped_at = now()`

## How the Enricher Works

`POST /enrich` triggers a background goroutine that:
1. Queries all places where `website IS NULL` (un-enriched)
2. For each place, calls Google Place Details API to fetch: website, phone numbers, opening hours, editorial summary, types, price level, ratings, UTC offset, Maps URL, icon
3. Updates the place record; sleeps 200ms between requests

---

## Makefile Targets

| Target | Description |
|--------|-------------|
| `build` | Compile to `bin/geopulse-backend` |
| `run` | Build and run with `DATABASE_URL` and `PORT` from env |
| `test` | `go test ./...` |
| `lint` | `golangci-lint run ./...` |
| `tidy` | `go mod tidy` |
| `typespec` | Compile `typespec/main.tsp` → `api/openapi.yaml` |
| `generate` | Run `typespec` then `oapi-codegen` → `internal/generated/server.gen.go` |
| `sqlc` | Run `sqlc generate` → `internal/store/` |
| `docker-up` | `docker compose up -d --wait` (blocks until postgres is healthy) |
| `docker-down` | Stop and remove containers |
| `venom` | Run E2E contract tests against a running server |

---

## Project Structure

```
backend/
├── cmd/main.go                  # Entry point — wiring only
├── internal/
│   ├── app/
│   │   ├── app.go               # Service interface + App struct
│   │   ├── scraper.go           # Google Places Nearby Search scraper
│   │   ├── enricher.go          # Google Place Details enricher
│   │   ├── csv.go               # CSV export and import logic
│   │   └── helpers.go           # String/type conversion helpers
│   ├── handler/handler.go       # HTTP handlers
│   ├── generated/server.gen.go  # Auto-generated — DO NOT EDIT
│   └── store/
│       ├── querier.go           # Auto-generated — DO NOT EDIT
│       ├── queries.sql.go       # Auto-generated — DO NOT EDIT
│       ├── models.go            # Auto-generated — DO NOT EDIT
│       └── db.go                # Auto-generated — DO NOT EDIT
├── api/openapi.yaml             # Generated from TypeSpec
├── typespec/main.tsp            # API definition source of truth
├── test/suite.yml               # Venom E2E contract tests
├── go.mod
├── Makefile
├── sqlc.yaml                    # schema points to ../bff/migrations/
└── oapi-codegen.yaml
```

---

## E2E Tests

Requires a running server and `venom` on `$PATH`.

```bash
make run &
make venom
```

Tests cover: liveness/readiness probes, scrape/enrich accepted (202), CSV export content-type.
