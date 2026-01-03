package repository

import (
	"errors"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestReportRepository_GetTotalItems_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewReportRepository(mockDB, zap.NewNop())

	mockDB.
		ExpectQuery(`SELECT COUNT\(\*\) FROM items`).
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(100))

	total, err := repo.GetTotalItems()
	require.NoError(t, err)
	require.Equal(t, 100, total)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestReportRepository_GetTotalItems_Error(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewReportRepository(mockDB, zap.NewNop())

	mockDB.
		ExpectQuery(`SELECT COUNT\(\*\) FROM items`).
		WillReturnError(errors.New("database error"))

	total, err := repo.GetTotalItems()
	require.Error(t, err)
	require.Equal(t, 0, total)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestReportRepository_GetLowStockItems_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewReportRepository(mockDB, zap.NewNop())

	mockDB.
		ExpectQuery(`SELECT COUNT\(\*\) FROM items WHERE stock < minimum_stock`).
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(15))

	total, err := repo.GetLowStockItems()
	require.NoError(t, err)
	require.Equal(t, 15, total)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestReportRepository_GetTotalSales_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewReportRepository(mockDB, zap.NewNop())

	mockDB.
		ExpectQuery(`SELECT COUNT\(\*\) FROM sales WHERE deleted_at IS NULL`).
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(50))

	total, err := repo.GetTotalSales()
	require.NoError(t, err)
	require.Equal(t, 50, total)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestReportRepository_GetTotalRevenue_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewReportRepository(mockDB, zap.NewNop())

	mockDB.
		ExpectQuery(`SELECT COALESCE\(SUM\(total_amount\), 0\) FROM sales WHERE deleted_at IS NULL`).
		WillReturnRows(pgxmock.NewRows([]string{"sum"}).AddRow(5000000.0))

	total, err := repo.GetTotalRevenue()
	require.NoError(t, err)
	require.Equal(t, 5000000.0, total)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestReportRepository_GetActiveUsers_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewReportRepository(mockDB, zap.NewNop())

	mockDB.
		ExpectQuery(`SELECT COUNT\(\*\) FROM users WHERE is_active = true`).
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(25))

	total, err := repo.GetActiveUsers()
	require.NoError(t, err)
	require.Equal(t, 25, total)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestReportRepository_GetTotalCategories_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewReportRepository(mockDB, zap.NewNop())

	mockDB.
		ExpectQuery(`SELECT COUNT\(\*\) FROM categories`).
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(10))

	total, err := repo.GetTotalCategories()
	require.NoError(t, err)
	require.Equal(t, 10, total)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestReportRepository_GetTotalWarehouses_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewReportRepository(mockDB, zap.NewNop())

	mockDB.
		ExpectQuery(`SELECT COUNT\(\*\) FROM warehouses`).
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(5))

	total, err := repo.GetTotalWarehouses()
	require.NoError(t, err)
	require.Equal(t, 5, total)

	require.NoError(t, mockDB.ExpectationsWereMet())
}
