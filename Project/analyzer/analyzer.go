package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type CodeAnalyzer struct {
	Errors  []string
	Imports []string
}

func (a *CodeAnalyzer) CheckFunctionNames(node ast.Node) {
	switch n := node.(type) {
	case *ast.FuncDecl:
		if !strings.HasPrefix(n.Name.Name, strings.ToLower(string(n.Name.Name[0]))) {
			a.Errors = append(a.Errors, fmt.Sprintf("Function name '%s' should start with lowercase", n.Name.Name))
		}
	}
}

func (a *CodeAnalyzer) CheckHardcodedCredentials(node ast.Node) {
	switch n := node.(type) {
	case *ast.AssignStmt:
		for _, expr := range n.Rhs {
			if lit, ok := expr.(*ast.BasicLit); ok && lit.Kind == token.STRING {
				value := strings.Trim(lit.Value, `"`)
				if strings.Contains(value, "password") || strings.Contains(value, "secret") {
					a.Errors = append(a.Errors, fmt.Sprintf("Potential hardcoded credential found: %s", value))
				}
			}
		}
	case *ast.ValueSpec:
		if len(n.Values) == 0 {
			return // No values assigned
		}
		for _, val := range n.Values {
			if lit, ok := val.(*ast.BasicLit); ok && lit.Kind == token.STRING {
				value := strings.Trim(lit.Value, `"`)
				if strings.Contains(value, "password") || strings.Contains(value, "secret") {
					a.Errors = append(a.Errors, fmt.Sprintf("Potential hardcoded credential found: %s", value))
				}
			}
		}
	}
}

func (a *CodeAnalyzer) ListImports(node ast.Node) {
	switch n := node.(type) {
	case *ast.GenDecl:
		if n.Tok == token.IMPORT {
			for _, spec := range n.Specs {
				importSpec := spec.(*ast.ImportSpec)
				a.Imports = append(a.Imports, importSpec.Path.Value)
			}
		}
	}
}

func (a *CodeAnalyzer) Analyze(code string) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", code, parser.AllErrors)
	if err != nil {
		a.Errors = append(a.Errors, fmt.Sprintf("Syntax error: %s", err))
		return
	}

	ast.Walk(a, file)
}

func (a *CodeAnalyzer) Visit(node ast.Node) ast.Visitor {
	a.CheckFunctionNames(node)
	a.ListImports(node)
	a.CheckHardcodedCredentials(node)
	return a
}
