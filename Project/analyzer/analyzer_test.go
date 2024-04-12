package analyzer_test

import (
	"go/ast"
	"go/token"
	"project/analyzer"
	"testing"
)

func TestCheckFunctionNames(t *testing.T) {
	codeAnalyzer := &analyzer.CodeAnalyzer{}
	funcDecl := &ast.FuncDecl{
		Name: &ast.Ident{Name: "TestFunction"},
	}
	codeAnalyzer.CheckFunctionNames(funcDecl)
	if len(codeAnalyzer.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(codeAnalyzer.Errors))
	}
}

func TestCheckHardcodedCredentials(t *testing.T) {
	codeAnalyzer := &analyzer.CodeAnalyzer{}
	assignStmt := &ast.AssignStmt{
		Rhs: []ast.Expr{
			&ast.BasicLit{Value: `"password"`},
		},
	}
	codeAnalyzer.CheckHardcodedCredentials(assignStmt)
	if len(codeAnalyzer.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(codeAnalyzer.Errors))
	}
}

func TestCheckSQLInjection(t *testing.T) {
	codeAnalyzer := &analyzer.CodeAnalyzer{}
	callExpr := &ast.CallExpr{
		Args: []ast.Expr{
			&ast.BasicLit{Value: `"SELECT * FROM users"`},
		},
	}
	codeAnalyzer.CheckSQLInjection(callExpr)
	if len(codeAnalyzer.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(codeAnalyzer.Errors))
	}
}

func TestListImports(t *testing.T) {
	codeAnalyzer := &analyzer.CodeAnalyzer{}
	genDecl := &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{Path: &ast.BasicLit{Value: `"fmt"`}},
		},
	}
	codeAnalyzer.ListImports(genDecl)
	if len(codeAnalyzer.Imports) != 1 {
		t.Errorf("Expected 1 import, got %d", len(codeAnalyzer.Imports))
	}
}

func TestAnalyze(t *testing.T) {
	codeAnalyzer := &analyzer.CodeAnalyzer{}
	code := `
	package main

import (
	"database/sql"
	"fmt"
)

func main() {

	secret := "secretKey"

	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer db.Close()

	userInput := "1 OR 1=1"
	query := fmt.Sprintf("INSERT INTO users (name) VALUES ('%s')", userInput) // SQL injection vulnerability
	_, err = db.Exec(query)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("User added successfully.")
}s
	`
	codeAnalyzer.Analyze(code)
	if len(codeAnalyzer.Errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(codeAnalyzer.Errors))
	}
	if len(codeAnalyzer.Imports) != 2 {
		t.Errorf("Expected 2 imports, got %d", len(codeAnalyzer.Imports))
	}
}
