package repository

import (
	"context"
	"project-app-inventory/database"

	"go.uber.org/zap"
)

type ReportRepository interface {
	GetTotalItems() (int, error)
	GetLowStockItems() (int, error)
	GetTotalSales() (int, error)
	GetTotalRevenue() (float64, error)
	GetActiveUsers() (int, error)
	GetTotalCategories() (int, error)
	GetTotalWarehouses() (int, error)
}

type reportRepository struct {
	db     database.PgxIface
	Logger *zap.Logger
}

func NewReportRepository(db database.PgxIface, log *zap.Logger) ReportRepository {
	return &reportRepository{db: db, Logger: log}
}

func (r *reportRepository) GetTotalItems() (int, error) {
	var total int
	query := `SELECT COUNT(*) FROM items`
	err := r.db.QueryRow(context.Background(), query).Scan(&total)
	if err != nil {
		if r.Logger != nil {
			r.Logger.Error("error getting total items", zap.Error(err))
		}
		return 0, err
	}
	return total, nil
}

func (r *reportRepository) GetLowStockItems() (int, error) {
	var total int
	query := `SELECT COUNT(*) FROM items WHERE stock < minimum_stock`
	err := r.db.QueryRow(context.Background(), query).Scan(&total)
	if err != nil {
		if r.Logger != nil {
			r.Logger.Error("error getting low stock items", zap.Error(err))
		}
		return 0, err
	}
	return total, nil
}

func (r *reportRepository) GetTotalSales() (int, error) {
	var total int
	query := `SELECT COUNT(*) FROM sales WHERE deleted_at IS NULL`
	err := r.db.QueryRow(context.Background(), query).Scan(&total)
	if err != nil {
		if r.Logger != nil {
			r.Logger.Error("error getting total sales", zap.Error(err))
		}
		return 0, err
	}
	return total, nil
}

func (r *reportRepository) GetTotalRevenue() (float64, error) {
	var total float64
	query := `SELECT COALESCE(SUM(total_amount), 0) FROM sales WHERE deleted_at IS NULL`
	err := r.db.QueryRow(context.Background(), query).Scan(&total)
	if err != nil {
		if r.Logger != nil {
			r.Logger.Error("error getting total revenue", zap.Error(err))
		}
		return 0, err
	}
	return total, nil
}

func (r *reportRepository) GetActiveUsers() (int, error) {
	var total int
	query := `SELECT COUNT(*) FROM users WHERE is_active = true`
	err := r.db.QueryRow(context.Background(), query).Scan(&total)
	if err != nil {
		if r.Logger != nil {
			r.Logger.Error("error getting active users", zap.Error(err))
		}
		return 0, err
	}
	return total, nil
}

func (r *reportRepository) GetTotalCategories() (int, error) {
	var total int
	query := `SELECT COUNT(*) FROM categories`
	err := r.db.QueryRow(context.Background(), query).Scan(&total)
	if err != nil {
		if r.Logger != nil {
			r.Logger.Error("error getting total categories", zap.Error(err))
		}
		return 0, err
	}
	return total, nil
}

func (r *reportRepository) GetTotalWarehouses() (int, error) {
	var total int
	query := `SELECT COUNT(*) FROM warehouses`
	err := r.db.QueryRow(context.Background(), query).Scan(&total)
	if err != nil {
		if r.Logger != nil {
			r.Logger.Error("error getting total warehouses", zap.Error(err))
		}
		return 0, err
	}
	return total, nil
}
