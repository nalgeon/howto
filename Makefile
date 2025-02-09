BUILD_TAG := $(shell git describe --tags --exact-match 2> /dev/null || git rev-parse --short HEAD)

.PHONY: build
build:
	@go build -ldflags "-X main.version=$(BUILD_TAG)" -o howto

.PHONY: test
test:
	@go test ./...

.PHONY: lint
lint:
	@go vet ./...
	@golangci-lint run --print-issued-lines=false --out-format=colored-line-number ./...
