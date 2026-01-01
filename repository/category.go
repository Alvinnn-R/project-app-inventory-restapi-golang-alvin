package repository

import (
	"context"
	"errors"
	"project-app-inventory/database"
	"project-app-inventory/model"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type CategoryRepository interface {
	Create(category *model.Category) error
	FindByID(id int) (*model.Category, error)
	FindByName(name string) (*model.Category, error)
	FindAll(page, limit int) ([]model.Category, int, error)
	Update(id int, data *model.Category) error
	Delete(id int) error
}

type categoryRepository struct {
	db     database.PgxIface
	Logger *zap.Logger
}

func NewCategoryRepository(db database.PgxIface, log *zap.Logger) CategoryRepository {
	return &categoryRepository{db: db, Logger: log}
}

func (r *categoryRepository) Create(category *model.Category) error {
	query := `
		INSERT INTO categories (name, description, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(context.Background(), query,
		category.Name, category.Description,
	).Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		r.Logger.Error("error creating category", zap.Error(err))
	}
	return err
}

func (r *categoryRepository) FindByID(id int) (*model.Category, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM categories 
		WHERE id = $1
	`
	var category model.Category
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&category.ID, &category.Name, &category.Description,
		&category.CreatedAt, &category.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		r.Logger.Error("error finding category by id", zap.Error(err))
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) FindByName(name string) (*model.Category, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM categories 
		WHERE name = $1
	`
	var category model.Category
	err := r.db.QueryRow(context.Background(), query, name).Scan(
		&category.ID, &category.Name, &category.Description,
		&category.CreatedAt, &category.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		r.Logger.Error("error finding category by name", zap.Error(err))
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) FindAll(page, limit int) ([]model.Category, int, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM categories`
	err := r.db.QueryRow(context.Background(), countQuery).Scan(&total)
	if err != nil {
		r.Logger.Error("error counting categories", zap.Error(err))
		return nil, 0, err
	}

	// Get data with pagination
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM categories
		ORDER BY name ASC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		r.Logger.Error("error querying categories", zap.Error(err))
		return nil, 0, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var category model.Category
		err := rows.Scan(
			&category.ID, &category.Name, &category.Description,
			&category.CreatedAt, &category.UpdatedAt,
		)
		if err != nil {
			r.Logger.Error("error scanning category", zap.Error(err))
			return nil, 0, err
		}
		categories = append(categories, category)
	}

	return categories, total, nil
}

func (r *categoryRepository) Update(id int, data *model.Category) error {
	query := `
		UPDATE categories
		SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $3
	`
	result, err := r.db.Exec(context.Background(), query,
		data.Name, data.Description, id,
	)
	if err != nil {
		r.Logger.Error("error updating category", zap.Error(err))
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("category not found")
	}
	return nil
}

func (r *categoryRepository) Delete(id int) error {
	query := `
		DELETE FROM categories 
		WHERE id = $1
	`
	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		r.Logger.Error("error deleting category", zap.Error(err))
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("category not found")
	}
	return nil
}
