package main

import (
	"context"
	"log"
	"reflect"
	"runtime"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func handlersMiddleware(fn server.ToolHandlerFunc) server.ToolHandlerFunc {
	pc := reflect.ValueOf(fn).Pointer()
	fullFuncName := runtime.FuncForPC(pc).Name()

	funcName := fullFuncName
	if idx := strings.LastIndex(funcName, "."); idx >= 0 {
		funcName = funcName[idx+1:]
	}

	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("'%s' called with params: %v", funcName, request.Params.Arguments)
		return fn(ctx, request)
	}
}

func handlerListEntries(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := request.Params.Arguments["path"].(string)

	entries, err := listEntries(path)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return mcp.NewToolResultErrorFromErr("", err), err
	}

	return mcp.NewToolResultText(entries), nil
}

func handlerReadFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := request.Params.Arguments["path"].(string)

	entries, err := readFile(path)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return mcp.NewToolResultErrorFromErr("", err), err
	}
	log.Printf("File sucessfully read from: %v\n", path)

	return mcp.NewToolResultText(entries), nil
}

func handlerWriteToFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := request.Params.Arguments["path"].(string)
	content := request.Params.Arguments["content"].(string)

	msg, err := writeToFile(content, path)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return mcp.NewToolResultErrorFromErr("", err), err
	}
	log.Printf("File sucessfully written at: %v\n", path)

	return mcp.NewToolResultText(msg), nil
}
