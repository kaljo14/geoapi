# GeoPulse BFF

Frontend-facing service. Runs on port **8080**.

## Service Role

Serves the Vue frontend with places data, GeoJSON metro shapes/stops, saturation scores, and heatmap tiles. The admin/seeding service lives at `../backend/`.

## Key Constraints

- **Owns all migrations.** `migrations/` here is the schema source of truth — `../backend/sqlc.yaml` also points here.
- **Never SELECT `location` or `geom` columns.** PostGIS returns WKB binary that pgx/v5 can't parse into Go types. Only `lat`/`lng` float columns are safe to read.
- **Spatial writes** use `ST_SetSRID(ST_MakePoint($lng, $lat), 4326)` inline in SQL — never pass geometry as a Go value.
- **`GetSaturation` and `GetHeatmapTiles` bypass sqlc.** These use complex PostGIS CTEs and are executed as raw `pgxpool.QueryRow`/`pgxpool.Query` calls directly in `app.go`.
- **GeoJSON assembly** for metro shapes/stops happens in `app.go`, not in SQL.

## Bootstrap Wiring

```
pgxpool.New → store.New(pool) → app.New(pool, querier, logger) →
handler.NewHandler(app, logger) → generated.NewStrictHandler(h, nil) →
generated.HandlerWithOptions(strict, ...) → http.Server
```

## Commands

```
make build        # compile
make run          # build + run (needs DATABASE_URL env var)
make test         # go test ./...
make lint         # golangci-lint
make sqlc         # regenerate store package
make generate     # typespec + oapi-codegen
make migrate-up   # apply migrations
make migrate-down # rollback migrations
make seed         # docker-up + migrate-up
make docker-up    # start postgres+martin via ../docker-compose.yml
make venom        # E2E contract tests (server must already be running)
```
