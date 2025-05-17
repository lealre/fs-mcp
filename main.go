package main

import (
	"log"

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
	mcpServer.AddTool(listEntriesTool, handlerListEntries)

	readFileTool := mcp.NewTool("read",
		mcp.WithDescription("Read the contents of a file at a given path"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Path to the file to be read"),
		),
	)
	mcpServer.AddTool(readFileTool, handlerReadFile)

	writeFileTool := mcp.NewTool("write",
		mcp.WithDescription("Create or overwrite a file with the given content"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Path to the file to write to"),
		),
		mcp.WithString("content",
			mcp.Required(),
			mcp.Description("Content to write to the file"),
		),
	)
	mcpServer.AddTool(writeFileTool, handlerWriteToFile)

	sseServer := server.NewSSEServer(mcpServer, server.WithBaseURL("http://localhost:8080"))
	log.Printf("SSE server listening on :8080")
	if err := sseServer.Start(":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}

}
