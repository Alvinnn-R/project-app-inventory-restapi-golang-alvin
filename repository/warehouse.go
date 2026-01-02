package repository

import (
	"context"
	"errors"
	"project-app-inventory/database"
	"project-app-inventory/model"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type WarehouseRepository interface {
	Create(warehouse *model.Warehouse) error
	FindByID(id int) (*model.Warehouse, error)
	FindByName(name string) (*model.Warehouse, error)
	FindAll(page, limit int) ([]model.Warehouse, int, error)
	Update(id int, data *model.Warehouse) error
	Delete(id int) error
}

type warehouseRepository struct {
	db     database.PgxIface
	Logger *zap.Logger
}

func NewWarehouseRepository(db database.PgxIface, log *zap.Logger) WarehouseRepository {
	return &warehouseRepository{db: db, Logger: log}
}

func (r *warehouseRepository) Create(warehouse *model.Warehouse) error {
	query := `
		INSERT INTO warehouses (name, location, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(context.Background(), query,
		warehouse.Name, warehouse.Location,
	).Scan(&warehouse.ID, &warehouse.CreatedAt, &warehouse.UpdatedAt)

	if err != nil {
		r.Logger.Error("error creating warehouse", zap.Error(err))
	}
	return err
}

func (r *warehouseRepository) FindByID(id int) (*model.Warehouse, error) {
	query := `
		SELECT id, name, location, created_at, updated_at
		FROM warehouses 
		WHERE id = $1
	`
	var warehouse model.Warehouse
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&warehouse.ID, &warehouse.Name, &warehouse.Location,
		&warehouse.CreatedAt, &warehouse.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		r.Logger.Error("error finding warehouse by id", zap.Error(err))
		return nil, err
	}
	return &warehouse, nil
}

func (r *warehouseRepository) FindByName(name string) (*model.Warehouse, error) {
	query := `
		SELECT id, name, location, created_at, updated_at
		FROM warehouses 
		WHERE name = $1
	`
	var warehouse model.Warehouse
	err := r.db.QueryRow(context.Background(), query, name).Scan(
		&warehouse.ID, &warehouse.Name, &warehouse.Location,
		&warehouse.CreatedAt, &warehouse.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		r.Logger.Error("error finding warehouse by name", zap.Error(err))
		return nil, err
	}
	return &warehouse, nil
}

func (r *warehouseRepository) FindAll(page, limit int) ([]model.Warehouse, int, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM warehouses`
	err := r.db.QueryRow(context.Background(), countQuery).Scan(&total)
	if err != nil {
		r.Logger.Error("error counting warehouses", zap.Error(err))
		return nil, 0, err
	}

	// Get data with pagination
	query := `
		SELECT id, name, location, created_at, updated_at
		FROM warehouses
		ORDER BY name ASC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		r.Logger.Error("error querying warehouses", zap.Error(err))
		return nil, 0, err
	}
	defer rows.Close()

	var warehouses []model.Warehouse
	for rows.Next() {
		var warehouse model.Warehouse
		err := rows.Scan(
			&warehouse.ID, &warehouse.Name, &warehouse.Location,
			&warehouse.CreatedAt, &warehouse.UpdatedAt,
		)
		if err != nil {
			r.Logger.Error("error scanning warehouse", zap.Error(err))
			return nil, 0, err
		}
		warehouses = append(warehouses, warehouse)
	}

	return warehouses, total, nil
}

func (r *warehouseRepository) Update(id int, data *model.Warehouse) error {
	query := `
		UPDATE warehouses
		SET name = $1, location = $2, updated_at = NOW()
		WHERE id = $3
	`
	result, err := r.db.Exec(context.Background(), query,
		data.Name, data.Location, id,
	)
	if err != nil {
		r.Logger.Error("error updating warehouse", zap.Error(err))
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("warehouse not found")
	}
	return nil
}

func (r *warehouseRepository) Delete(id int) error {
	query := `
		DELETE FROM warehouses 
		WHERE id = $1
	`
	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		r.Logger.Error("error deleting warehouse", zap.Error(err))
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("warehouse not found")
	}
	return nil
}
