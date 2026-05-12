.PHONY: build test race run tidy

build:
	go build -o bin/engine ./cmd/engine

test:
	go test ./...

race:
	go test -race -count=1 ./...

run:
	go run ./cmd/engine

tidy:
	go mod tidy
