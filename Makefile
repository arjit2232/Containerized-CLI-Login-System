.PHONY: build run test tidy docker-up docker-run docker-down fmt

build:
	go build -o cli ./cmd/cli

run: build
	./cli

test:
	go test ./...

tidy:
	go mod tidy

fmt:
	go fmt ./...

docker-up:
	docker compose -f docker/docker-compose.yml up --build -d db

docker-run:
	docker compose -f docker/docker-compose.yml run --rm cli

docker-down:
	docker compose -f docker/docker-compose.yml down
