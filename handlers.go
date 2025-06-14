package main

import (
	"context"
	"log"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type handlerCfg struct {
	baseDir       string
	dockerMode    bool
	volumeMapping *VolumeMapping
}

type VolumeMapping struct {
	HostPath      string
	ContainerPath string
}

func (h *handlerCfg) withSafePath(
	handler func(ctx context.Context, path string, request mcp.CallToolRequest) (*mcp.CallToolResult, error),
) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path := request.Params.Arguments["path"].(string)
		if !isSafePath(h.baseDir, path) {
			log.Printf("PATH NOT ALLOWED: path is outside of allowed base directory")
			return mcp.NewToolResultText("access denied: path is outside of allowed base directory"), nil
		}
		return handler(ctx, path, request)
	}
}

func (h *handlerCfg) withDockerPath(
	handler func(ctx context.Context, path string, request mcp.CallToolRequest) (*mcp.CallToolResult, error),
) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		hostPath := request.Params.Arguments["path"].(string)

		// Ensure the path is within the allowed host directory
		if !strings.HasPrefix(hostPath, h.volumeMapping.HostPath) {
			log.Printf("PATH NOT ALLOWED: %s is outside of %s", hostPath, h.volumeMapping.HostPath)
			return mcp.NewToolResultText("PATH NOT ALLOWED: path is outside of allowed directory"), nil
		}

		// Translate host path to container path
		relPath := strings.TrimPrefix(hostPath, h.volumeMapping.HostPath)
		containerPath := filepath.Join(h.volumeMapping.ContainerPath, filepath.Clean(relPath))

		log.Printf("Path translation:\nHost: %s\nContainer: %s", hostPath, containerPath)

		return handler(ctx, containerPath, request)
	}
}

func (h *handlerCfg) handlerListEntries(
	ctx context.Context, path string, request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {

	var depth float64 = 3
	if d, ok := request.Params.Arguments["depth"]; ok && d != nil {
		depth = d.(float64)
	}

	operationResult := listEntries(path, depth, "")
	if operationResult.Error != nil {
		log.Printf("ERROR: %v\n", operationResult.Error)
		return mcp.NewToolResultErrorFromErr("", operationResult.Error), operationResult.Error
	}

	if operationResult.Message != "" {
		log.Printf("WARNING: %v\n", operationResult.Message)
		return mcp.NewToolResultText(operationResult.Message), nil
	}

	return mcp.NewToolResultText(operationResult.Content), nil
}

func (h *handlerCfg) handlerReadFile(
	ctx context.Context, path string, request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	operationResult := readFile(path)
	if operationResult.Error != nil {
		log.Printf("ERROR: %v\n", operationResult.Error)
		return mcp.NewToolResultErrorFromErr("", operationResult.Error), operationResult.Error
	}

	if operationResult.Message != "" {
		log.Printf("WARNING: %v\n", operationResult.Message)
		return mcp.NewToolResultText(operationResult.Message), nil
	}

	log.Printf("File sucessfully read from: %v\n", path)

	return mcp.NewToolResultText(operationResult.Content), nil
}

func (h *handlerCfg) handlerWriteToFile(
	ctx context.Context, path string, request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	content := request.Params.Arguments["content"].(string)

	operationResult := writeToFile(content, path)
	if operationResult.Error != nil {
		log.Printf("ERROR: %v\n", operationResult.Error)
		return mcp.NewToolResultErrorFromErr("", operationResult.Error), operationResult.Error
	}

	if operationResult.Message != "" {
		log.Printf("WARNING: %v\n", operationResult.Message)
		return mcp.NewToolResultText(operationResult.Message), nil
	}

	log.Printf("File sucessfully written at: %v\n", path)

	return mcp.NewToolResultText(operationResult.Content), nil
}

func (h *handlerCfg) handlerGetFileInfo(
	ctx context.Context, path string, request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	operationResult := getFileInfo(path)
	if operationResult.Error != nil {
		log.Printf("ERROR: %v\n", operationResult.Error)
		return mcp.NewToolResultErrorFromErr("", operationResult.Error), operationResult.Error
	}

	if operationResult.Message != "" {
		log.Printf("WARNING: %v\n", operationResult.Message)
		return mcp.NewToolResultText(operationResult.Message), nil
	}

	log.Printf("Returning files info from file at: %v\n", path)

	return mcp.NewToolResultText(operationResult.Content), nil
}

func (h *handlerCfg) hadlerRenamePath(
	ctx context.Context, path string, request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	newPathFinalName := request.Params.Arguments["newPathFinalName"].(string)

	operationResult := renamePath(path, newPathFinalName)
	if operationResult.Error != nil {
		log.Printf("ERROR: %v\n", operationResult.Error)
		return mcp.NewToolResultErrorFromErr("", operationResult.Error), operationResult.Error
	}

	if operationResult.Message != "" {
		log.Printf("WARNING: %v\n", operationResult.Message)
		return mcp.NewToolResultText(operationResult.Message), nil
	}

	log.Printf("Returning files info from file at: %v\n", path)

	return mcp.NewToolResultText(operationResult.Content), nil
}

func (h *handlerCfg) hadlerCopyFileOrDir(
	ctx context.Context, path string, request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	destination := request.Params.Arguments["destination"].(string)

	if !h.dockerMode && !isSafePath(h.baseDir, destination) {
		log.Printf("PATH NOT ALLOWED: path is outside of allowed base directory")
		return mcp.NewToolResultText("access denied: path is outside of allowed base directory"), nil
	}

	operationResult := copyFileOrDir(path, destination)
	if operationResult.Error != nil {
		log.Printf("ERROR: %v\n", operationResult.Error)
		return mcp.NewToolResultErrorFromErr("", operationResult.Error), operationResult.Error
	}

	if operationResult.Message != "" {
		log.Printf("WARNING: %v\n", operationResult.Message)
		return mcp.NewToolResultText(operationResult.Message), nil
	}
	log.Printf("Returning files info from file at: %v\n", path)

	return mcp.NewToolResultText(operationResult.Content), nil
}

func handlersMiddleware(name string, fn server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("'%s' called with params: %v", name, request.Params.Arguments)
		return fn(ctx, request)
	}
}
