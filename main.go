package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	mcpServer := server.NewMCPServer(
		"fs-go-server",
		"1.0.0",
		// server.WithResourceCapabilities(true, true),
		// server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)

	listEntriesTool := mcp.NewTool("list",
		mcp.WithDescription("List entries for a given path"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Path to list all entries"),
		),
	)

	mcpServer.AddTool(listEntriesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path := request.Params.Arguments["operation"].(string)

		entries, err := listEntries(path)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("", err), err
		}

		return mcp.NewToolResultText(entries), nil
	})

	sseServer := server.NewSSEServer(mcpServer, server.WithBaseURL("http://localhost:8080"))
	log.Printf("SSE server listening on :8080")
	if err := sseServer.Start(":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}

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
