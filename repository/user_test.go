package repository

import (
	"errors"
	"project-app-inventory/model"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Create_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewUserRepository(mockDB)

	user := &model.User{
		Name:         "John Doe",
		Email:        "john@example.com",
		PasswordHash: "hashedpassword",
		RoleID:       2,
		IsActive:     true,
	}

	mockDB.
		ExpectQuery(`INSERT INTO users`).
		WithArgs(user.Name, user.Email, user.PasswordHash, user.RoleID, user.IsActive).
		WillReturnRows(pgxmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(1, time.Now(), time.Now()))

	err = repo.Create(user)
	require.NoError(t, err)
	require.Equal(t, 1, user.ID)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestUserRepository_Create_Error(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewUserRepository(mockDB)

	user := &model.User{
		Name:         "Jane Doe",
		Email:        "jane@example.com",
		PasswordHash: "hashedpassword",
		RoleID:       3,
		IsActive:     true,
	}

	mockDB.
		ExpectQuery(`INSERT INTO users`).
		WithArgs(user.Name, user.Email, user.PasswordHash, user.RoleID, user.IsActive).
		WillReturnError(errors.New("duplicate email"))

	err = repo.Create(user)
	require.Error(t, err)
	require.Equal(t, 0, user.ID)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestUserRepository_FindByID_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewUserRepository(mockDB)

	mockDB.
		ExpectQuery(`SELECT (.+) FROM users u LEFT JOIN roles`).
		WithArgs(1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "name", "email", "password_hash", "role_id", "role_name", "is_active", "created_at", "updated_at"}).
			AddRow(1, "John Doe", "john@example.com", "hashedpassword", 2, "admin", true, time.Now(), time.Now()))

	user, err := repo.FindByID(1)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, 1, user.ID)
	require.Equal(t, "John Doe", user.Name)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestUserRepository_FindByEmail_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewUserRepository(mockDB)

	mockDB.
		ExpectQuery(`SELECT (.+) FROM users u LEFT JOIN roles`).
		WithArgs("john@example.com").
		WillReturnRows(pgxmock.NewRows([]string{"id", "name", "email", "password_hash", "role_id", "role_name", "is_active", "created_at", "updated_at"}).
			AddRow(1, "John Doe", "john@example.com", "hashedpassword", 2, "admin", true, time.Now(), time.Now()))

	user, err := repo.FindByEmail("john@example.com")
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, "john@example.com", user.Email)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestUserRepository_Update_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewUserRepository(mockDB)

	user := &model.User{
		Name:         "John Updated",
		Email:        "john.updated@example.com",
		PasswordHash: "newhashedpassword",
		RoleID:       2,
		IsActive:     true,
	}

	mockDB.
		ExpectExec(`UPDATE users`).
		WithArgs(user.Name, user.Email, user.PasswordHash, user.RoleID, user.IsActive, 1).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.Update(1, user)
	require.NoError(t, err)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestUserRepository_Update_Error(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewUserRepository(mockDB)

	user := &model.User{
		Name:         "Test",
		Email:        "test@example.com",
		PasswordHash: "hash",
		RoleID:       2,
		IsActive:     true,
	}

	mockDB.
		ExpectExec(`UPDATE users`).
		WithArgs(user.Name, user.Email, user.PasswordHash, user.RoleID, user.IsActive, 999).
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	err = repo.Update(999, user)
	require.Error(t, err)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestUserRepository_Delete_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewUserRepository(mockDB)

	mockDB.
		ExpectExec(`DELETE FROM users`).
		WithArgs(1).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(1)
	require.NoError(t, err)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestUserRepository_Delete_Error(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewUserRepository(mockDB)

	mockDB.
		ExpectExec(`DELETE FROM users`).
		WithArgs(999).
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	err = repo.Delete(999)
	require.Error(t, err)

	require.NoError(t, mockDB.ExpectationsWereMet())
}
