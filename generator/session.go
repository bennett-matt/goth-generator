package generator

const sessionGoTemplate = `package session

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"{{.Module}}/internal/db"
)

type Store struct {
	queries *db.Queries
}

func NewStore(dbtx db.DBTX) *Store {
	return &Store{queries: db.New(dbtx)}
}

func (s *Store) Create(ctx context.Context, userID int64, expiresAt time.Time) (*db.Session, error) {
	id := uuid.New().String()
	sess, err := s.queries.CreateSession(ctx, db.CreateSessionParams{
		ID:        id,
		UserID:    userID,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

func (s *Store) Get(ctx context.Context, id string) (*db.Session, error) {
	sess, err := s.queries.GetSession(ctx, id)
	if err != nil {
		return nil, err
	}
	if time.Now().After(sess.ExpiresAt) {
		s.Delete(ctx, id)
		return nil, sql.ErrNoRows
	}
	return &sess, nil
}

func (s *Store) Delete(ctx context.Context, id string) error {
	return s.queries.DeleteSession(ctx, id)
}

func (s *Store) DeleteByUserID(ctx context.Context, userID int64) error {
	return s.queries.DeleteUserSessions(ctx, userID)
}
`

const sessionTestTemplate = `package session

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
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
	// Create a user for session tests
	if _, err := sqliteDB.Exec("INSERT INTO users (email, password_hash, name) VALUES ('u@test.com', 'hash', 'User')"); err != nil {
		t.Fatalf("insert user: %v", err)
	}
	return sqliteDB
}

func TestStore_CreateGetDelete(t *testing.T) {
	sqliteDB := setupTestDB(t)
	defer sqliteDB.Close()
	store := NewStore(sqliteDB)
	ctx := context.Background()

	sess, err := store.Create(ctx, 1, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if sess == nil || sess.ID == "" {
		t.Fatal("Create: got nil or empty ID")
	}

	got, err := store.Get(ctx, sess.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.ID != sess.ID || got.UserID != 1 {
		t.Errorf("Get: got ID=%q UserID=%d, want ID=%q UserID=1", got.ID, got.UserID, sess.ID)
	}

	if err := store.Delete(ctx, sess.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err = store.Get(ctx, sess.ID)
	if err == nil {
		t.Error("Get after Delete: expected error")
	}
}

func TestStore_DeleteByUserID(t *testing.T) {
	sqliteDB := setupTestDB(t)
	defer sqliteDB.Close()
	store := NewStore(sqliteDB)
	ctx := context.Background()

	sess, _ := store.Create(ctx, 1, time.Now().Add(time.Hour))
	store.DeleteByUserID(ctx, 1)
	_, err := store.Get(ctx, sess.ID)
	if err == nil {
		t.Error("Get after DeleteByUserID: expected error")
	}
}
`

func (g *Generator) generateSession() error {
	if !g.config.WithSessions {
		return nil
	}
	if err := g.writeTemplate(g.projectPath("internal/session/session.go"), sessionGoTemplate, g.config); err != nil {
		return err
	}
	if g.config.DBDriver == "sqlite" {
		if err := g.writeTemplate(g.projectPath("internal/session/session_test.go"), sessionTestTemplate, g.config); err != nil {
			return err
		}
	}
	return nil
}
