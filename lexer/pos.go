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

// (Zero-Indexed)
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

func (s Sel) LastOffset() int {
	return s.Offset + s.Count - 1
}

// selection does not include End
func (s Sel) End() int {
	return s.Offset + s.Count
}

func (p Pos) End() int {
	return p.Offset + 1
}
