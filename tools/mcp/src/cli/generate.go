package cli

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

func NewGenerationTool(template string) func(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, err := request.RequireString("name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		command := []string{"generate", template, name}
		res, err := executeCliCommand(command)

		isError := err != nil || res.Stderr != ""
		if isError {
			return mcp.NewToolResultError(res.Stderr), nil
		}

		return mcp.NewToolResultText(res.Stdout), nil
	}
}
