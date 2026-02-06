package generator

const dockerfileTemplate = `FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/server .

{{if eq .DBDriver "postgres"}}
# PostgreSQL client for migrations (optional)
RUN apk --no-cache add postgresql-client
{{end}}

EXPOSE {{.Port}}

CMD ["./server"]
`

const composeYamlTemplate = `services:
  app:
    build: .
    ports:
      - "{{.Port}}:{{.Port}}"
    environment:
      - PORT={{.Port}}
      {{if eq .DBDriver "postgres"}}
      - DATABASE_URL=postgres://user:password@db:5432/{{.Name}}?sslmode=disable
      {{else}}
      - DATABASE_URL=./{{.Name}}.db
      {{end}}
    depends_on:
      {{if eq .DBDriver "postgres"}}
      - db
      {{end}}
    volumes:
      - .:/app
    command: ./server

{{if eq .DBDriver "postgres"}}
  db:
    image: postgres:16-alpine
    environment:
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=password
      - POSTGRES_DB={{.Name}}
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
{{end}}
`

const dockerignoreTemplate = `# Git
.git
.gitignore

# Documentation
README.md
*.md

# Environment
.env
.env.local

# IDE
.vscode
.idea
*.swp
*.swo

# Build artifacts
*.exe
*.test
*.out

# Dependencies
vendor/

# Database
*.db
*.db-shm
*.db-wal

# Node
node_modules/
`

func (g *Generator) generateDocker() error {
	// Dockerfile
	if err := g.writeTemplate(g.projectPath("Dockerfile"), dockerfileTemplate, g.config); err != nil {
		return err
	}

	// docker-compose.yaml
	if err := g.writeTemplate(g.projectPath("compose.yaml"), composeYamlTemplate, g.config); err != nil {
		return err
	}

	// .dockerignore
	if err := g.writeFile(g.projectPath(".dockerignore"), dockerignoreTemplate); err != nil {
		return err
	}

	return nil
}
