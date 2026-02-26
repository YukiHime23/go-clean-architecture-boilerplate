.PHONY: run build test tidy docker-up docker-down migrate

APP_NAME=clean-anti
CMD_PATH=./cmd/api

## run: Run the application with live env
run:
	@go run $(CMD_PATH)/main.go

## build: Compile the application binary
build:
	@echo "Building $(APP_NAME)..."
	@go build -o bin/$(APP_NAME) $(CMD_PATH)/main.go
	@echo "Build complete: bin/$(APP_NAME)"

## test: Run all tests
test:
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

## tidy: Tidy and download Go modules
tidy:
	@go mod tidy

## docker-up: Start Docker Compose services
docker-up:
	@docker-compose up -d

## docker-down: Stop Docker Compose services
docker-down:
	@docker-compose down

## help: Print this help message
help:
	@echo "Usage: make [target]"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
