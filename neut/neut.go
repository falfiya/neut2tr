package neut

import "neutttr/lexer"

type NodeType int

const (
	// declarations
	ntdAlias NodeType = iota
	ntdAnnotation
	ntdSumType
	ntdTemplate
	// non-declarations
	ntIdentifier
	ntEssExpr
	ntFunction
	ntGeneric
)

type Node[D any] struct {
	typ NodeType
	extra D
}

func (n Node[any]) isDecl() bool {
	switch n.typ {
	case ntdAlias, ntdSumType, ntdAnnotation, ntdTemplate:
		return true
	}
	return false
}

type DeclNode = Node[any]
type ExprNode = Node[any]

type NodeExTargetAndType struct {
	target IdentifierNode
	typeExpr ExprNode
}

type NodeExSumType struct {
	target  IdentifierNode
	members []ExprNode
}

type NodeExTemplate struct {
	decl DeclNode
}

type NodeExIdentifier string
type IdentifierNode = Node[NodeExIdentifier]

type NodeExEssExpr struct {
	// either '(' or '['
	char byte
	members []ExprNode
}

type NodeExFunction struct {
	lhs []ExprNode
	rhs []ExprNode
}

type NodeExGeneric struct {
	params []IdentifierNode
	subExpr ExprNode
}

type Node2 interface {
	lexer.Sel
	
}

type IdentifierNode2 struct {
	name string
}
func (IdentifierNode2) typ() NodeType {
	return ntIdentifier
}
func (IdentifierNode2) isDecl() bool {
	return false
}

