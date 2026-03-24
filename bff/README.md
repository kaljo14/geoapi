# GeoPulse BFF

Frontend-facing API service for the GeoPulse platform. Serves place data, metro transit geometry, market saturation scores, and opportunity heatmap tiles to the Vue frontend.

**Port:** `8080`
**Database:** PostgreSQL + PostGIS (via `postgis/postgis:16-3.4`)

---

## Prerequisites

- Go 1.22+
- Docker (with Compose v2 — `docker compose` not `docker-compose`)
- [`golang-migrate`](https://github.com/golang-migrate/migrate) CLI
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

---

## Quick Start

```bash
# 1. Start postgres and apply all migrations
make seed

# 2. Run the server
make run
```

The server starts on `:8080`. Verify with:

```bash
curl localhost:8080/livez   # {"status":"ok"}
curl localhost:8080/readyz  # {"status":"ok","db":"ok"}
```

---

## Environment Variables

| Variable       | Default | Description                                                                                       |
| -------------- | ------- | ------------------------------------------------------------------------------------------------- |
| `DATABASE_URL` | —       | PostgreSQL connection string. If unset, server starts without a database (probes return degraded) |
| `PORT`         | `8080`  | HTTP listen port                                                                                  |

Example `.env`:

```
DATABASE_URL=postgres://geopulse:geopulse@localhost:5432/geopulse?sslmode=disable
PORT=8080
```

---

## API Endpoints

### Places

| Method   | Path                     | Query params      | Description                          |
| -------- | ------------------------ | ----------------- | ------------------------------------ |
| `GET`    | `/api/places`            | `category`, `tag` | List all places, optionally filtered |
| `POST`   | `/api/places`            | —                 | Create a place                       |
| `GET`    | `/api/places/{place_id}` | —                 | Get a single place                   |
| `PUT`    | `/api/places/{place_id}` | —                 | Update a place                       |
| `DELETE` | `/api/places/{place_id}` | —                 | Delete a place                       |

**Create/Update body:**

```json
{
  "name": "Test Cafe",
  "lat": 42.697,
  "lng": 23.321,
  "category": "cafe",
  "address": "...",
  "website": "...",
  "tags": "...",
  "rating": 4.5
}
```

### Metro Transit

| Method | Path                | Description                                     |
| ------ | ------------------- | ----------------------------------------------- |
| `GET`  | `/api/metro/shapes` | GeoJSON FeatureCollection of metro route shapes |
| `GET`  | `/api/metro/stops`  | GeoJSON FeatureCollection of metro stops        |

### Analytics

| Method | Path              | Query params                                                                         | Description                                    |
| ------ | ----------------- | ------------------------------------------------------------------------------------ | ---------------------------------------------- |
| `GET`  | `/api/saturation` | `lat`, `lng`, `radius` (default 500m), `category`                                    | Market saturation score (0–100) for a location |
| `GET`  | `/api/heatmap`    | `min_lat`, `min_lng`, `max_lat`, `max_lng`, `cell_size` (default 0.005°), `category` | Opportunity heatmap tiles for a bounding box   |

**Saturation response:**

```json
{ "competitor_count": 12, "density_per_km2": 5.3, "score": 26.5 }
```

**Heatmap response:** Array of `{ lat, lng, anchor_score, comp_penalty, score }` tiles.

### Probes

| Method | Path      | Description                               |
| ------ | --------- | ----------------------------------------- |
| `GET`  | `/livez`  | Always `200 {"status":"ok"}`              |
| `GET`  | `/readyz` | `200` if DB is reachable, `503` otherwise |

---

## Database Migrations

Migrations live in `migrations/` and are the **source of truth for the entire monorepo** — the `backend` service also points its sqlc schema here.

| #   | Name                  | Description                                                                           |
| --- | --------------------- | ------------------------------------------------------------------------------------- |
| 1   | `init_postgis`        | `CREATE EXTENSION IF NOT EXISTS postgis`                                              |
| 2   | `places`              | `places` table with spatial `location` column + GIST index                            |
| 3   | `gtfs_stops`          | `gtfs_stops` table for metro stops geometry                                           |
| 4   | `opportunity_heatmap` | Materialized view scoring locations 0–100 (60% metro proximity + 40% low competition) |

```bash
make migrate-up    # apply all pending migrations
make migrate-down  # rollback one step
```

---

## Makefile Targets

| Target         | Description                                                             |
| -------------- | ----------------------------------------------------------------------- |
| `build`        | Compile to `bin/geopulse-bff`                                           |
| `run`          | Build and run with `DATABASE_URL` and `PORT` from env                   |
| `test`         | `go test ./...`                                                         |
| `lint`         | `golangci-lint run ./...`                                               |
| `tidy`         | `go mod tidy`                                                           |
| `typespec`     | Compile `typespec/main.tsp` → `api/openapi.yaml`                        |
| `generate`     | Run `typespec` then `oapi-codegen` → `internal/generated/server.gen.go` |
| `sqlc`         | Run `sqlc generate` → `internal/store/`                                 |
| `migrate-up`   | Apply all pending migrations                                            |
| `migrate-down` | Rollback one migration                                                  |
| `docker-up`    | `docker compose up -d --wait` (blocks until postgres is healthy)        |
| `docker-down`  | Stop and remove containers                                              |
| `seed`         | `docker-up` + `migrate-up`                                              |
| `venom`        | Run E2E contract tests against a running server                         |

---

## Project Structure

```
bff/
├── cmd/main.go                  # Entry point — wiring only
├── internal/
│   ├── app/
│   │   ├── app.go               # Service interface + App struct
│   │   ├── places.go            # Place CRUD business logic
│   │   ├── metro.go             # Metro shapes/stops GeoJSON assembly
│   │   ├── analytics.go        # Saturation + heatmap scoring
│   │   ├── helpers.go           # Pointer helpers (pstr, pf64, …)
│   │   └── converters.go        # DB row → generated.Place conversions
│   ├── handler/handler.go       # HTTP handlers (oapi-codegen StrictServer)
│   ├── generated/server.gen.go  # Auto-generated — DO NOT EDIT
│   └── store/
│       ├── queries.sql          # Source of truth for all SQL
│       ├── store.go             # Extended Store interface (embeds Querier)
│       ├── analytics.go         # Raw PostGIS queries (saturation, heatmap)
│       ├── querier.go           # Auto-generated — DO NOT EDIT
│       ├── queries.sql.go       # Auto-generated — DO NOT EDIT
│       ├── models.go            # Auto-generated — DO NOT EDIT
│       └── db.go                # Auto-generated — DO NOT EDIT
├── api/openapi.yaml             # Generated from TypeSpec
├── typespec/main.tsp            # API definition source of truth
├── migrations/                  # Up/down SQL migration files
├── test/suite.yml               # Venom E2E contract tests
├── go.mod
├── Makefile
├── sqlc.yaml
└── oapi-codegen.yaml
```

---

## Code Generation

The generated files (`internal/generated/`, `internal/store/queries.sql.go`, etc.) are committed to the repo. Regenerate only when the API or SQL schema changes:

```bash
make typespec   # TypeSpec → api/openapi.yaml
make generate   # openapi.yaml → internal/generated/server.gen.go
make sqlc       # queries.sql + migrations → internal/store/
```

---

## E2E Tests

Requires a running server (`make run` in a separate terminal) and `venom` on `$PATH`.

```bash
make venom
```

Tests cover: liveness/readiness probes, place CRUD, category filter, metro GeoJSON, saturation score, heatmap tiles.
