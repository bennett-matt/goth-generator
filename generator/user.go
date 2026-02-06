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

const userTestTemplate = `package user

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupUserTestDB(t *testing.T) *sql.DB {
	t.Helper()
	sqliteDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	data, err := os.ReadFile("db/migrations/000001_initial_schema.up.sql")
	if err != nil {
		t.Skipf("migration file not found (run from project root): %v", err)
	}
	for _, stmt := range strings.Split(string(data), ";") {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}
		if _, err := sqliteDB.Exec(stmt); err != nil {
			t.Fatalf("exec schema: %v", err)
		}
	}
	return sqliteDB
}

func TestService_CreateGetByEmailGetByID(t *testing.T) {
	sqliteDB := setupUserTestDB(t)
	defer sqliteDB.Close()
	svc := NewService(sqliteDB)
	ctx := context.Background()

	u, err := svc.Create(ctx, "test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if u == nil || u.Email != "test@example.com" || u.Name != "Test User" {
		t.Fatalf("Create: got %+v", u)
	}

	byEmail, err := svc.GetByEmail(ctx, "test@example.com")
	if err != nil {
		t.Fatalf("GetByEmail: %v", err)
	}
	if byEmail.ID != u.ID {
		t.Errorf("GetByEmail: id %d != %d", byEmail.ID, u.ID)
	}

	byID, err := svc.GetByID(ctx, int(u.ID))
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if byID.Email != u.Email {
		t.Errorf("GetByID: email %q != %q", byID.Email, u.Email)
	}
}

func TestService_VerifyPassword(t *testing.T) {
	sqliteDB := setupUserTestDB(t)
	defer sqliteDB.Close()
	svc := NewService(sqliteDB)
	ctx := context.Background()

	u, err := svc.Create(ctx, "v@test.com", "secret456", "Verifier")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if err := svc.VerifyPassword(u, "secret456"); err != nil {
		t.Errorf("VerifyPassword(correct): %v", err)
	}
	if err := svc.VerifyPassword(u, "wrong"); err == nil {
		t.Error("VerifyPassword(wrong): expected error")
	}
}

func TestService_ListUsers(t *testing.T) {
	sqliteDB := setupUserTestDB(t)
	defer sqliteDB.Close()
	svc := NewService(sqliteDB)
	ctx := context.Background()

	_, _ = svc.Create(ctx, "a@test.com", "pass", "A")
	_, _ = svc.Create(ctx, "b@test.com", "pass", "B")

	list, err := svc.ListUsers(ctx)
	if err != nil {
		t.Fatalf("ListUsers: %v", err)
	}
	if len(list) < 2 {
		t.Errorf("ListUsers: got %d, want at least 2", len(list))
	}
}
`

func (g *Generator) generateUser() error {
	if !g.config.WithUsers && !g.config.WithAuth {
		return nil
	}
	if err := g.writeTemplate(g.projectPath("internal/user/user.go"), userGoTemplate, g.config); err != nil {
		return err
	}
	if g.config.DBDriver == "sqlite" {
		if err := g.writeTemplate(g.projectPath("internal/user/user_test.go"), userTestTemplate, g.config); err != nil {
			return err
		}
	}
	return nil
}
