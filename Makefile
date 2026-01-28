APP_NAME=worklayer

.PHONY: run build

# Run the API locally
run:
	go run cmd/server/main.go

# Build the binary
build:
	go build -o bin/$(APP_NAME) cmd/server/main.go
