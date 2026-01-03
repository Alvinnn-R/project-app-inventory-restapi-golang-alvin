package service

import (
	"errors"
	"project-app-inventory/model"
	"project-app-inventory/repository"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockItemRepository mocks ItemRepository interface
type MockItemRepository struct {
	mock.Mock
}

func (m *MockItemRepository) Create(item *model.Item) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockItemRepository) FindByID(id int) (*model.Item, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Item), args.Error(1)
}

func (m *MockItemRepository) FindBySKU(sku string) (*model.Item, error) {
	args := m.Called(sku)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Item), args.Error(1)
}

func (m *MockItemRepository) FindAll(page, limit int) ([]model.Item, int, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]model.Item), args.Int(1), args.Error(2)
}

func (m *MockItemRepository) FindLowStock(page, limit int) ([]model.Item, int, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]model.Item), args.Int(1), args.Error(2)
}

func (m *MockItemRepository) Update(id int, item *model.Item) error {
	args := m.Called(id, item)
	return args.Error(0)
}

func (m *MockItemRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// TestItemService_Create_Success tests successful item creation
func TestItemService_Create_Success(t *testing.T) {
	mockItemRepo := new(MockItemRepository)
	repo := repository.Repository{ItemRepo: mockItemRepo}
	service := NewItemService(repo)

	item := &model.Item{
		SKU:        "SKU001",
		Name:       "Test Item",
		CategoryID: 1,
		Stock:      100,
		Price:      10000,
	}

	mockItemRepo.On("FindBySKU", item.SKU).Return((*model.Item)(nil), nil)
	mockItemRepo.On("Create", item).Return(nil)

	err := service.Create(item)

	require.NoError(t, err)
	mockItemRepo.AssertExpectations(t)
}

// TestItemService_Create_SKUExists tests creation with existing SKU
func TestItemService_Create_SKUExists(t *testing.T) {
	mockItemRepo := new(MockItemRepository)
	repo := repository.Repository{ItemRepo: mockItemRepo}
	service := NewItemService(repo)

	item := &model.Item{SKU: "SKU001"}
	existingItem := &model.Item{ID: 1, SKU: "SKU001"}

	mockItemRepo.On("FindBySKU", item.SKU).Return(existingItem, nil)

	err := service.Create(item)

	require.Error(t, err)
	require.Equal(t, "SKU already exists", err.Error())
	mockItemRepo.AssertExpectations(t)
}

// TestItemService_Create_CheckError tests creation when check fails
func TestItemService_Create_CheckError(t *testing.T) {
	mockItemRepo := new(MockItemRepository)
	repo := repository.Repository{ItemRepo: mockItemRepo}
	service := NewItemService(repo)

	item := &model.Item{SKU: "SKU001"}

	mockItemRepo.On("FindBySKU", item.SKU).Return((*model.Item)(nil), errors.New("db error"))

	err := service.Create(item)

	require.Error(t, err)
	require.Equal(t, "failed to check SKU", err.Error())
	mockItemRepo.AssertExpectations(t)
}

// TestItemService_GetAllItems_Success tests getting all items
func TestItemService_GetAllItems_Success(t *testing.T) {
	mockItemRepo := new(MockItemRepository)
	repo := repository.Repository{ItemRepo: mockItemRepo}
	service := NewItemService(repo)

	items := []model.Item{
		{ID: 1, Name: "Item 1"},
		{ID: 2, Name: "Item 2"},
	}

	mockItemRepo.On("FindAll", 1, 10).Return(items, 2, nil)

	result, pagination, err := service.GetAllItems(1, 10)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, 2, len(*result))
	require.NotNil(t, pagination)
	require.Equal(t, 1, pagination.CurrentPage)
	require.Equal(t, 10, pagination.Limit)
	require.Equal(t, 2, pagination.TotalRecords)
	mockItemRepo.AssertExpectations(t)
}

// TestItemService_GetAllItems_Error tests error handling
func TestItemService_GetAllItems_Error(t *testing.T) {
	mockItemRepo := new(MockItemRepository)
	repo := repository.Repository{ItemRepo: mockItemRepo}
	service := NewItemService(repo)

	mockItemRepo.On("FindAll", 1, 10).Return([]model.Item{}, 0, errors.New("db error"))

	result, pagination, err := service.GetAllItems(1, 10)

	require.Error(t, err)
	require.Nil(t, result)
	require.Nil(t, pagination)
	mockItemRepo.AssertExpectations(t)
}

// TestItemService_GetLowStockItems_Success tests getting low stock items
func TestItemService_GetLowStockItems_Success(t *testing.T) {
	mockItemRepo := new(MockItemRepository)
	repo := repository.Repository{ItemRepo: mockItemRepo}
	service := NewItemService(repo)

	items := []model.Item{
		{ID: 1, Name: "Low Stock Item", Stock: 3},
	}

	mockItemRepo.On("FindLowStock", 1, 10).Return(items, 1, nil)

	result, pagination, err := service.GetLowStockItems(1, 10)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, 1, len(*result))
	require.NotNil(t, pagination)
	mockItemRepo.AssertExpectations(t)
}

// TestItemService_GetItemByID_Success tests getting item by ID
func TestItemService_GetItemByID_Success(t *testing.T) {
	mockItemRepo := new(MockItemRepository)
	repo := repository.Repository{ItemRepo: mockItemRepo}
	service := NewItemService(repo)

	item := &model.Item{ID: 1, Name: "Test Item"}

	mockItemRepo.On("FindByID", 1).Return(item, nil)

	result, err := service.GetItemByID(1)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "Test Item", result.Name)
	mockItemRepo.AssertExpectations(t)
}

// TestItemService_GetItemByID_NotFound tests item not found
func TestItemService_GetItemByID_NotFound(t *testing.T) {
	mockItemRepo := new(MockItemRepository)
	repo := repository.Repository{ItemRepo: mockItemRepo}
	service := NewItemService(repo)

	mockItemRepo.On("FindByID", 999).Return((*model.Item)(nil), nil)

	result, err := service.GetItemByID(999)

	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, "item not found", err.Error())
	mockItemRepo.AssertExpectations(t)
}

// TestItemService_Update_Success tests successful update
func TestItemService_Update_Success(t *testing.T) {
	mockItemRepo := new(MockItemRepository)
	repo := repository.Repository{ItemRepo: mockItemRepo}
	service := NewItemService(repo)

	existingItem := &model.Item{
		ID:         1,
		SKU:        "SKU001",
		Name:       "Old Name",
		CategoryID: 1,
		Stock:      50,
		Price:      10000,
	}

	updateData := &model.Item{
		Name:  "New Name",
		Stock: 100,
		Price: 15000,
	}

	mockItemRepo.On("FindByID", 1).Return(existingItem, nil)
	mockItemRepo.On("Update", 1, updateData).Return(nil)

	err := service.Update(1, updateData)

	require.NoError(t, err)
	require.Equal(t, "SKU001", updateData.SKU) // Should keep existing SKU
	mockItemRepo.AssertExpectations(t)
}

// TestItemService_Update_NotFound tests update with non-existent item
func TestItemService_Update_NotFound(t *testing.T) {
	mockItemRepo := new(MockItemRepository)
	repo := repository.Repository{ItemRepo: mockItemRepo}
	service := NewItemService(repo)

	updateData := &model.Item{Name: "New Name"}

	mockItemRepo.On("FindByID", 999).Return((*model.Item)(nil), nil)

	err := service.Update(999, updateData)

	require.Error(t, err)
	require.Equal(t, "item not found", err.Error())
	mockItemRepo.AssertExpectations(t)
}

// TestItemService_Delete_Success tests successful deletion
func TestItemService_Delete_Success(t *testing.T) {
	mockItemRepo := new(MockItemRepository)
	repo := repository.Repository{ItemRepo: mockItemRepo}
	service := NewItemService(repo)

	existingItem := &model.Item{ID: 1, Name: "Test Item"}

	mockItemRepo.On("FindByID", 1).Return(existingItem, nil)
	mockItemRepo.On("Delete", 1).Return(nil)

	err := service.Delete(1)

	require.NoError(t, err)
	mockItemRepo.AssertExpectations(t)
}

// TestItemService_Delete_NotFound tests deletion with non-existent item
func TestItemService_Delete_NotFound(t *testing.T) {
	mockItemRepo := new(MockItemRepository)
	repo := repository.Repository{ItemRepo: mockItemRepo}
	service := NewItemService(repo)

	mockItemRepo.On("FindByID", 999).Return((*model.Item)(nil), nil)

	err := service.Delete(999)

	require.Error(t, err)
	require.Equal(t, "item not found", err.Error())
	mockItemRepo.AssertExpectations(t)
}
