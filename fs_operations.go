package main

import (
	"fmt"
	"os"
	"path/filepath"
	"unicode/utf8"
)

func assertPath(path string) (os.FileInfo, error, bool) {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, nil, false
	}
	if err != nil {
		return nil, fmt.Errorf("error: %s", err), false
	}

	return fileInfo, nil, true
}

func listEntries(path string) (string, error) {
	info, err, exists := assertPath(path)
	if err != nil {
		return "", err
	}

	if !exists {
		return fmt.Sprintf("path not found at %s", path), nil
	}

	if !info.IsDir() {
		return "", fmt.Errorf("path is not a directory")
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return "", fmt.Errorf("error reading the directory: %s", err)
	}

	allEntries := ""
	for _, entry := range entries {
		pathType := "file"
		if entry.IsDir() {
			pathType = "directory"
		}
		allEntries += fmt.Sprintf("- %s (%s)\n", entry.Name(), pathType)
	}

	return allEntries, nil
}

func readFile(path string) (string, error) {
	info, err, exists := assertPath(path)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return "", fmt.Errorf("path is a directory, must be a file")
	}

	if !exists {
		return "", fmt.Errorf("path not found at %s", path)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("error reading the file: %s", err)
	}

	// Check if content is valid UTF-8 text
	if !utf8.Valid(content) {
		return "", fmt.Errorf("file is not valid UTF-8 text (likely binary)")
	}

	return string(content), nil
}

func writeToFile(content, path string) (string, error) {
	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return "", fmt.Errorf("could not create directory: %s", err)
		}
	}

	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return "", fmt.Errorf("path is a directory, must be a file")
	}

	err = os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return "", fmt.Errorf("could not write to file: %s", err)
	}

	return "file written successfully", nil
}
