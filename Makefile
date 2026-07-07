.DEFAULT_GOAL := ci
.PHONY: fmt lint vulncheck test ci

fmt:
	golangci-lint fmt

lint:
	golangci-lint run

vulncheck:
	govulncheck ./...

test:
	go test -v ./...

ci: vulncheck lint test
