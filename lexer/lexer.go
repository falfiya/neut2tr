package lexer

type Lexer struct {
	source string
	offset int
	lineAt int
	lineNo int
}

func New(source string) (lex Lexer) {
	lex.source = source
	return
}

func (lex Lexer) Offset() int {
	return lex.offset
}

func (lex Lexer) More() bool {
	return lex.offset < len(lex.source)
}

func (lex Lexer) Done() bool {
	return lex.offset >= len(lex.source)
}

// always call lex.More or lex.Done before calling this
func (lex Lexer) Current() byte {
	return lex.source[lex.offset]
}

// moves the lexer forward one byte
func (lex *Lexer) Bump() {
	switch lex.Current() {
	case '\n':
		lex.lineNo += 1
		// the next line starts after the \n char
		lex.lineAt = lex.offset + 1
	}
	lex.offset += 1
}
