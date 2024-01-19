package main

import (
	"fmt"
	"neutttr/racket"
	"os"

	// "github.com/davecgh/go-spew/spew"
)

func main() {
	dat, _ := os.ReadFile("./example/homework_10.rkt")
	complete, ts, diags := racket.Tokenize(dat)
	fmt.Printf("complete: %v\n", complete)
	// spew.Dump(ts)
	for _, d := range diags {
		fmt.Printf("DIAG: %+v\n", d)
	}
	fmt.Print(racket.Untokenize(ts))
}
