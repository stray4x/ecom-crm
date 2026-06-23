build:
	@go build -o bin/ecom cmd/api/main.go

run: build
	@./bin/ecom

dev:
	@air

migration:
	@migrate create -ext sql -dir migrations -seq $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

test:
	go test ./internal/...

test-verbose:
	go test -v ./internal/...

mocks:
	@mockery

%:
	@: