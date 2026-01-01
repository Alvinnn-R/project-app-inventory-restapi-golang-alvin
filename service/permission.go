package service

import (
	"project-app-inventory/repository"
)

type PermissionIface interface {
	Allowed(userID int, code string) (bool, error)
}

type permissionService struct {
	Repo repository.Repository
}

func NewPermissionService(repo repository.Repository) *permissionService {
	return &permissionService{Repo: repo}
}

func (permissionService *permissionService) Allowed(userID int, code string) (bool, error) {
	allowed, err := permissionService.Repo.PermissionRepository.Allowed(userID, code)
	if err != nil {
		return false, err
	}

	return allowed, nil
}
