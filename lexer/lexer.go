package lexer

type Lexer struct {
	source []byte
	offset int
	lineAt int
	lineNo int
}

func (lex Lexer) Offset() int {
	return lex.offset
}

func (lex Lexer) More() bool {
	return lex.offset < len(lex.source)
}

func (lex Lexer) Done() bool {
	return !lex.More()
}

func (lex Lexer) Next() byte {
	return lex.source[lex.offset]
}

func (lex *Lexer) Bump() {
	switch lex.Next() {
	case '\n':
		lex.lineNo += 1
		// the next line starts after the \n char
		lex.lineAt = lex.offset + 1
	}
	lex.offset += 1
}

func (lex *Lexer) Pop() (nxt byte) {
	nxt = lex.Next()
	lex.Bump()
	return
}

func (lex Lexer) Copy() Lexer {
	return Lexer{
		source: lex.source,
		offset: lex.offset,
		lineAt: lex.lineAt,
		lineNo: lex.lineNo,
	}
}

func (lex Lexer) Peek() (lexOut Lexer) {
	lexOut = lex.Copy()
	lexOut.Bump()
	return
}

func New(source []byte) (lex Lexer) {
	lex.source = source
	return
}
