package racket

import (
	"strings"
)

//go:generate stringer -type=TokenType
type TokenType int

const (
	TtQuote TokenType = iota
	TtPound
	TtSemi
	TtEscape
	TtPipe
	TtParenL
	TtParenR
	TtBracketL
	TtBracketR
	TtBraceL
	TtBraceR
	TtIndent
	TtLf
	TtString
	TtOther
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
	level DiagnosticLevel
	beg   Pos
	end   Pos
	msg   string
}

func Tokenize(s []byte) (good bool, out []Token, diags []Diagnostic) {
	// i is the current position and also the new "end" when committing tokens
	var i Pos
tokenLoop:
	for i < len(s) {
		// skip over spaces
		var beg Pos = i
		var typ TokenType
		switch s[i] {
		case '\'':
			typ = TtQuote
		case '#':
			typ = TtPound
		case ';':
			typ = TtSemi
		case '\\':
			typ = TtEscape
		case '|':
			typ = TtPipe
		case '(':
			typ = TtParenL
		case ')':
			typ = TtParenR
		case '[':
			typ = TtBracketL
		case ']':
			typ = TtBracketR
		case '{':
			typ = TtBraceL
		case '}':
			typ = TtBraceR
		case '\n':
			typ = TtLf
		case '	', ' ':
			typ = TtIndent
		default:
			// any control characters that weren't picked up before this are discarded
			if s[i] < '!' || s[i] == 127 {
				i += 1
				continue tokenLoop
			}
			goto stringOrOther // -----+
		}                     //      |
		i += 1                //      |
		goto commitToken      //      |
	stringOrOther:           //      |
		if s[i] == '"' {      // <----+
			i += 1
			typ = TtString
			escaped := false
			for {
				if i == len(s) {
					diags = append(diags, Diagnostic{
						level: DlFatal,
						beg:   beg,
						end:   i,
						msg:   "Expected \" but instead saw EOF",
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
						goto commitToken
					}
				}
				i += 1
			}
		} else {
			i += 1
			typ = TtOther
		}
	commitToken:
		// there are a few types of tokens we'd like to combine
		// - TtIndent
		// - TtOther
		// if this token is the same
		if typ == TtOther || typ == TtIndent {
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
		case TtQuote:
			sb.WriteByte('\'')
		case TtPound:
			sb.WriteByte('#')
		case TtSemi:
			sb.WriteByte(';')
		case TtEscape:
			sb.WriteByte('\\')
		case TtPipe:
			sb.WriteByte('|')
		case TtParenL:
			sb.WriteByte('(')
		case TtParenR:
			sb.WriteByte(')')
		case TtBracketL:
			sb.WriteByte('[')
		case TtBracketR:
			sb.WriteByte(']')
		case TtBraceL:
			sb.WriteByte('{')
		case TtBraceR:
			sb.WriteByte('}')
		case TtString:
			sb.WriteString("\"???\"")
		case TtLf:
			sb.WriteByte('\n')
		case TtIndent:
			sb.WriteByte(' ')
		default:
			sb.WriteString("_")
		}
	}
	return sb.String()
}
