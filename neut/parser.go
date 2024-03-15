package neut

import "neutttr/lexer"

type Node lexer.SelF

// returns [](nil | TemplateNode | DeclarationNode | AnnotationNode)
func Parse(tokens []token) []Node {
	var meaningful []Node
	// All valid NEU type info starts on the first token of a line.
	for len(tokens) > 0 {
		if maybeMeaningful := parseMeaningful(&tokens); maybeMeaningful != nil {
			meaningful = append(meaningful, maybeMeaningful)
			continue
		}
		// kill all tokens until another newline
	inner:
		for len(tokens) > 0 {
			_, isNewline := tokens[0].(NewlineToken)
			advance(&tokens)
			if isNewline {
				break inner
			}
		}
	}
	return meaningful
}

// returns nil | TemplateNode | DeclarationNode | AnnotationNode
func parseMeaningful(ts *[]token) Node {
	if temp := parseTemplate(ts); temp != nil {
		return temp
	}

	if decl := parseDeclaration(ts); decl != nil {
		return decl
	}

	return ParseAnnotation(ts)
}

type TemplateNode struct {
	lexer.Sel
	decl DeclarationNode
}

func parseTemplate(ts *[]token) *TemplateNode {
	tokens := *ts

	maybeTemplate := get[IdentifierToken](&tokens)
	if maybeTemplate == nil {
		return nil
	}
	if maybeTemplate.Name != "template" && maybeTemplate.Name != "template:" {
		return nil
	}

	decl := parseDeclaration(&tokens)
	if decl == nil {
		return nil
	}

	commit(ts, tokens)
	return &TemplateNode{maybeTemplate.Select(decl.End()), *decl}
}

type DeclarationNode struct {
	lexer.Sel
	// IdentifierNode | DeclarationGenericTargetNode
	Target Node
	// SumTypeNode | AliasNode
	Value Node
}

func parseDeclaration(ts *[]token) *DeclarationNode {
	tokens := *ts

	article1 := parseArticle(&tokens)
	if article1 == nil {
		return nil
	}

	var target Node
	if generic := parseGenericTarget(&tokens); generic != nil {
		target = generic
	} else if ident := get[IdentifierToken](&tokens); ident != nil {
		target = ident
	} else {
		return nil
	}

	maybeIs := get[IdentifierToken](&tokens)
	if maybeIs.Name != "is" {
		return nil
	}

	var value Node
	if sum := parseSumType(&tokens); sum != nil {
		value = sum
	} else if alias := parseAlias(&tokens); alias != nil {
		value = alias
	} else {
		return nil
	}

	commit(ts, tokens)
	sel := article1.Select(value.SelF().End())
	return &DeclarationNode{sel, target, value}
}

type GenericTargetNode struct {
	lexer.Sel
	target IdentifierToken
	params []IdentifierToken
}

func parseGenericTarget(ts *[]token) *GenericTargetNode {
	tokens := *ts
	// a declaration target looks like
	// (Listof X)
	// [Foo Xgeneric]
	// There should be at least four tokens
	if len(tokens) < 4 {
		return nil
	}

	open := get[SymbolToken](&tokens)
	var closeSymbol byte
	if open.Symbol == '(' {
		closeSymbol = ')'
	} else if open.Symbol == '[' {
		closeSymbol = ']'
	} else {
		return nil
	}

	target := get[IdentifierToken](&tokens)
	if target == nil {
		return nil
	}

	// need at least one generic parameter
	param1 := get[IdentifierToken](&tokens)
	if param1 == nil {
		return nil
	}

	var closingSymbolEnd int
	params := []IdentifierToken{*param1}
	for {
		if len(tokens) == 0 {
			return nil
		}
		current := tokens[0]
		advance(&tokens)
		switch v := current.(type) {
		case IdentifierToken:
			params = append(params, v)
		case SymbolToken:
			if v.Symbol == closeSymbol {
				closingSymbolEnd = v.End()
				goto commitGenericTarget
			} else {
				return nil
			}
		default:
			return nil
		}
	}

commitGenericTarget:
	commit(ts, tokens)
	return &GenericTargetNode{open.Select(closingSymbolEnd), *target, params}
}

type SumTypeNode struct {
	lexer.Sel
	one   IdentifierToken
	of    IdentifierToken
	terms []SumTypeTermNode
}

// ... one of
// - x
// - y
func parseSumType(ts *[]token) *SumTypeNode {
	tokens := *ts

	one := get[IdentifierToken](&tokens)
	if one == nil {
		return nil
	}
	if one.Name != "one" {
		return nil
	}
	of := get[IdentifierToken](&tokens)
	if of == nil {
		return nil
	}
	if of.Name != "of" {
		return nil
	}

	// must have at least two terms
	term1 := parseSumTypeTerm(&tokens)
	if term1 == nil {
		return nil
	}
	term2 := parseSumTypeTerm(&tokens)
	if term2 == nil {
		return nil
	}

	terms := []SumTypeTermNode{*term1, *term2}
	for {
		optionalTerm := parseSumTypeTerm(&tokens)
		if optionalTerm == nil {
			break
		} else {
			terms = append(terms, *optionalTerm)
		}
	}

	commit(ts, tokens)
	sel := one.Select(terms[len(terms)-1].End())
	return &SumTypeNode{sel, *one, *of, terms}
}

type SumTypeTermNode struct {
	lexer.Sel
	TypeNode TypeNode
}

func parseSumTypeTerm(ts *[]token) *SumTypeTermNode {
	tokens := *ts

	newline := get[NewlineToken](&tokens)
	if newline == nil {
		return nil
	}

	hyphen := get[IdentifierToken](&tokens)
	if hyphen == nil {
		return nil
	}
	if hyphen.Name != "-" {
		return nil
	}

	typeNode := ParseType(&tokens)
	if typeNode == nil {
		return nil
	}

	commit(ts, tokens)
	return &SumTypeTermNode{
		newline.Select(typeNode.SelF().End()),
		typeNode,
	}
}

type AliasNode struct {
	lexer.Sel
	Article  *ArticleNode
	TypeNode TypeNode
}

// ... (a|an) x
func parseAlias(ts *[]token) *AliasNode {
	tokens := *ts
	// doesn't matter if that's nil
	maybeArticle := parseArticle(&tokens)

	typeNode := ParseType(&tokens)
	if typeNode == nil {
		return nil
	}

	var sel lexer.Sel
	if maybeArticle == nil {
		sel = typeNode.SelF()
	} else {
		sel = maybeArticle.Select(typeNode.SelF().End())
	}

	commit(ts, tokens)
	return &AliasNode{
		sel,
		maybeArticle,
		typeNode,
	}
}

type ArticleNode struct {
	lexer.Sel
}

func parseArticle(ts *[]token) *ArticleNode {
	tokens := *ts
	maybeArticleToken := get[IdentifierToken](&tokens)
	if maybeArticleToken == nil {
		return nil
	}
	if maybeArticleToken.Name == "a" || maybeArticleToken.Name == "an" {
		commit(ts, tokens)
		return &ArticleNode{maybeArticleToken.Sel}
	} else {
		return nil
	}
}

type AnnotationNode struct {
	lexer.Sel
	target IdentifierToken
	// not a pointer
	typeNode TypeNode
}

// foo : bar
func ParseAnnotation(ts *[]token) *AnnotationNode {
	tokens := *ts

	target := get[IdentifierToken](&tokens)
	if target == nil {
		return nil
	}

	maybeColon := get[IdentifierToken](&tokens)
	if maybeColon == nil {
		return nil
	}
	if maybeColon.Name != ":" {
		return nil
	}

	var typeNode TypeNode
	if maybeFunctionInside := parseFunctionInside(&tokens); maybeFunctionInside != nil {
		typeNode = maybeFunctionInside
	} else if maybeTypeNode := ParseType(&tokens); maybeTypeNode != nil {
		typeNode = maybeTypeNode
	}

	if typeNode == nil {
		return nil
	}

	commit(ts, tokens)
	return &AnnotationNode{
		target.Select(typeNode.SelF().End()),
		*target,
		typeNode,
	}
}

// *****************************************************************************
// Type Parsing
// *****************************************************************************

// FunctionTypeNode | ListTypeNode | IdentifierToken | StringToken
type TypeNode Node

// or nil
func ParseType(ts *[]token) TypeNode {
	tokens := *ts

	if maybeFunctionTypeNode := ParseFunctionType(&tokens); maybeFunctionTypeNode != nil {
		commit(ts, tokens)
		return maybeFunctionTypeNode
	}

	if maybeListTypeNode := parseListType(&tokens); maybeListTypeNode != nil {
		commit(ts, tokens)
		return maybeListTypeNode
	}

	if maybeIdentifier := get[IdentifierToken](&tokens); maybeIdentifier != nil {
		commit(ts, tokens)
		return maybeIdentifier
	}

	if maybeString := get[StringToken](&tokens); maybeString != nil {
		commit(ts, tokens)
		return maybeString
	}

	return nil
}

type FunctionTypeNode struct {
	lexer.Sel
	// or nil
	maybeGeneric *FunctionGenericNode
	paramTypes   []TypeNode
	returnType   TypeNode
}

func ParseFunctionType(ts *[]token) *FunctionTypeNode {
	tokens := *ts

	maybeLeftBracket := get[SymbolToken](&tokens)
	if maybeLeftBracket == nil {
		return nil
	}
	if maybeLeftBracket.Symbol != '[' {
		return nil
	}

	inside := parseFunctionInside(&tokens)
	if inside == nil {
		return nil
	}

	maybeRightBracket := get[SymbolToken](&tokens)
	if maybeRightBracket == nil {
		return nil
	}
	if maybeRightBracket.Symbol != ']' {
		return nil
	}

	commit(ts, tokens)
	return inside
}

func parseFunctionInside(ts *[]token) *FunctionTypeNode {
	tokens := *ts

	maybeGeneric := parseFunctionGeneric(&tokens)

	var params []TypeNode
	maybeFirstType := ParseType(&tokens)
	if maybeFirstType == nil {
		return nil
	}
	params = append(params, maybeFirstType)

	for {
		if len(tokens) == 0 {
			return nil
		}
		if maybeArrow, ok := tokens[0].(IdentifierToken); ok {
			if maybeArrow.Name == "->" {
				advance(&tokens)
				break
			}
		}
		maybeType := ParseType(&tokens)
		if maybeType == nil {
			return nil
		}
		params = append(params, maybeType)
	}

	maybeReturnType := ParseType(&tokens)
	if maybeReturnType == nil {
		return nil
	}

	commit(ts, tokens)
	sel := maybeFirstType.SelF().Select(maybeReturnType.SelF().End())
	return &FunctionTypeNode{sel, maybeGeneric, params, maybeReturnType}
}

type FunctionGenericNode struct {
	lexer.Sel
	params []IdentifierToken
}

func parseFunctionGeneric(ts *[]token) *FunctionGenericNode {
	tokens := *ts

	open := get[SymbolToken](&tokens)
	if open == nil {
		return nil
	}
	if open.Symbol != '{' {
		return nil
	}

	var params []IdentifierToken
	var closingBraceEnd int
	for {
		if len(tokens) == 0 {
			return nil
		}
		switch v := tokens[0].(type) {
		case SymbolToken:
			// closing brace is the only allowed symbol
			if v.Symbol == '}' {
				closingBraceEnd = v.End()
				goto commitFunctionGeneric
			} else {
				return nil
			}
		case IdentifierToken:
			params = append(params, v)
		default:
			return nil
		}
	}

commitFunctionGeneric:
	commit(ts, tokens)
	return &FunctionGenericNode{
		open.Select(closingBraceEnd),
		params,
	}
}

type ListTypeNode struct {
	lexer.Sel
	members []TypeNode
}

func parseListType(ts *[]token) *ListTypeNode {
	tokens := *ts
	maybeLeft := get[SymbolToken](&tokens)

	if maybeLeft == nil {
		return nil
	}

	var endOffset int
	var members []TypeNode

	var close byte
	if maybeLeft.Symbol == '(' {
		close = ')'
		goto listInside
	}
	if maybeLeft.Symbol == '[' {
		close = ']'
		goto listInside
	}
	return nil

listInside:
	for {
		if len(tokens) == 0 {
			return nil
		}
		if maybeClosingSymbol, ok := tokens[0].(SymbolToken); ok {
			if maybeClosingSymbol.Symbol == close {
				advance(&tokens)
				endOffset = maybeClosingSymbol.End()
				goto listFinished
			}
		}
		maybeType := ParseType(&tokens)
		if maybeType == nil {
			return nil
		}
		members = append(members, maybeType)
	}

listFinished:
	commit(ts, tokens)
	return &ListTypeNode{
		Sel:     maybeLeft.Select(endOffset),
		members: members,
	}
}

// *****************************************************************************
// General Purpose Parser Functions
// *****************************************************************************

// maybe we need another one of these for getting and ignoring newlines
func get[T token](ts *[]token) *T {
	tokens := *ts
	if len(tokens) == 0 {
		return nil
	}
	maybeToken, ok := tokens[0].(T)
	if ok {
		advance(&tokens)
		commit(ts, tokens)
		return &maybeToken
	}
	return nil
}

func advance(ts *[]token) {
	tokens := *ts
	tokens = tokens[1:]
	commit(ts, tokens)
}

func commit(ts *[]token, newTokens []token) {
	*ts = newTokens
}
