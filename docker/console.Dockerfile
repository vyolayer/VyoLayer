FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-s -w' -o /out/console-service ./cmd/console

FROM alpine:3.22

WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /out/console-service /app/console-service

EXPOSE 50054

CMD ["/app/console-service"]
