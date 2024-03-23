run:
	go run ./cmd

build:
	go build -o bin/neut2tr.exe ./cmd

.PHONY: examples
examples: build
	bin/neut2tr examples/examples.rkt examples/examples.t.rkt

build-web: export GOOS=js
build-web: export GOARCH=wasm
build-web:
	go build -o bin/neut2tr.wasm ./web
