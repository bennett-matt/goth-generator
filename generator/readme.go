package generator

const readmeTemplate = `# {{.Name}}

A Go web application with authentication, user management, and session handling.

## Features

{{if .WithAuth}}- ✅ Authentication{{end}}
{{if .WithUsers}}- ✅ User Management{{end}}
{{if .WithSessions}}- ✅ Session Management{{end}}
- ✅ Health Check Endpoint
- ✅ Tailwind CSS & DaisyUI
- ✅ HTMX for dynamic interactions
- ✅ Templ Templates
- ✅ SQLC for Database Queries
- ✅ {{if eq .DBDriver "postgres"}}PostgreSQL{{else}}SQLite{{end}} Database

## Tech Stack

- **Language**: Go 1.23+
- **Web Framework**: httprouter
- **Templates**: Templ
- **Interactivity**: HTMX
- **Styling**: Tailwind CSS + DaisyUI
- **Database**: {{if eq .DBDriver "postgres"}}PostgreSQL{{else}}SQLite{{end}}
- **Query Builder**: SQLC
- **Authentication**: Session-based

## Getting Started

### Prerequisites

- Go 1.23 or later
- {{if eq .DBDriver "postgres"}}PostgreSQL 16+{{else}}SQLite{{end}}
- Make (optional)

### Installation

1. Clone the repository:
` + "```" + `bash
git clone <repository-url>
cd {{.Name}}
` + "```" + `

2. Install dependencies:
` + "```" + `bash
make setup
# or manually:
cp .env.example .env
` + "```" + `

3. Configure environment variables:
` + "```" + `bash
# Edit .env file with your settings
PORT={{.Port}}
DATABASE_URL={{if eq .DBDriver "postgres"}}postgres://user:password@localhost:5432/{{.Name}}?sslmode=disable{{else}}./{{.Name}}.db{{end}}
SECRET_KEY=your-secret-key-here
` + "```" + `

4. {{if eq .DBDriver "postgres"}}Start PostgreSQL database:
` + "```" + `bash
docker-compose up -d db
` + "```" + `
{{end}}

5. Run migrations:
` + "```" + `bash
# Migrations are run automatically on first start
# Or manually via: make migrate
` + "```" + `

6. Generate code and build CSS:
` + "```" + `bash
make sqlc  # Generate SQLC code
make templ # Generate Templ templates
npm install && npm run build:css  # Build Tailwind + DaisyUI (or 'make css' for watch mode)
` + "```" + `

7. Start the development server:
` + "```" + `bash
make dev
# or
go run ./cmd/server
` + "```" + `

The server will start on ` + "`" + `http://localhost:{{.Port}}` + "`" + `

## Project Structure

` + "```" + `
{{.Name}}/
├── cmd/
│   └── server/          # Main application entry point
├── internal/
│   ├── auth/            # Authentication logic
│   ├── database/        # Database connection and migrations
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # HTTP middleware
│   ├── models/          # Data models
│   ├── session/         # Session management
│   └── user/            # User service
├── web/
│   ├── static/          # Static assets (CSS, JS)
│   └── templates/       # Templ templates
├── db/
│   ├── migrations/      # SQL migrations
│   └── queries/         # SQLC queries
└── go.mod               # Go dependencies
` + "```" + `

## Development

### Running Tests

` + "```" + `bash
make test
` + "```" + `

### Building

` + "```" + `bash
make build
` + "```" + `

### Code Generation

- **SQLC**: ` + "`" + `make sqlc` + "`" + ` - Generates Go code from SQL queries
- **Templ**: ` + "`" + `make templ` + "`" + ` - Generates Go code from Templ templates

## Docker

### Build and Run with Docker Compose

` + "```" + `bash
docker-compose up --build
` + "```" + `

### Build Docker Image

` + "```" + `bash
docker build -t {{.Name}} .
` + "```" + `

## License

MIT
`

func (g *Generator) generateReadme() error {
	return g.writeTemplate(g.projectPath("README.md"), readmeTemplate, g.config)
}
