package repository

import (
	"errors"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

// TestPermissionRepository_Allowed_True tests permission check returning true
func TestPermissionRepository_Allowed_True(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewPermissionRepository(mockDB)

	mockDB.
		ExpectQuery(`WITH perm AS`).
		WithArgs(1, "manage_items").
		WillReturnRows(pgxmock.NewRows([]string{"allowed"}).AddRow(true))

	allowed, err := repo.Allowed(1, "manage_items")

	require.NoError(t, err)
	require.True(t, allowed)
	require.NoError(t, mockDB.ExpectationsWereMet())
}

// TestPermissionRepository_Allowed_False tests permission check returning false
func TestPermissionRepository_Allowed_False(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewPermissionRepository(mockDB)

	mockDB.
		ExpectQuery(`WITH perm AS`).
		WithArgs(1, "delete_items").
		WillReturnRows(pgxmock.NewRows([]string{"allowed"}).AddRow(false))

	allowed, err := repo.Allowed(1, "delete_items")

	require.NoError(t, err)
	require.False(t, allowed)
	require.NoError(t, mockDB.ExpectationsWereMet())
}

// TestPermissionRepository_Allowed_Error tests permission check error
func TestPermissionRepository_Allowed_Error(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewPermissionRepository(mockDB)

	mockDB.
		ExpectQuery(`WITH perm AS`).
		WithArgs(1, "invalid_permission").
		WillReturnError(errors.New("db error"))

	allowed, err := repo.Allowed(1, "invalid_permission")

	require.Error(t, err)
	require.False(t, allowed)
	require.NoError(t, mockDB.ExpectationsWereMet())
}
