package cli

import (
	"context"
	"os/exec"

	"github.com/mark3labs/mcp-go/mcp"
)

type CommandResult struct {
	Stdout string
	Stderr string
}

// The "command" argument should NOT include "2web" command or path itself.
// It should only include the subcommands and flags to be executed.
func NewExecuteTool(command ...string) func(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		res, err := executeCliCommand(command)

		isError := err != nil || res.Stderr != ""
		if isError {
			return mcp.NewToolResultError(res.Stderr), nil
		}

		return mcp.NewToolResultText(res.Stdout), nil
	}
}

// Executes a command from the 2web CLI and returns the output or an error.
func executeCliCommand(command []string) (CommandResult, error) {
	cliPath, err := twoWebCliPath()
	if err != nil {
		return CommandResult{}, err
	}

	// We always use the --no-cache for MCP executions to ensure fresh builds.
	command = append(command, "--no-cache")

	cmd := exec.Command(cliPath, command...)
	cmd.Stdin = nil

	stdoutStream, err := cmd.StdoutPipe()
	if err != nil {
		return CommandResult{}, err
	}

	stderrStream, err := cmd.StderrPipe()
	if err != nil {
		return CommandResult{}, err
	}

	if err := cmd.Start(); err != nil {
		return CommandResult{}, err
	}

	stdout := readStream(stdoutStream)
	stderr := readStream(stderrStream)
	cmd.Wait()

	return CommandResult{
		Stdout: stdout,
		Stderr: stderr,
	}, nil
}
