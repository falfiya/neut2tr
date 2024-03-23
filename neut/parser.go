package neut

import (
	"fmt"
	"neut2tr/lexer"
)

type Node interface {
	lexer.SelF
	Print() string
}

func Parse(s string) ([]Node, *TokenizerError) {
	tokens, te := Tokenize(s)
	if te != nil {
		return nil, te
	}
	return parseTokens(tokens), nil
}

// *****************************************************************************
// Parse Meta
// *****************************************************************************

// returns [](nil | TemplateNode | DeclarationNode | AnnotationNode)
func parseTokens(tokens []token) []Node {
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

	if anno := parseAnnotation(ts); anno != nil {
		return anno
	}

	return nil
}

// *****************************************************************************
// Higher-Level Syntax Elements
// *****************************************************************************

type TemplateNode struct {
	lexer.Sel
	// DeclarationNode | AnnotationNode
	Contents Node
}

func parseTemplate(ts *[]token) *TemplateNode {
	tokens := *ts

	maybeTemplate := get[IdentifierToken](&tokens)
	if maybeTemplate == nil {
		return nil
	}
	if maybeTemplate.CmpName != "template" && maybeTemplate.CmpName != "template:" {
		return nil
	}

	allowNewlines(&tokens)

	var contents Node
	if decl := parseDeclaration(&tokens); decl != nil {
		contents = decl
		goto commitTemplate
	}
	if anno := parseAnnotation(&tokens); anno != nil {
		contents = anno
		goto commitTemplate
	}
	return nil
commitTemplate:
	commit(ts, tokens)
	return &TemplateNode{maybeTemplate.Select(contents.SelF().End()), contents}
}

type DeclarationNode struct {
	lexer.Sel
	// IdentifierToken | GenericTargetNode
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
	if maybeIs == nil {
		return nil
	}

	if maybeIs.CmpName != "is" {
		return nil
	}

	var value Node
	if shouldParseSumType(&tokens) {
		sum := parseSumType(&tokens)
		if sum == nil {
			return nil
		}
		value = sum
	} else {
		alias := parseAlias(&tokens)
		if alias == nil {
			return nil
		}
		value = alias
	}

	commit(ts, tokens)
	sel := article1.Select(value.SelF().End())
	return &DeclarationNode{sel, target, value}
}

type GenericTargetNode struct {
	lexer.Sel
	Target IdentifierToken
	Params []IdentifierToken
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
	if open == nil {
		return nil
	}

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
	Terms []SumTypeTermNode
}

func shouldParseSumType(ts *[]token) bool {
	tokens := *ts

	one := get[IdentifierToken](&tokens)
	if one == nil {
		return false
	}
	if one.CmpName != "one" {
		return false
	}
	of := get[IdentifierToken](&tokens)
	if of == nil {
		return false
	}
	if of.CmpName != "of" && of.CmpName != "of:" {
		return false
	}

	commit(ts, tokens)
	return true
}

// - x
// - y
func parseSumType(ts *[]token) *SumTypeNode {
	tokens := *ts

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
	sel := terms[0].Select(terms[len(terms)-1].End())
	return &SumTypeNode{sel, terms}
}

type SumTypeTermNode struct {
	lexer.Sel
	Value TypeNode
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
	if hyphen.CmpName != "-" {
		return nil
	}

	typeNode := parseType(&tokens)
	if typeNode == nil {
		fmt.Print("died getting type\n")
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
	Value TypeNode
}

// ... (a|an) x
func parseAlias(ts *[]token) *AliasNode {
	tokens := *ts
	// doesn't matter if that's nil
	maybeArticle := parseArticle(&tokens)

	typeNode := parseType(&tokens)
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
	if maybeArticleToken.CmpName == "a" || maybeArticleToken.CmpName == "an" {
		commit(ts, tokens)
		return &ArticleNode{maybeArticleToken.Sel}
	} else {
		return nil
	}
}

type AnnotationNode struct {
	lexer.Sel
	Target IdentifierToken
	Value  TypeNode
}

// foo : bar
func parseAnnotation(ts *[]token) *AnnotationNode {
	tokens := *ts

	target := get[IdentifierToken](&tokens)
	if target == nil {
		return nil
	}

	maybeColon := get[IdentifierToken](&tokens)
	if maybeColon == nil {
		return nil
	}
	if maybeColon.CmpName != ":" {
		return nil
	}

	var typeNode TypeNode
	if maybeFunctionInside := parseFunctionInside(&tokens); maybeFunctionInside != nil {
		typeNode = maybeFunctionInside
	} else if maybeTypeNode := parseType(&tokens); maybeTypeNode != nil {
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
func parseType(ts *[]token) TypeNode {
	tokens := *ts

	if maybeQuoteNode := parseQuote(&tokens); maybeQuoteNode != nil {
		commit(ts, tokens)
		return maybeQuoteNode
	}

	if maybeFunctionTypeNode := parseFunction(&tokens); maybeFunctionTypeNode != nil {
		commit(ts, tokens)
		return maybeFunctionTypeNode
	}

	if maybeListTypeNode := parseList(&tokens); maybeListTypeNode != nil {
		commit(ts, tokens)
		return maybeListTypeNode
	}

	if maybeBoolean := parseBoolean(&tokens); maybeBoolean != nil {
		commit(ts, tokens)
		return maybeBoolean
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

type QuoteNode struct {
	lexer.Sel
	Value TypeNode
}

func parseQuote(ts *[]token) *QuoteNode {
	tokens := *ts

	allowNewlines(&tokens)

	maybeQuote := get[SymbolToken](&tokens)
	if maybeQuote == nil {
		return nil
	}

	if maybeQuote.Symbol != '\'' {
		return nil
	}

	maybeType := parseType(&tokens)
	if maybeType == nil {
		return nil
	}

	commit(ts, tokens)
	return &QuoteNode{maybeQuote.Select(maybeType.SelF().End()), maybeType}
}

type FunctionNode struct {
	lexer.Sel
	// or nil
	MaybeGeneric *FunctionGenericNode
	ParamTypes   []TypeNode
	ReturnType   TypeNode
}

func parseFunction(ts *[]token) *FunctionNode {
	tokens := *ts

	maybeLeftBracket := get[SymbolToken](&tokens)
	if maybeLeftBracket == nil {
		return nil
	}
	if maybeLeftBracket.Symbol != '[' {
		return nil
	}

	allowNewlines(&tokens)

	inside := parseFunctionInside(&tokens)
	if inside == nil {
		return nil
	}

	allowNewlines(&tokens)

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

func parseFunctionInside(ts *[]token) *FunctionNode {
	tokens := *ts

	var firstSel lexer.Sel

	maybeGeneric := parseFunctionGeneric(&tokens)
	if maybeGeneric != nil {
		allowNewlines(&tokens)
		firstSel = maybeGeneric.Sel
	}

	var params []TypeNode
	maybeFirstType := parseType(&tokens)
	if maybeFirstType == nil {
		return nil
	}
	if maybeGeneric == nil {
		firstSel = maybeFirstType.SelF()
	}
	params = append(params, maybeFirstType)

	for {
		if len(tokens) == 0 {
			return nil
		}
		if maybeArrow, ok := tokens[0].(IdentifierToken); ok {
			if maybeArrow.CmpName == "->" {
				advance(&tokens)
				break
			}
		}
		maybeType := parseType(&tokens)
		if maybeType == nil {
			return nil
		}
		params = append(params, maybeType)
	}

	allowNewlines(&tokens)

	maybeReturnType := parseType(&tokens)
	if maybeReturnType == nil {
		return nil
	}

	commit(ts, tokens)
	sel := firstSel.Select(maybeReturnType.SelF().End())
	return &FunctionNode{sel, maybeGeneric, params, maybeReturnType}
}

type FunctionGenericNode struct {
	lexer.Sel
	Params []IdentifierToken
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
		current := tokens[0]
		advance(&tokens)
		switch v := current.(type) {
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

type ListNode struct {
	lexer.Sel
	// one of ')' or ']'
	Symbol  byte
	Members []TypeNode
}

func parseList(ts *[]token) *ListNode {
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
		allowNewlines(&tokens)
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
		maybeType := parseType(&tokens)
		if maybeType == nil {
			return nil
		}
		members = append(members, maybeType)
	}

listFinished:
	commit(ts, tokens)
	return &ListNode{
		Sel:     maybeLeft.Select(endOffset),
		Symbol:  close,
		Members: members,
	}
}

type BooleanNode struct {
	lexer.Sel
	Value bool
}
func parseBoolean(ts *[]token) *BooleanNode {
	tokens := *ts

	maybeHash := get[SymbolToken](&tokens)
	if maybeHash == nil {
		return nil
	}
	if maybeHash.Symbol != '#' {
		return nil
	}

	maybeFT := get[IdentifierToken](&tokens)
	if maybeFT == nil {
		return nil
	}

	if maybeFT.CmpName == "f" {
		commit(ts, tokens)
		return &BooleanNode{
			maybeHash.Select(maybeFT.End()),
			false,
		}
	}

	if maybeFT.CmpName == "t" {
		commit(ts, tokens)
		return &BooleanNode{
			maybeHash.Select(maybeFT.End()),
			true,
		}
	}

	return nil
}

// *****************************************************************************
// General Purpose Parser Functions
// *****************************************************************************
func allowNewlines(ts *[]token) {
	for {
		maybeNewline := get[NewlineToken](ts)
		if maybeNewline == nil {
			break
		}
	}
}

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
