package neut

import (
	"neutttr/lexer"
	"strings"
)

type tokenType int

const (
	ttNewline    tokenType = iota
	ttString     tokenType = iota
	ttIdentifier tokenType = iota
)

type tokenizerError struct {
	lexer.Pos
	msg string
}

type token any

type newlineToken struct {
	lexer.Pos
}

type stringToken struct {
	lexer.Sel
}

type identifierToken struct {
	lexer.Sel
	name string
}

type symbolToken struct {
	lexer.Pos
	symbol byte
}

func isControlCharacter(c byte) bool {
	return c <= ' ' || c == 127
}

// https://docs.racket-lang.org/guide/symbols.html
// not allowed: ( ) [ ] { } " , ' ` ; | \
var disallowed = []byte("()[]{}\",'`;|\\")

func identifierDisallowed(c byte) bool {
	for _, d := range disallowed {
		if c == d {
			return false
		}
	}
	return true
}

func identifierAllowed(c byte) bool {
	return !isControlCharacter(c) && !identifierDisallowed(c)
}

// splits a string into tokens
// delimiters are any control characters, whitespace. this includes newlines
func tokenize(s string) (tokens []token, te *tokenizerError) {
	lex := lexer.New(s)
tokenLoop:
	for lex.More() {
		beg := lex.Pos()
		c := lex.Pop()
		if c == '\n' {
			tokens = append(tokens, newlineToken{beg})
			continue
		}
		if isControlCharacter(c) {
			continue
		}
		// we're going to push a token
		if c == '"' {
			escaped := false
		stringLoop:
			for lex.More() {
				c := lex.Pop()
				if escaped {
					if lex.More() {
						lex.Bump()
						continue stringLoop
					} else {
						te = &tokenizerError{
							lex.Pos(),
							"Expected character after string escape",
						}
						return
					}
				}
				if c == '"' {
					break stringLoop
				}
			}
			sel := lexer.Sel{
				Pos:   beg,
				Count: lex.Offset() - beg.Offset,
			}
			tokens = append(tokens, stringToken{sel})
		}
		if c == '#' {
			if !lex.More() {
				goto tokenizeSymbol
			}
			if lex.Next() != '%' {
				goto tokenizeSymbol
			}
			lex.Bump()
			// a hash is permitted in an identifier unless it's
			// - at the start
			// - AND not followed by a percent sign
			// that is to say that #notallowed
			// but #%thisisallowed
		}
		{
		identifierLoop:
			for {
				if !identifierAllowed(c) {
					break identifierLoop
				}
				if !lex.More() {
					break identifierLoop
				}
				c = lex.Pop()
			}
			sel := lexer.Sel{
				Pos:   beg,
				Count: lex.Offset() - beg.Offset,
			}
			text := strings.ToLower(sel.From(s))
			tokens = append(tokens, identifierToken{sel, text})
			continue tokenLoop
		}
	tokenizeSymbol:
		tokens = append(tokens, symbolToken{beg, c})
	}
	return
}
