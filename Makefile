run:
	go run .

build:
	go build -o neut2tr.exe .

.PHONY: examples
examples: build
	neut2tr examples/examples.rkt examples/examples.t.rkt

bin: bin-windows-amd64 bin-darwin-arm64 bin-darwin-amd64 bin-linux-amd64
	-

bin-windows-amd64: export GOOS=windows
bin-windows-amd64: export GOARCH=amd64
bin-windows-amd64:
	go build -o bin/neut2tr-windows-x64.exe

bin-darwin-arm64: export GOOS=darwin
bin-darwin-arm64: export GOARCH=arm64
bin-darwin-arm64:
	go build -o bin/neut2tr-darwin-arm64

bin-darwin-amd64: export GOOS=darwin
bin-darwin-amd64: export GOARCH=amd64
bin-darwin-amd64:
	go build -o bin/neut2tr-darwin-x64

bin-linux-amd64: export GOOS=linux
bin-linux-amd64: export GOARCH=amd64
bin-linux-amd64:
	go build -o bin/neut2tr-linux-x64
