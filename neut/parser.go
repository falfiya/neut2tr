package neut

import "neutttr/lexer"

type Node any

// returns [](nil | TemplateNode | DeclarationNode | AnnotationNode)
func parse(tokens []token) (meaningful []Node) {
	// All valid NEU type info starts on the first token of a line.
	for len(tokens) > 0 {
		maybeMeaningful := parseMeaningful(&tokens)
		if maybeMeaningful != nil {
			meaningful = append(meaningful, maybeMeaningful)
			continue
		} else {
			// If we encounter garbage at the start of a line, ignore the rest of the
			// line.
			for len(tokens) > 0 {
				_, isNewline := tokens[0].(newlineToken)
				tokens = tokens[1:]
				if isNewline {
					break
				}
			}
		}
	}
	return
}

// returns nil | TemplateNode | DeclarationNode | AnnotationNode
func parseMeaningful(ts *[]token) (node Node) {
	node = parseTemplate(ts)
	if node != nil {
		return
	}
	node = parseDeclaration(ts)
	if node != nil {
		return
	}
	node = parseAnnotation(ts)
	return
}

type TemplateNode struct {
	lexer.Sel
	decl DeclarationNode
}

func parseTemplate(ts *[]token) *TemplateNode {
	tokens := *ts
	maybeTemplate := getIdentifierToken(&tokens)
	if maybeTemplate == nil {
		return nil
	}

	if maybeTemplate.name != "template" && maybeTemplate.name != "template:" {
		return nil
	}

	decl := parseDeclaration(&tokens)
	if decl == nil {
		return nil
	}

	commit(ts, tokens)
	return &TemplateNode{
		sel: maybeTemplate.Select(),
		decl:     *decl,
	}
}

type DeclarationNode struct {
	lexer.Sel
	// IdentifierNode | DeclarationGenericTargetNode
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
	is := parseExactString(&tokens, "is")
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
	*ts = tokens
	return &DeclarationNode{
		NodeBase: newNodeBase(article1, rhs),
		target:   target,
		rhs:      rhs,
	}
}

type DeclarationGenericTargetNode struct {
	NodeBase
	target IdentifierNode
	params []IdentifierNode
}

func parseDeclarationTarget(ts *[]token) Node {
	generic := parseDeclarationGenericTarget(ts)
	if generic != nil {
		return generic
	}
	target := parseIdentifier(ts)
	if target == nil {
		return nil
	}
	return target
}

func parseDeclarationGenericTarget(ts *[]token) *DeclarationGenericTargetNode {
	tokens := *ts
	// a declaration target looks like
	// (Listof X)
	// There should be at least three tokens
	if len(tokens) < 3 {
		return nil
	}
	open := tokens[0]
	var closeText string
	if open.text == "(" {
		closeText = ")"
		goto genericDetected
	}
	if open.text == "[" {
		closeText = "]"
		goto genericDetected
	}
	// generic not detected
	return nil
genericDetected:
	target := tokens[1]
	tokens = tokens[2:]
	var current token
	var params []IdentifierNode
	for {
		if len(tokens) == 0 {
			return nil
		}
		current = tokens[0]
		tokens = tokens[1:]
		if current.text == closeText {
			// when we encounter the closing token, current will still exist
			// as the last token
			break
		}
		params = append(params, IdentifierNode{NodeBase{current.Sel}})
	}
	*ts = tokens
	return &DeclarationGenericTargetNode{
		NodeBase: newNodeBase2(open, current),
		target:   IdentifierNode{NodeBase{target.Sel}},
		params:   params,
	}
}

type SumTypeTerm struct {
	NodeBase
	terms []SumTypeElementNode
}

// ... one of
// - x
// - y
func parseSumTypeRhs(ts *[]token) *SumTypeTerm {
	// somewhat amusingly, this doesn't check for newlines in between terms.
	// shhhhh don't tell!
	tokens := *ts
	one := parseExactString(&tokens, "one")
	if one == nil {
		return nil
	}
	of := parseExactString(&tokens, "of")
	if of == nil {
		return nil
	}
	term1 := parseSumTypeTerm(&tokens)
	if term1 == nil {
		return nil
	}
	term2 := parseSumTypeTerm(&tokens)
	if term2 == nil {
		return nil
	}
	terms := []SumTypeElementNode{*term1, *term2}
	for {
		optionalTerm := parseSumTypeTerm(&tokens)
		if optionalTerm == nil {
			break
		} else {
			terms = append(terms, *optionalTerm)
		}
	}
	*ts = tokens
	return &SumTypeTerm{
		NodeBase: newNodeBase(one, terms[len(terms)-1]),
		terms:    terms,
	}
}

type SumTypeElementNode struct {
	NodeBase
	typeNode TypeNode
}

func parseSumTypeTerm(ts *[]token) *SumTypeElementNode {
	tokens := *ts
	entry := parseExactString(&tokens, "-")
	if entry == nil {
		return nil
	}

	typeNode := parseType(&tokens)
	if typeNode == nil {
		return nil
	}

	*ts = tokens
	return &SumTypeElementNode{
		NodeBase: newNodeBase(entry, typeNode),
	}
}

type AliasRhsNode struct {
	NodeBase
	typeNode TypeNode
}

// ... (a|an) x
func parseAliasRhs(ts *[]token) *AliasRhsNode {
	tokens := *ts
	// doesn't matter if that's nil
	maybeArticle := parseArticle(&tokens)

	typeNode := parseType(&tokens)
	if typeNode == nil {
		return nil
	}

	var nodeBase NodeBase
	if maybeArticle == nil {
		nodeBase = typeNode.Base()
	} else {
		nodeBase = newNodeBase(maybeArticle, typeNode)
	}

	*ts = tokens
	return &AliasRhsNode{
		NodeBase: nodeBase,
		typeNode: typeNode,
	}
}

type AnnotationNode struct {
	NodeBase
	target   IdentifierNode
	typeNode TypeNode
}

// foo : bar
func parseAnnotation(ts *[]token) *AnnotationNode {
	tokens := *ts

	target := parseIdentifier(&tokens)
	if target == nil {
		return nil
	}

	colon := parseExactString(&tokens, ":")
	if colon == nil {
		return nil
	}

	var typeNode Node
	typeNode = parseFunctionInside(&tokens)
	if typeNode != nil {
		// typeNode : FunctionTypeNode
		goto typeResolved
	}
	typeNode = parseType(&tokens)
	if typeNode != nil {
		// typeNode : TypeNode
		goto typeResolved
	}
	return nil
typeResolved:
	*ts = tokens
	return &AnnotationNode{
		NodeBase: newNodeBase3(target, typeNode),
		target:   *target,
		typeNode: typeNode,
	}
}

// FunctionTypeNode | ListTypeNode | IdentifierNode
type TypeNode Node
// or nil
func parseType(ts *[]token) Node {
	tokens := *ts

	maybeFunctionTypeNode := parseFunctionType()
}

type FunctionTypeNode struct {
	// or nil
	generic *FunctionGenericNode
	input []TypeNode
	output TypeNode
}

func parseFunctionType(ts *[]token) *FunctionTypeNode {
	tokens := *ts

	left := parseExactString(&tokens, "[")
	if left == nil {
		return nil
	}

	inside := parseFunctionInside(&tokens)
	if inside == nil {
		return nil
	}

	right := parseExactString(&tokens, "]")
	if right == nil {
		return nil
	}

	*ts = tokens
	return inside
}

func parseFunctionInside(ts *[]token) *FunctionTypeNode {
	tokens := *ts

	maybeGeneric := parseFunctionGeneric(&tokens)
	for {
		parseExactString(&tokens, "->")
	}
}

type FunctionGenericNode struct {
	NodeBase
	params []IdentifierNode
}

func parseFunctionGeneric(ts *[]token) *FunctionGenericNode {
	tokens := *ts

	open := parseExactString(&tokens, "{")
	if open == nil {
		return nil
	}

	var current token
	var params []IdentifierNode
	for {
		if len(tokens) == 0 {
			return nil
		}
		current = tokens[0]
		tokens = tokens[1:]
		if current.text == "}" {
			break
		} else {
			params = append(params, IdentifierNode{NodeBase{current.Sel}})
		}
	}

	*ts = tokens
	return &FunctionGenericNode{
		NodeBase: newNodeBase(open, params[len(params)-1]),
		params: params,
	}
}

type ListTypeNode struct {
	lexer.Sel
	members []TypeNode
}

func parseListType(ts *[]token) *ListTypeNode {
	tokens := *ts
	maybeLeft := get[symbolToken](&tokens)

	if maybeLeft == nil {
		return nil
	}

	var endOffset int
	var members []TypeNode

	var close byte
	if maybeLeft.symbol == '(' {
		close = ')'
		goto listInside
	}
	if maybeLeft.symbol == '[' {
		close = ']'
		goto listInside
	}
	return nil

listInside:
	for {
		maybeRight := get[symbolToken](&tokens)
		if maybeRight != nil && maybeRight.symbol == close {
			endOffset = maybeRight.End()
			break listInside
		}
		maybeType := parseType(&tokens)
		if maybeType == nil {
			return nil
		}
		members = append(members, maybeType)
	}
	commit(ts, tokens)
	return &ListTypeNode{
		Sel: maybeLeft.SelectTill(endOffset),
		members: members,
	}
}

type IdentifierNode struct {
	lexer.Sel
}

func parseIdentifier(ts *[]token) *IdentifierNode {
	tokens := *ts
	if len(tokens) == 0 {
		return nil
	}
	current, ok := tokens[0].(identifierToken)
	if ok {
		commit(ts, tokens[1:])
		return &IdentifierNode{current.Sel}
	}
	return nil
}

type ArticleNode struct {
	lexer.Sel
}

func parseArticle(ts *[]token) *ArticleNode {
	tokens := *ts
	maybeArticleToken := get[identifierToken](&tokens)
	if maybeArticleToken.name == "a" || maybeArticleToken.name == "an" {
		commit(ts, tokens)
		return &ArticleNode{maybeArticleToken.Sel}
	} else {
		return nil
	}
}

func get[T token](ts *[]token) *T {
	tokens := *ts
	if len(tokens) == 0 {
		return nil
	}
	maybeIdentifierToken, ok := tokens[0].(T)
	if ok {
		tokens = tokens[1:]
		commit(ts, tokens)
		return &maybeIdentifierToken
	}
	return nil
}

func commit(ts *[]token, newTokens []token) {
	*ts = newTokens
}
