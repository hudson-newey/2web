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

	newPageTool := mcp.NewTool("new-page",
		mcp.WithDescription("Create a new page in the current project"),
		mcp.WithString(
			"name",
			mcp.Required(),
			mcp.Description("The name of the page"),
		),
	)

	newComponentTool := mcp.NewTool("new-component",
		mcp.WithDescription("Create a component in the current project"),
		mcp.WithString(
			"name",
			mcp.Required(),
			mcp.Description("The name of the component"),
		),
	)

	newServiceTool := mcp.NewTool("new-service",
		mcp.WithDescription("Create a service in the current project"),
		mcp.WithString(
			"name",
			mcp.Required(),
			mcp.Description("The name of the service"),
		),
	)

	newModelTool := mcp.NewTool("new-model",
		mcp.WithDescription("Create a model in the current project"),
		mcp.WithString(
			"name",
			mcp.Required(),
			mcp.Description("The name of the model"),
		),
	)

	newEnumTool := mcp.NewTool("new-enum",
		mcp.WithDescription("Create an enum in the current project"),
		mcp.WithString(
			"name",
			mcp.Required(),
			mcp.Description("The name of the enum"),
		),
	)

	newInterfaceTool := mcp.NewTool("new-interface",
		mcp.WithDescription("Create an interface in the current project"),
		mcp.WithString(
			"name",
			mcp.Required(),
			mcp.Description("The name of the interface"),
		),
	)

	mcpServer.AddTool(lintTool, cli.NewExecuteTool("lint"))
	mcpServer.AddTool(formatTool, cli.NewExecuteTool("format"))
	mcpServer.AddTool(testTool, cli.NewExecuteTool("test"))

	mcpServer.AddTool(newPageTool, cli.NewGenerationTool("page"))
	mcpServer.AddTool(newComponentTool, cli.NewGenerationTool("component"))
	mcpServer.AddTool(newServiceTool, cli.NewGenerationTool("service"))

	mcpServer.AddTool(newModelTool, cli.NewGenerationTool("model"))
	mcpServer.AddTool(newEnumTool, cli.NewGenerationTool("enum"))
	mcpServer.AddTool(newInterfaceTool, cli.NewGenerationTool("interface"))

	// Start the stdio server
	if err := server.ServeStdio(mcpServer); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
