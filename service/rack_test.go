package service

import (
	"errors"
	"project-app-inventory/model"
	"project-app-inventory/repository"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockRackRepository mocks RackRepository interface
type MockRackRepository struct {
	mock.Mock
}

func (m *MockRackRepository) Create(rack *model.Rack) error {
	args := m.Called(rack)
	return args.Error(0)
}

func (m *MockRackRepository) FindByID(id int) (*model.Rack, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Rack), args.Error(1)
}

func (m *MockRackRepository) FindByWarehouseAndCode(warehouseID int, code string) (*model.Rack, error) {
	args := m.Called(warehouseID, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Rack), args.Error(1)
}

func (m *MockRackRepository) FindAll(page, limit int) ([]model.Rack, int, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]model.Rack), args.Int(1), args.Error(2)
}

func (m *MockRackRepository) FindByWarehouseID(warehouseID, page, limit int) ([]model.Rack, int, error) {
	args := m.Called(warehouseID, page, limit)
	return args.Get(0).([]model.Rack), args.Int(1), args.Error(2)
}

func (m *MockRackRepository) Update(id int, rack *model.Rack) error {
	args := m.Called(id, rack)
	return args.Error(0)
}

func (m *MockRackRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// TestRackService_Create_Success tests successful rack creation
func TestRackService_Create_Success(t *testing.T) {
	mockRackRepo := new(MockRackRepository)
	repo := repository.Repository{RackRepo: mockRackRepo}
	service := NewRackService(repo)

	rack := &model.Rack{
		WarehouseID: 1,
		Code:        "A1",
	}

	mockRackRepo.On("FindByWarehouseAndCode", 1, "A1").Return((*model.Rack)(nil), nil)
	mockRackRepo.On("Create", rack).Return(nil)

	err := service.Create(rack)

	require.NoError(t, err)
	mockRackRepo.AssertExpectations(t)
}

// TestRackService_Create_CodeExists tests creation with existing code
func TestRackService_Create_CodeExists(t *testing.T) {
	mockRackRepo := new(MockRackRepository)
	repo := repository.Repository{RackRepo: mockRackRepo}
	service := NewRackService(repo)

	rack := &model.Rack{WarehouseID: 1, Code: "A1"}
	existingRack := &model.Rack{ID: 1, WarehouseID: 1, Code: "A1"}

	mockRackRepo.On("FindByWarehouseAndCode", 1, "A1").Return(existingRack, nil)

	err := service.Create(rack)

	require.Error(t, err)
	require.Equal(t, "rack code already exists in this warehouse", err.Error())
	mockRackRepo.AssertExpectations(t)
}

// TestRackService_Create_CheckError tests creation when check fails
func TestRackService_Create_CheckError(t *testing.T) {
	mockRackRepo := new(MockRackRepository)
	repo := repository.Repository{RackRepo: mockRackRepo}
	service := NewRackService(repo)

	rack := &model.Rack{WarehouseID: 1, Code: "A1"}

	mockRackRepo.On("FindByWarehouseAndCode", 1, "A1").Return((*model.Rack)(nil), errors.New("db error"))

	err := service.Create(rack)

	require.Error(t, err)
	require.Equal(t, "failed to check rack code", err.Error())
	mockRackRepo.AssertExpectations(t)
}

// TestRackService_GetAllRacks_Success tests getting all racks
func TestRackService_GetAllRacks_Success(t *testing.T) {
	mockRackRepo := new(MockRackRepository)
	repo := repository.Repository{RackRepo: mockRackRepo}
	service := NewRackService(repo)

	racks := []model.Rack{
		{ID: 1, Code: "A1", WarehouseID: 1},
		{ID: 2, Code: "A2", WarehouseID: 1},
	}

	mockRackRepo.On("FindAll", 1, 10).Return(racks, 2, nil)

	result, pagination, err := service.GetAllRacks(1, 10)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, 2, len(*result))
	require.NotNil(t, pagination)
	require.Equal(t, 1, pagination.CurrentPage)
	require.Equal(t, 10, pagination.Limit)
	mockRackRepo.AssertExpectations(t)
}

// TestRackService_GetAllRacks_Error tests error handling
func TestRackService_GetAllRacks_Error(t *testing.T) {
	mockRackRepo := new(MockRackRepository)
	repo := repository.Repository{RackRepo: mockRackRepo}
	service := NewRackService(repo)

	mockRackRepo.On("FindAll", 1, 10).Return([]model.Rack{}, 0, errors.New("db error"))

	result, pagination, err := service.GetAllRacks(1, 10)

	require.Error(t, err)
	require.Nil(t, result)
	require.Nil(t, pagination)
	mockRackRepo.AssertExpectations(t)
}

// TestRackService_GetRacksByWarehouse_Success tests getting racks by warehouse
func TestRackService_GetRacksByWarehouse_Success(t *testing.T) {
	mockRackRepo := new(MockRackRepository)
	repo := repository.Repository{RackRepo: mockRackRepo}
	service := NewRackService(repo)

	racks := []model.Rack{
		{ID: 1, Code: "A1", WarehouseID: 1},
	}

	mockRackRepo.On("FindByWarehouseID", 1, 1, 10).Return(racks, 1, nil)

	result, pagination, err := service.GetRacksByWarehouse(1, 1, 10)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, 1, len(*result))
	require.NotNil(t, pagination)
	mockRackRepo.AssertExpectations(t)
}

// TestRackService_GetRackByID_Success tests getting rack by ID
func TestRackService_GetRackByID_Success(t *testing.T) {
	mockRackRepo := new(MockRackRepository)
	repo := repository.Repository{RackRepo: mockRackRepo}
	service := NewRackService(repo)

	rack := &model.Rack{ID: 1, Code: "A1", WarehouseID: 1}

	mockRackRepo.On("FindByID", 1).Return(rack, nil)

	result, err := service.GetRackByID(1)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "A1", result.Code)
	mockRackRepo.AssertExpectations(t)
}

// TestRackService_GetRackByID_NotFound tests rack not found
func TestRackService_GetRackByID_NotFound(t *testing.T) {
	mockRackRepo := new(MockRackRepository)
	repo := repository.Repository{RackRepo: mockRackRepo}
	service := NewRackService(repo)

	mockRackRepo.On("FindByID", 999).Return((*model.Rack)(nil), nil)

	result, err := service.GetRackByID(999)

	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, "rack not found", err.Error())
	mockRackRepo.AssertExpectations(t)
}

// TestRackService_Update_Success tests successful update
func TestRackService_Update_Success(t *testing.T) {
	mockRackRepo := new(MockRackRepository)
	repo := repository.Repository{RackRepo: mockRackRepo}
	service := NewRackService(repo)

	existingRack := &model.Rack{
		ID:          1,
		WarehouseID: 1,
		Code:        "A1",
	}

	updateData := &model.Rack{
		Code: "A2",
	}

	mockRackRepo.On("FindByID", 1).Return(existingRack, nil)
	mockRackRepo.On("FindByWarehouseAndCode", 1, "A2").Return((*model.Rack)(nil), nil)
	mockRackRepo.On("Update", 1, updateData).Return(nil)

	err := service.Update(1, updateData)

	require.NoError(t, err)
	require.Equal(t, 1, updateData.WarehouseID) // Should keep existing WarehouseID
	mockRackRepo.AssertExpectations(t)
}

// TestRackService_Update_CodeExists tests update with existing code
func TestRackService_Update_CodeExists(t *testing.T) {
	mockRackRepo := new(MockRackRepository)
	repo := repository.Repository{RackRepo: mockRackRepo}
	service := NewRackService(repo)

	existingRack := &model.Rack{ID: 1, WarehouseID: 1, Code: "A1"}
	otherRack := &model.Rack{ID: 2, WarehouseID: 1, Code: "A2"}
	updateData := &model.Rack{Code: "A2"}

	mockRackRepo.On("FindByID", 1).Return(existingRack, nil)
	mockRackRepo.On("FindByWarehouseAndCode", 1, "A2").Return(otherRack, nil)

	err := service.Update(1, updateData)

	require.Error(t, err)
	require.Equal(t, "rack code already exists in this warehouse", err.Error())
	mockRackRepo.AssertExpectations(t)
}

// TestRackService_Update_NotFound tests update with non-existent rack
func TestRackService_Update_NotFound(t *testing.T) {
	mockRackRepo := new(MockRackRepository)
	repo := repository.Repository{RackRepo: mockRackRepo}
	service := NewRackService(repo)

	updateData := &model.Rack{Code: "A1"}

	mockRackRepo.On("FindByID", 999).Return((*model.Rack)(nil), nil)

	err := service.Update(999, updateData)

	require.Error(t, err)
	require.Equal(t, "rack not found", err.Error())
	mockRackRepo.AssertExpectations(t)
}

// TestRackService_Delete_Success tests successful deletion
func TestRackService_Delete_Success(t *testing.T) {
	mockRackRepo := new(MockRackRepository)
	repo := repository.Repository{RackRepo: mockRackRepo}
	service := NewRackService(repo)

	existingRack := &model.Rack{ID: 1, Code: "A1"}

	mockRackRepo.On("FindByID", 1).Return(existingRack, nil)
	mockRackRepo.On("Delete", 1).Return(nil)

	err := service.Delete(1)

	require.NoError(t, err)
	mockRackRepo.AssertExpectations(t)
}

// TestRackService_Delete_NotFound tests deletion with non-existent rack
func TestRackService_Delete_NotFound(t *testing.T) {
	mockRackRepo := new(MockRackRepository)
	repo := repository.Repository{RackRepo: mockRackRepo}
	service := NewRackService(repo)

	mockRackRepo.On("FindByID", 999).Return((*model.Rack)(nil), nil)

	err := service.Delete(999)

	require.Error(t, err)
	require.Equal(t, "rack not found", err.Error())
	mockRackRepo.AssertExpectations(t)
}
