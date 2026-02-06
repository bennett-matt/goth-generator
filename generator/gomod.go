package generator

import (
	"fmt"
	"os/exec"
)

func (g *Generator) generateGoMod() error {
	projectDir := g.projectPath("")

	// Run go mod init with module path (or app name if not provided)
	module := g.config.Module
	if module == "" {
		module = g.config.Name
	}

	initCmd := exec.Command("go", "mod", "init", module)
	initCmd.Dir = projectDir
	if out, err := initCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go mod init: %w: %s", err, out)
	}

	// Run go mod tidy to add dependencies from generated code
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = projectDir
	if out, err := tidyCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go mod tidy: %w: %s", err, out)
	}

	return nil
}
