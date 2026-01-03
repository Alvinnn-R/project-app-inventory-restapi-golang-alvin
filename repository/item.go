package repository

import (
	"context"
	"errors"
	"project-app-inventory/database"
	"project-app-inventory/model"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type ItemRepository interface {
	Create(item *model.Item) error
	FindByID(id int) (*model.Item, error)
	FindBySKU(sku string) (*model.Item, error)
	FindAll(page, limit int) ([]model.Item, int, error)
	FindLowStock(page, limit int) ([]model.Item, int, error)
	Update(id int, data *model.Item) error
	Delete(id int) error
}

type itemRepository struct {
	db     database.PgxIface
	Logger *zap.Logger
}

func NewItemRepository(db database.PgxIface, log *zap.Logger) ItemRepository {
	return &itemRepository{db: db, Logger: log}
}

func (r *itemRepository) Create(item *model.Item) error {
	query := `
		INSERT INTO items (sku, name, category_id, rack_id, stock, minimum_stock, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(context.Background(), query,
		item.SKU, item.Name, item.CategoryID, item.RackID,
		item.Stock, item.MinimumStock, item.Price,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)

	if err != nil {
		r.Logger.Error("error creating item", zap.Error(err))
	}
	return err
}

func (r *itemRepository) FindByID(id int) (*model.Item, error) {
	query := `
		SELECT id, sku, name, category_id, rack_id, stock, minimum_stock, price, 
		       created_at, updated_at
		FROM items 
		WHERE id = $1
	`
	var item model.Item
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&item.ID, &item.SKU, &item.Name, &item.CategoryID, &item.RackID,
		&item.Stock, &item.MinimumStock, &item.Price,
		&item.CreatedAt, &item.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		r.Logger.Error("error finding item by id", zap.Error(err))
		return nil, err
	}
	return &item, nil
}

func (r *itemRepository) FindBySKU(sku string) (*model.Item, error) {
	query := `
		SELECT id, sku, name, category_id, rack_id, stock, minimum_stock, price, 
		       created_at, updated_at
		FROM items 
		WHERE sku = $1
	`
	var item model.Item
	err := r.db.QueryRow(context.Background(), query, sku).Scan(
		&item.ID, &item.SKU, &item.Name, &item.CategoryID, &item.RackID,
		&item.Stock, &item.MinimumStock, &item.Price,
		&item.CreatedAt, &item.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		r.Logger.Error("error finding item by sku", zap.Error(err))
		return nil, err
	}
	return &item, nil
}

func (r *itemRepository) FindAll(page, limit int) ([]model.Item, int, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM items`
	err := r.db.QueryRow(context.Background(), countQuery).Scan(&total)
	if err != nil {
		r.Logger.Error("error counting items", zap.Error(err))
		return nil, 0, err
	}

	// Get data with pagination
	query := `
		SELECT id, sku, name, category_id, rack_id, stock, minimum_stock, price, 
		       created_at, updated_at
		FROM items
		ORDER BY name ASC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		r.Logger.Error("error querying items", zap.Error(err))
		return nil, 0, err
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var item model.Item
		err := rows.Scan(
			&item.ID, &item.SKU, &item.Name, &item.CategoryID, &item.RackID,
			&item.Stock, &item.MinimumStock, &item.Price,
			&item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			r.Logger.Error("error scanning item", zap.Error(err))
			return nil, 0, err
		}
		items = append(items, item)
	}

	return items, total, nil
}

func (r *itemRepository) FindLowStock(page, limit int) ([]model.Item, int, error) {
	offset := (page - 1) * limit

	// Get total count of low stock items
	var total int
	countQuery := `SELECT COUNT(*) FROM items WHERE stock < minimum_stock`
	err := r.db.QueryRow(context.Background(), countQuery).Scan(&total)
	if err != nil {
		r.Logger.Error("error counting low stock items", zap.Error(err))
		return nil, 0, err
	}

	// Get data with pagination
	query := `
		SELECT id, sku, name, category_id, rack_id, stock, minimum_stock, price, created_at, updated_at
		FROM items
		WHERE stock < minimum_stock
		ORDER BY stock ASC, name ASC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		r.Logger.Error("error querying low stock items", zap.Error(err))
		return nil, 0, err
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var item model.Item
		err := rows.Scan(
			&item.ID, &item.SKU, &item.Name, &item.CategoryID, &item.RackID,
			&item.Stock, &item.MinimumStock, &item.Price,
			&item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			r.Logger.Error("error scanning low stock item", zap.Error(err))
			return nil, 0, err
		}
		items = append(items, item)
	}

	return items, total, nil
}

func (r *itemRepository) Update(id int, data *model.Item) error {
	query := `
		UPDATE items
		SET sku = $1, name = $2, category_id = $3, rack_id = $4,
		    stock = $5, minimum_stock = $6, price = $7, updated_at = NOW()
		WHERE id = $8
	`
	result, err := r.db.Exec(context.Background(), query,
		data.SKU, data.Name, data.CategoryID, data.RackID,
		data.Stock, data.MinimumStock, data.Price, id,
	)
	if err != nil {
		r.Logger.Error("error updating item", zap.Error(err))
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("item not found")
	}
	return nil
}

func (r *itemRepository) Delete(id int) error {
	query := `
		DELETE FROM items 
		WHERE id = $1
	`
	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		r.Logger.Error("error deleting item", zap.Error(err))
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("item not found")
	}
	return nil
}
