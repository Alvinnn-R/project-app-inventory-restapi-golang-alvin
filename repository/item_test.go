package repository

import (
	"errors"
	"project-app-inventory/model"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestItemRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	logger, _ := zap.NewDevelopment()
	repo := NewItemRepository(mock, logger)

	item := &model.Item{
		SKU:          "TEST-001",
		Name:         "Test Item",
		CategoryID:   1,
		RackID:       1,
		Stock:        10,
		MinimumStock: 5,
		Price:        100000,
	}

	t.Run("Success", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(1, time.Now(), time.Now())

		mock.ExpectQuery("INSERT INTO items").
			WithArgs(item.SKU, item.Name, item.CategoryID, item.RackID,
				item.Stock, item.MinimumStock, item.Price).
			WillReturnRows(rows)

		err := repo.Create(item)
		assert.NoError(t, err)
		assert.Equal(t, 1, item.ID)
	})

	t.Run("Error - Database Error", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO items").
			WithArgs(item.SKU, item.Name, item.CategoryID, item.RackID,
				item.Stock, item.MinimumStock, item.Price).
			WillReturnError(errors.New("database error"))

		err := repo.Create(item)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
	})
}

func TestItemRepository_FindByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	logger, _ := zap.NewDevelopment()
	repo := NewItemRepository(mock, logger)

	t.Run("Success - Item Found", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{
			"id", "sku", "name", "category_id", "rack_id",
			"stock", "minimum_stock", "price", "created_at", "updated_at",
		}).AddRow(
			1, "TEST-001", "Test Item", 1, 1,
			10, 5, 100000.0, time.Now(), time.Now(),
		)

		mock.ExpectQuery("SELECT (.+) FROM items WHERE id").
			WithArgs(1).
			WillReturnRows(rows)

		item, err := repo.FindByID(1)
		assert.NoError(t, err)
		assert.NotNil(t, item)
		assert.Equal(t, "TEST-001", item.SKU)
		assert.Equal(t, "Test Item", item.Name)
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM items WHERE id").
			WithArgs(999).
			WillReturnError(pgx.ErrNoRows)

		item, err := repo.FindByID(999)
		assert.NoError(t, err)
		assert.Nil(t, item)
	})
}

func TestItemRepository_FindLowStock(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	logger, _ := zap.NewDevelopment()
	repo := NewItemRepository(mock, logger)

	t.Run("Success - Low Stock Items Found", func(t *testing.T) {
		// Mock count query
		countRows := pgxmock.NewRows([]string{"count"}).AddRow(2)
		mock.ExpectQuery("SELECT COUNT(.+) FROM items WHERE stock < minimum_stock").
			WillReturnRows(countRows)

		// Mock data query
		dataRows := pgxmock.NewRows([]string{
			"id", "sku", "name", "category_id", "rack_id",
			"stock", "minimum_stock", "price", "created_at", "updated_at",
		}).
			AddRow(1, "LOW-001", "Low Stock Item 1", 1, 1, 2, 5, 50000.0, time.Now(), time.Now()).
			AddRow(2, "LOW-002", "Low Stock Item 2", 1, 1, 3, 10, 75000.0, time.Now(), time.Now())

		mock.ExpectQuery("SELECT (.+) FROM items WHERE stock < minimum_stock").
			WithArgs(10, 0).
			WillReturnRows(dataRows)

		items, total, err := repo.FindLowStock(1, 10)
		assert.NoError(t, err)
		assert.Equal(t, 2, total)
		assert.Len(t, items, 2)
		assert.Equal(t, "LOW-001", items[0].SKU)
	})

	t.Run("Success - No Low Stock Items", func(t *testing.T) {
		countRows := pgxmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery("SELECT COUNT(.+) FROM items WHERE stock < minimum_stock").
			WillReturnRows(countRows)

		dataRows := pgxmock.NewRows([]string{
			"id", "sku", "name", "category_id", "rack_id",
			"stock", "minimum_stock", "price", "created_at", "updated_at",
		})
		mock.ExpectQuery("SELECT (.+) FROM items WHERE stock < minimum_stock").
			WithArgs(10, 0).
			WillReturnRows(dataRows)

		items, total, err := repo.FindLowStock(1, 10)
		assert.NoError(t, err)
		assert.Equal(t, 0, total)
		assert.Len(t, items, 0)
	})
}

func TestItemRepository_Update(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	logger, _ := zap.NewDevelopment()
	repo := NewItemRepository(mock, logger)

	item := &model.Item{
		SKU:          "TEST-001-UPDATED",
		Name:         "Updated Item",
		CategoryID:   1,
		RackID:       1,
		Stock:        20,
		MinimumStock: 5,
		Price:        150000,
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("UPDATE items").
			WithArgs(item.SKU, item.Name, item.CategoryID, item.RackID,
				item.Stock, item.MinimumStock, item.Price, 1).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := repo.Update(1, item)
		assert.NoError(t, err)
	})

	t.Run("Error - Item Not Found", func(t *testing.T) {
		mock.ExpectExec("UPDATE items").
			WithArgs(item.SKU, item.Name, item.CategoryID, item.RackID,
				item.Stock, item.MinimumStock, item.Price, 999).
			WillReturnResult(pgxmock.NewResult("UPDATE", 0))

		err := repo.Update(999, item)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestItemRepository_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	logger, _ := zap.NewDevelopment()
	repo := NewItemRepository(mock, logger)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM items").
			WithArgs(1).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		err := repo.Delete(1)
		assert.NoError(t, err)
	})

	t.Run("Error - Item Not Found", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM items").
			WithArgs(999).
			WillReturnResult(pgxmock.NewResult("DELETE", 0))

		err := repo.Delete(999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}
