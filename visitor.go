package strchecker

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

// treeVisitor carries the package name and file name
// for passing it to the imports map, and the fileSet for
// retrieving the token.Position.
type treeVisitor struct {
	p           *Parser
	fileSet     *token.FileSet
	packageName string
	fileName    string
}

// Visit browses the AST tree for strings that could be potentially
// replaced by constants.
// A map of existing constants is built as well (-match-constant).
func (v *treeVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return v
	}
	switch t := node.(type) {
	// Scan for constants in an attempt to match strings with existing constants
	case *ast.GenDecl:
		if t.Tok != token.CONST {
			return v
		}
		v.processConst(t)
	// foo := "moo"
	case *ast.AssignStmt:

	// if foo == "moo"
	case *ast.BinaryExpr:
		v.processBinaryExpr(t)

	// case "foo":
	case *ast.CaseClause:
		v.processCaseExpr(t)

	// return "boo"
	case *ast.ReturnStmt:
		v.processReturnExpr(t)

	// fn("http://")
	case *ast.CallExpr:
		v.processCallExpr(t)
	}

	return v
}

func (v *treeVisitor) processConst(t *ast.GenDecl) {
	for _, spec := range t.Specs {
		val := spec.(*ast.ValueSpec)
		for _, str := range val.Values {
			v.processExpr(str, Const)
		}
	}
}
func (v *treeVisitor) processAssignStmt(t *ast.AssignStmt) {
	for _, rhs := range t.Rhs {
		v.processExpr(rhs, Assignment)
	}
}
func (v *treeVisitor) processCaseExpr(t *ast.CaseClause) {
	for _, item := range t.List {
		v.processExpr(item, Case)
	}
}
func (v *treeVisitor) processBinaryExpr(t *ast.BinaryExpr) {
	v.processExpr(t.X, Binary)
	v.processExpr(t.Y, Binary)
}

func (v *treeVisitor) processCallExpr(t *ast.CallExpr) {
	for _, item := range t.Args {
		v.processExpr(item, Call)
	}
}
func (v *treeVisitor) processReturnExpr(t *ast.ReturnStmt) {
	for _, item := range t.Results {
		v.processExpr(item, Return)
	}
}
func (v *treeVisitor) processExpr(t ast.Expr, typ Type) {
	if bl, ok := t.(*ast.BasicLit); ok {
		v.addString(bl.Value, bl.Pos(), typ)
	}
}

// addString adds a string in the map along with its position in the tree.
func (v *treeVisitor) addString(str string, pos token.Pos, typ Type) {
	ok, excluded := v.p.excludeTypes[typ]
	if ok && excluded {
		return
	}
	// Drop quotes if any
	if strings.HasPrefix(str, `"`) || strings.HasPrefix(str, "`") {
		str, _ = strconv.Unquote(str)
	}

	// Ignore empty strings
	if len(str) == 0 {
		return
	}
	//check string is an invalid string
	if !v.p.invalidStrReg.MatchString(str) {
		return
	}
	v.p.strs = append(v.p.strs, InvalidString{
		Str:         str,
		Position:    v.fileSet.Position(pos),
		packageName: v.packageName,
	})
}
