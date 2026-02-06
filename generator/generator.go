package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Generator struct {
	config *Config
}

func New(config *Config) *Generator {
	return &Generator{
		config: config,
	}
}

func (g *Generator) Generate() error {
	if err := os.MkdirAll(g.config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	steps := []struct {
		name string
		fn   func() error
	}{
		{"project structure", g.generateStructure},
		{"main.go", g.generateMain},
		{"database", g.generateDatabase},
		{"models", g.generateModels},
		{"session service", g.generateSession},
		{"user service", g.generateUser},
		{"handlers", g.generateHandlers},
		{"middleware", g.generateMiddleware},
		{"templates", g.generateTemplates},
		{"static files", g.generateStatic},
		{"sqlc config", g.generateSQLC},
		{"config files", g.generateConfig},
		{"docker files", g.generateDocker},
		{"makefile", g.generateMakefile},
		{"readme", g.generateReadme},
		{"go mod", g.generateGoMod},
	}

	for _, step := range steps {
		if err := step.fn(); err != nil {
			return fmt.Errorf("failed to generate %s: %w", step.name, err)
		}
	}

	return nil
}

func (g *Generator) projectPath(paths ...string) string {
	allPaths := append([]string{g.config.OutputDir}, paths...)
	return filepath.Join(allPaths...)
}

func (g *Generator) writeFile(path string, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func (g *Generator) writeTemplate(path string, tmpl string, data interface{}) error {
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return err
	}

	var buf strings.Builder
	if err := t.Execute(&buf, data); err != nil {
		return err
	}

	return g.writeFile(path, buf.String())
}
