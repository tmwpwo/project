package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func WalkTraversal(path string, info os.FileInfo, err error) error {
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
