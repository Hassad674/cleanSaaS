.PHONY: run build test migrate seed mock

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

test:
	go test ./... -v

test-unit:
	go test ./internal/domain/... ./internal/app/... ./pkg/... -v

test-integration:
	go test ./test/... -v

migrate:
	go run cmd/migrate/main.go

seed:
	go run cmd/seed/main.go

mock:
	@echo "Generate mocks from port interfaces"

lint:
	golangci-lint run ./...

tidy:
	go mod tidy
