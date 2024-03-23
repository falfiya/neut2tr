package neut

import (
	"neut2tr/lexer"
	"neut2tr/util"
	"strings"
)

type TokenizerError struct {
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
	CmpName string
}

type SymbolToken struct {
	lexer.Pos
	Symbol byte
}

// https://docs.racket-lang.org/guide/symbols.html
// not allowed: ( ) [ ] { } " , ' ` ; | \
var disallowed = []byte("()[]{}\",'`;|\\")

func identifierAllowed(c byte) bool {
	if util.IsControlCharacter(c) {
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
func Tokenize(s string) (tokens []token, te *TokenizerError) {
	lex := lexer.New(s)
tokenLoop:
	for lex.More() {
		beg := lex.Pos()
		c := lex.Current()
		lex.Bump()
		if c == '\n' {
			tokens = append(tokens, NewlineToken{beg})
			continue tokenLoop
		}
		if util.IsControlCharacter(c) {
			continue tokenLoop
		}
		// we're going to push a token
		if c == '"' {
		stringLoop:
			for lex.More() {
				c := lex.Current()
				lex.Bump()
				if c == '\\' {
					if lex.Done() {
						te = &TokenizerError{
							lex.Pos(),
							"Expected character after string escape",
						}
						return
					}
					lex.Bump()
					continue stringLoop
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
			if lex.Done() {
				goto tokenizeSymbol
			}
			// A hash is permitted in an identifier.
			// But if it's at the start, it must be followed by a percent sign!
			// That is to say that this is #not-allowed but #%this-is-allowed
			if lex.Current() != '%' {
				// In this case, we know what we're looking at is not an identifier
				goto tokenizeSymbol
			}
		}
		if identifierAllowed(c) {
			for {
				if lex.Done() {
					goto commitIdentifier
				}
				c = lex.Current()
				if identifierAllowed(c) {
					lex.Bump()
					continue
				} else {
					goto commitIdentifier
				}
			}
		commitIdentifier:
			name := s[beg.Offset:lex.Offset()]
			cmpName := strings.ToLower(name)
			sel := lexer.Sel{
				Pos:   beg,
				Count: lex.Offset() - beg.Offset,
			}
			tokens = append(tokens, IdentifierToken{sel, name, cmpName})
			continue tokenLoop
		}
	tokenizeSymbol:
		tokens = append(tokens, SymbolToken{beg, c})
	}
	return
}
