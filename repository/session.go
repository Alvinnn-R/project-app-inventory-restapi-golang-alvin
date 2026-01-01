package repository

import (
	"context"
	"database/sql"
	"project-app-inventory/database"
	"project-app-inventory/model"
)

type SessionRepository interface {
	Create(session *model.Session) error
	FindByToken(token string) (*model.Session, error)
	RevokeByToken(token string) error
	DeleteExpiredSessions() error
}

type sessionRepositoryImpl struct {
	db database.PgxIface
}

func NewSessionRepository(db database.PgxIface) SessionRepository {
	return &sessionRepositoryImpl{db: db}
}

func (r *sessionRepositoryImpl) Create(session *model.Session) error {
	query := `
		INSERT INTO sessions (user_id, token, expired_at, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id
	`
	return r.db.QueryRow(context.Background(), query, session.UserID, session.Token, session.ExpiredAt).Scan(&session.ID)
}

func (r *sessionRepositoryImpl) FindByToken(token string) (*model.Session, error) {
	query := `
		SELECT id, user_id, token, expired_at, revoked_at, created_at
		FROM sessions
		WHERE token = $1 AND revoked_at IS NULL AND expired_at > NOW()
	`
	var session model.Session
	err := r.db.QueryRow(context.Background(), query, token).Scan(
		&session.ID, &session.UserID, &session.Token, &session.ExpiredAt, &session.RevokedAt, &session.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // session tidak ditemukan atau sudah expired/revoked
	}

	return &session, err
}

func (r *sessionRepositoryImpl) RevokeByToken(token string) error {
	query := `
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE token = $1 AND revoked_at IS NULL
	`
	_, err := r.db.Exec(context.Background(), query, token)
	return err
}

func (r *sessionRepositoryImpl) DeleteExpiredSessions() error {
	query := `
		DELETE FROM sessions
		WHERE expired_at < NOW() OR revoked_at IS NOT NULL
	`
	_, err := r.db.Exec(context.Background(), query)
	return err
}
