.PHONY: build test release clean fmt
build: bin/gogoferd bin/gophermap-example bin/readme-example

release: fmt test bin/gogoferd

bin/gogoferd: .bin-stamp *.go gogoferd/*.go
    go build -o bin/gogoferd ./gogoferd

bin/gophermap-example: .bin-stamp *.go ./cmd/gophermap-examplee/*.go
    go build -o bin/readme-example ./cmd/gophermap-example

bin/readme-example: .bin-stamp *.go ./cmd/readme-example/*.go
    go build -o bin/readme-example ./cmd/readme-example

.bin-stamp:
    mkdir bin
    touch .bin-stamp

fmt:
    go fmt ./...

test:
    go test ./...

clean:
    rm -f .*-stamp
    rm -rf ./bin