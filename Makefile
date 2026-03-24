.PHONY: build run test lint clean tidy generate mock typespec sqlc migrate-up migrate-down docker-up docker-down seed venom

SERVICE_NAME := geopulse
DATABASE_URL ?= postgres://geopulse:geopulse@localhost:5432/geopulse?sslmode=disable

build:
	go build -o bin/$(SERVICE_NAME) ./cmd/

run: build
	DATABASE_URL=$(DATABASE_URL) PORT=8080 ./bin/$(SERVICE_NAME)

test:
	go test ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/

tidy:
	go mod tidy

# --- Code generation ---

generate: typespec
	oapi-codegen -config oapi-codegen.yaml api/openapi.yaml

mock:
	mockery

typespec:
	cd typespec && ./node_modules/.bin/tsp compile .

sqlc:
	sqlc generate

# --- Database ---

migrate-up:
	migrate -path migrations -database "$$DATABASE_URL" up

migrate-down:
	migrate -path migrations -database "$$DATABASE_URL" down

# --- Docker ---

docker-up:
	docker compose up -d --wait

docker-down:
	docker compose down

seed: docker-up
	DATABASE_URL=$(DATABASE_URL) $(MAKE) migrate-up

# --- Tests ---

venom:
	rm -f venom.log venom.*.log
	venom run test/suite.yml
