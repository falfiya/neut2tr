package racket

import "strings"

type Comment struct {
	IsLineComment bool
	Beg           int
	End           int
	LineNum       int
	CharNum       int
	KnownIndent   []byte
}

func (c Comment) InnerText(src string) string {
	if c.IsLineComment {
		return src[c.Beg+1 : c.End]
	} else {
		return src[c.Beg+2 : c.End-2]
	}
}
func (c Comment) GetIndent() string {
	var sb strings.Builder
	for _, i := range c.KnownIndent {
		sb.WriteByte(i)
	}
	missingIndent := c.CharNum - 1 - len(c.KnownIndent)
	for missingIndent > 0 {
		sb.WriteByte(' ')
		missingIndent -= 1
	}
	return sb.String()
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
	Beg   int
	End   int
	Msg   string
}

func Comments(s []byte) (good bool, comments []Comment, diags []Diagnostic) {
	i := 0
	// lineCounter specific variables
	lineCounter := 0
	lastLf := 0
	shouldIndent := true
	indentSoFar := make([]byte, 0)
	nextLine := func() {
		lastLf = i
		lineCounter += 1
		indentSoFar = make([]byte, 0)
		shouldIndent = true
		i += 1
	}
	for i < len(s) {
		beg := i
		begCharNum := beg - lastLf
		c := s[i]
		switch c {
		case '\n':
			nextLine()
		case '	', ' ':
			// keep track of the current indent
			// indents can only consist of tabs and spaces and this branch sees
			// them all.
			var prev byte
			prevIdx := i - 1
			if prevIdx < 0 {
				prev = '\n'
			} else {
				prev = s[prevIdx]
			}
			i += 1
			if shouldIndent {
				switch prev {
				case '\n', '	', ' ':
					indentSoFar = append(indentSoFar, c)
				default:
					// something broke the indent chain
					shouldIndent = false
				}
			}
		case '"':
			i += 1
			escaped := false
		stringLoop:
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
				c2 := s[i]
				if escaped {
					if c2 == '\n' {
						nextLine()
					} else {
						i += 1
					}
					escaped = false
				} else {
					switch c2 {
					case '\n':
						nextLine()
					case '\\':
						i += 1
						escaped = true
					case '"':
						i += 1
						break stringLoop
					default:
						i += 1
					}
				}
			}
		case '#':
			savedIndentSoFar := indentSoFar
			i += 1
			if i < len(s) {
				c2 := s[i]
				switch c2 {
				case ';':
					// we don't care about s-expression comments
					// eat the tokens
					i += 1
					continue
				case '\\': // #\
					// eat the \ plus one more character
					i += 2
				case '|': // #|
					i += 1
					// just want to point out that yes we could've used nicey wicey
					// recursion (owo) but a simple counter will suffice
					nestLevel := 1
					for nestLevel > 0 {
						// eat every token that isn't either
						// 	another #|
						// 	a closing |#
						if !(i < len(s)) {
							diags = append(diags, Diagnostic{
								Level: DlFatal,
								Beg:   beg,
								End:   i,
								Msg:   "Expected |# but instead saw EOF",
							})
							return
						}
						switch s[i] {
						case '\n':
							nextLine()
						case '#': // first byte is #
							i += 1
							if i < len(s) && s[i] == '|' {
								// the second byte is |
								// we have another block comment
								i += 1
								nestLevel += 1
							}
						case '|': // first byte |
							i += 1
							if i < len(s) && s[i] == '#' {
								// the second byte is #
								// that's a closing block comment
								i += 1
								nestLevel -= 1
							}
						default:
							i += 1
						}
					}
					comments = append(comments, Comment{
						IsLineComment: false,
						Beg:           beg,
						End:           i,
						LineNum:       lineCounter,
						CharNum:       begCharNum,
						KnownIndent:   savedIndentSoFar,
					})
				}
			}
		case ';':
			i += 1
			savedIndentSoFar := indentSoFar
		lineComment:
			for i < len(s) {
				if s[i] == '\n' {
					nextLine()
					break lineComment
				}
				i += 1
			}
			comments = append(comments, Comment{
				IsLineComment: true,
				Beg:           beg,
				End:           i,
				LineNum:       lineCounter,
				CharNum:       begCharNum,
				KnownIndent:   savedIndentSoFar,
			})
		default:
			i += 1
		}
	}
	good = true
	return
}
