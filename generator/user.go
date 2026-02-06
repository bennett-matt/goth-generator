package generator

const userGoTemplate = `package user

import (
	"context"

	"golang.org/x/crypto/bcrypt"
	"{{.Module}}/internal/db"
)

type Service struct {
	queries *db.Queries
}

func NewService(dbtx db.DBTX) *Service {
	return &Service{queries: db.New(dbtx)}
}

func (s *Service) Create(ctx context.Context, email, password, name string) (*db.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		Email:        email,
		PasswordHash: string(hashedPassword),
		Name:         name,
	})
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*db.User, error) {
	u, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Service) GetByID(ctx context.Context, id int) (*db.User, error) {
	u, err := s.queries.GetUser(ctx, int64(id))
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Service) ListUsers(ctx context.Context) ([]db.User, error) {
	return s.queries.ListUsers(ctx)
}

func (s *Service) VerifyPassword(user *db.User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
}
`

func (g *Generator) generateUser() error {
	if !g.config.WithUsers || !g.config.WithAuth {
		return nil
	}
	return g.writeTemplate(g.projectPath("internal/user/user.go"), userGoTemplate, g.config)
}
