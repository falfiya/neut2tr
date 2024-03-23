package rewrite

import (
	"fmt"
	"log"
	. "neut2tr/neut"
	"strings"
)

// None of these are user errors
func Rewrite(node Node) string {
	if node == nil {
		log.Panicf("Cannot rewrite nil!")
	}
	// pointer or no pointer?
	// don't you love when the parser internals are exposed to another module?
	// when the
	// when the abstraction is leaky
	switch n := node.(type) {
	case *TemplateNode:
		return fmt.Sprintf("#;%s", Rewrite(n.Contents))
	case *DeclarationNode:
		return fmt.Sprintf("(define-type %s %s)", Rewrite(n.Target), Rewrite(n.Value))
	case *AliasNode:
		return Rewrite(n.Value)
	case *AnnotationNode:
		return fmt.Sprintf("(: %s %s)", Rewrite(n.Target), Rewrite(n.Value))
	case *GenericTargetNode:
		return fmt.Sprintf("(%s %s)", Rewrite(n.Target), RewriteList(n.Params))
	case *SumTypeNode:
		return fmt.Sprintf("(U %s)", RewriteList(n.Terms))
	case SumTypeTermNode:
		return Rewrite(n.Value)
	case *QuoteNode:
		return fmt.Sprintf("'%s", Rewrite(n.Value))
	case *FunctionNode:
		if len(n.ParamTypes) == 0 {
			log.Panicf("SANITY: Function cannot be nullary!")
		}
		nonGenericFunction := fmt.Sprintf("(-> %s %s)", RewriteList(n.ParamTypes), Rewrite(n.ReturnType))
		if n.MaybeGeneric == nil {
			return nonGenericFunction
		} else {
			return fmt.Sprintf("(All %s %s)", Rewrite(n.MaybeGeneric), nonGenericFunction)
		}
	case *FunctionGenericNode:
		if len(n.Params) == 0 {
			log.Panicf("SANITY: Must have at least 1 generic parameter!")
			// because if you're going to specify that there are generic parameters
			// using the syntax, you had better at least use it.
			// that means something like this:
			// {} Number -> Number
			// is not allowed
		}
		return fmt.Sprintf("(%s)", RewriteList(n.Params))
	case *ListNode:
		switch n.Symbol {
		case ')':
			if len(n.Members) == 0 {
				return "()"
			}
			return fmt.Sprintf("(%s)", RewriteList(n.Members))
		case ']':
			if len(n.Members) == 0 {
				return "[]"
			}
			return fmt.Sprintf("[%s]", RewriteList(n.Members))
		default:
			log.Panicf("SANITY: ListNode can only be ')' or ']' but was %v!", n.Symbol)
		}
	case *BooleanNode:
		if n.Value {
			return "#t"
		} else {
			return "#f"
		}
	case *IdentifierToken:
		return n.Name
	// ALL ACCORDING TO ROB PIKE'S VISION
	case IdentifierToken:
		return n.Name
	case *StringToken:
		return n.Literal
	default:
		log.Panicf("Unhandled type %T!", node)
	}
	panic("Unreachable")
}

func RewriteList[T Node](list []T) string {
	if len(list) == 0 {
		log.Panicf("RewriteList does not handle the empty slice case!")
	}
	var sb strings.Builder
	sb.WriteString(Rewrite(list[0]))
	for _, t := range list[1:] {
		sb.WriteByte(' ')
		sb.WriteString(Rewrite(t))
	}
	return sb.String()
}
