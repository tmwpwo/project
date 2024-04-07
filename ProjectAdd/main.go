package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Specify the directory path
	// Define a flag for the directory path
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

func walkTraversal(path string, info os.FileInfo, err error) error {
	if err == nil && !info.IsDir() {
		extension := filepath.Ext(path)
		fileName := strings.TrimSuffix(info.Name(), extension)
		fmt.Printf("File: %s, Extension: %s\n", fileName, extension)
	}
	return nil
}
