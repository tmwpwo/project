package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
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

func (a *CodeAnalyzer) CheckSQLInjection(node ast.Node) {
	sqlInjectionPattern := regexp.MustCompile(`(?:SELECT|INSERT|UPDATE|DELETE)\s+.+\s+(?:FROM|INTO|VALUES)\s*\(`)
	switch n := node.(type) {
	case *ast.CallExpr:
		// Check if it's a function call to database query/execution function
		// For demonstration, assuming it's db.Query() or db.Exec()
		if len(n.Args) > 0 {
			for _, arg := range n.Args {
				// Check if the argument is a basic literal
				if lit, ok := arg.(*ast.BasicLit); ok {
					if sqlInjectionPattern.MatchString(lit.Value) {
						a.Errors = append(a.Errors, "Potential SQL injection vulnerability found")
					}
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
	a.CheckSQLInjection(node)
	return a
}
