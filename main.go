package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Print("Usage: go run main.go <path>")
		os.Exit(1)
	}

	fmt.Println("Printing all args passed")
	for i, arg := range os.Args {
		fmt.Printf("Arg %d: %s\n", i, arg)
	}
	path := os.Args[1]

	listEntries(path)

}

func listEntries(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("path not found")

	}
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory")
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("error reading the directory: %s", err)
	}

	for _, entry := range entries {
		pathType := "file"
		if entry.IsDir() {
			pathType = "directory"
		}
		fmt.Printf("- %s (%s)\n", entry.Name(), pathType)
	}

	return nil

}
