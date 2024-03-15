package main

import (
	"fmt"
	"neut2tr/neut"

	// "neut2tr/racket"
	// "neut2tr/util"
	// "os"
	"github.com/cockroachdb/errors"
	"github.com/fatih/color"
)

func main() {
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
