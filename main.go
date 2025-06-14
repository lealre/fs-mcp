package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func createMCPServer(handlerCfg *handlerCfg, pathMiddleware func(handlerFunc) server.ToolHandlerFunc) *server.MCPServer {
	mcpServer := server.NewMCPServer(
		"fs-mcp-server",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)

	// Define all tools
	tools := []struct {
		name        string
		description string
		params      []mcp.ToolOption
		handler     func(ctx context.Context, path string, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
	}{
		{
			name:        "listEntries",
			description: "List entries at a given path",
			params: []mcp.ToolOption{
				mcp.WithString("path",
					mcp.Required(),
					mcp.Description("Path for which to list all entries"),
				),
				mcp.WithNumber("depth",
					mcp.Description("Depth of the directory tree (default is 3)"),
				),
			},
			handler: handlerCfg.handlerListEntries,
		},
		{
			name:        "readFromFile",
			description: "Read the contents of a file at a given path",
			params: []mcp.ToolOption{
				mcp.WithString("path",
					mcp.Required(),
					mcp.Description("Path to the file to be read"),
				),
			},
			handler: handlerCfg.handlerReadFile,
		},
		{
			name:        "writeToFile",
			description: "Create or overwrite a file with the given content",
			params: []mcp.ToolOption{
				mcp.WithString("path",
					mcp.Required(),
					mcp.Description("Path to the file to write to"),
				),
				mcp.WithString("content",
					mcp.Required(),
					mcp.Description("Content to write to the file"),
				),
			},
			handler: handlerCfg.handlerWriteToFile,
		},
		{
			name: "getFileInfo",
			description: "Retrieve file information including size, last modified time, " +
				"detected MIME type, and file permissions",
			params: []mcp.ToolOption{
				mcp.WithString("path",
					mcp.Required(),
					mcp.Description("Path to the file to retrieve information from"),
				),
			},
			handler: handlerCfg.handlerGetFileInfo,
		},
		{
			name:        "renamePath",
			description: "Renames a file or directory to a new name",
			params: []mcp.ToolOption{
				mcp.WithString("path",
					mcp.Required(),
					mcp.Description("Path to the file or directory to be renamed"),
				),
				mcp.WithString("newPathFinalName",
					mcp.Required(),
					mcp.Description("New name for the file or directory (just the name, not the full path)"),
				),
			},
			handler: handlerCfg.hadlerRenamePath,
		},
		{
			name:        "copyFileOrDir",
			description: "Copies a file or directory to a new location",
			params: []mcp.ToolOption{
				mcp.WithString("path",
					mcp.Required(),
					mcp.Description("Path to the file or directory to be copied"),
				),
				mcp.WithString("destination",
					mcp.Required(),
					mcp.Description("Destination path where the file or directory will be copied"),
				),
			},
			handler: handlerCfg.hadlerCopyFileOrDir,
		},
	}

	// Add all tools to the server
	for _, tool := range tools {
		t := mcp.NewTool(tool.name,
			append([]mcp.ToolOption{
				mcp.WithDescription(tool.description),
			}, tool.params...)...,
		)
		mcpServer.AddTool(
			t,
			handlersMiddleware(tool.name, pathMiddleware(tool.handler)),
		)
	}

	return mcpServer
}

func fileSystemMCP(handlerCfg *handlerCfg) *server.MCPServer {
	return createMCPServer(handlerCfg, handlerCfg.withSafePath)
}

func fileSystemDockerMCP(handlerCfg *handlerCfg) *server.MCPServer {
	return createMCPServer(handlerCfg, handlerCfg.withDockerPath)
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

	dockerMode = os.Getenv("FS_MCP_DOCKER_MODE") == "true"

	// directory resolution
	finalDir := ""
	switch {
	case dockerMode && volumeMapping != "":
		parts := strings.Split(volumeMapping, ":")
		if len(parts) != 2 {
			log.Fatal("Invalid volume format. Use hostPath:containerPath")
		}
		finalDir = parts[1]
	case dockerMode:
		finalDir = "/baseDir"
	case dir != "":
		finalDir = dir
	default:
		log.Fatal("\nError: -dir is required")
		flag.Usage()
	}

	_, err, exists := assertPath(finalDir)
	if err != nil {
		log.Fatalf("Error reading the base path: %v", err)
	}
	if !exists {
		log.Fatalf("Base path not found: %v", finalDir)
	}

	if dockerMode && volumeMapping == "" {
		log.Println("Warning: Docker mode is enabled but no volume mapping is specified")
	}

	// Choose the appropriate MCP server based on dockerMode
	var mcpServer *server.MCPServer
	handlerCfg := &handlerCfg{baseDir: finalDir, dockerMode: dockerMode}
	if dockerMode {
		if dockerMode && volumeMapping != "" {
			parts := strings.Split(volumeMapping, ":")
			if len(parts) == 2 {
				handlerCfg.volumeMapping = &VolumeMapping{
					HostPath:      parts[0],
					ContainerPath: parts[1],
				}
			}
		}
		mcpServer = fileSystemDockerMCP(handlerCfg)
	} else {
		mcpServer = fileSystemMCP(handlerCfg)
	}

	// Start the server based on transport
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
