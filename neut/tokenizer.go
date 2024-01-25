package neut

import (
	"neutttr/lexer"
	"neutttr/util/ascii"
	"strings"
)

type token struct {
	lexer.Sel
	text string
}
func (t token) isWord() bool {
	return t.Count > 1
}

func tokenize(s string) (tokens []token) {
	lex := lexer.New(s)
	for lex.More() {
		beg := lex.Pos()
		c := lex.Pop()
		if ascii.IsAlphanumericOr_(c) {
			for lex.More() && ascii.IsAlphanumericOr_(lex.Next()) {
				lex.Bump()
			}
			sel := lexer.Sel {
				Pos: beg,
				Count: lex.Offset() - beg.Offset,
			}
			tokens = append(tokens, token{sel, strings.ToLower(sel.From(s))})
		} else if c == '\n' || !ascii.IsControlCharacter(c) {
			sel := lexer.Sel{
				Pos: lex.Pos(),
				Count: 1,
			}
			tokens = append(tokens, token{sel, strings.ToLower(sel.From(s))})
		}
	}
	return
}
