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

func TestWarehouseRepository_Create_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewWarehouseRepository(mockDB, zap.NewNop())

	warehouse := &model.Warehouse{
		Name:     "Main Warehouse",
		Location: "Jakarta",
	}

	mockDB.
		ExpectQuery(`INSERT INTO warehouses`).
		WithArgs(warehouse.Name, warehouse.Location).
		WillReturnRows(pgxmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(5, time.Now(), time.Now()))

	err = repo.Create(warehouse)
	require.NoError(t, err)
	require.Equal(t, 5, warehouse.ID)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestWarehouseRepository_Create_Error(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewWarehouseRepository(mockDB, zap.NewNop())

	warehouse := &model.Warehouse{
		Name:     "Test Warehouse",
		Location: "Bandung",
	}

	mockDB.
		ExpectQuery(`INSERT INTO warehouses`).
		WithArgs(warehouse.Name, warehouse.Location).
		WillReturnError(errors.New("database error"))

	err = repo.Create(warehouse)
	require.Error(t, err)
	require.Equal(t, 0, warehouse.ID)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestWarehouseRepository_FindByID_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewWarehouseRepository(mockDB, zap.NewNop())

	mockDB.
		ExpectQuery(`SELECT (.+) FROM warehouses WHERE id`).
		WithArgs(1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "name", "location", "created_at", "updated_at"}).
			AddRow(1, "Main Warehouse", "Jakarta", time.Now(), time.Now()))

	warehouse, err := repo.FindByID(1)
	require.NoError(t, err)
	require.NotNil(t, warehouse)
	require.Equal(t, 1, warehouse.ID)
	require.Equal(t, "Main Warehouse", warehouse.Name)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestWarehouseRepository_Update_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewWarehouseRepository(mockDB, zap.NewNop())

	warehouse := &model.Warehouse{
		Name:     "Updated Warehouse",
		Location: "Surabaya",
	}

	mockDB.
		ExpectExec(`UPDATE warehouses`).
		WithArgs(warehouse.Name, warehouse.Location, 1).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.Update(1, warehouse)
	require.NoError(t, err)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestWarehouseRepository_Update_Error(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewWarehouseRepository(mockDB, zap.NewNop())

	warehouse := &model.Warehouse{
		Name:     "Test",
		Location: "Test",
	}

	mockDB.
		ExpectExec(`UPDATE warehouses`).
		WithArgs(warehouse.Name, warehouse.Location, 999).
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	err = repo.Update(999, warehouse)
	require.Error(t, err)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestWarehouseRepository_Delete_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewWarehouseRepository(mockDB, zap.NewNop())

	mockDB.
		ExpectExec(`DELETE FROM warehouses`).
		WithArgs(1).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(1)
	require.NoError(t, err)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestWarehouseRepository_Delete_Error(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewWarehouseRepository(mockDB, zap.NewNop())

	mockDB.
		ExpectExec(`DELETE FROM warehouses`).
		WithArgs(999).
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	err = repo.Delete(999)
	require.Error(t, err)

	require.NoError(t, mockDB.ExpectationsWereMet())
}
