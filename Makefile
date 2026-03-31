.PHONY: build run stop restart migrate-up migrate-down martin-restart \
        test lint clean tidy generate mock typespec sqlc \
        docker-up docker-down seed venom

SERVICE_NAME  := geopulse
DATABASE_URL  ?= postgres://geopulse:geopulse@localhost:5432/geopulse?sslmode=disable

# Same URL but with host.docker.internal so the migrate container can reach
# the postgres container that is port-forwarded to the Mac host.
MIGRATE_DB_URL := postgres://geopulse:geopulse@host.docker.internal:5432/geopulse?sslmode=disable
MIGRATE_IMAGE  := migrate/migrate:v4.18.1

# ---------------------------------------------------------------------------
# Build & Run
# ---------------------------------------------------------------------------

build:
	go build -o bin/$(SERVICE_NAME) ./cmd/

# Kill any previous instance on :8080 before starting, so you never get
# "address already in use" after closing a terminal.
run: build
	@lsof -ti :8080 | xargs kill -9 2>/dev/null || true
	DATABASE_URL=$(DATABASE_URL) PORT=8080 ./bin/$(SERVICE_NAME)

stop:
	@lsof -ti :8080 | xargs kill -9 2>/dev/null && echo "server stopped" || echo "nothing running on :8080"

restart: stop run

# ---------------------------------------------------------------------------
# Tests / Lint
# ---------------------------------------------------------------------------

test:
	go test ./...

lint:
	golangci-lint run ./...

# ---------------------------------------------------------------------------
# Code generation
# ---------------------------------------------------------------------------

clean:
	rm -rf bin/

tidy:
	go mod tidy

generate: typespec
	oapi-codegen -config oapi-codegen.yaml api/openapi.yaml

mock:
	mockery

typespec:
	cd typespec && ./node_modules/.bin/tsp compile .

sqlc:
	sqlc generate

# ---------------------------------------------------------------------------
# Database migrations  (no local `migrate` binary needed — runs in Docker)
# ---------------------------------------------------------------------------

migrate-up:
	docker run --rm \
		-v "$(CURDIR)/migrations:/migrations" \
		$(MIGRATE_IMAGE) \
		-path=/migrations \
		-database "$(MIGRATE_DB_URL)" \
		up

migrate-down:
	docker run --rm \
		-v "$(CURDIR)/migrations:/migrations" \
		$(MIGRATE_IMAGE) \
		-path=/migrations \
		-database "$(MIGRATE_DB_URL)" \
		down 1

# Restart Martin so it re-discovers any new tables/views created by migrations.
martin-restart:
	docker compose restart martin

# ---------------------------------------------------------------------------
# Docker
# ---------------------------------------------------------------------------

docker-up:
	docker compose up -d --wait

docker-down:
	docker compose down

# Full fresh start: bring up postgres+martin, apply all migrations.
seed: docker-up migrate-up martin-restart

# ---------------------------------------------------------------------------
# E2E tests
# ---------------------------------------------------------------------------

venom:
	rm -f venom.log venom.*.log
	venom run test/suite.yml
