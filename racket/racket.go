package racket

import (
	"neutttr/lexer"
)

type ParseError struct {
	lexer.Pos
	Msg   string
}

type Comment struct {
	lexer.Sel
	IsLineComment bool
}

func Comments(s string) (comments []Comment, pe *ParseError) {
	lex := lexer.New(s)
	for lex.More() {
		beg := lex.Pos()
		c := lex.Next()
		lex.Bump()
		switch c {
		case '"':
			escaped := false
		stringLoop:
			for {
				if lex.Done() {
					pe = new(ParseError)
					pe.Pos = beg
					pe.Msg = "Expected \" but instead saw EOF"
					return
				}
				c := lex.Next()
				lex.Bump()
				if escaped {
					escaped = false
				} else {
					switch c {
					case '\\':
						escaped = true
					case '"':
						break stringLoop
					}
				}
			}
		case '#':
			if lex.More() {
				c := lex.Next()
				lex.Bump()
				switch c {
				case '\\': // #\
					// eat the \ plus one more character
					lex.Bump()
				case '|': // #|
					// just want to point out that yes we could've used
					// <3 nicey wicey recursion <3
					// but a simple counter will suffice
					nestLevel := 1
					for nestLevel > 0 {
						// eat every token that isn't either
						// 	another #|
						// 	a closing |#
						if lex.Done() {
							pe = new(ParseError)
							pe.Pos = lex.Pos()
							pe.Msg = "Expected |# but instead saw EOF"
							return
						}
						c2 := lex.Next()
						lex.Bump()
						switch c2 {
						case '#': // first byte is #
							if lex.More() && lex.Next() == '|' {
								// the second byte is |
								// we have another block comment
								lex.Bump()
								nestLevel += 1
							}
						case '|': // first byte |
							if lex.More() && lex.Next() == '#' {
								// the second byte is #
								// that's a closing block comment
								lex.Bump()
								nestLevel -= 1
							}
						}
					}
					comments = append(comments, Comment{
						IsLineComment: false,
						Sel: beg.Select(lex.Offset()),
					})
				}
			}
		case ';':
			for lex.More() && lex.Next() != '\n' {
				lex.Bump()
			}
			comments = append(comments, Comment{
				IsLineComment: true,
				Sel: beg.Select(lex.Offset()),
			})
		}
	}
	return
}
