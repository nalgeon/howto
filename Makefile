.PHONY: build
build:
	@go build -o howto

.PHONY: test
test:
	@go test ./...

.PHONY: lint
lint:
	@go vet ./...
	@golangci-lint run --print-issued-lines=false --out-format=colored-line-number ./...
