package neut

import (
	"neutttr/lexer"
	"strings"
)

type tokenizerError struct {
	lexer.Pos
	Msg string
}

type token any

type NewlineToken struct {
	lexer.Pos
}

type StringToken struct {
	lexer.Sel
	Literal string
}

// numbers are considered identifiers
type IdentifierToken struct {
	lexer.Sel
	Name string
}

type SymbolToken struct {
	lexer.Pos
	Symbol byte
}

func isControlCharacter(c byte) bool {
	return c <= ' ' || c == 127
}

// https://docs.racket-lang.org/guide/symbols.html
// not allowed: ( ) [ ] { } " , ' ` ; | \
var disallowed = []byte("()[]{}\",'`;|\\")

func identifierAllowed(c byte) bool {
	if isControlCharacter(c) {
		return false
	}
	for _, d := range disallowed {
		if c == d {
			return false
		}
	}
	return true
}

// splits a string into tokens
// delimiters are any control characters, whitespace. this includes newlines
func Tokenize(s string) (tokens []token, te *tokenizerError) {
	lex := lexer.New(s)
tokenLoop:
	for lex.More() {
		beg := lex.Pos()
		c := lex.Next()
		lex.Bump()
		if c == '\n' {
			tokens = append(tokens, NewlineToken{beg})
			continue tokenLoop
		}
		if isControlCharacter(c) {
			continue tokenLoop
		}
		// we're going to push a token
		if c == '"' {
		stringLoop:
			for lex.More() {
				c := lex.Next()
				lex.Bump()
				if c == '\\' {
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
			tokens = append(tokens, StringToken{sel, s[beg.Offset:lex.Offset()]})
			continue tokenLoop
		}
		if c == '#' {
			if lex.More() && lex.Next() == '%' {
				// A hash is permitted in an identifier.
				// But if it's at the start, it must be followed by a percent sign!
				// That is to say that this is #not-allowed but #%this-is-allowed

				// Let the code flow normally into the identifierLoop
			} else {
				// Identifier cannot begin with a hash!
				goto tokenizeSymbol
			}
		}
		if identifierAllowed(c) {
			for {
				if !lex.More() {
					goto commitIdentifier
				}
				c = lex.Next()
				if identifierAllowed(c) {
					lex.Bump()
					continue
				} else {
					goto commitIdentifier
				}
			}
		commitIdentifier:
			text := strings.ToLower(s[beg.Offset:lex.Offset()])
			sel := lexer.Sel{
				Pos:   beg,
				Count: lex.Offset() - beg.Offset,
			}
			tokens = append(tokens, IdentifierToken{sel, text})
			continue tokenLoop
		}
	tokenizeSymbol:
		tokens = append(tokens, SymbolToken{beg, c})
	}
	return
}
