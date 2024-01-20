package main

import (
	"fmt"
	"neutttr/racket"
	"neutttr/util"
	"os"

	"github.com/fatih/color"
)

func main() {
	dat, _ := os.ReadFile("./examples/awful.rkt")
	comments, err := racket.Comments(dat)
	if err != nil {
		fmt.Printf("Error: %+v\n", *err)
	}
	pos := 0
	for _, c := range comments {
		fmt.Print(string(dat[pos:c.Offset]))
		pos = c.Offset
		color.HiBlack(string(dat[c.Offset:c.Offset + c.Count]))
		pos += c.Count
		if c.IsLineComment {
			fmt.Print(string(util.IndentOnly(c.StartOfLine(dat))))
			color.Yellow("; ^ that is a comment")
		}
	}
	fmt.Print(string(dat[pos:]))
}
