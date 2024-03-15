package main

import (
	"fmt"
	"log"
	"neut2tr/neut"
	"neut2tr/racket"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/fatih/color"
)

func main() {
	testCommentFinder()
	// testParser()
}

func testCommentFinder() {
	dat, err := os.ReadFile("examples/homework_10.rkt")
	file := string(dat)
	if err != nil {
		log.Fatal(err)
	}
	comments, err2 := racket.Comments(file)
	if err2 != nil {
		log.Fatal(err)
	}
	for i, c := range comments {
		fmt.Printf("%d: %s\n", i, file[c.Offset: c.End()])
	}
}

func testParser() {
	tokens, err := neut.Tokenize("A (Listof X) is one of\n-'()\n-(cons X [Listof X])")
	if err != nil {
		errors.Errorf("Error at line %d, char %d:\n%s", err.LineNo, err.CharNo(), err.Msg)
	}
	for _, t := range tokens {
		switch v := t.(type) {
		case neut.NewlineToken:
			fmt.Printf("%s ", color.HiBlackString("\\n"))
		case neut.StringToken:
			fmt.Printf("%s ", color.HiGreenString(v.Literal))
		case neut.IdentifierToken:
			fmt.Printf("%s ", v.Name)
		case neut.SymbolToken:
			fmt.Printf("%s ", color.YellowString(string(v.Symbol)))
		}
	}
	fmt.Print("\n---\n")
	meaningful := neut.Parse(tokens)
	for _, m := range meaningful {
		fmt.Print(m.Print())
	}
}

