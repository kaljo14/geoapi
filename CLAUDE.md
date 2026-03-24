# GeoPulse

Unified service for geospatial place intelligence. Runs on port **8080**.

## Service Role

Serves the Vue frontend with places data, GeoJSON metro shapes/stops, saturation scores, and heatmap tiles. Also handles Google Places scraping, enrichment, and CSV import/export.

## Key Constraints

- **Never SELECT `location` or `geom` columns.** PostGIS returns WKB binary that pgx/v5 can't parse into Go types. Only `lat`/`lng` float columns are safe to read.
- **Spatial writes** use `ST_SetSRID(ST_MakePoint($lng, $lat), 4326)` inline in SQL — never pass geometry as a Go value.
- **`GetSaturation` and `GetHeatmapTiles` bypass sqlc.** These use complex PostGIS CTEs and are executed as raw `pgxpool.QueryRow`/`pgxpool.Query` calls in `internal/store/analytics.go`.
- **GeoJSON assembly** for metro shapes/stops happens in `internal/app/metro.go`, not in SQL.
- **CSV export routes bypass oapi-codegen.** `GET /api/places/export`, `GET /api/places/export-simple`, and `POST /api/places/import` are registered manually on the chi router.
- **Scrape/enrich are fire-and-forget goroutines.** Both endpoints return 202 immediately.

## Bootstrap Wiring

```
pgxpool.New → store.New(pool) → app.New(store, pool, logger) →
handler.NewHandler(app, logger) → generated.NewStrictHandler(h, nil) →
generated.HandlerWithOptions(strict, ChiServerOptions{BaseRouter: r}) →
r.Get/Post (manual CSV routes) → http.Server
```

## Commands

```
make build        # compile
make run          # build + run (needs DATABASE_URL env var)
make test         # go test ./...
make lint         # golangci-lint
make sqlc         # regenerate store package
make generate     # typespec + oapi-codegen
make mock         # regenerate mocks
make migrate-up   # apply migrations
make migrate-down # rollback migrations
make seed         # docker-up + migrate-up
make docker-up    # start postgres+martin via docker-compose.yml
make venom        # E2E contract tests (server must already be running)
```
