.PHONY: dev up down db-gen
.PHONY: test test-clean

up:
	docker-compose -f deploy/docker-compose.yml up -d

down:
	docker-compose -f deploy/docker-compose.yml down

db-gen:
	sqlc generate

run-api:
	go run main.go

run-web:
	cd web && npm run dev

test:
	@echo "Running all tests..."
	go test -race -v ./...
