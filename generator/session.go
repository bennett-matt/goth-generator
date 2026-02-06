package generator

const sessionGoTemplate = `package session

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"{{.Module}}/internal/models"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Create(ctx context.Context, userID int, expiresAt time.Time) (*models.Session, error) {
	id := uuid.New().String()
	session := &models.Session{
		ID:        id,
		UserID:    userID,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	query := "INSERT INTO sessions (id, user_id, expires_at) VALUES ($1, $2, $3)"
	_, err := s.db.ExecContext(ctx, query, session.ID, session.UserID, session.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Store) Get(ctx context.Context, id string) (*models.Session, error) {
	var session models.Session
	query := "SELECT id, user_id, expires_at, created_at FROM sessions WHERE id = $1"
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&session.ID,
		&session.UserID,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		s.Delete(ctx, id)
		return nil, sql.ErrNoRows
	}

	return &session, nil
}

func (s *Store) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM sessions WHERE id = $1"
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

func (s *Store) DeleteByUserID(ctx context.Context, userID int) error {
	query := "DELETE FROM sessions WHERE user_id = $1"
	_, err := s.db.ExecContext(ctx, query, userID)
	return err
}
`

func (g *Generator) generateSession() error {
	if !g.config.WithSessions {
		return nil
	}
	return g.writeTemplate(g.projectPath("internal/session/session.go"), sessionGoTemplate, g.config)
}
