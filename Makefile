.PHONY: test
test:
	@go test -v -race -timeout 30s ./...

.PHONY: statictest
lint:
	@golangci-lint run --no-config --disable-all -E govet