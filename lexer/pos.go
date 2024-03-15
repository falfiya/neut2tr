package lexer

type Pos struct {
	// In bytes
	Offset int
	// Line number (Zero-Indexed)
	LineNo int
	// Offset of the first character of the line
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

func (p Pos) Select(notIncludedOffset int) Sel {
	return Sel {
		Pos: p,
		Count: notIncludedOffset - p.Offset,
	}
}

// selection does not include End
func (s Sel) End() int {
	return s.Offset + s.Count
}

// the offset after this position
func (p Pos) End() int {
	return p.Offset + 1
}

type SelF interface {
	SelF() Sel
}

func (s Sel) SelF() Sel {
	return s
}
