package service

import (
	"errors"
	"project-app-inventory/model"
	"project-app-inventory/repository"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockWarehouseRepository mocks WarehouseRepository interface
type MockWarehouseRepository struct {
	mock.Mock
}

func (m *MockWarehouseRepository) Create(warehouse *model.Warehouse) error {
	args := m.Called(warehouse)
	return args.Error(0)
}

func (m *MockWarehouseRepository) FindByID(id int) (*model.Warehouse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Warehouse), args.Error(1)
}

func (m *MockWarehouseRepository) FindByName(name string) (*model.Warehouse, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Warehouse), args.Error(1)
}

func (m *MockWarehouseRepository) FindAll(page, limit int) ([]model.Warehouse, int, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]model.Warehouse), args.Int(1), args.Error(2)
}

func (m *MockWarehouseRepository) Update(id int, warehouse *model.Warehouse) error {
	args := m.Called(id, warehouse)
	return args.Error(0)
}

func (m *MockWarehouseRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// TestWarehouseService_Create_Success tests successful warehouse creation
func TestWarehouseService_Create_Success(t *testing.T) {
	mockWarehouseRepo := new(MockWarehouseRepository)
	repo := repository.Repository{WarehouseRepo: mockWarehouseRepo}
	service := NewWarehouseService(repo)

	warehouse := &model.Warehouse{
		Name:     "Main Warehouse",
		Location: "Jakarta",
	}

	mockWarehouseRepo.On("FindByName", warehouse.Name).Return((*model.Warehouse)(nil), nil)
	mockWarehouseRepo.On("Create", warehouse).Return(nil)

	err := service.Create(warehouse)

	require.NoError(t, err)
	mockWarehouseRepo.AssertExpectations(t)
}

// TestWarehouseService_Create_NameExists tests creation with existing name
func TestWarehouseService_Create_NameExists(t *testing.T) {
	mockWarehouseRepo := new(MockWarehouseRepository)
	repo := repository.Repository{WarehouseRepo: mockWarehouseRepo}
	service := NewWarehouseService(repo)

	warehouse := &model.Warehouse{Name: "Main Warehouse"}
	existingWarehouse := &model.Warehouse{ID: 1, Name: "Main Warehouse"}

	mockWarehouseRepo.On("FindByName", warehouse.Name).Return(existingWarehouse, nil)

	err := service.Create(warehouse)

	require.Error(t, err)
	require.Equal(t, "warehouse name already exists", err.Error())
	mockWarehouseRepo.AssertExpectations(t)
}

// TestWarehouseService_Create_CheckError tests creation when check fails
func TestWarehouseService_Create_CheckError(t *testing.T) {
	mockWarehouseRepo := new(MockWarehouseRepository)
	repo := repository.Repository{WarehouseRepo: mockWarehouseRepo}
	service := NewWarehouseService(repo)

	warehouse := &model.Warehouse{Name: "Main Warehouse"}

	mockWarehouseRepo.On("FindByName", warehouse.Name).Return((*model.Warehouse)(nil), errors.New("db error"))

	err := service.Create(warehouse)

	require.Error(t, err)
	require.Equal(t, "failed to check warehouse name", err.Error())
	mockWarehouseRepo.AssertExpectations(t)
}

// TestWarehouseService_GetAllWarehouses_Success tests getting all warehouses
func TestWarehouseService_GetAllWarehouses_Success(t *testing.T) {
	mockWarehouseRepo := new(MockWarehouseRepository)
	repo := repository.Repository{WarehouseRepo: mockWarehouseRepo}
	service := NewWarehouseService(repo)

	warehouses := []model.Warehouse{
		{ID: 1, Name: "Main Warehouse"},
		{ID: 2, Name: "Secondary Warehouse"},
	}

	mockWarehouseRepo.On("FindAll", 1, 10).Return(warehouses, 2, nil)

	result, pagination, err := service.GetAllWarehouses(1, 10)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, 2, len(*result))
	require.NotNil(t, pagination)
	require.Equal(t, 1, pagination.CurrentPage)
	require.Equal(t, 10, pagination.Limit)
	mockWarehouseRepo.AssertExpectations(t)
}

// TestWarehouseService_GetAllWarehouses_Error tests error handling
func TestWarehouseService_GetAllWarehouses_Error(t *testing.T) {
	mockWarehouseRepo := new(MockWarehouseRepository)
	repo := repository.Repository{WarehouseRepo: mockWarehouseRepo}
	service := NewWarehouseService(repo)

	mockWarehouseRepo.On("FindAll", 1, 10).Return([]model.Warehouse{}, 0, errors.New("db error"))

	result, pagination, err := service.GetAllWarehouses(1, 10)

	require.Error(t, err)
	require.Nil(t, result)
	require.Nil(t, pagination)
	mockWarehouseRepo.AssertExpectations(t)
}

// TestWarehouseService_GetWarehouseByID_Success tests getting warehouse by ID
func TestWarehouseService_GetWarehouseByID_Success(t *testing.T) {
	mockWarehouseRepo := new(MockWarehouseRepository)
	repo := repository.Repository{WarehouseRepo: mockWarehouseRepo}
	service := NewWarehouseService(repo)

	warehouse := &model.Warehouse{ID: 1, Name: "Main Warehouse"}

	mockWarehouseRepo.On("FindByID", 1).Return(warehouse, nil)

	result, err := service.GetWarehouseByID(1)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "Main Warehouse", result.Name)
	mockWarehouseRepo.AssertExpectations(t)
}

// TestWarehouseService_GetWarehouseByID_NotFound tests warehouse not found
func TestWarehouseService_GetWarehouseByID_NotFound(t *testing.T) {
	mockWarehouseRepo := new(MockWarehouseRepository)
	repo := repository.Repository{WarehouseRepo: mockWarehouseRepo}
	service := NewWarehouseService(repo)

	mockWarehouseRepo.On("FindByID", 999).Return((*model.Warehouse)(nil), nil)

	result, err := service.GetWarehouseByID(999)

	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, "warehouse not found", err.Error())
	mockWarehouseRepo.AssertExpectations(t)
}

// TestWarehouseService_Update_Success tests successful update
func TestWarehouseService_Update_Success(t *testing.T) {
	mockWarehouseRepo := new(MockWarehouseRepository)
	repo := repository.Repository{WarehouseRepo: mockWarehouseRepo}
	service := NewWarehouseService(repo)

	existingWarehouse := &model.Warehouse{
		ID:       1,
		Name:     "Main Warehouse",
		Location: "Jakarta",
	}

	updateData := &model.Warehouse{
		Name:     "Updated Warehouse",
		Location: "Bandung",
	}

	mockWarehouseRepo.On("FindByID", 1).Return(existingWarehouse, nil)
	mockWarehouseRepo.On("FindByName", "Updated Warehouse").Return((*model.Warehouse)(nil), nil)
	mockWarehouseRepo.On("Update", 1, updateData).Return(nil)

	err := service.Update(1, updateData)

	require.NoError(t, err)
	mockWarehouseRepo.AssertExpectations(t)
}

// TestWarehouseService_Update_NameExists tests update with existing name
func TestWarehouseService_Update_NameExists(t *testing.T) {
	mockWarehouseRepo := new(MockWarehouseRepository)
	repo := repository.Repository{WarehouseRepo: mockWarehouseRepo}
	service := NewWarehouseService(repo)

	existingWarehouse := &model.Warehouse{ID: 1, Name: "Main Warehouse"}
	otherWarehouse := &model.Warehouse{ID: 2, Name: "Secondary Warehouse"}
	updateData := &model.Warehouse{Name: "Secondary Warehouse"}

	mockWarehouseRepo.On("FindByID", 1).Return(existingWarehouse, nil)
	mockWarehouseRepo.On("FindByName", "Secondary Warehouse").Return(otherWarehouse, nil)

	err := service.Update(1, updateData)

	require.Error(t, err)
	require.Equal(t, "warehouse name already exists", err.Error())
	mockWarehouseRepo.AssertExpectations(t)
}

// TestWarehouseService_Update_NotFound tests update with non-existent warehouse
func TestWarehouseService_Update_NotFound(t *testing.T) {
	mockWarehouseRepo := new(MockWarehouseRepository)
	repo := repository.Repository{WarehouseRepo: mockWarehouseRepo}
	service := NewWarehouseService(repo)

	updateData := &model.Warehouse{Name: "New Name"}

	mockWarehouseRepo.On("FindByID", 999).Return((*model.Warehouse)(nil), nil)

	err := service.Update(999, updateData)

	require.Error(t, err)
	require.Equal(t, "warehouse not found", err.Error())
	mockWarehouseRepo.AssertExpectations(t)
}

// TestWarehouseService_Delete_Success tests successful deletion
func TestWarehouseService_Delete_Success(t *testing.T) {
	mockWarehouseRepo := new(MockWarehouseRepository)
	repo := repository.Repository{WarehouseRepo: mockWarehouseRepo}
	service := NewWarehouseService(repo)

	existingWarehouse := &model.Warehouse{ID: 1, Name: "Main Warehouse"}

	mockWarehouseRepo.On("FindByID", 1).Return(existingWarehouse, nil)
	mockWarehouseRepo.On("Delete", 1).Return(nil)

	err := service.Delete(1)

	require.NoError(t, err)
	mockWarehouseRepo.AssertExpectations(t)
}

// TestWarehouseService_Delete_NotFound tests deletion with non-existent warehouse
func TestWarehouseService_Delete_NotFound(t *testing.T) {
	mockWarehouseRepo := new(MockWarehouseRepository)
	repo := repository.Repository{WarehouseRepo: mockWarehouseRepo}
	service := NewWarehouseService(repo)

	mockWarehouseRepo.On("FindByID", 999).Return((*model.Warehouse)(nil), nil)

	err := service.Delete(999)

	require.Error(t, err)
	require.Equal(t, "warehouse not found", err.Error())
	mockWarehouseRepo.AssertExpectations(t)
}
