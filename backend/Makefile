.PHONY: build run test lint clean tidy generate typespec sqlc docker-up docker-down venom

SERVICE_NAME := geopulse-backend
DATABASE_URL ?= postgres://geopulse:geopulse@localhost:5432/geopulse?sslmode=disable

build:
	go build -o bin/$(SERVICE_NAME) ./cmd/

run: build
	DATABASE_URL=$(DATABASE_URL) PORT=8081 ./bin/$(SERVICE_NAME)

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

typespec:
	cd typespec && ./node_modules/.bin/tsp compile .

sqlc:
	sqlc generate

# --- Docker ---

docker-up:
	docker compose -f ../docker-compose.yml up -d --wait

docker-down:
	docker compose -f ../docker-compose.yml down

# --- Tests ---

venom:
	rm -f venom.log venom.*.log
	venom run test/suite.yml
