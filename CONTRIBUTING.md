# Contributing to wtx

Thank you for your interest in contributing to wtx!

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Git
- Make (optional, but recommended)

### Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/yourusername/wtx.git
   cd wtx
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

4. Build:
   ```bash
   make build
   # or
   go build -o bin/wtx ./cmd/wtx
   ```

5. Run tests:
   ```bash
   make test
   # or
   go test ./...
   ```

## Project Structure

```
wtx/
├── cmd/wtx/              # CLI entry point
├── internal/
│   ├── git/              # Git operations
│   ├── tui/              # Bubble Tea UI
│   ├── editor/           # Editor adapters
│   ├── metadata/         # Data storage
│   ├── config/           # Configuration
│   ├── process/          # Process management
│   └── util/             # Utilities
├── plugins/
│   ├── vscode/           # VS Code extension
│   └── claude-code/      # Claude Code MCP server
└── scripts/              # Build and install scripts
```

## Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Keep functions small (<50 lines)
- Write meaningful commit messages
- Add tests for new features

## Testing

```bash
# Run all tests
make test

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/git
```

## Adding a New Editor

1. Create a new file in `internal/editor/`:
   ```go
   // internal/editor/myeditor.go
   package editor

   type MyEditor struct{}

   func (e *MyEditor) Name() string {
       return "myeditor"
   }

   func (e *MyEditor) IsAvailable() bool {
       return commandExists("myeditor")
   }

   func (e *MyEditor) Open(path string, reuse bool) error {
       cmd := exec.Command("myeditor", path)
       return cmd.Start()
   }
   ```

2. Add to editor list in `editor.go`:
   ```go
   editors := []Editor{
       &VSCode{},
       &Cursor{},
       &MyEditor{}, // Add here
       // ...
   }
   ```

3. Add to `GetEditor` switch:
   ```go
   case "myeditor":
       return &MyEditor{}
   ```

4. Test thoroughly
5. Submit PR with documentation

## Commit Messages

Follow conventional commits:

- `feat: add new command`
- `fix: resolve path issue`
- `docs: update README`
- `refactor: simplify git operations`
- `test: add metadata tests`

## Pull Request Process

1. Update README if adding features
2. Add tests for new functionality
3. Ensure all tests pass
4. Update CHANGELOG.md
5. Submit PR with clear description

## Feature Requests

Open an issue with:
- Clear use case
- Expected behavior
- Proposed solution (optional)

## Bug Reports

Include:
- wtx version (`wtx version`)
- OS and version
- Steps to reproduce
- Expected vs actual behavior
- Error messages

## Questions?

- Open a discussion on GitHub
- Check existing issues
- Read the documentation

## Code of Conduct

Be respectful and constructive. We're all here to make wtx better!

## License

By contributing, you agree your contributions will be licensed under MIT License.
