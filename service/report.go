package service

import (
	"project-app-inventory/dto"
	"project-app-inventory/repository"
)

type ReportService interface {
	GetSummary() (*dto.ReportSummaryResponse, error)
}

type reportService struct {
	Repo *repository.Repository
}

func NewReportService(repo *repository.Repository) ReportService {
	return &reportService{Repo: repo}
}

func (s *reportService) GetSummary() (*dto.ReportSummaryResponse, error) {
	// Get all metrics in parallel for better performance
	var report dto.ReportSummaryResponse

	// Total items
	totalItems, err := s.Repo.ReportRepo.GetTotalItems()
	if err != nil {
		return nil, err
	}
	report.TotalItems = totalItems

	// Low stock items
	lowStockItems, err := s.Repo.ReportRepo.GetLowStockItems()
	if err != nil {
		return nil, err
	}
	report.LowStockItems = lowStockItems

	// Total sales
	totalSales, err := s.Repo.ReportRepo.GetTotalSales()
	if err != nil {
		return nil, err
	}
	report.TotalSales = totalSales

	// Total revenue
	totalRevenue, err := s.Repo.ReportRepo.GetTotalRevenue()
	if err != nil {
		return nil, err
	}
	report.TotalRevenue = totalRevenue

	// Active users
	activeUsers, err := s.Repo.ReportRepo.GetActiveUsers()
	if err != nil {
		return nil, err
	}
	report.ActiveUsers = activeUsers

	// Total categories
	totalCategories, err := s.Repo.ReportRepo.GetTotalCategories()
	if err != nil {
		return nil, err
	}
	report.TotalCategories = totalCategories

	// Total warehouses
	totalWarehouses, err := s.Repo.ReportRepo.GetTotalWarehouses()
	if err != nil {
		return nil, err
	}
	report.TotalWarehouses = totalWarehouses

	return &report, nil
}
