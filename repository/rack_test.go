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

func TestRackRepository_Create_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewRackRepository(mockDB, zap.NewNop())

	desc := "Rack A1 for electronics"
	rack := &model.Rack{
		WarehouseID: 1,
		Code:        "A1",
		Description: &desc,
	}

	mockDB.
		ExpectQuery(`INSERT INTO racks`).
		WithArgs(rack.WarehouseID, rack.Code, rack.Description).
		WillReturnRows(pgxmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(10, time.Now(), time.Now()))

	err = repo.Create(rack)
	require.NoError(t, err)
	require.Equal(t, 10, rack.ID)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestRackRepository_Create_Error(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewRackRepository(mockDB, zap.NewNop())

	desc := "Rack A1"
	rack := &model.Rack{
		WarehouseID: 1,
		Code:        "A1",
		Description: &desc,
	}

	mockDB.
		ExpectQuery(`INSERT INTO racks`).
		WithArgs(rack.WarehouseID, rack.Code, rack.Description).
		WillReturnError(errors.New("database error"))

	err = repo.Create(rack)
	require.Error(t, err)
	require.Equal(t, 0, rack.ID)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestRackRepository_FindByID_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewRackRepository(mockDB, zap.NewNop())

	desc := "Rack A1"
	mockDB.
		ExpectQuery(`SELECT (.+) FROM racks WHERE id`).
		WithArgs(1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "warehouse_id", "code", "description", "created_at", "updated_at"}).
			AddRow(1, 1, "A1", &desc, time.Now(), time.Now()))

	rack, err := repo.FindByID(1)
	require.NoError(t, err)
	require.NotNil(t, rack)
	require.Equal(t, 1, rack.ID)
	require.Equal(t, "A1", rack.Code)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestRackRepository_Update_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewRackRepository(mockDB, zap.NewNop())

	desc := "Updated description"
	rack := &model.Rack{
		WarehouseID: 1,
		Code:        "A1-Updated",
		Description: &desc,
	}

	mockDB.
		ExpectExec(`UPDATE racks`).
		WithArgs(rack.WarehouseID, rack.Code, rack.Description, 1).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.Update(1, rack)
	require.NoError(t, err)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestRackRepository_Update_Error(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewRackRepository(mockDB, zap.NewNop())

	desc := "Test"
	rack := &model.Rack{
		WarehouseID: 1,
		Code:        "A1",
		Description: &desc,
	}

	mockDB.
		ExpectExec(`UPDATE racks`).
		WithArgs(rack.WarehouseID, rack.Code, rack.Description, 999).
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	err = repo.Update(999, rack)
	require.Error(t, err)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestRackRepository_Delete_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewRackRepository(mockDB, zap.NewNop())

	mockDB.
		ExpectExec(`DELETE FROM racks`).
		WithArgs(1).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(1)
	require.NoError(t, err)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestRackRepository_Delete_Error(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewRackRepository(mockDB, zap.NewNop())

	mockDB.
		ExpectExec(`DELETE FROM racks`).
		WithArgs(999).
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	err = repo.Delete(999)
	require.Error(t, err)

	require.NoError(t, mockDB.ExpectationsWereMet())
}
