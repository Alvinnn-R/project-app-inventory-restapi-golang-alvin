package service

import (
	"errors"
	"project-app-inventory/dto"
	"project-app-inventory/model"
	"project-app-inventory/repository"
	"project-app-inventory/utils"
)

type UserService interface {
	Create(user *model.User) error
	GetAllUsers(page, limit int) (*[]model.User, *dto.Pagination, error)
	GetUserByID(id int) (model.User, error)
	GetUserByIDDetailed(id int) (*model.User, error)
	Update(id int, data *model.User) error
	Delete(id int) error
}

type userService struct {
	Repo repository.Repository
}

func NewUserService(repo repository.Repository) UserService {
	return &userService{Repo: repo}
}

func (s *userService) Create(user *model.User) error {
	// Check if email already exists
	existingUser, err := s.Repo.UserRepo.FindByEmail(user.Email)
	if err != nil {
		return errors.New("failed to check email")
	}
	if existingUser != nil {
		return errors.New("email already exists")
	}

	return s.Repo.UserRepo.Create(user)
}

func (s *userService) GetAllUsers(page, limit int) (*[]model.User, *dto.Pagination, error) {
	users, total, err := s.Repo.UserRepo.FindAll(page, limit)
	if err != nil {
		return nil, nil, err
	}

	pagination := dto.Pagination{
		CurrentPage:  page,
		Limit:        limit,
		TotalPages:   utils.TotalPage(limit, int64(total)),
		TotalRecords: total,
	}
	return &users, &pagination, nil
}

func (s *userService) GetUserByID(id int) (model.User, error) {
	return s.Repo.UserRepo.GetUserByID(id)
}

func (s *userService) GetUserByIDDetailed(id int) (*model.User, error) {
	user, err := s.Repo.UserRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *userService) Update(id int, data *model.User) error {
	// Check if user exists
	existingUser, err := s.Repo.UserRepo.FindByID(id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	// If fields are empty, keep existing values
	if data.Name == "" {
		data.Name = existingUser.Name
	}
	if data.Email == "" {
		data.Email = existingUser.Email
	}
	if data.PasswordHash == "" {
		data.PasswordHash = existingUser.PasswordHash
	}
	if data.RoleID == 0 {
		data.RoleID = existingUser.RoleID
	}

	// Check if email is being changed and if new email already exists
	if data.Email != existingUser.Email {
		emailExists, err := s.Repo.UserRepo.FindByEmail(data.Email)
		if err != nil {
			return errors.New("failed to check email")
		}
		if emailExists != nil {
			return errors.New("email already exists")
		}
	}

	return s.Repo.UserRepo.Update(id, data)
}

func (s *userService) Delete(id int) error {
	// Check if user exists
	existingUser, err := s.Repo.UserRepo.FindByID(id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	return s.Repo.UserRepo.Delete(id)
}
