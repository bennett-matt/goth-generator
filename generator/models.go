package generator

const modelsGoTemplate = `package models

import "time"

type User struct {
	ID           int       ` + "`json:\"id\"`" + `
	Email        string    ` + "`json:\"email\"`" + `
	PasswordHash string    ` + "`json:\"-\"`" + `
	Name         string    ` + "`json:\"name\"`" + `
	CreatedAt    time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt    time.Time ` + "`json:\"updated_at\"`" + `
}

{{if .WithSessions}}
type Session struct {
	ID        string    ` + "`json:\"id\"`" + `
	UserID    int       ` + "`json:\"user_id\"`" + `
	ExpiresAt time.Time ` + "`json:\"expires_at\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
}
{{end}}
`

func (g *Generator) generateModels() error {
	return g.writeTemplate(g.projectPath("internal/models/models.go"), modelsGoTemplate, g.config)
}
