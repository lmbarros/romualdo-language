package main

import (
	"strconv"
	"strings"

	"gitlab.com/stackedboxes/romulang/romulangc/ast"
)

type ASTPrinter struct {
	indentLevel int
	result      string
}

func (ap *ASTPrinter) Visit(node ast.Node) {
	switch n := node.(type) {
	case *ast.SourceFile:
		ap.result += indent(ap.indentLevel) + "SourceFile [" + n.Namespace +
			"]\n"
	case *ast.Declaration:
		ap.result += indent(ap.indentLevel) + "Declaration\n"
	case *ast.Storyworld:
		ap.result += indent(ap.indentLevel) + "Storyworld\n"
	case *ast.StoryworldBlock:
		ap.result += indent(ap.indentLevel) + "StoryworldBlock\n"
	case *ast.Meta:
		ap.result += indent(ap.indentLevel) + "Meta\n"
	case *ast.MetaEntry:
		ap.result += indent(ap.indentLevel) + "MetaEntry (" + *n.Name + " = " +
			*n.Value + ")\n"
	case *ast.Vars:
		ap.result += indent(ap.indentLevel) + "Vars\n"
	case *ast.VarDecl:
		ap.result += indent(ap.indentLevel) + "VarDecl (" + *n.Name + ": " +
			*n.Type + " = " + *n.InitialValue + ")\n"
	case *ast.Passage:
		ap.result += indent(ap.indentLevel) + "Passage (" + *n.Name + "@" +
			strconv.Itoa(*n.Version) + "(): " + *n.ReturnType + "\n"
	case *ast.Assignment:
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
