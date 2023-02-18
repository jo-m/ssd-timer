.PHONY: all run test lint check format bench

all: $(wildcard *.go)
	go build

run: $(wildcard *.go)
	go run $(wildcard *.go) -a localhost:8080 -p abc123

test:
	# Tags bounds,noasm,safe are added for safety checks in Gonum.
	go test -count 1 -tags bounds,noasm,safe -race -v ./...

lint:
	gofmt -l .; test -z "$$(gofmt -l .)"
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

check: lint test

format:
	gofmt -w -s .

bench:
	go test -v -bench=. ./...
