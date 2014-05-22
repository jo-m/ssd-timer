all: $(wildcard *.go)
	go build

run: $(wildcard *.go)
	go run $(wildcard *.go) -a localhost:8080 -p abc123
