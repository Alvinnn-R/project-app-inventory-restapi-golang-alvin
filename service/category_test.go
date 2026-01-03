package service

import (
	"errors"
	"project-app-inventory/model"
	"project-app-inventory/repository"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockCategoryRepository mocks CategoryRepository interface
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(category *model.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) FindByID(id int) (*model.Category, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Category), args.Error(1)
}

func (m *MockCategoryRepository) FindByName(name string) (*model.Category, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Category), args.Error(1)
}

func (m *MockCategoryRepository) FindAll(page, limit int) ([]model.Category, int, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]model.Category), args.Int(1), args.Error(2)
}

func (m *MockCategoryRepository) Update(id int, category *model.Category) error {
	args := m.Called(id, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// TestCategoryService_Create_Success tests successful category creation
func TestCategoryService_Create_Success(t *testing.T) {
	mockCategoryRepo := new(MockCategoryRepository)
	repo := repository.Repository{CategoryRepo: mockCategoryRepo}
	service := NewCategoryService(repo)

	category := &model.Category{
		Name: "Electronics",
	}

	mockCategoryRepo.On("FindByName", category.Name).Return((*model.Category)(nil), nil)
	mockCategoryRepo.On("Create", category).Return(nil)

	err := service.Create(category)

	require.NoError(t, err)
	mockCategoryRepo.AssertExpectations(t)
}

// TestCategoryService_Create_NameExists tests creation with existing name
func TestCategoryService_Create_NameExists(t *testing.T) {
	mockCategoryRepo := new(MockCategoryRepository)
	repo := repository.Repository{CategoryRepo: mockCategoryRepo}
	service := NewCategoryService(repo)

	category := &model.Category{Name: "Electronics"}
	existingCategory := &model.Category{ID: 1, Name: "Electronics"}

	mockCategoryRepo.On("FindByName", category.Name).Return(existingCategory, nil)

	err := service.Create(category)

	require.Error(t, err)
	require.Equal(t, "category name already exists", err.Error())
	mockCategoryRepo.AssertExpectations(t)
}

// TestCategoryService_Create_CheckError tests creation when check fails
func TestCategoryService_Create_CheckError(t *testing.T) {
	mockCategoryRepo := new(MockCategoryRepository)
	repo := repository.Repository{CategoryRepo: mockCategoryRepo}
	service := NewCategoryService(repo)

	category := &model.Category{Name: "Electronics"}

	mockCategoryRepo.On("FindByName", category.Name).Return((*model.Category)(nil), errors.New("db error"))

	err := service.Create(category)

	require.Error(t, err)
	require.Equal(t, "failed to check category name", err.Error())
	mockCategoryRepo.AssertExpectations(t)
}

// TestCategoryService_GetAllCategories_Success tests getting all categories
func TestCategoryService_GetAllCategories_Success(t *testing.T) {
	mockCategoryRepo := new(MockCategoryRepository)
	repo := repository.Repository{CategoryRepo: mockCategoryRepo}
	service := NewCategoryService(repo)

	categories := []model.Category{
		{ID: 1, Name: "Electronics"},
		{ID: 2, Name: "Furniture"},
	}

	mockCategoryRepo.On("FindAll", 1, 10).Return(categories, 2, nil)

	result, pagination, err := service.GetAllCategories(1, 10)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, 2, len(*result))
	require.NotNil(t, pagination)
	require.Equal(t, 1, pagination.CurrentPage)
	require.Equal(t, 10, pagination.Limit)
	require.Equal(t, 2, pagination.TotalRecords)
	mockCategoryRepo.AssertExpectations(t)
}

// TestCategoryService_GetAllCategories_Error tests error handling
func TestCategoryService_GetAllCategories_Error(t *testing.T) {
	mockCategoryRepo := new(MockCategoryRepository)
	repo := repository.Repository{CategoryRepo: mockCategoryRepo}
	service := NewCategoryService(repo)

	mockCategoryRepo.On("FindAll", 1, 10).Return([]model.Category{}, 0, errors.New("db error"))

	result, pagination, err := service.GetAllCategories(1, 10)

	require.Error(t, err)
	require.Nil(t, result)
	require.Nil(t, pagination)
	mockCategoryRepo.AssertExpectations(t)
}

// TestCategoryService_GetCategoryByID_Success tests getting category by ID
func TestCategoryService_GetCategoryByID_Success(t *testing.T) {
	mockCategoryRepo := new(MockCategoryRepository)
	repo := repository.Repository{CategoryRepo: mockCategoryRepo}
	service := NewCategoryService(repo)

	category := &model.Category{ID: 1, Name: "Electronics"}

	mockCategoryRepo.On("FindByID", 1).Return(category, nil)

	result, err := service.GetCategoryByID(1)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "Electronics", result.Name)
	mockCategoryRepo.AssertExpectations(t)
}

// TestCategoryService_GetCategoryByID_NotFound tests category not found
func TestCategoryService_GetCategoryByID_NotFound(t *testing.T) {
	mockCategoryRepo := new(MockCategoryRepository)
	repo := repository.Repository{CategoryRepo: mockCategoryRepo}
	service := NewCategoryService(repo)

	mockCategoryRepo.On("FindByID", 999).Return((*model.Category)(nil), nil)

	result, err := service.GetCategoryByID(999)

	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, "category not found", err.Error())
	mockCategoryRepo.AssertExpectations(t)
}

// TestCategoryService_Update_Success tests successful update
func TestCategoryService_Update_Success(t *testing.T) {
	mockCategoryRepo := new(MockCategoryRepository)
	repo := repository.Repository{CategoryRepo: mockCategoryRepo}
	service := NewCategoryService(repo)

	desc := "Updated description"
	existingCategory := &model.Category{
		ID:   1,
		Name: "Electronics",
	}

	updateData := &model.Category{
		Name:        "Electronics Updated",
		Description: &desc,
	}

	mockCategoryRepo.On("FindByID", 1).Return(existingCategory, nil)
	mockCategoryRepo.On("FindByName", "Electronics Updated").Return((*model.Category)(nil), nil)
	mockCategoryRepo.On("Update", 1, updateData).Return(nil)

	err := service.Update(1, updateData)

	require.NoError(t, err)
	mockCategoryRepo.AssertExpectations(t)
}

// TestCategoryService_Update_NameExists tests update with existing name
func TestCategoryService_Update_NameExists(t *testing.T) {
	mockCategoryRepo := new(MockCategoryRepository)
	repo := repository.Repository{CategoryRepo: mockCategoryRepo}
	service := NewCategoryService(repo)

	existingCategory := &model.Category{ID: 1, Name: "Electronics"}
	otherCategory := &model.Category{ID: 2, Name: "Furniture"}
	updateData := &model.Category{Name: "Furniture"}

	mockCategoryRepo.On("FindByID", 1).Return(existingCategory, nil)
	mockCategoryRepo.On("FindByName", "Furniture").Return(otherCategory, nil)

	err := service.Update(1, updateData)

	require.Error(t, err)
	require.Equal(t, "category name already exists", err.Error())
	mockCategoryRepo.AssertExpectations(t)
}

// TestCategoryService_Update_NotFound tests update with non-existent category
func TestCategoryService_Update_NotFound(t *testing.T) {
	mockCategoryRepo := new(MockCategoryRepository)
	repo := repository.Repository{CategoryRepo: mockCategoryRepo}
	service := NewCategoryService(repo)

	updateData := &model.Category{Name: "New Name"}

	mockCategoryRepo.On("FindByID", 999).Return((*model.Category)(nil), nil)

	err := service.Update(999, updateData)

	require.Error(t, err)
	require.Equal(t, "category not found", err.Error())
	mockCategoryRepo.AssertExpectations(t)
}

// TestCategoryService_Delete_Success tests successful deletion
func TestCategoryService_Delete_Success(t *testing.T) {
	mockCategoryRepo := new(MockCategoryRepository)
	repo := repository.Repository{CategoryRepo: mockCategoryRepo}
	service := NewCategoryService(repo)

	existingCategory := &model.Category{ID: 1, Name: "Electronics"}

	mockCategoryRepo.On("FindByID", 1).Return(existingCategory, nil)
	mockCategoryRepo.On("Delete", 1).Return(nil)

	err := service.Delete(1)

	require.NoError(t, err)
	mockCategoryRepo.AssertExpectations(t)
}

// TestCategoryService_Delete_NotFound tests deletion with non-existent category
func TestCategoryService_Delete_NotFound(t *testing.T) {
	mockCategoryRepo := new(MockCategoryRepository)
	repo := repository.Repository{CategoryRepo: mockCategoryRepo}
	service := NewCategoryService(repo)

	mockCategoryRepo.On("FindByID", 999).Return((*model.Category)(nil), nil)

	err := service.Delete(999)

	require.Error(t, err)
	require.Equal(t, "category not found", err.Error())
	mockCategoryRepo.AssertExpectations(t)
}
