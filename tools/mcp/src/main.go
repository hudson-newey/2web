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

	lintTool := mcp.NewTool("lint",
		mcp.WithDescription("Lint the current project"),
	)

	formatTool := mcp.NewTool("format",
		mcp.WithDescription("Format the current project"),
	)

	testTool := mcp.NewTool("test",
		mcp.WithDescription("Test the current project"),
	)

	// newPageTool := mcp.NewTool("new-page",
	// 	mcp.WithDescription("Create a new page in the current project"),
	// 	mcp.WithString(
	// 		"name",
	// 		mcp.Required(),
	// 		mcp.Description("The name of the new page"),
	// 	),
	// )

	mcpServer.AddTool(lintTool, cli.ExecuteTool("lint"))
	mcpServer.AddTool(formatTool, cli.ExecuteTool("format"))
	mcpServer.AddTool(testTool, cli.ExecuteTool("test"))

	// Start the stdio server
	if err := server.ServeStdio(mcpServer); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
