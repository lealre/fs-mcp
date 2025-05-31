package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func fileSystemMCP(dir string) *server.MCPServer {

	mcpServer := server.NewMCPServer(
		"fs-mcp-server",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)

	h := &handler{baseDir: dir}

	listEntriesTool := mcp.NewTool("listEntries",
		mcp.WithDescription("List entries at a given path"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Path for which to list all entries"),
		),
		mcp.WithNumber("depth",
			mcp.Description("Depth of the directory tree (default is 3)"),
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
		mcp.WithDescription(
			"Retrieve file information including size, last modified time, "+
				"detected MIME type, and file permissions",
		),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Path to the file to retrieve information from"),
		),
	)
	renamePath := mcp.NewTool("renamePath",
		mcp.WithDescription("Renames a file or directory to a new name"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Path to the file or directory to be renamed"),
		),
		mcp.WithString("newPathFinalName",
			mcp.Required(),
			mcp.Description("New name for the file or directory (just the name, not the full path)"),
		),
	)
	copyFileOrDir := mcp.NewTool("copyFileOrDir",
		mcp.WithDescription("Copies a file or directory to a new location"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Path to the file or directory to be copied"),
		),
		mcp.WithString("destination",
			mcp.Required(),
			mcp.Description("Destination path where the file or directory will be copied"),
		),
	)

	mcpServer.AddTool(
		listEntriesTool,
		handlersMiddleware("listEntries", h.withSafePath(h.handlerListEntries)),
	)
	mcpServer.AddTool(
		readFileTool,
		handlersMiddleware("readFromFile", h.withSafePath(h.handlerReadFile)),
	)
	mcpServer.AddTool(
		writeFileTool,
		handlersMiddleware("writeToFile", h.withSafePath(h.handlerWriteToFile)),
	)
	mcpServer.AddTool(
		getFileInfo,
		handlersMiddleware("getFileInfo", h.withSafePath(h.handlerGetFileInfo)),
	)
	mcpServer.AddTool(
		renamePath,
		handlersMiddleware("renamePath", h.withSafePath(h.hadlerRenamePath)),
	)
	mcpServer.AddTool(
		copyFileOrDir,
		handlersMiddleware("copyFileOrDir", h.withSafePath(h.hadlerCopyFileOrDir)),
	)

	return mcpServer
}

func main() {
	var port int
	var dir string

	flag.IntVar(&port, "port", 8080, "Port to listen on (optional)")
	flag.StringVar(&dir, "dir", "", "Directory to serve")

	flag.Parse()

	if dir == "" {
		log.Println("Error: --dir must be provided.")
		flag.Usage()
		os.Exit(1)
	}

	_, err, exists := assertPath(dir)
	if err != nil {
		log.Fatalf("Error reading the base path: %v", err)
	}
	if !exists {
		log.Fatalf("Base path not found: %v", dir)
	}

	mcpServer := fileSystemMCP(dir)

	addr := fmt.Sprintf(":%d", port)
	sseServer := server.NewSSEServer(mcpServer, server.WithBaseURL("http://localhost"+addr))
	log.Printf("SSE server listening on %s", addr)
	if err := sseServer.Start(addr); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
