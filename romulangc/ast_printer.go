package main

import (
	"strconv"
	"strings"
)

type ASTPrinter struct {
	indentLevel int
	result      string
}

func (ap *ASTPrinter) Visit(node interface{}) {
	switch n := node.(type) {
	case *SourceFile:
		ap.result += indent(ap.indentLevel) + "SourceFile [" + n.Namespace +
			"]\n"
	case *Declaration:
		ap.result += indent(ap.indentLevel) + "Declaration\n"
	case *Storyworld:
		ap.result += indent(ap.indentLevel) + "Storyworld\n"
	case *StoryworldBlock:
		ap.result += indent(ap.indentLevel) + "StoryworldBlock\n"
	case *Meta:
		ap.result += indent(ap.indentLevel) + "Meta\n"
	case *MetaEntry:
		ap.result += indent(ap.indentLevel) + "MetaEntry (" + *n.Name + " = " +
			*n.Value + ")\n"
	case *Vars:
		ap.result += indent(ap.indentLevel) + "Vars\n"
	case *VarDecl:
		ap.result += indent(ap.indentLevel) + "VarDecl (" + *n.Name + ": " +
			*n.Type + " = " + *n.InitialValue + ")\n"
	case *Passage:
		ap.result += indent(ap.indentLevel) + "Passage (" + *n.Name + "@" +
			strconv.Itoa(*n.Version) + "(): " + *n.ReturnType + "\n"
	case *Assignment:
		ap.result += indent(ap.indentLevel) + "Assignment (" + *n.Var + " = " +
			*n.Value + ")\n"
	}
	ap.indentLevel++
}

func (ap *ASTPrinter) Leave(node interface{}) {
	ap.indentLevel--
}

// indent returns a string good for indenting code level levels deep.
func indent(level int) string {
	return strings.Repeat("\t", level)
}
