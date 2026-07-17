.DEFAULT_GOAL := ci
.PHONY: fmt lint vulncheck tidy test ci

fmt:
	golangci-lint fmt

lint:
	golangci-lint run

vulncheck:
	govulncheck ./...

tidy:
	go mod tidy -diff

test:
	go test -v ./...

ci: vulncheck tidy lint test
