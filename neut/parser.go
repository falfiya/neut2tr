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
type DeclarationNode struct {
	NodeBase
	// DeclarationTargetNode or DeclarationGenericTargetNode
	target Node
	rhs    Node
}

func parseDeclaration(ts *[]token) *DeclarationNode {
	tokens := *ts
	article1 := parseArticle(&tokens)
	if article1 == nil {
		return nil
	}
	target := parseDeclarationTarget(&tokens)
	if target == nil {
		return nil
	}
	is := parseString(&tokens, "is")
	if is == nil {
		return nil
	}

	var rhs Node
	rhs = parseSumTypeRhs(&tokens)
	if rhs != nil {
		goto commitDeclaration
	}
	rhs = parseAliasRhs(&tokens)
	if rhs != nil {
		goto commitDeclaration
	}
	// neither worked
	return nil
commitDeclaration:
	ts = &tokens
	startPos := article1.Pos
	endOffset := rhs.Base().LastOffset()
	count := endOffset - startPos.Offset
	sel := lexer.Sel{Pos: startPos, Count: count}
	return &DeclarationNode{
		NodeBase: NodeBase{sel},
		target: target,
		rhs: rhs,
	}
}

type DeclarationTargetNode struct {
	
}

type DeclarationGenericTargetNode struct {
	
}

func parseDeclarationTarget(ts *[]token) Node {
	tokens := *ts
	firstToken := parseWord(&tokens)
	if firstToken == nil {
		return nil
	}
	if firstToken.text == "(" {
		
	}
}

type SumTypeRhsNode struct {
	NodeBase
	members []SumTypeElementNode
}

// ... is one of
// - x
// - y
func parseSumTypeRhs(ts *[]token) *SumTypeRhsNode {

}

type SumTypeElementNode struct {
	
}
func parseSumTypeElement(ts *[]token) *SumTypeElementNode {
	tokens := *ts
	entry := parseString(&tokens, "-")
	if entry == nil {
		return nil
	}
	
	if typeExpr == nil {
		return nil
	}
	ts = &tokens
	var startPos lexer.Pos
	if maybeArticle == nil {
		startPos = typeExpr.Base().Pos
	} else {
		startPos = maybeArticle.Pos
	}
	endOffset := typeExpr.Base().LastOffset()
	count := endOffset - startPos.Offset
	sel := lexer.Sel{Pos: startPos, Count: count}
	return &AliasRhsNode{
		NodeBase: NodeBase{sel},
		typeExpr: typeExpr,
	}
}

type AliasRhsNode struct {
	NodeBase
	typeExpr Node
}

// ... (a|an) x
func parseAliasRhs(ts *[]token) *AliasRhsNode {
	tokens := *ts
	// doesn't matter if that's nil
	maybeArticle := parseArticle(&tokens)
	typeExpr := parseTypeExpr(&tokens)
	if typeExpr == nil {
		return nil
	}
	ts = &tokens
	var startPos lexer.Pos
	if maybeArticle == nil {
		startPos = typeExpr.Base().Pos
	} else {
		startPos = maybeArticle.Pos
	}
	endOffset := typeExpr.Base().LastOffset()
	count := endOffset - startPos.Offset
	sel := lexer.Sel{Pos: startPos, Count: count}
	return &AliasRhsNode{
		NodeBase: NodeBase{sel},
		typeExpr: typeExpr,
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
	colon := parseString(&tokens, ":")
	if colon == nil {
		return nil
	}
	typeExpr := parseTypeExpr(&tokens)
	if typeExpr == nil {
		return nil
	}
	ts = &tokens
	startPos := target.Pos
	endOffset := typeExpr.Base().LastOffset()
	count := endOffset - startPos.Offset
	sel := lexer.Sel{Pos: startPos, Count: count}
	return &AnnotationNode{
		NodeBase: NodeBase{sel},
		target:   *target,
		typeExpr: typeExpr,
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
func parseTypeExpr(ts *[]token) Node {

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

func parseString(ts *[]token, s string) *token {
	tokens := *ts
	if len(tokens) == 0 {
		return nil
	}
	current := tokens[0]
	if current.text != c {
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
