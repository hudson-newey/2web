package cli

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

func ExecuteTool(command ...string) func(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		res, err := Execute(command...)

		isError := err != nil || res.Stderr != ""
		if isError {
			return mcp.NewToolResultError(res.Stderr), nil
		}

		return mcp.NewToolResultText(res.Stdout), nil
	}
}

// The "command" argument should NOT include "2web" command or path itself.
// It should only include the subcommands and flags to be executed.
func Execute(command ...string) (CommandResult, error) {
	return executeCliCommand(command)
}
