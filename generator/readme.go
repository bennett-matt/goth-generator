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

- **Language**: Go {{.GoVersion}}+
- **Web Framework**: httprouter
- **Templates**: Templ
- **Interactivity**: HTMX
- **Styling**: Tailwind CSS + DaisyUI
- **Database**: {{if eq .DBDriver "postgres"}}PostgreSQL{{else}}SQLite{{end}}
- **Query Builder**: SQLC
- **Authentication**: Session-based

## Getting Started

### Prerequisites

- Go {{.GoVersion}} or later
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

5. Run migrations (or let the server run them on startup):
` + "```" + `bash
make migrate-up    # Run migrations up
make migrate-down  # Rollback one migration
make migrate-create name=add_feature  # Create a new migration
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
│   ├── server/          # Main application entry point
│   └── migrate/         # Migration CLI (up, down, create)
├── internal/
│   ├── auth/            # Authentication logic
│   ├── database/        # Database connection and migrations
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # HTTP middleware
│   ├── session/         # Session management
│   └── user/            # User service
├── web/
│   ├── static/          # Static assets (CSS, JS)
│   └── templates/       # Templ templates
├── db/
│   ├── migrations/      # SQL migrations (up/down)
│   ├── schema/          # Schema for SQLC
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

### Database Migrations

- **Up**: ` + "`" + `make migrate-up` + "`" + ` - Apply all pending migrations
- **Down**: ` + "`" + `make migrate-down` + "`" + ` - Rollback the last migration
- **Create**: ` + "`" + `make migrate-create name=add_users_table` + "`" + ` - Create a new migration

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
