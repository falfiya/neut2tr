package neut

import (
	. "neut2tr/util"
)

func (t TemplateNode) Print() string {
	return "" +
		"TemplateNode {\n" +
		Reindent(t.Contents.Print(), 1) +
		"}\n"
}

func (d DeclarationNode) Print() string {
	return "" +
		"DeclarationNode {\n" +
		Reindent(d.Target.Print(), 1) +
		Reindent(d.Value.Print(), 1) +
		"}\n"
}

func (g GenericTargetNode) Print() string {
	params := "<" + g.Params[0].CmpName
	for _, p := range g.Params[1:] {
		params += ", " + p.Print()
	}
	params += ">"
	return "GenericTargetNode { " + g.Target.Print() + params + " }\n"
}

func (s SumTypeNode) Print() string {
	out := "SumTypeNode {\n"
	for _, s := range s.Terms {
		out += Reindent(s.Print(), 1)
	}
	out += "}\n"
	return out
}

func (s SumTypeTermNode) Print() string {
	return "" +
		"SumTypeTermNode {\n" +
		Reindent(s.Value.Print(), 1) +
		"}\n"
}

func (a AliasNode) Print() string {
	return "" +
		"AliasNode {\n" +
		Reindent(a.Value.Print(), 1) +
		"}\n"
}

func (a AnnotationNode) Print() string {
	return "" +
		"AnnotationNode {\n" +
		"   Target { " + a.Target.Print() + " }\n" +
		Reindent(a.Value.Print(), 1) +
		"}\n"
}

func (q QuoteNode) Print() string {
	return "" +
		"QuoteNode {\n" +
		Reindent(q.Value.Print(), 1) +
		"}\n"
}

func (f FunctionNode) Print() string {
	out := "FunctionNode {\n"
	if f.MaybeGeneric != nil {
		out += Reindent(f.MaybeGeneric.Print(), 1)
	}
	out += "   Params {\n"
	for _, param := range f.ParamTypes {
		out += Reindent(param.Print(), 2)
	}
	out += "   }\n"
	out += "   Return {\n"
	out += Reindent(f.ReturnType.Print(), 2)
	out += "   }\n"
	out += "}\n"
	return out
}

func (f FunctionGenericNode) Print() string {
	out := "FunctionGenericNode { " + f.Params[0].Print()
	for _, param := range f.Params[1:] {
		out += ", " + param.Print()
	}
	out += " }\n"
	return out
}

func (l ListNode) Print() string {
	out := "ListNode {\n"
	for _, member := range l.Members {
		out += Reindent(member.Print(), 1)
	}
	out += "}\n"
	return out
}

func (b BooleanNode) Print() string {
	if b.Value {
		return "#t"
	} else {
		return "#f"
	}
}

func (i IdentifierToken) Print() string {
	return "#%" + i.Name
}

func (s StringToken) Print() string {
	return s.Literal
}
