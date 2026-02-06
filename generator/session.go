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

func (g *Generator) generateSession() error {
	if !g.config.WithSessions {
		return nil
	}
	return g.writeTemplate(g.projectPath("internal/session/session.go"), sessionGoTemplate, g.config)
}
