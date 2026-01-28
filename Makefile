# Simple variables
APP_NAME=worklayer
DB_URL=postgres://worklayer_user:worklayer_password@localhost:4444/worklayer_db?sslmode=disable

.PHONY: run build docker-up docker-start docker-stop docker-down migrate

# Run the API locally
run:
	go run cmd/server/main.go

# Build the binary
build:
	go build -o bin/$(APP_NAME) cmd/server/main.go

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

migrate:
	go run cmd/migrate/main.go