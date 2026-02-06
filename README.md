# Goth Generator

Standalone Go project generator for web apps with users, sessions, auth, Tailwind, DaisyUI, Templ, SQLC, and HTMX.

## Installation

```bash
go install github.com/bennett-matt/goth-generator@latest
```

Or from source:

```bash
go build -o goth-generate .
```

## Usage

```bash
goth-generate -name myapp -module github.com/user/myapp -output ./projects
```

### Options

- `-name` (required): Project name
- `-module`: Go module path (defaults to app name)
- `-output`: Output directory (default: current directory)
- `-db`: Database driver - `postgres` or `sqlite` (default: `postgres`)
- `-port`: Server port (default: `8080`)
- `-auth`: Include authentication (default: `true`)
- `-users`: Include user management (default: `true`)
- `-sessions`: Include session management (default: `true`)

## Portability

This is a self-contained Go project. To move to another repo:

1. Copy the entire `goth-generator/` directory
2. Update the module path in `go.mod` to match your repo
3. Build and run
