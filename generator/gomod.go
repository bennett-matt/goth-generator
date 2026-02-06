package generator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (g *Generator) generateGoMod() error {
	projectDir := g.projectPath("")

	// Run go mod init with module path (or app name if not provided)
	module := g.config.Module
	if module == "" {
		module = g.config.Name
	}
	// Strip URL scheme - Go expects github.com/user/repo, not https://github.com/user/repo
	for _, prefix := range []string{"https://", "http://"} {
		if strings.HasPrefix(module, prefix) {
			module = strings.TrimPrefix(module, prefix)
			break
		}
	}

	initCmd := exec.Command("go", "mod", "init", module)
	initCmd.Dir = projectDir
	if out, err := initCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go mod init: %w: %s", err, out)
	}

	// Generate templ components (creates .templ.go from .templ files)
	templCmd := exec.Command("go", "run", "github.com/a-h/templ/cmd/templ@latest", "generate")
	templCmd.Dir = projectDir
	if out, err := templCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("templ generate: %w: %s (run 'make templ' manually if needed)", err, out)
	}

	// Fix duplicate import in templ-generated files (templ compiler bug)
	if err := fixTemplImports(projectDir); err != nil {
		return fmt.Errorf("fix templ imports: %w", err)
	}

	// Run go mod tidy to add dependencies from generated code
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = projectDir
	if out, err := tidyCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go mod tidy: %w: %s", err, out)
	}

	return nil
}

// fixTemplImports removes duplicate "github.com/a-h/templ" imports from templ-generated files.
func fixTemplImports(projectDir string) error {
	templDir := filepath.Join(projectDir, "web", "templates")
	entries, err := os.ReadDir(templDir)
	if err != nil {
		return err
	}
	duplicateImport := "\n\nimport \"github.com/a-h/templ\"\n"
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), "_templ.go") {
			continue
		}
		path := filepath.Join(templDir, e.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		content := string(data)
		if strings.Count(content, `import "github.com/a-h/templ"`) > 1 {
			content = strings.Replace(content, duplicateImport, "\n", 1)
			if err := os.WriteFile(path, []byte(content), 0644); err != nil {
				return err
			}
		}
	}
	return nil
}
