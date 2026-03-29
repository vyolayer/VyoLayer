# Simple variables
APP_NAME=vyolayer
DB_URL=postgres://vyolayer_user:vyolayer_password@localhost:4444/vyolayer_db?sslmode=disable
AIR_BIN?=air

.PHONY: run run-all air-install dev-gateway dev-account dev-all build docs docker-up docker-start docker-stop docker-down migrate seed

# Run the API locally
run:
	go run cmd/server/main.go	

# Run contionusly
run-all:
	@echo "Starting services..."
	@go run cmd/account-service/main.go &
	@go run cmd/gateway/main.go

# Install air for live reload
air-install:
	@go install github.com/air-verse/air@latest

# Run gateway with live reload
dev-gateway:
	@$(AIR_BIN) -c internal/gateway/.air.gateway.toml

# Run account service with live reload
dev-account:
	@$(AIR_BIN) -c internal/account/.air.account.toml

# Run iam service with live reload
dev-iam:
	@$(AIR_BIN) -c internal/iam/.air.iam.toml

# Run gateway + account, iam service with live reload
dev-all:
	@echo "Starting services with Air..."
	@trap 'kill 0' INT TERM EXIT; \
	$(AIR_BIN) -c internal/account/.air.account.toml & \
	$(AIR_BIN) -c internal/gateway/.air.gateway.toml & \
	$(AIR_BIN) -c internal/iam/.air.iam.toml & \
	wait

# Build the binary
build:
	go build -o bin/$(APP_NAME) cmd/server/main.go

# Docs generate
docs:
	@echo "Installing swag if needed..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Generating Swagger documentation..."
	@$(shell go env GOPATH)/bin/swag init -g internal/app/server.go

# Gateway docs generate	
docs-gateway:
	@echo "Installing swag if needed..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Generating Gateway Swagger docs in docs/gateway..."
	@mkdir -p docs/gateway
	@$(shell go env GOPATH)/bin/swag init \
		-g internal/gateway/server/server.go \
		-o docs/gateway \
		--parseInternal \
		--exclude internal/app

# Start the database container
docker-up:
	docker-compose -f docker/docker-compose.dev.yml up -d

# Start the database container
docker-start:
	docker-compose -f docker/docker-compose.dev.yml start

# Stop the database container
docker-stop:
	docker-compose -f docker/docker-compose.dev.yml stop

# Stop the database container
docker-down:
	docker-compose -f docker/docker-compose.dev.yml down

# Migrate database
migrate:
	go run cmd/migrate/main.go

# Seed database
seed:
	go run cmd/seed/main.go
