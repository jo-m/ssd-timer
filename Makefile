.PHONY: all run lint check format

all: $(wildcard *.go)
	go build

run: $(wildcard *.go)
	go run $(wildcard *.go) -a localhost:8080 -p abc123

lint:
	gofmt -l .; test -z "$$(gofmt -l .)"
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

check: lint test

format:
	gofmt -w -s .
