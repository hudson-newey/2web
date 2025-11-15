package main

import (
	"fmt"

	"github.com/hudson-newey/2web-mcp/src/cli"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	mcpServer := server.NewMCPServer(
		"2Web MCP",
		"0.1.0",
		server.WithToolCapabilities(false),
	)

	lintTool := mcp.NewTool("lint-project",
		mcp.WithDescription("Lint the current project"),
	)

	formatTool := mcp.NewTool("format-project",
		mcp.WithDescription("Format the current project"),
	)

	testTool := mcp.NewTool("test-project",
		mcp.WithDescription("Test the current project"),
	)

	mcpServer.AddTool(lintTool, cli.ExecuteCallback("lint"))
	mcpServer.AddTool(formatTool, cli.ExecuteCallback("format"))
	mcpServer.AddTool(testTool, cli.ExecuteCallback("test"))

	mcpServer.AddTool(cli.NewGenerationTool("page"))
	mcpServer.AddTool(cli.NewGenerationTool("component"))
	mcpServer.AddTool(cli.NewGenerationTool("service"))

	mcpServer.AddTool(cli.NewGenerationTool("model"))
	mcpServer.AddTool(cli.NewGenerationTool("enum"))
	mcpServer.AddTool(cli.NewGenerationTool("interface"))

	if err := server.ServeStdio(mcpServer); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
