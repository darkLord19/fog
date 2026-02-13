# wtx MCP Server (Experimental)

Model Context Protocol (MCP) server that exposes `wtx` worktree operations to Claude Code.

This is a thin wrapper around the `wtx` CLI:
- receives MCP tool calls
- executes corresponding `wtx` commands
- returns results to the MCP client

## Requirements

- `wtx` installed and available in PATH
- Claude Code with MCP support

## Install (From Source)

```bash
cd plugins/claude-code
npm install
npm run build
npm link
```

## Configure Claude Code

Add to `~/.claude/config.json`:

```json
{
  "mcpServers": {
    "wtx": {
      "command": "wtx-mcp-server"
    }
  }
}
```

## License

See repository root `LICENSE` (AGPL-3.0-or-later).

