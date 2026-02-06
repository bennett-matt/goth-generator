package generator

const databaseGoTemplate = `package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
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

func Migrate(db *sql.DB) error {
	ctx := context.Background()
	
	{{if eq .DBDriver "postgres"}}
	query := "` + "`" + `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS sessions (
		id VARCHAR(255) PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		expires_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	` + "`" + `"
	{{else if eq .DBDriver "sqlite"}}
	query := "` + "`" + `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		name TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		expires_at DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	` + "`" + `"
	{{end}}

	if _, err := db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
`

func (g *Generator) generateDatabase() error {
	return g.writeTemplate(g.projectPath("internal/database/database.go"), databaseGoTemplate, g.config)
}
