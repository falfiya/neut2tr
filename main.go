package main

import (
	"fmt"
	"neutttr/racket"
	"os"
	"github.com/fatih/color"
)

func main() {
	dat, _ := os.ReadFile("./racket/awful.rkt")
	complete, comments, diags := racket.Comments(dat)
	fmt.Printf("complete: %v\n", complete)
	pos := 0
	for _, c := range comments {
		fmt.Print(string(dat[pos:c.Beg]))
		color.HiBlack(string(dat[c.Beg:c.End]))
		pos = c.End
	}
	for _, d := range diags {
		fmt.Printf("DIAG: %+v\n", d)
	}
}
