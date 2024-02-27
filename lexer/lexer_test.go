package lexer

import "testing"

func TestConstructor(t *testing.T) {
	lex := New("foobar")
	if lex.offset != 0 {
		t.Errorf("lexer should start at offset 0. instead it was %d", lex.offset)
	}
	if lex.lineNo != 0 {
		t.Errorf("lexer must start at line 0. instead it was %d", lex.lineNo)
	}
	if lex.lineAt != 0 {
		t.Errorf("The first line starts at 0. instead it was %d", lex.lineAt)
	}
}
