vendor:
	go mod vendor
	go mod tidy

check:
	golangci-lint run -E gofmt ./...

build:
	go build -mod=vendor

install:
	go install -mod=vendor
