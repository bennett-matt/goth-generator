package generator

const databaseGoTemplate = `package database

import (
	"database/sql"
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

func New() (*sql.DB, error) {
	{{if eq .DBDriver "postgres"}}
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://user:password@localhost:5432/{{.Name}}?sslmode=disable"
	}
	db, err := sql.Open("pgx", dsn)
	{{else if eq .DBDriver "sqlite"}}
	dbPath := os.Getenv("DATABASE_URL")
	if dbPath == "" {
		dbPath = "./{{.Name}}.db"
	}
	db, err := sql.Open("sqlite3", dbPath)
	{{end}}
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func MigrateUp(db *sql.DB) error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		{{if eq .DBDriver "postgres"}}
		dbURL = "postgres://user:password@localhost:5432/{{.Name}}?sslmode=disable"
		{{else if eq .DBDriver "sqlite"}}
		dbURL = "./{{.Name}}.db"
		{{end}}
	}

	{{if eq .DBDriver "sqlite"}}
	if !strings.HasPrefix(dbURL, "sqlite3://") && !strings.HasPrefix(dbURL, "file://") {
		dbURL = "sqlite3://" + dbURL
	}
	{{end}}

	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = "db/migrations"
	}
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to resolve migrations path: %w", err)
	}

	m, err := migrate.New(
		"file://"+absPath,
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
`

func (g *Generator) generateDatabase() error {
	return g.writeTemplate(g.projectPath("internal/database/database.go"), databaseGoTemplate, g.config)
}
