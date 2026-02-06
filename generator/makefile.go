package generator

const makefileTemplate = `.PHONY: dev build run test clean migrate sqlc templ css

# Development
dev:
	@templ generate 2>/dev/null || true
	@echo "Building CSS..."
	@npm run build:css:once
	@echo "Starting development server and CSS watcher..."
	@(make css &) && go run ./cmd/server

# Build CSS (Tailwind + DaisyUI) - watch mode, also started by make dev
css:
	@npm run build:css

# Build
build:
	@echo "Building application..."
	@go build -o bin/server ./cmd/server

# Run
run: build
	@./bin/server

# Test
test:
	@go test -v ./...

# Clean
clean:
	@rm -rf bin/
	@go clean

# Database migrations
migrate:
	@echo "Running migrations..."
	@go run ./cmd/migrate || echo "Migration command not implemented yet"

# Generate SQLC code
sqlc:
	@echo "Generating SQLC code..."
	@sqlc generate

# Generate Templ code
templ:
	@echo "Generating Templ code..."
	@templ generate

# Install dependencies
deps:
	@go mod tidy
	@go mod download

# Setup development environment
setup:
	@cp .env.example .env
	@echo "Generating Templ code..."
	@templ generate
	@$(MAKE) deps
	@npm install
	@echo "✅ Setup complete! Edit .env file with your configuration."
	@echo "   Run 'make dev' to start the server and CSS watcher."
`

func (g *Generator) generateMakefile() error {
	return g.writeTemplate(g.projectPath("Makefile"), makefileTemplate, g.config)
}
