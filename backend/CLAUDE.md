# GeoPulse Backend

Admin and data-seeding service. Runs on port **8081**.

## Service Role

Handles Google Places scraping, enrichment, and CSV import/export. Not frontend-facing — the BFF (`../bff/`) serves the Vue frontend.

## Key Constraints

- **No migrations here.** `sqlc.yaml` points `schema:` at `../bff/migrations/`. All schema changes go in `bff/`.
- **Never SELECT `location` or `geom` columns.** PostGIS returns WKB binary that pgx/v5 can't parse into Go types. Only `lat`/`lng` float columns are safe to read.
- **CSV export routes bypass oapi-codegen.** `GET /api/places/export` and `GET /api/places/export-simple` are registered manually on the chi router (after `generated.HandlerWithOptions()`) because they stream directly to `http.ResponseWriter`.
- **Scrape/enrich are fire-and-forget goroutines.** Both endpoints return 202 immediately. Errors from missing `GOOGLE_API_KEY` are logged as `Warn`, not `Error`.

## Bootstrap Wiring

```
pgxpool.New → store.New(pool) → app.New(querier, pinger, logger) →
handler.NewHandler(app, logger) → generated.NewStrictHandler(h, nil) →
generated.HandlerWithOptions(strict, ...) → http.Server
```

## Commands

```
make build       # compile
make run         # build + run (needs DATABASE_URL env var)
make test        # go test ./...
make lint        # golangci-lint
make sqlc        # regenerate store package (schema from ../bff/migrations/)
make generate    # typespec + oapi-codegen
make docker-up   # start postgres+martin via ../docker-compose.yml
make venom       # E2E contract tests (server must already be running)
```
