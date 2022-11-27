.PHONY: test
test:
	@go test -v -race -timeout 30s ./...

.PHONY: statictest
statictest:
	@go vet -vettool="$(shell which statictest)" ./...