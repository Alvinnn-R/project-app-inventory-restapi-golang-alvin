package service

import (
	"errors"
	"project-app-inventory/dto"
	"project-app-inventory/model"
	"project-app-inventory/repository"
	"project-app-inventory/utils"
	"time"
)

type AuthService interface {
	Login(email, password string) (*dto.LoginResponse, error)
	Logout(token string) error
	ValidateToken(token string) (*model.User, error)
}

type authService struct {
	Repo repository.Repository
}

func NewAuthService(repo repository.Repository) AuthService {
	return &authService{Repo: repo}
}

func (s *authService) Login(email, password string) (*dto.LoginResponse, error) {
	// Find user by email
	user, err := s.Repo.UserRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("failed to find user")
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Check password
	if !utils.CheckPassword(password, user.PasswordHash) {
		return nil, errors.New("incorrect password")
	}

	// Generate UUID token
	token := utils.GenerateUUIDToken()

	// Create session with 24 hours expiration
	session := &model.Session{
		UserID:    user.ID,
		Token:     token,
		ExpiredAt: time.Now().Add(24 * time.Hour),
	}

	err = s.Repo.SessionRepo.Create(session)
	if err != nil {
		return nil, errors.New("failed to create session")
	}

	// Prepare response
	response := &dto.LoginResponse{
		Token:     token,
		ExpiredAt: session.ExpiredAt,
		User: dto.UserInfo{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			RoleID:   user.RoleID,
			RoleName: user.RoleName,
			IsActive: user.IsActive,
		},
	}

	return response, nil
}

func (s *authService) Logout(token string) error {
	// Revoke session
	err := s.Repo.SessionRepo.RevokeByToken(token)
	if err != nil {
		return errors.New("failed to revoke session")
	}
	return nil
}

func (s *authService) ValidateToken(token string) (*model.User, error) {
	// Find session by token
	session, err := s.Repo.SessionRepo.FindByToken(token)
	if err != nil {
		return nil, errors.New("failed to validate token")
	}
	if session == nil {
		return nil, errors.New("invalid or expired token")
	}

	// Find user by ID
	user, err := s.Repo.UserRepo.FindByID(session.UserID)
	if err != nil {
		return nil, errors.New("failed to find user")
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
