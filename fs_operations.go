package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"
)

// isSafePath checks if the given path is within the base directory and does not contain directory traversal
func isSafePath(base, target string) bool {
	absBase, err := filepath.Abs(base)
	if err != nil {
		return false
	}
	absTarget, err := filepath.Abs(target)
	if err != nil {
		return false
	}
	return strings.HasPrefix(absTarget, absBase)
}

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

func listEntries(path string, depth float64, prefix string) (string, error) {
	info, err, exists := assertPath(path)
	if err != nil {
		return "", err
	}
	if !exists {
		return fmt.Sprintf("path not found at %s", path), nil
	}
	if !info.IsDir() {
		return "path is not a directory", nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return "", fmt.Errorf("error reading the directory: %s", err)
	}

	allEntries := ""
	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())
		pathType := "file"
		if entry.IsDir() {
			pathType = "directory"
		}
		allEntries += fmt.Sprintf("%s- %s (%s)\n", prefix, entry.Name(), pathType)
		if entry.IsDir() && depth != 0 {
			subEntries, err := listEntries(entryPath, depth-1, prefix+"  ")
			if err != nil {
				return "", err
			}
			allEntries += subEntries
		}
	}
	return allEntries, nil
}

func readFile(path string) (string, error) {
	info, err, exists := assertPath(path)
	if err != nil {
		return "", err
	}
	if !exists {
		return fmt.Sprintf("path not found at %s", path), nil
	}
	if info.IsDir() {
		return "path is a directory, must be a file", nil
	}

	content, err := os.ReadFile(path) // #nosec G304
	if err != nil {
		return "", fmt.Errorf("error reading the file: %s", err)
	}

	// Check if content is valid UTF-8 text
	if !utf8.Valid(content) {
		return "file is not valid UTF-8 text (likely binary)", nil
	}

	return string(content), nil
}

func writeToFile(content, path string) (string, error) {
	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0750)
		if err != nil {
			return "", fmt.Errorf("could not create directory: %s", err)
		}
	}

	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return "path is a directory, must be a file", nil
	}

	err = os.WriteFile(path, []byte(content), 0600)
	if err != nil {
		return "", fmt.Errorf("could not write to file: %s", err)
	}

	return "file written successfully", nil
}

func getFileInfo(path string) (string, error) {
	info, err, exists := assertPath(path)
	if err != nil {
		return "", err
	}
	if !exists {
		return fmt.Sprintf("path not found at %s", path), nil
	}
	if info.IsDir() {
		return "path is a directory, must be a file", nil
	}

	mimetype, err := getMimeType(path)
	if err != nil {
		return "", err
	}

	perms := info.Mode().String()
	modTime := info.ModTime().Format(time.RFC3339)

	fileInfo := fmt.Sprintf(
		"File: %s\n"+
			"Size: %d bytes\n"+
			"Permissions: %s\n"+
			"Last Modified: %s\n"+
			"MIME Type: %s\n",
		path,
		info.Size(),
		perms,
		modTime,
		mimetype,
	)

	return fileInfo, nil
}

func getMimeType(path string) (string, error) {
	file, err := os.Open(path) // #nosec G304
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read the first 512 bytes for MIME detection
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}

	mimeType := http.DetectContentType(buffer)
	return mimeType, nil
}

func renamePath(path, newName string) (string, error) {
	_, err, exists := assertPath(path)
	if err != nil {
		return "", err
	}

	if !exists {
		return fmt.Sprintf("path not found at %s", path), nil
	}

	fileDir := filepath.Dir(path)
	newPathName := filepath.Join(fileDir, newName)

	// Check if new name already exists
	if _, err := os.Stat(newPathName); err == nil {
		return fmt.Sprintf("target path %s already exists", newPathName), nil
	}

	os.Rename(path, newPathName)

	return newPathName, nil
}
