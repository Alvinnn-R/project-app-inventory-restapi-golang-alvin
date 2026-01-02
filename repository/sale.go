package repository

import (
	"context"
	"errors"
	"project-app-inventory/database"
	"project-app-inventory/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type SaleRepository interface {
	Create(sale *model.Sale, items []model.SaleItem) error
	FindByID(id int) (*model.Sale, error)
	FindSaleItems(saleID int) ([]model.SaleItem, error)
	FindAll(page, limit int) ([]model.Sale, int, error)
	Update(id int, sale *model.Sale, items []model.SaleItem) error
	Delete(id int) error
}

type saleRepository struct {
	db     database.PgxIface
	Logger *zap.Logger
}

func NewSaleRepository(db database.PgxIface, log *zap.Logger) SaleRepository {
	return &saleRepository{db: db, Logger: log}
}

func (r *saleRepository) Create(sale *model.Sale, items []model.SaleItem) error {
	// Type assert to get pool for transaction
	pool, ok := r.db.(*pgxpool.Pool)
	if !ok {
		return errors.New("database connection does not support transactions")
	}

	// Begin transaction
	tx, err := pool.Begin(context.Background())
	if err != nil {
		r.Logger.Error("error beginning transaction", zap.Error(err))
		return err
	}
	defer tx.Rollback(context.Background())

	// Insert sale
	query := `
		INSERT INTO sales (user_id, total_amount, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(context.Background(), query,
		sale.UserID, sale.TotalAmount,
	).Scan(&sale.ID, &sale.CreatedAt, &sale.UpdatedAt)

	if err != nil {
		r.Logger.Error("error creating sale", zap.Error(err))
		return err
	}

	// Insert sale items
	itemQuery := `
		INSERT INTO sale_items (sale_id, item_id, quantity, price_at_sale, subtotal)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	for i := range items {
		items[i].SaleID = sale.ID
		err = tx.QueryRow(context.Background(), itemQuery,
			items[i].SaleID, items[i].ItemID, items[i].Quantity,
			items[i].PriceAtSale, items[i].Subtotal,
		).Scan(&items[i].ID)

		if err != nil {
			r.Logger.Error("error creating sale item", zap.Error(err))
			return err
		}

		// Update item stock
		updateStockQuery := `
			UPDATE items
			SET stock = stock - $1, updated_at = NOW()
			WHERE id = $2 AND stock >= $3
		`
		result, err := tx.Exec(context.Background(), updateStockQuery,
			items[i].Quantity, items[i].ItemID, items[i].Quantity,
		)
		if err != nil {
			r.Logger.Error("error updating item stock", zap.Error(err))
			return err
		}

		if result.RowsAffected() == 0 {
			return errors.New("insufficient stock for item")
		}
	}

	// Commit transaction
	err = tx.Commit(context.Background())
	if err != nil {
		r.Logger.Error("error committing transaction", zap.Error(err))
		return err
	}

	return nil
}

func (r *saleRepository) FindByID(id int) (*model.Sale, error) {
	query := `
		SELECT s.id, s.user_id, s.total_amount, s.created_at, s.updated_at, s.deleted_at
		FROM sales s
		WHERE s.id = $1 AND s.deleted_at IS NULL
	`
	var sale model.Sale
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&sale.ID, &sale.UserID, &sale.TotalAmount, &sale.CreatedAt, &sale.UpdatedAt, &sale.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		r.Logger.Error("error finding sale by id", zap.Error(err))
		return nil, err
	}
	return &sale, nil
}

func (r *saleRepository) FindSaleItems(saleID int) ([]model.SaleItem, error) {
	query := `
		SELECT id, sale_id, item_id, quantity, price_at_sale, subtotal
		FROM sale_items
		WHERE sale_id = $1
		ORDER BY id ASC
	`
	rows, err := r.db.Query(context.Background(), query, saleID)
	if err != nil {
		r.Logger.Error("error querying sale items", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var items []model.SaleItem
	for rows.Next() {
		var item model.SaleItem
		err := rows.Scan(
			&item.ID, &item.SaleID, &item.ItemID,
			&item.Quantity, &item.PriceAtSale, &item.Subtotal,
		)
		if err != nil {
			r.Logger.Error("error scanning sale item", zap.Error(err))
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *saleRepository) FindAll(page, limit int) ([]model.Sale, int, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM sales WHERE deleted_at IS NULL`
	err := r.db.QueryRow(context.Background(), countQuery).Scan(&total)
	if err != nil {
		r.Logger.Error("error counting sales", zap.Error(err))
		return nil, 0, err
	}

	// Get data with pagination
	query := `
		SELECT id, user_id, total_amount, created_at, updated_at, deleted_at
		FROM sales
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		r.Logger.Error("error querying sales", zap.Error(err))
		return nil, 0, err
	}
	defer rows.Close()

	var sales []model.Sale
	for rows.Next() {
		var sale model.Sale
		err := rows.Scan(
			&sale.ID, &sale.UserID, &sale.TotalAmount, &sale.CreatedAt, &sale.UpdatedAt, &sale.DeletedAt,
		)
		if err != nil {
			r.Logger.Error("error scanning sale", zap.Error(err))
			return nil, 0, err
		}
		sales = append(sales, sale)
	}

	return sales, total, nil
}

func (r *saleRepository) Update(id int, sale *model.Sale, items []model.SaleItem) error {
	// Type assert to get pool for transaction
	pool, ok := r.db.(*pgxpool.Pool)
	if !ok {
		return errors.New("database connection does not support transactions")
	}

	// Begin transaction
	tx, err := pool.Begin(context.Background())
	if err != nil {
		r.Logger.Error("error beginning transaction", zap.Error(err))
		return err
	}
	defer tx.Rollback(context.Background())

	// Get old sale items to return stock
	oldItems, err := r.FindSaleItems(id)
	if err != nil {
		r.Logger.Error("error finding old sale items", zap.Error(err))
		return err
	}

	// Return stock from old items
	for _, oldItem := range oldItems {
		returnStockQuery := `
			UPDATE items
			SET stock = stock + $1, updated_at = NOW()
			WHERE id = $2
		`
		_, err := tx.Exec(context.Background(), returnStockQuery,
			oldItem.Quantity, oldItem.ItemID,
		)
		if err != nil {
			r.Logger.Error("error returning stock", zap.Error(err))
			return err
		}
	}

	// Delete old sale items
	deleteItemsQuery := `DELETE FROM sale_items WHERE sale_id = $1`
	_, err = tx.Exec(context.Background(), deleteItemsQuery, id)
	if err != nil {
		r.Logger.Error("error deleting old sale items", zap.Error(err))
		return err
	}

	// Update sale total amount
	updateSaleQuery := `
		UPDATE sales
		SET total_amount = $1, updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
	`
	result, err := tx.Exec(context.Background(), updateSaleQuery, sale.TotalAmount, id)
	if err != nil {
		r.Logger.Error("error updating sale", zap.Error(err))
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("sale not found")
	}

	// Insert new sale items
	itemQuery := `
		INSERT INTO sale_items (sale_id, item_id, quantity, price_at_sale, subtotal)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	for i := range items {
		items[i].SaleID = id
		err = tx.QueryRow(context.Background(), itemQuery,
			items[i].SaleID, items[i].ItemID, items[i].Quantity,
			items[i].PriceAtSale, items[i].Subtotal,
		).Scan(&items[i].ID)

		if err != nil {
			r.Logger.Error("error creating new sale item", zap.Error(err))
			return err
		}

		// Reduce stock for new items
		updateStockQuery := `
			UPDATE items
			SET stock = stock - $1, updated_at = NOW()
			WHERE id = $2 AND stock >= $3
		`
		result, err := tx.Exec(context.Background(), updateStockQuery,
			items[i].Quantity, items[i].ItemID, items[i].Quantity,
		)
		if err != nil {
			r.Logger.Error("error updating item stock", zap.Error(err))
			return err
		}

		if result.RowsAffected() == 0 {
			return errors.New("insufficient stock for item")
		}
	}

	// Commit transaction
	err = tx.Commit(context.Background())
	if err != nil {
		r.Logger.Error("error committing transaction", zap.Error(err))
		return err
	}

	return nil
}

func (r *saleRepository) Delete(id int) error {
	// Type assert to get pool for transaction
	pool, ok := r.db.(*pgxpool.Pool)
	if !ok {
		return errors.New("database connection does not support transactions")
	}

	// Begin transaction
	tx, err := pool.Begin(context.Background())
	if err != nil {
		r.Logger.Error("error beginning transaction", zap.Error(err))
		return err
	}
	defer tx.Rollback(context.Background())

	// Get sale items to return stock
	saleItems, err := r.FindSaleItems(id)
	if err != nil {
		r.Logger.Error("error finding sale items", zap.Error(err))
		return err
	}

	// Return stock for each item
	for _, item := range saleItems {
		returnStockQuery := `
			UPDATE items
			SET stock = stock + $1, updated_at = NOW()
			WHERE id = $2
		`
		_, err := tx.Exec(context.Background(), returnStockQuery,
			item.Quantity, item.ItemID,
		)
		if err != nil {
			r.Logger.Error("error returning stock", zap.Error(err))
			return err
		}
	}

	// Delete sale items first
	deleteItemsQuery := `DELETE FROM sale_items WHERE sale_id = $1`
	_, err = tx.Exec(context.Background(), deleteItemsQuery, id)
	if err != nil {
		r.Logger.Error("error deleting sale items", zap.Error(err))
		return err
	}

	// Soft delete sale
	deleteSaleQuery := `UPDATE sales SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	result, err := tx.Exec(context.Background(), deleteSaleQuery, id)
	if err != nil {
		r.Logger.Error("error deleting sale", zap.Error(err))
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("sale not found")
	}

	// Commit transaction
	err = tx.Commit(context.Background())
	if err != nil {
		r.Logger.Error("error committing transaction", zap.Error(err))
		return err
	}

	return nil
}
