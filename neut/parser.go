package neut

import "neutttr/lexer"

type Node = any

type CommentNode struct {
	children []Node
}

func parse(s string) any {
	tokens := tokenize(s)
	i := 0
	for i < len(tokens) {
		tok := tokens[i]
		i += 1
		switch tok.text {
		case "a", "an":
			if i < len(tokens) {
				tok2 := tokens[i]
				if tok2.typ == ttWord {
					
					continue
				} else {
					panic("This isn't valid syntax")
				}
			}
		}
	}
}

// declarations
type AliasNode struct {
	lexer.Sel
	target   IdentifierNode
	typeExpr Node
}

type AnnotationNode struct {
	lexer.Sel
	target   IdentifierNode
	typeExpr Node
}

type SumTypeNode struct {
	lexer.Sel
	target  IdentifierNode
	members []Node
}

type TemplateNode struct {
	lexer.Sel
	sub Node
}

func parseTemplate(ts *[]token) *TemplateNode {
	tokens := *ts
	
}

// Expression Nodes
type IdentifierNode struct {
	lexer.Sel
	name string
}

type SExprNode struct {
	lexer.Sel
	// either '(' or '['
	char    byte
	members []Node
}

type FunctionNode struct {
	lexer.Sel
	lhs []Node
	rhs []Node
}

type GenericNode struct {
	lexer.Sel
	params []Node
	sub    Node
}

func parseIdentifier(tokens *[]token) *IdentifierNode {
	ts := *tokens
	t := ts[0]
	if t.typ == ttWord {
		rest := ts[1:]
		tokens = &rest
		return &IdentifierNode{
			t.Sel,
			t.text,
		}
	} else {
		return nil;
	}
}
