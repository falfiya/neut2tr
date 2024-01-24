package lexer

type Pos struct {
	// In bytes
	Offset int
	// Line number (Zero-Indexed)
	LineNo int
	// Start of the line
	LineAt int
}

type Sel struct {
	Pos
	Count int
}

func (lex Lexer) Pos() Pos {
	return Pos{
		Offset: lex.offset,
		LineNo: lex.lineNo,
		LineAt: lex.lineAt,
	}
}

func (p Pos) CharNo() int {
	return p.Offset - p.LineAt
}

func (p Pos) Select(count int) Sel {
	return Sel {
		Pos: p,
		Count: count,
	}
}

func (p Pos) SelectTill(offset int) Sel {
	return Sel {
		Pos: p,
		Count: offset - p.Offset,
	}
}

func (p Pos) StartOfLine(source string) string {
	return source[p.LineAt:p.Offset]
}

func (p Pos) Before(source string) string {
	return source[:p.Offset]
}

func (s Sel) From(source string) string {
	return source[s.Offset:s.Offset + s.Count]
}

func (s Sel) After(source string) string {
	return source[s.Offset + s.Count:]
}
