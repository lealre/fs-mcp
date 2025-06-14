package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func fileSystemMCP(handlerCfg *handlerCfg) *server.MCPServer {

	mcpServer := server.NewMCPServer(
		"fs-mcp-server",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)

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
		handlersMiddleware("listEntries", handlerCfg.withSafePath(handlerCfg.handlerListEntries)),
	)
	mcpServer.AddTool(
		readFileTool,
		handlersMiddleware("readFromFile", handlerCfg.withSafePath(handlerCfg.handlerReadFile)),
	)
	mcpServer.AddTool(
		writeFileTool,
		handlersMiddleware("writeToFile", handlerCfg.withSafePath(handlerCfg.handlerWriteToFile)),
	)
	mcpServer.AddTool(
		getFileInfo,
		handlersMiddleware("getFileInfo", handlerCfg.withSafePath(handlerCfg.handlerGetFileInfo)),
	)
	mcpServer.AddTool(
		renamePath,
		handlersMiddleware("renamePath", handlerCfg.withSafePath(handlerCfg.hadlerRenamePath)),
	)
	mcpServer.AddTool(
		copyFileOrDir,
		handlersMiddleware("copyFileOrDir", handlerCfg.withSafePath(handlerCfg.hadlerCopyFileOrDir)),
	)

	return mcpServer
}

func fileSystemDockerMCP(handlerCfg *handlerCfg) *server.MCPServer {

	mcpServer := server.NewMCPServer(
		"fs-mcp-server",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)

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
		handlersMiddleware("listEntries", handlerCfg.withDockerPath(handlerCfg.handlerListEntries)),
	)
	mcpServer.AddTool(
		readFileTool,
		handlersMiddleware("readFromFile", handlerCfg.withDockerPath(handlerCfg.handlerReadFile)),
	)
	mcpServer.AddTool(
		writeFileTool,
		handlersMiddleware("writeToFile", handlerCfg.withDockerPath(handlerCfg.handlerWriteToFile)),
	)
	mcpServer.AddTool(
		getFileInfo,
		handlersMiddleware("getFileInfo", handlerCfg.withDockerPath(handlerCfg.handlerGetFileInfo)),
	)
	mcpServer.AddTool(
		renamePath,
		handlersMiddleware("renamePath", handlerCfg.withDockerPath(handlerCfg.hadlerRenamePath)),
	)
	mcpServer.AddTool(
		copyFileOrDir,
		handlersMiddleware("copyFileOrDir", handlerCfg.withDockerPath(handlerCfg.hadlerCopyFileOrDir)),
	)

	return mcpServer
}

func main() {

	var port int
	var dir string
	var transport string
	var dockerMode bool
	var volumeMapping string

	flag.IntVar(&port, "port", 8080, "Port to listen on (optional)")
	flag.StringVar(&dir, "dir", "", "Directory to serve")
	flag.StringVar(&transport, "t", "stdio", "Transport type (stdio or http)")
	flag.BoolVar(&dockerMode, "docker", false, "Enable Docker mode with volume mapping")
	flag.StringVar(&volumeMapping, "volume", "", "Volume mapping in format 'hostPath:containerPath'")

	flag.Parse()
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `fs-mcp - A simple filesystem MCP

Usage:
	fs-mcp --dir <directory> [--port <port>] [-t <transport>]

Options:
`)
		flag.PrintDefaults()
	}

	if dir == "" {
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

	if dockerMode {
		// TODO: if dockerMode is false, warns that the --volume wont be used
		// Create the config
		handlerCfg := &handlerCfg{baseDir: dir, dockerMode: true}
		if dockerMode && volumeMapping != "" {
			parts := strings.Split(volumeMapping, ":")
			if len(parts) == 2 {
				handlerCfg.volumeMapping = &VolumeMapping{
					HostPath:      parts[0],
					ContainerPath: parts[1],
				}
			}
		}

		// Start the server
		mcpServer := fileSystemDockerMCP(handlerCfg)
		addr := fmt.Sprintf(":%d", port)
		sseServer := server.NewSSEServer(mcpServer, server.WithBaseURL("http://localhost"+addr))
		log.Printf("SSE server listening inside Docker on %s", addr)
		if err := sseServer.Start(addr); err != nil {
			log.Fatalf("Server error: %v", err)
		}
		return
	}

	handlerCfg := &handlerCfg{baseDir: dir, dockerMode: false}

	mcpServer := fileSystemMCP(handlerCfg)
	if transport == "http" {
		addr := fmt.Sprintf(":%d", port)
		sseServer := server.NewSSEServer(mcpServer, server.WithBaseURL("http://localhost"+addr))
		log.Printf("SSE server listening on %s", addr)
		if err := sseServer.Start(addr); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	} else {
		if err := server.ServeStdio(mcpServer); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}
}
