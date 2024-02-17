package neut

import (
	"neutttr/lexer"
	"strings"
)

type token struct {
	lexer.Sel
	text string
}

func isDelimiter(c byte) bool {
	return c <= ' ' || c == 127
}

// splits a string into tokens
// delimiters are any control characters, whitespace. this includes newlines
func tokenize(s string) (tokens []token) {
	panic("This doesn't split things correctly. {} should not be one token!")
	lex := lexer.New(s)
	for lex.More() {
		beg := lex.Pos()
		c := lex.Pop()
		if isDelimiter(c) {
			continue
		}
		// we're going to push a token
		tokenLoop: for lex.More() {
			c := lex.Pop()
			if !isDelimiter(c) {
				continue tokenLoop
			}
			// ignore c
			sel := lexer.Sel {
				Pos: beg,
				Count: lex.Offset() - beg.Offset,
			}
			text := strings.ToLower(sel.From(s))
			tokens = append(tokens, token{sel, text})
			break
		}
	}
	return
}
