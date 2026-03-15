.PHONY: run build test migrate-up migrate-down migrate-down-all migrate-status migrate-force seed mock lint tidy

run:
	set -a && [ -f .env ] && . ./.env; set +a; go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

test:
	go test ./... -v

test-unit:
	go test ./internal/domain/... ./internal/app/... ./pkg/... -v

test-integration:
	go test ./test/... -v

# Migrations
migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down

migrate-down-all:
	go run cmd/migrate/main.go down-all

migrate-status:
	go run cmd/migrate/main.go status

migrate-force:
	go run cmd/migrate/main.go force $(VERSION)

# Seed
seed:
	go run cmd/seed/main.go

# Tools
mock:
	@echo "Generate mocks from port interfaces"

lint:
	golangci-lint run ./...

tidy:
	go mod tidy
