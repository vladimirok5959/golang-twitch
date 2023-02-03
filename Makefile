default: test

clean:
	go clean -testcache ./...

test:
	go test ./...

lint:
	golangci-lint run --disable=structcheck

tidy:
	go mod tidy

build:
	go build -o bin/cli cmd/cli/main.go

.PHONY: default clean test lint tidy build
