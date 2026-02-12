# wtx - Git Worktree Manager

> "Cmd+Tab for your Git branches"

Fast, keyboard-driven workspace switcher for Git worktrees.

## Features

- 🚀 **Instant switching** - Switch between worktrees in <2 seconds
- 🎨 **Beautiful TUI** - Interactive fuzzy finder with status indicators
- 🔧 **Multi-editor** - VS Code, Cursor, Neovim, Claude Code support
- 🌐 **Dev server tracking** - Manage ports and processes per worktree
- 🔒 **Safe by default** - Never lose uncommitted work
- ⌨️ **Keyboard-first** - Everything accessible without mouse

## Quick Start

```bash
# Install
go install github.com/yourusername/wtx/cmd/wtx@latest

# Switch workspace (interactive)
wtx

# List all worktrees
wtx list

# Create new worktree
wtx add feature-branch

# Open in editor
wtx open feature-branch
```

## Status

🚧 **In Development** - v0.1.0 coming soon

## License

MIT
