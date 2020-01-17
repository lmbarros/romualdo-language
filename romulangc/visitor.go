package main

// A Visitor has all the methods needed to traverse a Romualdo parse tree.
type Visitor interface {
	Visit(node interface{})
	Leave(node interface{})
}

func (sf *SourceFile) Walk(v Visitor) {
	v.Visit(sf)

	for _, d := range sf.Declarations {
		d.Walk(v)
	}

	v.Leave(sf)
}

func (d *Declaration) Walk(v Visitor) {
	v.Visit(d)

	switch {
	case d.Storyworld != nil:
		d.Storyworld.Walk(v)

	case d.Passage != nil:
		d.Passage.Walk(v)
	}

	v.Leave(d)
}

func (s *Storyworld) Walk(v Visitor) {
	v.Visit(s)

	for _, b := range s.StoryworldBlocks {
		b.Walk(v)
	}

	v.Leave(s)
}

func (s *StoryworldBlock) Walk(v Visitor) {
	v.Visit(s)

	switch {
	case s.Meta != nil:
		s.Meta.Walk(v)

	case s.Vars != nil:
		s.Vars.Walk(v)
	}

	v.Leave(s)
}

func (m *Meta) Walk(v Visitor) {
	v.Visit(m)

	for _, m := range m.MetaEntries {
		m.Walk(v)
	}

	v.Leave(m)
}

func (m *MetaEntry) Walk(v Visitor) {
	v.Visit(m)
	v.Leave(m)
}

func (vs *Vars) Walk(v Visitor) {
	v.Visit(vs)

	for _, vd := range vs.VarDecls {
		vd.Walk(v)
	}

	v.Leave(vs)
}

func (vd *VarDecl) Walk(v Visitor) {
	v.Visit(vd)
	v.Leave(vd)
}

func (p *Passage) Walk(v Visitor) {
	v.Visit(p)

	for _, s := range p.Body {
		s.Walk(v)
	}

	v.Leave(p)
}

func (a *Assignment) Walk(v Visitor) {
	v.Visit(a)
	v.Leave(a)
}
