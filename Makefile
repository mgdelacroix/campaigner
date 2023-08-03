.PHONY: vendor

vendor:
	go mod vendor
	go mod tidy

check:
	golangci-lint run -E gofmt ./...

build:
	go build

install:
	go install
