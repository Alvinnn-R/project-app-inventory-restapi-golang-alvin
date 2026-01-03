package service

import (
	"errors"
	"project-app-inventory/model"
	"project-app-inventory/repository"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockUserRepository mocks UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id int) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(id int) (model.User, error) {
	args := m.Called(id)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindAll(page, limit int) ([]model.User, int, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]model.User), args.Int(1), args.Error(2)
}

func (m *MockUserRepository) Update(id int, user *model.User) error {
	args := m.Called(id, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) FindAllStudents() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

// TestUserService_Create_Success tests successful user creation
func TestUserService_Create_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	repo := repository.Repository{UserRepo: mockUserRepo}
	service := NewUserService(repo)

	user := &model.User{
		Name:         "John Doe",
		Email:        "john@example.com",
		PasswordHash: "hashedpassword",
		RoleID:       2,
	}

	mockUserRepo.On("FindByEmail", user.Email).Return((*model.User)(nil), nil)
	mockUserRepo.On("Create", user).Return(nil)

	err := service.Create(user)

	require.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

// TestUserService_Create_EmailExists tests creation with existing email
func TestUserService_Create_EmailExists(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	repo := repository.Repository{UserRepo: mockUserRepo}
	service := NewUserService(repo)

	user := &model.User{Email: "john@example.com"}
	existingUser := &model.User{ID: 1, Email: "john@example.com"}

	mockUserRepo.On("FindByEmail", user.Email).Return(existingUser, nil)

	err := service.Create(user)

	require.Error(t, err)
	require.Equal(t, "email already exists", err.Error())
	mockUserRepo.AssertExpectations(t)
}

// TestUserService_Create_CheckError tests creation when check fails
func TestUserService_Create_CheckError(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	repo := repository.Repository{UserRepo: mockUserRepo}
	service := NewUserService(repo)

	user := &model.User{Email: "john@example.com"}

	mockUserRepo.On("FindByEmail", user.Email).Return((*model.User)(nil), errors.New("db error"))

	err := service.Create(user)

	require.Error(t, err)
	require.Equal(t, "failed to check email", err.Error())
	mockUserRepo.AssertExpectations(t)
}

// TestUserService_GetAllUsers_Success tests getting all users
func TestUserService_GetAllUsers_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	repo := repository.Repository{UserRepo: mockUserRepo}
	service := NewUserService(repo)

	users := []model.User{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
		{ID: 2, Name: "Jane Doe", Email: "jane@example.com"},
	}

	mockUserRepo.On("FindAll", 1, 10).Return(users, 2, nil)

	result, pagination, err := service.GetAllUsers(1, 10)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, 2, len(*result))
	require.NotNil(t, pagination)
	require.Equal(t, 1, pagination.CurrentPage)
	require.Equal(t, 10, pagination.Limit)
	mockUserRepo.AssertExpectations(t)
}

// TestUserService_GetAllUsers_Error tests error handling
func TestUserService_GetAllUsers_Error(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	repo := repository.Repository{UserRepo: mockUserRepo}
	service := NewUserService(repo)

	mockUserRepo.On("FindAll", 1, 10).Return([]model.User{}, 0, errors.New("db error"))

	result, pagination, err := service.GetAllUsers(1, 10)

	require.Error(t, err)
	require.Nil(t, result)
	require.Nil(t, pagination)
	mockUserRepo.AssertExpectations(t)
}

// TestUserService_GetUserByID_Success tests getting user by ID
func TestUserService_GetUserByID_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	repo := repository.Repository{UserRepo: mockUserRepo}
	service := NewUserService(repo)

	user := model.User{ID: 1, Name: "John Doe"}

	mockUserRepo.On("GetUserByID", 1).Return(user, nil)

	result, err := service.GetUserByID(1)

	require.NoError(t, err)
	require.Equal(t, "John Doe", result.Name)
	mockUserRepo.AssertExpectations(t)
}

// TestUserService_GetUserByIDDetailed_Success tests getting detailed user by ID
func TestUserService_GetUserByIDDetailed_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	repo := repository.Repository{UserRepo: mockUserRepo}
	service := NewUserService(repo)

	user := &model.User{ID: 1, Name: "John Doe"}

	mockUserRepo.On("FindByID", 1).Return(user, nil)

	result, err := service.GetUserByIDDetailed(1)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "John Doe", result.Name)
	mockUserRepo.AssertExpectations(t)
}

// TestUserService_GetUserByIDDetailed_NotFound tests user not found
func TestUserService_GetUserByIDDetailed_NotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	repo := repository.Repository{UserRepo: mockUserRepo}
	service := NewUserService(repo)

	mockUserRepo.On("FindByID", 999).Return((*model.User)(nil), nil)

	result, err := service.GetUserByIDDetailed(999)

	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, "user not found", err.Error())
	mockUserRepo.AssertExpectations(t)
}

// TestUserService_Update_Success tests successful update
func TestUserService_Update_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	repo := repository.Repository{UserRepo: mockUserRepo}
	service := NewUserService(repo)

	existingUser := &model.User{
		ID:           1,
		Name:         "John Doe",
		Email:        "john@example.com",
		PasswordHash: "oldpassword",
		RoleID:       2,
	}

	updateData := &model.User{
		Name:  "John Updated",
		Email: "john.updated@example.com",
	}

	mockUserRepo.On("FindByID", 1).Return(existingUser, nil)
	mockUserRepo.On("FindByEmail", "john.updated@example.com").Return((*model.User)(nil), nil)
	mockUserRepo.On("Update", 1, updateData).Return(nil)

	err := service.Update(1, updateData)

	require.NoError(t, err)
	require.Equal(t, "oldpassword", updateData.PasswordHash) // Should keep existing password
	require.Equal(t, 2, updateData.RoleID)                   // Should keep existing role
	mockUserRepo.AssertExpectations(t)
}

// TestUserService_Update_EmailExists tests update with existing email
func TestUserService_Update_EmailExists(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	repo := repository.Repository{UserRepo: mockUserRepo}
	service := NewUserService(repo)

	existingUser := &model.User{ID: 1, Email: "john@example.com"}
	otherUser := &model.User{ID: 2, Email: "jane@example.com"}
	updateData := &model.User{Email: "jane@example.com"}

	mockUserRepo.On("FindByID", 1).Return(existingUser, nil)
	mockUserRepo.On("FindByEmail", "jane@example.com").Return(otherUser, nil)

	err := service.Update(1, updateData)

	require.Error(t, err)
	require.Equal(t, "email already exists", err.Error())
	mockUserRepo.AssertExpectations(t)
}

// TestUserService_Update_NotFound tests update with non-existent user
func TestUserService_Update_NotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	repo := repository.Repository{UserRepo: mockUserRepo}
	service := NewUserService(repo)

	updateData := &model.User{Name: "New Name"}

	mockUserRepo.On("FindByID", 999).Return((*model.User)(nil), nil)

	err := service.Update(999, updateData)

	require.Error(t, err)
	require.Equal(t, "user not found", err.Error())
	mockUserRepo.AssertExpectations(t)
}

// TestUserService_Delete_Success tests successful deletion
func TestUserService_Delete_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	repo := repository.Repository{UserRepo: mockUserRepo}
	service := NewUserService(repo)

	existingUser := &model.User{ID: 1, Name: "John Doe"}

	mockUserRepo.On("FindByID", 1).Return(existingUser, nil)
	mockUserRepo.On("Delete", 1).Return(nil)

	err := service.Delete(1)

	require.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

// TestUserService_Delete_NotFound tests deletion with non-existent user
func TestUserService_Delete_NotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	repo := repository.Repository{UserRepo: mockUserRepo}
	service := NewUserService(repo)

	mockUserRepo.On("FindByID", 999).Return((*model.User)(nil), nil)

	err := service.Delete(999)

	require.Error(t, err)
	require.Equal(t, "user not found", err.Error())
	mockUserRepo.AssertExpectations(t)
}
