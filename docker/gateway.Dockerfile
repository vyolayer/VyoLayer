FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-s -w' -o /out/gateway ./cmd/gateway

FROM alpine:3.22

WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /out/gateway /app/gateway

EXPOSE 8080

CMD ["/app/gateway"]
