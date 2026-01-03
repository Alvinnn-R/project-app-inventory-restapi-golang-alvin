package repository

import (
	"errors"
	"project-app-inventory/model"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCategoryRepository_Create_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewCategoryRepository(mockDB, zap.NewNop())

	desc := "Electronics and gadgets"
	category := &model.Category{
		Name:        "Electronics",
		Description: &desc,
	}

	mockDB.
		ExpectQuery(`INSERT INTO categories`).
		WithArgs(category.Name, category.Description).
		WillReturnRows(pgxmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(1, time.Now(), time.Now()))

	err = repo.Create(category)
	require.NoError(t, err)
	require.Equal(t, 1, category.ID)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestCategoryRepository_Create_Error(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewCategoryRepository(mockDB, zap.NewNop())

	category := &model.Category{
		Name:        "Books",
		Description: nil,
	}

	mockDB.
		ExpectQuery(`INSERT INTO categories`).
		WithArgs(category.Name, category.Description).
		WillReturnError(errors.New("database error"))

	err = repo.Create(category)
	require.Error(t, err)
	require.Equal(t, 0, category.ID)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestCategoryRepository_FindByID_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewCategoryRepository(mockDB, zap.NewNop())

	desc := "Test description"
	mockDB.
		ExpectQuery(`SELECT (.+) FROM categories WHERE id`).
		WithArgs(1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at"}).
			AddRow(1, "Electronics", &desc, time.Now(), time.Now()))

	category, err := repo.FindByID(1)
	require.NoError(t, err)
	require.NotNil(t, category)
	require.Equal(t, 1, category.ID)
	require.Equal(t, "Electronics", category.Name)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestCategoryRepository_FindByName_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewCategoryRepository(mockDB, zap.NewNop())

	mockDB.
		ExpectQuery(`SELECT (.+) FROM categories WHERE name`).
		WithArgs("Electronics").
		WillReturnRows(pgxmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at"}).
			AddRow(1, "Electronics", nil, time.Now(), time.Now()))

	category, err := repo.FindByName("Electronics")
	require.NoError(t, err)
	require.NotNil(t, category)
	require.Equal(t, "Electronics", category.Name)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestCategoryRepository_Update_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewCategoryRepository(mockDB, zap.NewNop())

	desc := "Updated description"
	category := &model.Category{
		Name:        "Electronics Updated",
		Description: &desc,
	}

	mockDB.
		ExpectExec(`UPDATE categories`).
		WithArgs(category.Name, category.Description, 1).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.Update(1, category)
	require.NoError(t, err)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestCategoryRepository_Update_Error(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewCategoryRepository(mockDB, zap.NewNop())

	category := &model.Category{
		Name:        "Test",
		Description: nil,
	}

	mockDB.
		ExpectExec(`UPDATE categories`).
		WithArgs(category.Name, category.Description, 999).
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	err = repo.Update(999, category)
	require.Error(t, err)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestCategoryRepository_Delete_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewCategoryRepository(mockDB, zap.NewNop())

	mockDB.
		ExpectExec(`DELETE FROM categories`).
		WithArgs(1).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(1)
	require.NoError(t, err)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestCategoryRepository_Delete_Error(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewCategoryRepository(mockDB, zap.NewNop())

	mockDB.
		ExpectExec(`DELETE FROM categories`).
		WithArgs(999).
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	err = repo.Delete(999)
	require.Error(t, err)

	require.NoError(t, mockDB.ExpectationsWereMet())
}
