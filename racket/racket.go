package racket

import (
	"strings"
)

//go:generate stringer -type=TokenType
type TokenType int

const (
	ttString TokenType = iota
	ttPound
	ttPipe
	ttQuote
	ttSemi
	ttEscape
	ttParenL
	ttParenR
	ttBracketL
	ttBracketR
	ttBraceL
	ttBraceR
	ttLf
	ttWhiteSpace
	ttOther
)

type Pos = int
type Token struct {
	beg Pos
	end Pos
	typ TokenType
}

//go:generate stringer -type=DiagnosticLevel
type DiagnosticLevel int

const (
	DlInfo DiagnosticLevel = iota
	DlWarn
	DlFatal
)

type Diagnostic struct {
	Level DiagnosticLevel
	Beg   Pos
	End   Pos
	Msg   string
}

func isControlChar(b byte) bool {
	return b < '!' || b == 127
}

func Tokenize(s []byte) (good bool, out []Token, diags []Diagnostic) {
	// i is the current position and also the new "end" when committing tokens
	var i Pos
tokenLoop:
	for i < len(s) {
		beg := i
		var typ TokenType
		if s[i] == '"' {
			i += 1
			typ = ttString
			escaped := false
			for {
				if !(i < len(s)) {
					diags = append(diags, Diagnostic{
						Level: DlFatal,
						Beg:   beg,
						End:   i,
						Msg:   "Expected \" but instead saw EOF",
					})
					return
				}
				if escaped {
					escaped = false
				} else {
					switch s[i] {
					case '\\':
						escaped = true
					case '"':
						i += 1
						goto commitToken // <--- this is the only exit path
					}
				}
				i += 1
			}
		}
		switch s[i] {
		case '#':
			typ = ttPound
		case '\'':
			typ = ttQuote
		case ';':
			typ = ttSemi
		case '\\':
			typ = ttEscape
		case '(':
			typ = ttParenL
		case ')':
			typ = ttParenR
		case '[':
			typ = ttBracketL
		case ']':
			typ = ttBracketR
		case '{':
			typ = ttBraceL
		case '|':
			typ = ttPipe
		case '}':
			typ = ttBraceR
		case '\n':
			typ = ttLf
		case '	', ' ':
			typ = ttWhiteSpace
		default:
			// control characters that weren't picked up before this are discarded
			if isControlChar(s[i]) {
				i += 1
				continue tokenLoop
			}
			typ = ttOther
		}
		i += 1
	commitToken:
		// there are a few types of tokens we'd like to combine
		// - ttIndent
		// - ttOther
		if typ == ttOther || typ == ttWhiteSpace {
			// if this token's type is the same as the previous one,
			// just combine them
			lastIdx := len(out) - 1
			if 0 <= lastIdx {
				last := out[lastIdx]
				if last.typ == typ {
					last.end = i
					continue
				}
			}
		}

		out = append(out, Token{beg, i, typ})
	}
	good = true
	return
}

func Untokenize(ts []Token) string {
	var sb strings.Builder
	for _, token := range ts {
		switch token.typ {
		case ttString:
			sb.WriteString("\"abc\"")
		case ttQuote:
			sb.WriteByte('\'')
		case ttPound:
			sb.WriteByte('#')
		case ttSemi:
			sb.WriteByte(';')
		case ttEscape:
			sb.WriteByte('\\')
		case ttPipe:
			sb.WriteByte('|')
		case ttParenL:
			sb.WriteByte('(')
		case ttParenR:
			sb.WriteByte(')')
		case ttBracketL:
			sb.WriteByte('[')
		case ttBracketR:
			sb.WriteByte(']')
		case ttBraceL:
			sb.WriteByte('{')
		case ttBraceR:
			sb.WriteByte('}')
		case ttLf:
			sb.WriteByte('\n')
		case ttWhiteSpace:
			sb.WriteByte(' ')
		default:
			sb.WriteString("_")
		}
	}
	return sb.String()
}

type Comment struct {
	Line bool
	Beg  Pos
	End  Pos
}

func Comments(s []byte) (good bool, comments []Comment, diags []Diagnostic) {
	var ts []Token
	good, ts, diags = Tokenize(s)
	var i Pos
	for i < len(ts) {
		t := ts[i]
		beg := i
		switch t.typ {
		case ttPound:
			i += 1
			if i < len(ts) {
				t2 := ts[i]
				switch t2.typ {
				case ttSemi:
					// we don't care about s-expression comments
					// eat the tokens
					i += 1
					continue
				case ttEscape: // #\
					i += 1
					if !(i < len(ts)) {
						diags = append(diags, Diagnostic{
							Level: DlFatal,
							Beg:   beg,
							End:   i,
							Msg:   "A character after #\\ but instead saw EOF",
						})
					}
					t3 := ts[i]
					switch t3.typ {
					case 
					}
				case ttPipe: // #|
					// just want to point out that yes we could've used nicey wicey
					// recursion (owo) but a simple counter will suffice
					nestLevel := 1
					for nestLevel > 0 {
						// eat every token that isn't either
						// 	another #|
						// 	a closing |#
						if !(i < len(ts)) {
							diags = append(diags, Diagnostic{
								Level: DlFatal,
								Beg:   beg,
								End:   i,
								Msg:   "Expected |# but instead saw EOF",
							})
						}
						i += 1
						t3 := ts[i]
						switch t3.typ {
						case ttPound:
							if i < len(ts) {
								// normally here we'd i += 1
								// but imagine this scenario:
								// ##|
								// 345
								// even if t4 isn't a pipe, we still want it to become
								// t3 next time
								t4 := ts[i]
								if t4.typ == ttPipe {
									// only consume it if it means something
									// in this case, another nest level
									i += 1
									nestLevel += 1
								}
							}
						case ttPipe:
							if i < len(ts) {
								t4 := ts[i]
								if t4.typ == ttPound {
									i += 1
									nestLevel -= 1
								}
							}
						}
					}
					comments = append(comments, Comment{
						Line: false,
						Beg: beg,
						End: i,
					})
				}
			}
		case ttPipe:
		case ttParenL:
		case ttParenR:
		case ttBracketL:
		case ttBracketR:
		case ttBraceL:
		case ttBraceR:
		case ttLf:
		case ttWhiteSpace:
		case ttOther:
		}
	}
	return
}
