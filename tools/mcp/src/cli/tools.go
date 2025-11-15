package cli

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

type ToolCallback = func(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error)
