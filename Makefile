PHONY: setup lint fmt test build

setup:
	go get golang.org/x/tools/cmd/goimports

lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.42.0 golangci-lint run -E golint,gofmt,goimports

fmt:
	go fmt ./...
	goimports -l -w .

test:
	go test ./...

build:
	go build -ldflags='-s -w -X github.com/kobtea/mixcli/cmd.Version=$(shell cat VERSION)'
