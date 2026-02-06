package generator

import (
	"fmt"
	"os/exec"
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

	// Run go mod tidy to add dependencies from generated code
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = projectDir
	if out, err := tidyCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go mod tidy: %w: %s", err, out)
	}

	return nil
}
