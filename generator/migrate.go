package generator

const migrationUpTemplate = `-- Users table
CREATE TABLE IF NOT EXISTS users (
	{{if eq .DBDriver "postgres"}}
	id SERIAL PRIMARY KEY,
	email VARCHAR(255) UNIQUE NOT NULL,
	password_hash VARCHAR(255) NOT NULL,
	name VARCHAR(255) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	{{else if eq .DBDriver "sqlite"}}
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	email TEXT UNIQUE NOT NULL,
	password_hash TEXT NOT NULL,
	name TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	{{end}}
);

{{if .WithSessions}}
-- Sessions table
CREATE TABLE IF NOT EXISTS sessions (
	{{if eq .DBDriver "postgres"}}
	id VARCHAR(255) PRIMARY KEY,
	user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
	expires_at TIMESTAMP NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	{{else if eq .DBDriver "sqlite"}}
	id TEXT PRIMARY KEY,
	user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
	expires_at DATETIME NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	{{end}}
);
{{end}}
`

const migrationDownTemplate = `{{if .WithSessions}}
DROP TABLE IF EXISTS sessions;
{{end}}
DROP TABLE IF EXISTS users;
`

const migrateMainTemplate = `package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	{{if eq .DBDriver "postgres"}}
	_ "github.com/jackc/pgx/v5/stdlib"
	{{else if eq .DBDriver "sqlite"}}
	_ "github.com/mattn/go-sqlite3"
	{{end}}
)

func main() {
	up := flag.Bool("up", false, "run migrations up")
	down := flag.Bool("down", false, "run migrations down")
	create := flag.String("create", "", "create a new migration (usage: -create migration_name)")
	flag.Parse()

	migrationsPath := "db/migrations"
	if path := os.Getenv("MIGRATIONS_PATH"); path != "" {
		migrationsPath = path
	}

	if *create != "" {
		if err := createMigration(migrationsPath, *create); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating migration: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if !*up && !*down {
		fmt.Fprintln(os.Stderr, "Usage: go run ./cmd/migrate -up | -down | -create <name>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		{{if eq .DBDriver "postgres"}}
		dbURL = "postgres://user:password@localhost:5432/{{.Name}}?sslmode=disable"
		{{else if eq .DBDriver "sqlite"}}
		dbURL = "./{{.Name}}.db"
		{{end}}
	}

	{{if eq .DBDriver "sqlite"}}
	// golang-migrate requires sqlite3:// prefix for file paths
	if !strings.HasPrefix(dbURL, "sqlite3://") && !strings.HasPrefix(dbURL, "file://") {
		dbURL = "sqlite3://" + dbURL
	}
	{{end}}

	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving migrations path: %v\n", err)
		os.Exit(1)
	}

	m, err := migrate.New(
		"file://"+absPath,
		{{if eq .DBDriver "postgres"}}dbURL{{else}}dbURL{{end}},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer m.Close()

	if *up {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			fmt.Fprintf(os.Stderr, "Error running migrations up: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migrations up complete")
	} else if *down {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			fmt.Fprintf(os.Stderr, "Error running migrations down: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migrations down complete")
	}
}

func createMigration(migrationsPath, name string) error {
	if err := os.MkdirAll(migrationsPath, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(migrationsPath)
	if err != nil {
		return err
	}

	nextVersion := 1
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		var v int
		if _, err := fmt.Sscanf(e.Name(), "%d_", &v); err == nil && v >= nextVersion {
			nextVersion = v + 1
		}
	}

	base := fmt.Sprintf("%06d_%s", nextVersion, strings.ReplaceAll(name, " ", "_"))
	upPath := filepath.Join(migrationsPath, base+".up.sql")
	downPath := filepath.Join(migrationsPath, base+".down.sql")

	if err := os.WriteFile(upPath, []byte("-- Add your migration here\n"), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(downPath, []byte("-- Add your rollback here\n"), 0644); err != nil {
		return err
	}

	fmt.Printf("Created migrations:\n  %s\n  %s\n", upPath, downPath)
	return nil
}
`

func (g *Generator) generateMigrate() error {
	// cmd/migrate/main.go
	if err := g.writeTemplate(g.projectPath("cmd/migrate/main.go"), migrateMainTemplate, g.config); err != nil {
		return err
	}

	// Initial migration files
	if err := g.writeTemplate(g.projectPath("db/migrations/000001_initial_schema.up.sql"), migrationUpTemplate, g.config); err != nil {
		return err
	}
	if err := g.writeTemplate(g.projectPath("db/migrations/000001_initial_schema.down.sql"), migrationDownTemplate, g.config); err != nil {
		return err
	}

	return nil
}
