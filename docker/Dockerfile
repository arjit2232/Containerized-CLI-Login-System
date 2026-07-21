# ---- build stage ----
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o cli ./cmd/cli

# ---- run stage ----
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/cli .
COPY migrations ./migrations
ENTRYPOINT ["./cli"]
