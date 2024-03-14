package main

import (
	"fmt"
	"neutttr/neut"

	// "neutttr/racket"
	// "neutttr/util"
	// "os"
	"github.com/cockroachdb/errors"
	"github.com/davecgh/go-spew/spew"
	"github.com/fatih/color"
)

func main() {
	tokens, err := neut.Tokenize("fppbar : [hello -> (make-something with \"a string\")]")
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
	spew.Dump(neut.ParseAnnotation(&tokens))
}
