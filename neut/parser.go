package neut

import "neutttr/lexer"

type Node interface {
	Base() NodeBase
}

type NodeBase struct {
	lexer.Sel
}

func (n NodeBase) Base() NodeBase {
	return n
}

type CommentNode struct {
	NodeBase
	children []Node
}

// declarations
type SumTypeNode struct {
	lexer.Sel
	target  IdentifierNode
	members []Node
}

func parseSumType(ts *[]token) *SumTypeNode {
	
}

type AliasNode struct {
	NodeBase
	target   IdentifierNode
	typeExpr Node
}

// (a|an) <word> is (a|an)
func parseAlias(ts *[]token) *AliasNode {
	tokens := *ts
	article1 := parseArticle(&tokens)
	if article1 == nil {
		return nil
	}
	identifier := parseIdentifier(&tokens)
	if identifier == nil {
		return nil
	}
	is := parseWord(&tokens)
	if is == nil {
		return nil
	}
	if is.text != "is" {
		return nil
	}
	maybeArticle2 := parseArticle(&tokens)
	// doesn't matter if that's nil
	_ = maybeArticle2
	typeExpr := parseTypeExpr(&tokens)
	if typeExpr == nil {
		return nil
	}
	ts = &tokens
	startPos := article1.Pos
	endOffset := (*typeExpr).Base().LastOffset()
	count := endOffset - startPos.Offset
	sel := lexer.Sel{Pos: startPos, Count: count}
	return &AliasNode{
		NodeBase: NodeBase{sel},
		target:   *identifier,
		typeExpr: *typeExpr,
	}
}

type AnnotationNode struct {
	NodeBase
	target   IdentifierNode
	typeExpr Node
}

// foo : bar
func parseAnnotation(ts *[]token) *AnnotationNode {
	tokens := *ts
	target := parseIdentifier(&tokens)
	if target == nil {
		return nil
	}
	colon := parseSymbol(&tokens, ':')
	if colon == nil {
		return nil
	}
	typeExpr := parseTypeExpr(&tokens)
	if typeExpr == nil {
		return nil
	}
	ts = &tokens
	startPos := target.Pos
	endOffset := (*typeExpr).Base().LastOffset()
	count := endOffset - startPos.Offset
	sel := lexer.Sel{Pos: startPos, Count: count}
	return &AnnotationNode{
		NodeBase: NodeBase{sel},
		target: *target,
		typeExpr: *typeExpr,
	}
}

type TemplateNode struct {
	lexer.Sel
	sub Node
}

func parseTemplate(ts *[]token) *TemplateNode {
	tokens := *ts
	
}

// Expression Nodes
func parseTypeExpr(ts *[]token) *Node {

}

type IdentifierNode struct {
	lexer.Sel
}

func parseIdentifier(ts *[]token) *IdentifierNode {
	tokens := *ts
	word := parseWord(&tokens)
	return &IdentifierNode{
		Sel: word.Sel,
	}
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

func parseSymbol(ts *[]token, c byte) *token {
	tokens := *ts
	if len(tokens) == 0 {
		return nil
	}
	current := tokens[0]
	if current.isWord() {
		return nil
	}
	if current.text[0] != c {
		return nil
	}
	rest := tokens[1:]
	ts = &rest
	return &current
}

func parseWord(ts *[]token) *token {
	tokens := *ts
	if len(tokens) == 0 {
		return nil
	}
	current := tokens[0]
	if !current.isWord() {
		return nil
	}
	rest := tokens[1:]
	ts = &rest
	return &current
}

func parseArticle(ts *[]token) *token {
	tokens := *ts
	word := parseWord(&tokens)
	if word == nil {
		return nil
	}
	if word.text != "a" {
		return nil
	}
	if word.text != "an" {
		return nil
	}
	ts = &tokens
	return word
}
