package main

import (
	"context"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type handler struct {
	baseDir string
}

func (s *handler) withSafePath(
	handler func(ctx context.Context, path string, request mcp.CallToolRequest) (*mcp.CallToolResult, error),
) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path := request.Params.Arguments["path"].(string)
		if !isSafePath(s.baseDir, path) {
			log.Printf("PATH NOT ALLOWED: path is outside of allowed base directory")
			return mcp.NewToolResultText("access denied: path is outside of allowed base directory"), nil
		}
		return handler(ctx, path, request)
	}
}

func (s *handler) handlerListEntries(
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

func (s *handler) handlerReadFile(
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

func (s *handler) handlerWriteToFile(
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

func (s *handler) handlerGetFileInfo(
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

func (s *handler) hadlerRenamePath(
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

func (s *handler) hadlerCopyFileOrDir(
	ctx context.Context, path string, request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	destination := request.Params.Arguments["destination"].(string)

	if !isSafePath(s.baseDir, destination) {
		log.Printf("PATH NOT ALLOWED: path is outside of allowed base directory")
		return mcp.NewToolResultText("access denied: path is outside of allowed base directory"), nil
	}

	msg, err := copyFileOrDir(path, destination)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return mcp.NewToolResultErrorFromErr("", err), err
	}
	log.Printf("Returning files info from file at: %v\n", path)

	return mcp.NewToolResultText(msg), nil
}

func handlersMiddleware(name string, fn server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("'%s' called with params: %v", name, request.Params.Arguments)
		return fn(ctx, request)
	}
}
