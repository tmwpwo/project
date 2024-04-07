package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
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

func walkTraversal(path string, info os.FileInfo, err error) error {
	if err == nil && !info.IsDir() {
		extension := filepath.Ext(path)
		if extension == ".go" {
			fileName := strings.TrimSuffix(info.Name(), extension)
			fmt.Printf("File: %s, Extension: %s\n", fileName, extension)

			// Read the content of the file and analyze
			code, err := os.ReadFile(path)
			if err != nil {
				fmt.Printf("Error reading file %s: %s\n", path, err)
				return nil
			}

			analyzer := &CodeAnalyzer{}
			analyzer.Analyze(string(code))

			if len(analyzer.Errors) > 0 {
				for _, err := range analyzer.Errors {
					fmt.Println("\t", err)
				}
			} else {
				fmt.Println("No issues found in", fileName)
			}

			fmt.Println("List of imports in", fileName+":")
			for _, imp := range analyzer.Imports {
				fmt.Println("\t", imp)
			}
		}
	}
	return nil
}

func main() {
	directoryPtr := flag.String("directory", "", "the directory path to walk")
	flag.Parse()

	// If directory flag is not provided, use the current working directory
	directory := *directoryPtr
	if directory == "" {
		var err error
		directory, err = os.Getwd()
		if err != nil {
			fmt.Println("Error getting current directory:", err)
			os.Exit(1)
		}
	}

	err := filepath.Walk(directory, walkTraversal)

	if err != nil {
		fmt.Println("Error:", err)
	}
}
