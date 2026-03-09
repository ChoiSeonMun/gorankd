.PHONY: proto-gen build run test lint

proto-gen:
	bash scripts/proto-gen.sh

build:
	go build -o bin/gorankd ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./... -v

lint:
	golangci-lint run ./...
