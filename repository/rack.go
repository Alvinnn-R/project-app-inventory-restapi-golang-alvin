package repository

import (
	"context"
	"errors"
	"project-app-inventory/database"
	"project-app-inventory/model"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type RackRepository interface {
	Create(rack *model.Rack) error
	FindByID(id int) (*model.Rack, error)
	FindByWarehouseAndCode(warehouseID int, code string) (*model.Rack, error)
	FindAll(page, limit int) ([]model.Rack, int, error)
	FindByWarehouseID(warehouseID, page, limit int) ([]model.Rack, int, error)
	Update(id int, data *model.Rack) error
	Delete(id int) error
}

type rackRepository struct {
	db     database.PgxIface
	Logger *zap.Logger
}

func NewRackRepository(db database.PgxIface, log *zap.Logger) RackRepository {
	return &rackRepository{db: db, Logger: log}
}

func (r *rackRepository) Create(rack *model.Rack) error {
	query := `
		INSERT INTO racks (warehouse_id, code, description, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(context.Background(), query,
		rack.WarehouseID, rack.Code, rack.Description,
	).Scan(&rack.ID, &rack.CreatedAt, &rack.UpdatedAt)

	if err != nil {
		r.Logger.Error("error creating rack", zap.Error(err))
	}
	return err
}

func (r *rackRepository) FindByID(id int) (*model.Rack, error) {
	query := `
		SELECT id, warehouse_id, code, description, created_at, updated_at
		FROM racks 
		WHERE id = $1
	`
	var rack model.Rack
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&rack.ID, &rack.WarehouseID, &rack.Code, &rack.Description,
		&rack.CreatedAt, &rack.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		r.Logger.Error("error finding rack by id", zap.Error(err))
		return nil, err
	}
	return &rack, nil
}

func (r *rackRepository) FindByWarehouseAndCode(warehouseID int, code string) (*model.Rack, error) {
	query := `
		SELECT id, warehouse_id, code, description, created_at, updated_at
		FROM racks 
		WHERE warehouse_id = $1 AND code = $2
	`
	var rack model.Rack
	err := r.db.QueryRow(context.Background(), query, warehouseID, code).Scan(
		&rack.ID, &rack.WarehouseID, &rack.Code, &rack.Description,
		&rack.CreatedAt, &rack.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		r.Logger.Error("error finding rack by warehouse and code", zap.Error(err))
		return nil, err
	}
	return &rack, nil
}

func (r *rackRepository) FindAll(page, limit int) ([]model.Rack, int, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM racks`
	err := r.db.QueryRow(context.Background(), countQuery).Scan(&total)
	if err != nil {
		r.Logger.Error("error counting racks", zap.Error(err))
		return nil, 0, err
	}

	// Get data with pagination
	query := `
		SELECT id, warehouse_id, code, description, created_at, updated_at
		FROM racks
		ORDER BY warehouse_id ASC, code ASC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		r.Logger.Error("error querying racks", zap.Error(err))
		return nil, 0, err
	}
	defer rows.Close()

	var racks []model.Rack
	for rows.Next() {
		var rack model.Rack
		err := rows.Scan(
			&rack.ID, &rack.WarehouseID, &rack.Code, &rack.Description,
			&rack.CreatedAt, &rack.UpdatedAt,
		)
		if err != nil {
			r.Logger.Error("error scanning rack", zap.Error(err))
			return nil, 0, err
		}
		racks = append(racks, rack)
	}

	return racks, total, nil
}

func (r *rackRepository) FindByWarehouseID(warehouseID, page, limit int) ([]model.Rack, int, error) {
	offset := (page - 1) * limit

	// Get total count for warehouse
	var total int
	countQuery := `SELECT COUNT(*) FROM racks WHERE warehouse_id = $1`
	err := r.db.QueryRow(context.Background(), countQuery, warehouseID).Scan(&total)
	if err != nil {
		r.Logger.Error("error counting racks by warehouse", zap.Error(err))
		return nil, 0, err
	}

	// Get data with pagination
	query := `
		SELECT id, warehouse_id, code, description, created_at, updated_at
		FROM racks
		WHERE warehouse_id = $1
		ORDER BY code ASC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(context.Background(), query, warehouseID, limit, offset)
	if err != nil {
		r.Logger.Error("error querying racks by warehouse", zap.Error(err))
		return nil, 0, err
	}
	defer rows.Close()

	var racks []model.Rack
	for rows.Next() {
		var rack model.Rack
		err := rows.Scan(
			&rack.ID, &rack.WarehouseID, &rack.Code, &rack.Description,
			&rack.CreatedAt, &rack.UpdatedAt,
		)
		if err != nil {
			r.Logger.Error("error scanning rack", zap.Error(err))
			return nil, 0, err
		}
		racks = append(racks, rack)
	}

	return racks, total, nil
}

func (r *rackRepository) Update(id int, data *model.Rack) error {
	query := `
		UPDATE racks
		SET warehouse_id = $1, code = $2, description = $3, updated_at = NOW()
		WHERE id = $4
	`
	result, err := r.db.Exec(context.Background(), query,
		data.WarehouseID, data.Code, data.Description, id,
	)
	if err != nil {
		r.Logger.Error("error updating rack", zap.Error(err))
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("rack not found")
	}
	return nil
}

func (r *rackRepository) Delete(id int) error {
	query := `
		DELETE FROM racks 
		WHERE id = $1
	`
	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		r.Logger.Error("error deleting rack", zap.Error(err))
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("rack not found")
	}
	return nil
}
