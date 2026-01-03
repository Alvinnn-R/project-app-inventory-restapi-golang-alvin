package repository

import (
	"errors"
	"project-app-inventory/model"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

// TestSessionRepository_Create_Success tests successful session creation
func TestSessionRepository_Create_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewSessionRepository(mockDB)
	expiredAt := time.Now().Add(24 * time.Hour)

	session := &model.Session{
		UserID:    1,
		Token:     "test-token-123",
		ExpiredAt: expiredAt,
	}

	mockDB.
		ExpectQuery(`INSERT INTO sessions`).
		WithArgs(session.UserID, session.Token, session.ExpiredAt).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))

	err = repo.Create(session)

	require.NoError(t, err)
	require.Equal(t, 1, session.ID)
	require.NoError(t, mockDB.ExpectationsWereMet())
}

// TestSessionRepository_Create_Error tests session creation error
func TestSessionRepository_Create_Error(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewSessionRepository(mockDB)
	session := &model.Session{UserID: 1, Token: "test-token"}

	mockDB.
		ExpectQuery(`INSERT INTO sessions`).
		WithArgs(session.UserID, session.Token, session.ExpiredAt).
		WillReturnError(errors.New("db error"))

	err = repo.Create(session)

	require.Error(t, err)
	require.NoError(t, mockDB.ExpectationsWereMet())
}

// TestSessionRepository_FindByToken_Success tests finding session by token
func TestSessionRepository_FindByToken_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewSessionRepository(mockDB)
	expiredAt := time.Now().Add(24 * time.Hour)
	createdAt := time.Now()

	mockDB.
		ExpectQuery(`SELECT (.+) FROM sessions WHERE token`).
		WithArgs("test-token").
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "user_id", "token", "expired_at", "revoked_at", "created_at",
		}).AddRow(1, 1, "test-token", expiredAt, nil, createdAt))

	session, err := repo.FindByToken("test-token")

	require.NoError(t, err)
	require.NotNil(t, session)
	require.Equal(t, "test-token", session.Token)
	require.Equal(t, 1, session.UserID)
	require.NoError(t, mockDB.ExpectationsWereMet())
}

// TestSessionRepository_FindByToken_NotFound tests session not found
func TestSessionRepository_FindByToken_NotFound(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewSessionRepository(mockDB)

	mockDB.
		ExpectQuery(`SELECT (.+) FROM sessions WHERE token`).
		WithArgs("invalid-token").
		WillReturnError(pgx.ErrNoRows)

	session, err := repo.FindByToken("invalid-token")

	require.NoError(t, err)
	require.Nil(t, session)
	require.NoError(t, mockDB.ExpectationsWereMet())
}

// TestSessionRepository_RevokeByToken_Success tests revoking session
func TestSessionRepository_RevokeByToken_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewSessionRepository(mockDB)

	mockDB.
		ExpectExec(`UPDATE sessions SET revoked_at`).
		WithArgs("test-token").
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.RevokeByToken("test-token")

	require.NoError(t, err)
	require.NoError(t, mockDB.ExpectationsWereMet())
}

// TestSessionRepository_RevokeByToken_Error tests revoke error
func TestSessionRepository_RevokeByToken_Error(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewSessionRepository(mockDB)

	mockDB.
		ExpectExec(`UPDATE sessions SET revoked_at`).
		WithArgs("test-token").
		WillReturnError(errors.New("db error"))

	err = repo.RevokeByToken("test-token")

	require.Error(t, err)
	require.NoError(t, mockDB.ExpectationsWereMet())
}

// TestSessionRepository_DeleteExpiredSessions_Success tests deleting expired sessions
func TestSessionRepository_DeleteExpiredSessions_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewSessionRepository(mockDB)

	mockDB.
		ExpectExec(`DELETE FROM sessions WHERE expired_at`).
		WillReturnResult(pgxmock.NewResult("DELETE", 5))

	err = repo.DeleteExpiredSessions()

	require.NoError(t, err)
	require.NoError(t, mockDB.ExpectationsWereMet())
}

// TestSessionRepository_DeleteExpiredSessions_Error tests deletion error
func TestSessionRepository_DeleteExpiredSessions_Error(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewSessionRepository(mockDB)

	mockDB.
		ExpectExec(`DELETE FROM sessions WHERE expired_at`).
		WillReturnError(errors.New("db error"))

	err = repo.DeleteExpiredSessions()

	require.Error(t, err)
	require.NoError(t, mockDB.ExpectationsWereMet())
}
