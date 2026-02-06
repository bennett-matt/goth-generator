package generator

const envTemplate = `PORT={{.Port}}
DATABASE_URL={{if eq .DBDriver "postgres"}}postgres://user:password@localhost:5432/{{.Name}}?sslmode=disable{{else}}./{{.Name}}.db{{end}}
SECRET_KEY=change-this-secret-key-in-production
`

const gitignoreTemplate = `# If you prefer the allow list template instead of the deny list, see community template:
# https://github.com/github/gitignore/blob/main/community/Golang/Go.AllowList.gitignore
#
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with ` + "`go test -c`" + `
*.test

# Code coverage profiles and other test artifacts
*.out
coverage.*
*.coverprofile
profile.cov

# Dependency directories (remove the comment below to include it)
# vendor/

# Go workspace file
go.work
go.work.sum

# env file
.env

# Database
*.db
*.db-shm
*.db-wal

# Node modules
node_modules/

# Build outputs
dist/
build/

# Editor/IDE
# .idea/
# .vscode/
`

func (g *Generator) generateConfig() error {
	// .env file
	if err := g.writeTemplate(g.projectPath(".env.example"), envTemplate, g.config); err != nil {
		return err
	}

	// .gitignore
	if err := g.writeFile(g.projectPath(".gitignore"), gitignoreTemplate); err != nil {
		return err
	}

	return nil
}
