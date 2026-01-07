.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	go build -o branch ./...

.PHONY: install
install:
	go install ./...

.PHONY: install-lint
install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: format
format:
	go fmt ./...

.PHONY: clean