package main

import (
	"fmt"
	"neutttr/racket"
	"neutttr/util"
	"os"

	"github.com/fatih/color"
)

func main() {
	bytes, _ := os.ReadFile("./examples/awful.rkt")
	src := string(bytes)
	comments, err := racket.Comments(src)
	if err != nil {
		fmt.Printf("Error: %+v\n", *err)
	}
	pos := 0
	for _, c := range comments {
		fmt.Print(src[pos:c.Offset])
		pos = c.Offset
		color.HiBlack(src[c.Offset:c.Offset + c.Count])
		pos += c.Count
		if c.IsLineComment {
			fmt.Print(util.IndentOnly(c.StartOfLine(src)))
			color.Yellow("; ^ that is a comment")
		}
	}
	fmt.Print(src[pos:])
}
