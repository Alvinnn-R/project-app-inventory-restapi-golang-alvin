package service

import "project-app-inventory/repository"

type Service struct {
	AssignmentService AssignmentService
	SubmissionService SubmissionService
	UserService       UserService
	AuthService       AuthService
	PermissionService PermissionIface
	ItemService       ItemService
	CategoryService   CategoryService
	RackService       RackService
}

func NewService(repo repository.Repository) Service {
	return Service{
		AssignmentService: NewAssignmentService(repo),
		SubmissionService: NewSubmissionService(repo),
		UserService:       NewUserService(repo),
		AuthService:       NewAuthService(repo),
		PermissionService: NewPermissionService(repo),
		ItemService:       NewItemService(repo),
		CategoryService:   NewCategoryService(repo),
		RackService:       NewRackService(repo),
	}
}
