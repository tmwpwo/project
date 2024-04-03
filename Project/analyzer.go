package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

type CodeAnalyzer struct {
	Errors []string
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
	return a
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(os.Args)
		fmt.Println("Usage: go run analyzer.go <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]
	code, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}

	analyzer := &CodeAnalyzer{}
	analyzer.Analyze(string(code))

	if len(analyzer.Errors) > 0 {
		for _, err := range analyzer.Errors {
			fmt.Println(err)
		}
	} else {
		fmt.Println("No issues found.")
	}

	fmt.Println("List of all imports:")
	for _, imp := range analyzer.Imports {
		fmt.Println(imp)
	}
}
