# Simple variables
APP_NAME=vyolayer
DB_URL=postgres://vyolayer_user:vyolayer_password@localhost:4444/vyolayer_db?sslmode=disable

.PHONY: run build docs docker-up docker-start docker-stop docker-down migrate seed

# Run the API locally
run:
	go run cmd/server/main.go	

# Build the binary
build:
	go build -o bin/$(APP_NAME) cmd/server/main.go

# Docs generate
docs:
	@echo "Installing swag if needed..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Generating Swagger documentation..."
	@$(shell go env GOPATH)/bin/swag init -g internal/app/server.go


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
