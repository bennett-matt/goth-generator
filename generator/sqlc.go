package generator

const sqlcYamlTemplate = `version: "2"
sql:
  - engine: "{{if eq .DBDriver "postgres"}}postgresql{{else}}sqlite{{end}}"
    queries: "db/queries"
    schema: "db/schema"
    gen:
      go:
        package: "db"
        out: "internal/db"
        sql_package: "database/sql"
        emit_interface: true
        emit_json_tags: true
        overrides:
          - column: "users.id"
            go_type: "int64"
          - column: "users.password_hash"
            go_struct_tag: 'json:"-"'
          - column: "sessions.user_id"
            go_type: "int64"
`

const schemaSqlTemplate = `-- Users table
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

const queriesSqlTemplate = `-- name: GetUser :one
SELECT * FROM users
WHERE id = {{if eq .DBDriver "postgres"}}$1{{else}}?1{{end}} LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = {{if eq .DBDriver "postgres"}}$1{{else}}?1{{end}} LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC;

-- name: CreateUser :one
INSERT INTO users (email, password_hash, name)
VALUES ({{if eq .DBDriver "postgres"}}$1, $2, $3{{else}}?1, ?2, ?3{{end}})
RETURNING *;

{{if .WithSessions}}
-- name: GetSession :one
SELECT * FROM sessions
WHERE id = {{if eq .DBDriver "postgres"}}$1{{else}}?1{{end}} LIMIT 1;

-- name: CreateSession :one
INSERT INTO sessions (id, user_id, expires_at)
VALUES ({{if eq .DBDriver "postgres"}}$1, $2, $3{{else}}?1, ?2, ?3{{end}})
RETURNING *;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = {{if eq .DBDriver "postgres"}}$1{{else}}?1{{end}};

-- name: DeleteUserSessions :exec
DELETE FROM sessions
WHERE user_id = {{if eq .DBDriver "postgres"}}$1{{else}}?1{{end}};
{{end}}
`

func (g *Generator) generateSQLC() error {
	// SQLC config
	if err := g.writeTemplate(g.projectPath("sqlc.yaml"), sqlcYamlTemplate, g.config); err != nil {
		return err
	}

	// Schema SQL (for sqlc - mirrors migration 000001)
	if err := g.writeTemplate(g.projectPath("db/schema/schema.sql"), schemaSqlTemplate, g.config); err != nil {
		return err
	}

	// Queries SQL
	if err := g.writeTemplate(g.projectPath("db/queries/users.sql"), queriesSqlTemplate, g.config); err != nil {
		return err
	}

	return nil
}
