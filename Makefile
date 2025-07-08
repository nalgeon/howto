all: build lint vet test

build:
	@go build -o howto
	@echo "✓ build"

lint:
	@golangci-lint run ./...
	@echo "✓ lint"

test:
	@go test ./...
	@echo "✓ test"

vet:
	@go vet ./...
	@echo "✓ vet"
