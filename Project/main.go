package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"project/travel"
)

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

	err := filepath.Walk(directory, travel.WalkTraversal)

	if err != nil {
		fmt.Println("Error:", err)
	}
}
