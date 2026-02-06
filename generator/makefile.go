package generator

const makefileTemplate = `.PHONY: dev build run test clean migrate sqlc templ css

# Development
dev:
	@go run github.com/a-h/templ/cmd/templ@latest generate 2>/dev/null || true
	@for f in web/templates/*_templ.go; do [ -f "$$f" ] && perl -i -0pe 's/(import templruntime "github\.com\/a-h\/templ\/runtime")\n\nimport "github\.com\/a-h\/templ"\n/\1\n/g' "$$f"; done
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

# Generate Templ code (fixes templ compiler duplicate-import bug)
templ:
	@echo "Generating Templ code..."
	@go run github.com/a-h/templ/cmd/templ@latest generate
	@for f in web/templates/*_templ.go; do [ -f "$$f" ] && perl -i -0pe 's/(import templruntime "github\.com\/a-h\/templ\/runtime")\n\nimport "github\.com\/a-h\/templ"\n/\1\n/g' "$$f"; done

# Install dependencies
deps:
	@go mod tidy
	@go mod download

# Setup development environment
setup:
	@cp .env.example .env
	@echo "Generating Templ code..."
	@go run github.com/a-h/templ/cmd/templ@latest generate
	@for f in web/templates/*_templ.go; do [ -f "$$f" ] && perl -i -0pe 's/(import templruntime "github\.com\/a-h\/templ\/runtime")\n\nimport "github\.com\/a-h\/templ"\n/\1\n/g' "$$f"; done
	@echo "Generating SQLC code..."
	@sqlc generate
	@$(MAKE) deps
	@npm install
	@echo "✅ Setup complete! Edit .env file with your configuration."
	@echo "   Run 'make dev' to start the server and CSS watcher."
`

func (g *Generator) generateMakefile() error {
	return g.writeTemplate(g.projectPath("Makefile"), makefileTemplate, g.config)
}
