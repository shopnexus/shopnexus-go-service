# Set default goal to prevent running `make` without a target
.DEFAULT_GOAL := prevent-naked-make

# Prevent naked `make` (running `make` with no target)
prevent-naked-make:
	@echo "Error: Please specify a make target. See available targets below."
	@$(MAKE) help
	@false

# Help target - Displays available commands
help:
	@echo "Shopnexus Go Service - Available Commands"
	@echo "======================================"
	@echo
	@echo "Development:"
	@echo "  make dev          - Run development server using Air"
	@echo "  make dev2         - Run development server using Nodemon"
	@echo
	@echo "Code Generation:"
	@echo "  make proto        - Generate protobuf code"
	@echo "  make sqlc         - Generate SQL code"
	@echo "  make init-migrate - Initialize database migration"
	@echo "  make generate     - Generate all code (proto + sqlc)"
	@echo
	@echo "Build & Run:"
	@echo "  make build        - Build the binary"
	@echo "  make run          - Build and run the binary"
	@echo
	@echo "Testing:"
	@echo "  make test         - Run tests"
	@echo "  make test-coverage- Run tests with coverage report"
	@echo
	@echo "Code Quality:"
	@echo "  make fmt          - Format Go code"
	@echo "  make vet          - Run Go vet"
	@echo "  make lint         - Run linter"
	@echo
	@echo "Maintenance:"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make deps         - Download dependencies"

# Go related variables
BINARY_NAME=shopnexus-service
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")

# Development
dev:
	air .air.toml -h

dev2:
	nodemon --watch . --ext go --exec "go build ./cmd/main.go && ./main || exit 1" --watch "src"

# Code generation
proto:
	cd proto && \
	protoc --go_out=../gen/pb \
		   --go-grpc_out=../gen/pb \
		   --go_opt=paths=source_relative \
		   --go-grpc_opt=paths=source_relative \
		   ./*.proto && \
	cd ..

sqlc: init-migrate
	sqlc generate

init-migrate:
	npx prisma migrate diff --from-empty --to-schema-datamodel prisma/schema.prisma --script > prisma/migrations/0_init/migration.sql

# Build and run
build:
	go build -o $(BINARY_NAME) ./cmd/main.go

run: build
	./$(BINARY_NAME)

# Testing
test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Cleanup
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f coverage.out

# Dependencies
deps:
	go mod download
	go mod tidy
