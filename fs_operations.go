package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"
)

type OperationResult struct {
	Content string
	Message string
	Error   error
}

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

func listEntries(path string, depth float64, prefix string) OperationResult {
	info, err, exists := assertPath(path)
	if err != nil {
		return OperationResult{Error: err}
	}
	if !exists {
		return OperationResult{Message: fmt.Sprintf("path not found at %s", path)}
	}
	if !info.IsDir() {
		return OperationResult{Message: "path is not a directory"}
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return OperationResult{Error: fmt.Errorf("error reading the directory: %s", err)}
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
			operationResult := listEntries(entryPath, depth-1, prefix+"  ")
			if operationResult.Error != nil {
				return OperationResult{Error: fmt.Errorf("error reading subDirectory: %s", err)}
			}
			allEntries += operationResult.Content
		}
	}
	return OperationResult{Content: allEntries}
}

func readFile(path string) OperationResult {
	info, err, exists := assertPath(path)
	if err != nil {
		return OperationResult{Error: err}
	}
	if !exists {
		return OperationResult{Message: fmt.Sprintf("path not found at %s", path)}
	}
	if info.IsDir() {
		return OperationResult{Message: "path is a directory, must be a file"}
	}

	content, err := os.ReadFile(path) // #nosec G304
	if err != nil {
		return OperationResult{Error: fmt.Errorf("error reading the file: %s", err)}
	}

	// Check if content is valid UTF-8 text
	if !utf8.Valid(content) {
		return OperationResult{Message: "file is not valid UTF-8 text (likely binary)"}
	}

	return OperationResult{Content: string(content)}
}

func writeToFile(content, path string) OperationResult {
	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0750)
		if err != nil {
			return OperationResult{Error: fmt.Errorf("could not create directory: %s", err)}
		}
	}

	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return OperationResult{Message: "path is a directory, must be a file"}
	}

	err = os.WriteFile(path, []byte(content), 0600)
	if err != nil {
		return OperationResult{Error: fmt.Errorf("could not write to file: %s", err)}
	}

	return OperationResult{Content: "file written successfully"}
}

func getFileInfo(path string) OperationResult {
	info, err, exists := assertPath(path)
	if err != nil {
		return OperationResult{Error: err}
	}
	if !exists {
		return OperationResult{Message: fmt.Sprintf("path not found at %s", path)}
	}
	if info.IsDir() {
		return OperationResult{Message: "path is a directory, must be a file"}
	}

	mimetype, err := getMimeType(path)
	if err != nil {
		return OperationResult{Error: err}
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

	return OperationResult{Content: fileInfo}
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

func renamePath(path, newName string) OperationResult {
	_, err, exists := assertPath(path)
	if err != nil {
		return OperationResult{Error: err}
	}

	if !exists {
		return OperationResult{Message: fmt.Sprintf("path not found at %s", path)}
	}

	fileDir := filepath.Dir(path)
	newPathName := filepath.Join(fileDir, newName)

	// Check if new name already exists
	if _, err := os.Stat(newPathName); err == nil {
		return OperationResult{Message: fmt.Sprintf("target path %s already exists", newPathName)}
	}

	err = os.Rename(path, newPathName)
	if err != nil {
		return OperationResult{Error: err}
	}

	return OperationResult{Content: newPathName}
}

func copyFileOrDir(path, dst string) OperationResult {
	fileInfo, err, exists := assertPath(path)
	if err != nil {
		return OperationResult{Error: err}
	}
	if !exists {
		return OperationResult{Message: fmt.Sprintf("path not found at %s", path)}
	}

	if fileInfo.IsDir() {
		return copyDir(path, dst)
	}
	return copyFile(path, dst)
}

func copyFile(path, destination string) OperationResult {
	sourceFile, err := os.Open(path)
	if err != nil {
		return OperationResult{Error: err}
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destination)
	if err != nil {
		return OperationResult{Error: err}
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return OperationResult{Error: err}
	}

	err = destFile.Sync()
	if err != nil {
		return OperationResult{Error: err}
	}

	return OperationResult{Content: "File copied to destination"}
}

func copyDir(path, dst string) OperationResult {
	pathInfo, err := os.Stat(path)
	if err != nil {
		return OperationResult{Error: err}
	}

	// Create the destination directory
	if err := os.MkdirAll(dst, pathInfo.Mode()); err != nil {
		return OperationResult{Error: err}
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return OperationResult{Error: err}
	}

	for _, entry := range entries {
		srcPath := filepath.Join(path, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if operationResult := copyDir(srcPath, dstPath); operationResult.Error != nil {
				return operationResult
			}
		} else {
			if operationResult := copyFile(srcPath, dstPath); operationResult.Error != nil {
				return operationResult
			}
		}
	}

	return OperationResult{Error: err}
}
