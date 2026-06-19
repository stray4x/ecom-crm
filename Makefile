build:
	@go build -o bin/ecom cmd/api/main.go

run: build
	@./bin/ecom

dev:
	@air