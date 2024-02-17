package neut

import "neutttr/lexer"

type sel = lexer.Sel
type pos = lexer.Pos

type Node interface {
	Base() NodeBase
}

type NodeBase struct {
	lexer.Sel
}

// function overloading
func newNodeBase(start *token, end Node) NodeBase {
	return _newNodeBase(start.Pos, end.Base().Sel)
}
func newNodeBase2(start token, end token) NodeBase {
	return _newNodeBase(start.Pos, end.Sel)
}
func newNodeBase3(start Node, end Node) NodeBase {
	return _newNodeBase(start.Base().Pos, end.Base().Sel)
}
func _newNodeBase(start lexer.Pos, end lexer.Sel) NodeBase {
	count := end.Offset + end.Count - start.Offset
	sel := lexer.Sel{Pos: start, Count: count}
	return NodeBase{sel}
}

func (n NodeBase) Base() NodeBase {
	return n
}
