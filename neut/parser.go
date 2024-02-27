package neut

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
	NodeBase
	decl DeclarationNode
}

func parseTemplate(ts *[]token) *TemplateNode {
	tokens := *ts
	template := parseString(&tokens, "template")
	if template == nil {
		return nil
	}
	// maybe colon
	_ = parseString(&tokens, ":")
	decl := parseDeclaration(&tokens)
	if decl == nil {
		return nil
	}
	*ts = tokens
	return &TemplateNode{
		NodeBase: newNodeBase(template, decl),
		decl:     *decl,
	}
}

type DeclarationNode struct {
	NodeBase
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
	one := parseString(&tokens, "one")
	if one == nil {
		return nil
	}
	of := parseString(&tokens, "of")
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
	entry := parseString(&tokens, "-")
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

	colon := parseString(&tokens, ":")
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
type TypeNode = Node

// Expression Nodes
func parseType(ts *[]token) Node {
	
}

type FunctionTypeNode struct {
	// or nil
	generic *FunctionGenericNode
	input []TypeNode
	output TypeNode
}

func parseFunction(ts *[]token) *FunctionTypeNode {
	tokens := *ts

	left := parseString(&tokens, "[")
	if left == nil {
		return nil
	}

	inside := parseFunctionInside(&tokens)
	if inside == nil {
		return nil
	}

	right := parseString(&tokens, "]")
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
		parseString(&tokens, "->")
	}
}

type FunctionGenericNode struct {
	NodeBase
	params []IdentifierNode
}

func parseFunctionGeneric(ts *[]token) *FunctionGenericNode {
	tokens := *ts

	open := parseString(&tokens, "{")
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

type IdentifierNode struct {
	NodeBase
}

func parseIdentifier(ts *[]token) *IdentifierNode {
	tokens := *ts
	if len(tokens) == 0 {
		return nil
	}
	current, ok := tokens[0].(identifierToken)
	if ok {
		*ts = tokens[1:]
		return &IdentifierNode{NodeBase{current.Sel}}
	}
	return nil
}

// Returns nil if the token did not match s.
// Modifies ts on success.
func parseString(ts *[]token, s string) *token {
	tokens := *ts
	if len(tokens) == 0 {
		return nil
	}
	current := tokens[0]
	if current.text != s {
		return nil
	}
	rest := tokens[1:]
	ts = &rest
	return &current
}

func parseArticle(ts *[]token) *identifierToken {
	tokens := *ts
	if len(tokens) == 0 {
		return nil
	}
	current := tokens[0]
	switch t := current.(type) {
	case identifierToken:
		if t.name == "a" || t.name == "an" {
			*ts = tokens[1:]
			return &t
		}
	}
	return nil
}
