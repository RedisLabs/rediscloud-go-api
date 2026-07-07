.DEFAULT_GOAL := ci
.PHONY: fmt lint test ci

fmt:
	golangci-lint fmt

lint:
	golangci-lint run

test:
	go test -v ./...

ci: lint test
