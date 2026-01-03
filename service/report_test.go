package service

import (
	"errors"
	"project-app-inventory/repository"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockReportRepository mocks ReportRepository interface
type MockReportRepository struct {
	mock.Mock
}

func (m *MockReportRepository) GetTotalItems() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockReportRepository) GetLowStockItems() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockReportRepository) GetTotalSales() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockReportRepository) GetTotalRevenue() (float64, error) {
	args := m.Called()
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockReportRepository) GetActiveUsers() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockReportRepository) GetTotalCategories() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockReportRepository) GetTotalWarehouses() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

// TestReportService_GetSummary_Success tests getting report summary successfully
func TestReportService_GetSummary_Success(t *testing.T) {
	mockReportRepo := new(MockReportRepository)
	repo := repository.Repository{ReportRepo: mockReportRepo}
	service := NewReportService(&repo)

	mockReportRepo.On("GetTotalItems").Return(100, nil)
	mockReportRepo.On("GetLowStockItems").Return(5, nil)
	mockReportRepo.On("GetTotalSales").Return(50, nil)
	mockReportRepo.On("GetTotalRevenue").Return(1000000.0, nil)
	mockReportRepo.On("GetActiveUsers").Return(25, nil)
	mockReportRepo.On("GetTotalCategories").Return(10, nil)
	mockReportRepo.On("GetTotalWarehouses").Return(3, nil)

	result, err := service.GetSummary()

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, 100, result.TotalItems)
	require.Equal(t, 5, result.LowStockItems)
	require.Equal(t, 50, result.TotalSales)
	require.Equal(t, 1000000.0, result.TotalRevenue)
	require.Equal(t, 25, result.ActiveUsers)
	require.Equal(t, 10, result.TotalCategories)
	require.Equal(t, 3, result.TotalWarehouses)
	mockReportRepo.AssertExpectations(t)
}

// TestReportService_GetSummary_TotalItemsError tests error on GetTotalItems
func TestReportService_GetSummary_TotalItemsError(t *testing.T) {
	mockReportRepo := new(MockReportRepository)
	repo := repository.Repository{ReportRepo: mockReportRepo}
	service := NewReportService(&repo)

	mockReportRepo.On("GetTotalItems").Return(0, errors.New("db error"))

	result, err := service.GetSummary()

	require.Error(t, err)
	require.Nil(t, result)
	mockReportRepo.AssertExpectations(t)
}

// TestReportService_GetSummary_LowStockError tests error on GetLowStockItems
func TestReportService_GetSummary_LowStockError(t *testing.T) {
	mockReportRepo := new(MockReportRepository)
	repo := repository.Repository{ReportRepo: mockReportRepo}
	service := NewReportService(&repo)

	mockReportRepo.On("GetTotalItems").Return(100, nil)
	mockReportRepo.On("GetLowStockItems").Return(0, errors.New("db error"))

	result, err := service.GetSummary()

	require.Error(t, err)
	require.Nil(t, result)
	mockReportRepo.AssertExpectations(t)
}

// TestReportService_GetSummary_TotalSalesError tests error on GetTotalSales
func TestReportService_GetSummary_TotalSalesError(t *testing.T) {
	mockReportRepo := new(MockReportRepository)
	repo := repository.Repository{ReportRepo: mockReportRepo}
	service := NewReportService(&repo)

	mockReportRepo.On("GetTotalItems").Return(100, nil)
	mockReportRepo.On("GetLowStockItems").Return(5, nil)
	mockReportRepo.On("GetTotalSales").Return(0, errors.New("db error"))

	result, err := service.GetSummary()

	require.Error(t, err)
	require.Nil(t, result)
	mockReportRepo.AssertExpectations(t)
}

// TestReportService_GetSummary_TotalRevenueError tests error on GetTotalRevenue
func TestReportService_GetSummary_TotalRevenueError(t *testing.T) {
	mockReportRepo := new(MockReportRepository)
	repo := repository.Repository{ReportRepo: mockReportRepo}
	service := NewReportService(&repo)

	mockReportRepo.On("GetTotalItems").Return(100, nil)
	mockReportRepo.On("GetLowStockItems").Return(5, nil)
	mockReportRepo.On("GetTotalSales").Return(50, nil)
	mockReportRepo.On("GetTotalRevenue").Return(0.0, errors.New("db error"))

	result, err := service.GetSummary()

	require.Error(t, err)
	require.Nil(t, result)
	mockReportRepo.AssertExpectations(t)
}

// TestReportService_GetSummary_ActiveUsersError tests error on GetActiveUsers
func TestReportService_GetSummary_ActiveUsersError(t *testing.T) {
	mockReportRepo := new(MockReportRepository)
	repo := repository.Repository{ReportRepo: mockReportRepo}
	service := NewReportService(&repo)

	mockReportRepo.On("GetTotalItems").Return(100, nil)
	mockReportRepo.On("GetLowStockItems").Return(5, nil)
	mockReportRepo.On("GetTotalSales").Return(50, nil)
	mockReportRepo.On("GetTotalRevenue").Return(1000000.0, nil)
	mockReportRepo.On("GetActiveUsers").Return(0, errors.New("db error"))

	result, err := service.GetSummary()

	require.Error(t, err)
	require.Nil(t, result)
	mockReportRepo.AssertExpectations(t)
}

// TestReportService_GetSummary_TotalCategoriesError tests error on GetTotalCategories
func TestReportService_GetSummary_TotalCategoriesError(t *testing.T) {
	mockReportRepo := new(MockReportRepository)
	repo := repository.Repository{ReportRepo: mockReportRepo}
	service := NewReportService(&repo)

	mockReportRepo.On("GetTotalItems").Return(100, nil)
	mockReportRepo.On("GetLowStockItems").Return(5, nil)
	mockReportRepo.On("GetTotalSales").Return(50, nil)
	mockReportRepo.On("GetTotalRevenue").Return(1000000.0, nil)
	mockReportRepo.On("GetActiveUsers").Return(25, nil)
	mockReportRepo.On("GetTotalCategories").Return(0, errors.New("db error"))

	result, err := service.GetSummary()

	require.Error(t, err)
	require.Nil(t, result)
	mockReportRepo.AssertExpectations(t)
}

// TestReportService_GetSummary_TotalWarehousesError tests error on GetTotalWarehouses
func TestReportService_GetSummary_TotalWarehousesError(t *testing.T) {
	mockReportRepo := new(MockReportRepository)
	repo := repository.Repository{ReportRepo: mockReportRepo}
	service := NewReportService(&repo)

	mockReportRepo.On("GetTotalItems").Return(100, nil)
	mockReportRepo.On("GetLowStockItems").Return(5, nil)
	mockReportRepo.On("GetTotalSales").Return(50, nil)
	mockReportRepo.On("GetTotalRevenue").Return(1000000.0, nil)
	mockReportRepo.On("GetActiveUsers").Return(25, nil)
	mockReportRepo.On("GetTotalCategories").Return(10, nil)
	mockReportRepo.On("GetTotalWarehouses").Return(0, errors.New("db error"))

	result, err := service.GetSummary()

	require.Error(t, err)
	require.Nil(t, result)
	mockReportRepo.AssertExpectations(t)
}
