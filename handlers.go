package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"unicode/utf8"

	"github.com/mark3labs/mcp-go/mcp"
)

// List entries
func handlerListEntries(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := request.Params.Arguments["path"].(string)
	log.Printf("'handlerListEntries' called with path %s", path)

	entries, err := listEntries(path)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("", err), err
	}

	return mcp.NewToolResultText(entries), nil
}

func listEntries(path string) (string, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("path not found")
	}
	if err != nil {
		return "", fmt.Errorf("error: %s", err)
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

// Read file
func handlerReadFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := request.Params.Arguments["path"].(string)
	log.Printf("'handlerReadFile' called with path %s", path)

	entries, err := readFile(path)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("", err), err
	}

	return mcp.NewToolResultText(entries), nil
}

func readFile(path string) (string, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("path not found")
	}
	if err != nil {
		return "", fmt.Errorf("error: %s", err)
	}
	if info.IsDir() {
		return "", fmt.Errorf("path is a directory, must be a file")
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

// Read file
func handlerWriteToFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := request.Params.Arguments["path"].(string)
	content := request.Params.Arguments["content"].(string)

	// Log the function name dynamically
	pc, _, _, _ := runtime.Caller(0)
	funcName := runtime.FuncForPC(pc).Name()
	log.Printf("'%s' called with path %s", funcName, path)

	entries, err := writeToFile(content, path)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("", err), err
	}

	return mcp.NewToolResultText(entries), nil
}

func writeToFile(content, path string) (string, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			return "", fmt.Errorf("could not create file: %s", err)
		}
		return "file created and written", nil
	}
	if err != nil {
		return "", fmt.Errorf("error: %s", err)
	}

	if info.IsDir() {
		return "", fmt.Errorf("path is a directory, must be a file")
	}

	// Write the content (overwrite)
	err = os.WriteFile(path, []byte(content), info.Mode())
	if err != nil {
		return "", fmt.Errorf("error writing to file: %s", err)
	}

	return "file written successfully", nil
}
