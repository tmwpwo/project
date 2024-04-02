package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

func main() {
    // Specify the directory path
    directory := "/home/tmwpwo/Project"

    // Walk through the directory and list all files
    err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
        // Check if it's a regular file
        if err == nil && !info.IsDir() {
            // Get the file extension
            extension := filepath.Ext(path)
            // Extract the file name without extension
            fileName := strings.TrimSuffix(info.Name(), extension)
            // Print file name along with its extension
            fmt.Printf("File: %s, Extension: %s\n", fileName, extension)
        }
        return nil
    })

    if err != nil {
        fmt.Println("Error:", err)
    }
}
