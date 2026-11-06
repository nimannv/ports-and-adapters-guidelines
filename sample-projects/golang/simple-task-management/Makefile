.PHONY: test test-integration run-api run-cli tidy docker-build docker-up docker-down

tidy:
	go mod tidy

test:
	go test -v ./internal/core/... ./internal/adapters/inbound/cli/...

test-integration:
	go test -v -tags=integration ./test/integration/...

run-api:
	go run ./bootstrap/api/main.go

run-cli:
	go run ./bootstrap/cli/main.go $(ARGS)

docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down

