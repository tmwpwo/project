package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"project/analyzer"
	"strings"
)

func WalkTraversal(directory string, info os.FileInfo, err error, resultChan chan<- error) {

	if err != nil {
		resultChan <- fmt.Errorf("error accessing directory %s: %v", directory, err)
		return
	}

	if info.IsDir() {
		return
	}

	if filepath.Ext(directory) != ".go" || strings.Contains(info.Name(), "_test") {
		return
	}

	fileName := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))

	fmt.Printf("File: %s\n", fileName)

	code, err := os.ReadFile(directory)
	if err != nil {
		resultChan <- fmt.Errorf("error reading file %s: %v", directory, err)
		return
	}

	// Analyze the file content
	analyzer := &analyzer.CodeAnalyzer{}
	analyzer.Analyze(string(code))
	if len(analyzer.Errors) > 0 {
		for _, err := range analyzer.Errors {
			resultChan <- fmt.Errorf(" %s: %v", fileName, err)
		}
	} else {
		fmt.Println("No issues found in", fileName)
	}
	fmt.Println("List of imports in", fileName+":")

	for _, imp := range analyzer.Imports {
		fmt.Println("\t", imp)
	}

	resultChan <- nil
}

func main() {
	directoryPtr := flag.String("directory", "", "the directory path to walk")
	flag.Parse()

	directory := *directoryPtr
	if directory == "" {
		var err error
		directory, err = os.Getwd()
		if err != nil {
			fmt.Println("Error getting current directory:", err)
			os.Exit(1)
		}
	}

	resultChan := make(chan error)

	go func() {
		for err := range resultChan {
			if err != nil {
				fmt.Println("potential vulnerability:", err)
			}
		}
	}()

	// Walk the directory concurrently
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		WalkTraversal(path, info, err, resultChan)
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the directory: %s\n", err)
	}

	close(resultChan)
}
