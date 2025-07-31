# Stage 1: Build
FROM golang:1.24.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

# Stage 2: Runtime
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./app"]
