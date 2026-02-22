.PHONY: fmt lint test setup

setup:
	@command -v lefthook >/dev/null || (echo "Install lefthook: brew install lefthook" && exit 1)
	lefthook install

fmt:
	gofumpt -w ./

lint:
	golangci-lint run ./...

test:
	go test ./...
