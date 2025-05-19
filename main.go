package main

import (
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func fileSystemMCP() *server.MCPServer {

	mcpServer := server.NewMCPServer(
		"fs-mcp-server",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)

	listEntriesTool := mcp.NewTool("listEntries",
		mcp.WithDescription("List entries for a given path"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Path to list all entries"),
		),
	)
	readFileTool := mcp.NewTool("readFromFile",
		mcp.WithDescription("Read the contents of a file at a given path"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Path to the file to be read"),
		),
	)
	writeFileTool := mcp.NewTool("writeToFile",
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
	getFileInfo := mcp.NewTool("getFileInfo",
		mcp.WithDescription("Retrieve file information including size, last modified time, detected MIME type, and file permissions"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Path to the file to retrieve information from"),
		),
	)

	mcpServer.AddTool(listEntriesTool, handlersMiddleware(handlerListEntries))
	mcpServer.AddTool(readFileTool, handlersMiddleware(handlerReadFile))
	mcpServer.AddTool(writeFileTool, handlersMiddleware(handlerWriteToFile))
	mcpServer.AddTool(getFileInfo, handlersMiddleware(handlerGetFileInfo))

	return mcpServer
}

func main() {
	mcpServer := fileSystemMCP()

	sseServer := server.NewSSEServer(mcpServer, server.WithBaseURL("http://localhost:8080"))
	log.Printf("SSE server listening on :8080")
	if err := sseServer.Start(":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}

}
