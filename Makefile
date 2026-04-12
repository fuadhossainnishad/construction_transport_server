run:
	go run ./cmd/api

dev:
	go run ./cmd/api

build:
	go build -o bin/app ./cmd/api

test:
	go test ./...

fmt:
	go fmt ./...