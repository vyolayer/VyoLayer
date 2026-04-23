FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-s -w' -o /out/apikeys-service ./cmd/apikey

FROM alpine:3.22

WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /out/apikeys-service /app/apikeys-service

EXPOSE 50055

CMD ["/app/apikeys-service"]
