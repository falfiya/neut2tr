package racket

import (
	"neut2tr/lexer"
)

type CommentType int

const (
	CtLine CommentType = iota
	CtMultiLine
	CtBlock
)

type ParseError struct {
	lexer.Pos
	Msg string
}

type Comment struct {
	lexer.Sel
	Ct   CommentType
	Text string
}

func ExtractComments(s string) (comments []Comment, pe *ParseError) {
	lex := lexer.New(s)
	for lex.More() {
		beg := lex.Pos()
		c := lex.Current()
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
				c := lex.Current()
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
				c := lex.Current()
				lex.Bump()
				switch c {
				case '\\': // #\
					// eat the \ plus one more character
					lex.Bump()
				case '|': // #|
					textStart := lex.Offset()
					nestLevel := 1
					for nestLevel > 0 {
						// eat every token that isn't either
						// 	another #|
						// 	a closing |#
						potentialError := &ParseError{
							Pos: lex.Pos(),
							Msg: "Unexpected end of block comment!",
						}

						if lex.Done() {
							pe = potentialError
							return
						}
						c2 := lex.Current()
						lex.Bump()

						if lex.Done() {
							pe = potentialError
							return
						}
						c3 := lex.Current()
						switch c2 {
						case '#': // first byte is #
							if c3 == '|' {
								// the second byte is |
								// we have another block comment
								lex.Bump()
								nestLevel += 1
							}
						case '|': // first byte |
							if c3 == '#' {
								// the second byte is #
								// that's a closing block comment
								lex.Bump()
								nestLevel -= 1
							}
						}
					}
					// #| blah blah |# (some code)
					//                ^ lex.Offset()
					//              ^ desired textEnd
					textEnd := lex.Offset() - 2
					text := s[textStart:textEnd]
					comments = append(comments, Comment{beg.Select(textEnd), CtBlock, text})
				}
			}
		case ';':
			textStart := lex.Offset()
			ct := CtLine
			workingText := ""
			for lex.More() {
				if lex.Current() == '\n' {
					lex.Bump()
					workingText += s[textStart:lex.Offset()]
					if lex.Done() {
						goto commitComment
					}
					if lex.Current() != ';' {
						goto commitComment
					}
					ct = CtMultiLine
					lex.Bump()
					textStart = lex.Offset()
				} else {
					lex.Bump()
				}
			}
		commitComment:
			comments = append(comments, Comment{beg.Select(lex.Offset()), ct, workingText})
		}
	}
	return
}
