package main

import (
	"strings"
)

func (sf *SourceFile) PrintAST() string {
	return sf.toString(0)
}

func (sf *SourceFile) toString(level int) string {
	result := indent(level) + "SourceFile [" + sf.Namespace + "]\n"

	for _, d := range sf.Declarations {
		result += d.toString(level + 1)
	}

	return result
}

func (d *Declaration) toString(level int) string {
	result := indent(level) + "Declaration\n"

	if d.Storyworld != nil {
		result += d.Storyworld.toString(level + 1)
	}

	return result
}

func (s *Storyworld) toString(level int) string {
	result := indent(level) + "Storyworld\n"

	for _, b := range s.StoryworldBlocks {
		result += b.toString(level + 1)
	}

	return result
}

func (s *StoryworldBlock) toString(level int) string {
	result := indent(level) + "StoryworldBlock\n"

	switch {
	case s.Meta != nil:
		result += s.Meta.toString(level + 1)

	case s.Vars != nil:
		result += s.Vars.toString(level + 1)
	}

	return result
}

func (m *Meta) toString(level int) string {
	result := indent(level) + "Meta\n"

	for _, m := range m.MetaEntries {
		result += m.toString(level + 1)
	}

	return result
}

func (m *MetaEntry) toString(level int) string {
	return indent(level) + "MetaEntry (" + *m.Name + " = " + *m.Value + ")\n"
}

func (v *Vars) toString(level int) string {
	result := indent(level) + "Vars\n"

	for _, vd := range v.VarDecls {
		result += vd.toString(level + 1)
	}

	return result
}

func (vd *VarDecl) toString(level int) string {
	return indent(level) + "VarDecl (" + *vd.Name + ": " + *vd.Type + " = " +
		*vd.InitialValue + ")\n"
}

// indent returns a string good for indenting code level levels deep.
func indent(level int) string {
	return strings.Repeat("\t", level)
}
