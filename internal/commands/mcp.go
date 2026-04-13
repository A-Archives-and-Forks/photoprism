package commands

import (
	"context"
	"log/slog"
	"os"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/urfave/cli/v2"

	"github.com/photoprism/photoprism/internal/mcp"
)

// MCPCommands configures the Model Context Protocol (MCP) command group.
var MCPCommands = &cli.Command{
	Name:  "mcp",
	Usage: "Shows the Model Context Protocol (MCP) server subcommands",
	Subcommands: []*cli.Command{
		MCPServeCommand,
	},
}

// MCPServeCommand starts the MCP server over the stdio transport.
var MCPServeCommand = &cli.Command{
	Name:   "serve",
	Usage:  "Starts the internal MCP server via stdio for development and testing",
	Action: mcpServeAction,
}

// mcpServeAction starts the MCP server using the stdio transport.
func mcpServeAction(ctx *cli.Context) error {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	logger.Info("starting mcp server", "transport", "stdio", "tools", 2, "resources", 2)

	return runMCPServer(context.Background(), ctx, &sdkmcp.StdioTransport{})
}

// runMCPServer builds an MCP server from the CLI metadata and runs it over the
// given transport until the context is cancelled or the transport closes. The
// transport is a parameter so tests can substitute an in-memory transport for
// the stdio one the CLI Action uses in production.
func runMCPServer(ctx context.Context, appCtx *cli.Context, transport sdkmcp.Transport) error {
	implementation := &sdkmcp.Implementation{
		Name:    "photoprism-mcp",
		Version: mcpAppMetadata(appCtx, "Version", "development"),
	}
	edition := mcpAppMetadata(appCtx, "Edition", "unknown")

	return mcp.NewServer(implementation, edition).Run(ctx, transport)
}

// mcpAppMetadata returns the named string entry from the CLI app metadata,
// falling back to the supplied default if it is missing or not a string.
func mcpAppMetadata(ctx *cli.Context, key, fallback string) string {
	if value, ok := ctx.App.Metadata[key].(string); ok && value != "" {
		return value
	}

	return fallback
}
