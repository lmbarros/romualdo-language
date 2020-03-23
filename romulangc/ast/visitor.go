package ast

// A Visitor has all the methods needed to traverse a Romualdo parse tree.
type Visitor interface {
	Visit(node Node)
	Leave(node Node)
}

func (n *SourceFile) Walk(v Visitor) {
	v.Visit(n)

	for _, d := range n.Declarations {
		d.Walk(v)
	}

	v.Leave(n)
}

func (n *Declaration) Walk(v Visitor) {
	v.Visit(n)

	switch {
	case n.MetaBlock != nil:
		n.MetaBlock.Walk(v)
	case n.VarsBlock != nil:
		n.VarsBlock.Walk(v)
	case n.AliasDecl != nil:
		n.AliasDecl.Walk(v)
	case n.EnumDecl != nil:
		n.EnumDecl.Walk(v)
	case n.StructDecl != nil:
		n.StructDecl.Walk(v)
	case n.FunctionDecl != nil:
		n.FunctionDecl.Walk(v)
	case n.PassageDecl != nil:
		n.PassageDecl.Walk(v)
	}

	v.Leave(n)
}

func (n *MetaBlock) Walk(v Visitor) {
	v.Visit(n)
	for _, e := range n.Entries {
		e.Walk(v)
	}
	v.Leave(n)
}

func (n *VarDecl) Walk(v Visitor) {
	v.Visit(n)
	v.Leave(n)
}

func (n *VarsBlock) Walk(v Visitor) {
	v.Visit(n)

	for _, vd := range n.Vars {
		vd.Walk(v)
	}

	v.Leave(n)
}

func (n *AliasDecl) Walk(v Visitor) {
	v.Visit(n)
	n.Type.Walk(v)
	v.Leave(n)
}

func (n *EnumDecl) Walk(v Visitor) {
	v.Visit(n)
	v.Leave(n)
}

func (n *StructDecl) Walk(v Visitor) {
	v.Visit(n)

	for _, m := range n.Members {
		m.Walk(v)
	}

	v.Leave(n)
}

func (n *FunctionDecl) Walk(v Visitor) {
	v.Visit(n)

	if n.Parameters != nil {
		n.Parameters.Walk(v)
	}
	n.ReturnType.Walk(v)

	for _, s := range n.Body {
		s.Walk(v)
	}

	v.Leave(n)
}

func (n *Parameters) Walk(v Visitor) {
	v.Visit(n)

	n.First.Walk(v)

	for _, p := range n.Rest {
		p.Walk(v)
	}

	v.Leave(n)
}

func (n *Parameter) Walk(v Visitor) {
	v.Visit(n)
	n.Type.Walk(v)
	v.Leave(n)
}

func (n *PassageDecl) Walk(v Visitor) {
	v.Visit(n)

	n.Parameters.Walk(v)
	n.ReturnType.Walk(v)

	for _, s := range n.Body {
		s.Walk(v)
	}

	v.Leave(n)
}

func (n *Type) Walk(v Visitor) {
	v.Visit(n)

	switch {
	case n.Array != nil:
		n.Array.Walk(v)
	case n.Function != nil:
		n.Function.Walk(v)
	case n.QualifiedIdentifier != nil:
		n.QualifiedIdentifier.Walk(v)
	}

	v.Leave(n)
}

func (n *FunctionType) Walk(v Visitor) {
	v.Visit(n)
	n.ParameterTypes.Walk(v)
	n.ReturnType.Walk(v)
	v.Leave(n)
}

func (n *TypeList) Walk(v Visitor) {
	v.Visit(n)

	n.First.Walk(v)

	for _, t := range n.Rest {
		t.Walk(v)
	}

	v.Leave(n)
}

func (n *QualifiedIdentifier) Walk(v Visitor) {
	v.Visit(n)
	v.Leave(n)
}

func (n *Statement) Walk(v Visitor) {
	v.Visit(n)

	switch {
	case n.Expression != nil:
		n.Expression.Walk(v)
	case n.WhileStmt != nil:
		n.WhileStmt.Walk(v)
	case n.IfStmt != nil:
		n.IfStmt.Walk(v)
	case n.ReturnStmt != nil:
		n.ReturnStmt.Walk(v)
	case n.GotoStmt != nil:
		n.GotoStmt.Walk(v)
	case n.SayStmt != nil:
		n.SayStmt.Walk(v)
	}

	v.Leave(n)
}

func (n *WhileStmt) Walk(v Visitor) {
	v.Visit(n)

	n.Condition.Walk(v)

	for _, s := range n.Body {
		s.Walk(v)
	}

	v.Leave(n)
}

func (n *IfStmt) Walk(v Visitor) {
	v.Visit(n)

	n.Condition.Walk(v)

	for _, s := range n.Then {
		s.Walk(v)
	}

	for _, ei := range n.Elseifs {
		ei.Walk(v)
	}

	for _, s := range n.Else {
		s.Walk(v)
	}

	v.Leave(n)
}

func (n *Elseif) Walk(v Visitor) {
	v.Visit(n)
	n.Condition.Walk(v)
	for _, s := range n.Body {
		s.Walk(v)
	}
	v.Leave(n)
}

func (n *ReturnStmt) Walk(v Visitor) {
	v.Visit(n)
	n.What.Walk(v)
	v.Leave(n)
}

func (n *GotoStmt) Walk(v Visitor) {
	v.Visit(n)
	n.Passage.Walk(v)
	n.Arguments.Walk(v)
	v.Leave(n)
}

func (n *SayStmt) Walk(v Visitor) {
	v.Visit(n)
	n.What.Walk(v)
	v.Leave(n)
}

func (n *Arguments) Walk(v Visitor) {
	v.Visit(n)
	n.First.Walk(v)
	for _, a := range n.Rest {
		a.Walk(v)
	}
	v.Leave(n)
}

func (n *Expression) Walk(v Visitor) {
	v.Visit(n)
	n.Expression.Walk(v)
	v.Leave(n)
}

func (n *Assignment) Walk(v Visitor) {
	v.Visit(n)

	if n.QualifiedIdentifier != nil {
		n.QualifiedIdentifier.Walk(v)
		n.Value.Walk(v)
	} else {
		n.Next.Walk(v)
	}

	v.Leave(n)
}

func (n *LogicOr) Walk(v Visitor) {
	v.Visit(n)

	n.First.Walk(v)
	for _, e := range n.Rest {
		e.Walk(v)
	}

	v.Leave(n)
}

func (n *LogicAnd) Walk(v Visitor) {
	v.Visit(n)

	n.First.Walk(v)
	for _, e := range n.Rest {
		e.Walk(v)
	}

	v.Leave(n)
}

func (n *Equality) Walk(v Visitor) {
	v.Visit(n)

	n.First.Walk(v)
	for _, e := range n.Rest {
		e.Walk(v)
	}

	v.Leave(n)
}

func (n *Comparison) Walk(v Visitor) {
	v.Visit(n)

	n.First.Walk(v)
	for _, e := range n.Rest {
		e.Walk(v)
	}

	v.Leave(n)
}

func (n *Addition) Walk(v Visitor) {
	v.Visit(n)

	n.First.Walk(v)
	for _, e := range n.Rest {
		e.Walk(v)
	}

	v.Leave(n)
}

func (n *Multiplication) Walk(v Visitor) {
	v.Visit(n)

	n.First.Walk(v)
	for _, e := range n.Rest {
		e.Walk(v)
	}

	v.Leave(n)
}

func (n *Exponentiation) Walk(v Visitor) {
	v.Visit(n)

	n.Base.Walk(v)
	n.Exponent.Walk(v)

	v.Leave(n)
}

func (n *Unary) Walk(v Visitor) {
	v.Visit(n)

	switch {
	case n.Operand != nil:
		n.Operand.Walk(v)
	case n.Call != nil:
		n.Call.Walk(v)
	default:
		panic("Unary should have a non-nil member")
	}

	v.Leave(n)
}

func (n *Call) Walk(v Visitor) {
	v.Visit(n)

	n.Primary.Walk(v)

	for _, c := range n.Complements {
		c.Walk(v)
	}

	v.Leave(n)
}

func (n *CallComplement) Walk(v Visitor) {
	v.Visit(n)

	switch {
	case n.Arguments != nil:
		n.Arguments.Walk(v)
	case n.QualifiedIdentifier != nil:
		n.QualifiedIdentifier.Walk(v)
	case n.Index != nil:
		n.QualifiedIdentifier.Walk(v)
	default:
		panic("CallComponent should have a non-nil member")
	}

	v.Leave(n)
}

func (n *Primary) Walk(v Visitor) {
	v.Visit(n)

	switch {
	case n.Array != nil:
		n.Array.Walk(v)
	case n.Map != nil:
		n.Map.Walk(v)
	case n.QualifiedIdentifier != nil:
		n.QualifiedIdentifier.Walk(v)
	case n.ParenthesizedExpression != nil:
		n.ParenthesizedExpression.Walk(v)
	case n.Gosub != nil:
		n.Gosub.Walk(v)
	case n.Listen != nil:
		n.Listen.Walk(v)
	}

	v.Leave(n)
}

func (n *ArrayLiteral) Walk(v Visitor) {
	v.Visit(n)

	n.First.Walk(v)
	for _, e := range n.Rest {
		e.Walk(v)
	}

	v.Leave(n)
}

func (n *MapLiteral) Walk(v Visitor) {
	n.First.Walk(v)

	for _, e := range n.Rest {
		e.Walk(v)
	}

	v.Leave(n)
}

func (n *MapEntry) Walk(v Visitor) {
	v.Visit(n)
	n.Value.Walk(v)
	v.Leave(n)
}

func (n *Gosub) Walk(v Visitor) {
	v.Visit(n)
	n.QualifiedIdentifier.Walk(v)
	n.Arguments.Walk(v)
	v.Leave(n)
}
