package generator

import (
	"os"
)

func (g *Generator) generateStructure() error {
	dirs := []string{
		"cmd/server",
		"internal/database",
		"internal/db",
		"internal/handlers",
		"internal/middleware",
		"internal/auth",
		"internal/session",
		"internal/user",
		"web/templates",
		"web/static/css",
		"web/static/js",
		"cmd/migrate",
		"db/migrations",
		"db/schema",
		"db/queries",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(g.projectPath(dir), 0755); err != nil {
			return err
		}
	}

	return nil
}
