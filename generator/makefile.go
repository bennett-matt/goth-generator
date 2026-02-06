package generator

const makefileTemplate = `.PHONY: dev build run test clean migrate sqlc templ css

# Development
dev:
	@echo "Starting development server..."
	@go run ./cmd/server

# Build CSS (Tailwind + DaisyUI) - run in separate terminal for watch mode
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
setup: deps
	@cp .env.example .env
	@npm install
	@echo "✅ Setup complete! Edit .env file with your configuration."
	@echo "   Run 'make css' in another terminal to build/watch Tailwind CSS."
`

func (g *Generator) generateMakefile() error {
	return g.writeTemplate(g.projectPath("Makefile"), makefileTemplate, g.config)
}
