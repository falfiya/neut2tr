package neut

import "strings"

func reindent(s string, times int) string {
	indent := strings.Repeat("   ", times)
	lines := strings.Split(s, "\n")
	out := ""
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		out += indent + line + "\n"
	}
	return out
}

func (t TemplateNode) Print() string {
	return "" +
		"TemplateNode {\n" +
		reindent(t.decl.Print(), 1) +
		"}\n"
}

func (d DeclarationNode) Print() string {
	return "" +
		"DeclarationNode {\n" +
		reindent(d.Target.Print(), 1) +
		reindent(d.Value.Print(), 1) +
		"}\n"
}

func (g GenericTargetNode) Print() string {
	params := "<" + g.params[0].Name
	for _, p := range g.params[1:] {
		params += ", " + p.Print()
	}
	params += ">"
	return "GenericTargetNode { " + g.target.Print() + params + " }\n"
}

func (s SumTypeNode) Print() string {
	out := "SumTypeNode {\n"
	for _, s := range s.terms {
		out += reindent(s.Print(), 1)
	}
	out += "}\n"
	return out
}

func (s SumTypeTermNode) Print() string {
	return "" +
		"SumTypeTermNode {\n" +
		reindent(s.TypeNode.Print(), 1) +
		"}\n"
}

func (a AliasNode) Print() string {
	return "" +
		"AliasNode {\n" +
		reindent(a.TypeNode.Print(), 1) +
		"}\n"
}

func (a AnnotationNode) Print() string {
	return "" +
		"AnnotationNode {\n" +
		"   Target { " + a.target.Print() + " }\n" +
		reindent(a.value.Print(), 1) +
		"}\n"
}

func (q QuoteNode) Print() string {
	return "" +
		"QuoteNode {\n" +
		reindent(q.typeNode.Print(), 1) +
		"}\n"
}

func (f FunctionNode) Print() string {
	out := "FunctionNode {\n"
	if f.maybeGeneric != nil {
		out += reindent(f.maybeGeneric.Print(), 1)
	}
	out += "   Params {\n"
	for _, param := range f.paramTypes {
		out += reindent(param.Print(), 2)
	}
	out += "   }\n"
	out += "   Return {\n"
	out += reindent(f.returnType.Print(), 2)
	out += "   }\n"
	out += "}\n"
	return out
}

func (f FunctionGenericNode) Print() string {
	out := "FunctionGenericNode { " + f.params[0].Print()
	for _, param := range f.params[1:] {
		out += ", " + param.Print()
	}
	out += " }\n"
	return out
}

func (l ListNode) Print() string {
	out := "ListNode {\n"
	for _, member := range l.members {
		out += reindent(member.Print(), 1)
	}
	out += "}\n"
	return out
}

func (i IdentifierToken) Print() string {
	return "#%" + i.Name
}

func (s StringToken) Print() string {
	return s.Literal
}
