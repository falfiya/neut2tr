package main

import (
	"fmt"
	"neutttr/racket"
	"os"
	"github.com/fatih/color"
)

func main() {
	dat, _ := os.ReadFile("./examples/homework_10.rkt")
	complete, comments, diags := racket.Comments(dat)
	fmt.Printf("complete: %v\n", complete)
	for _, d := range diags {
		fmt.Printf("DIAG: %+v\n", d)
	}
	pos := 0
	for _, c := range comments {
		fmt.Print(string(dat[pos:c.Beg]))
		color.HiBlack(string(dat[c.Beg:c.End]))
		if c.IsLineComment {
			fmt.Print(c.GetIndent())
			color.Yellow("; ooga booba")
		}
		pos = c.End
	}
	fmt.Print(string(dat[pos:]))
}
