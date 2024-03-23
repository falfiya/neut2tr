run:
	go run ./cmd

build:
	go build -o bin/neut2tr.exe ./cmd

build-wasm: export GOOS=js
build-wasm: export GOARCH=wasm
build-wasm:
	go build -o bin/neut2tr.wasm ./wasm
