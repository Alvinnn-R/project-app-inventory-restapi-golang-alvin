package repository

import (
	"project-app-inventory/database"

	"go.uber.org/zap"
)

type Repository struct {
	AssignmentRepo       AssignmentRepository
	SubmissionRepo       SubmissionRepo
	UserRepo             UserRepository
	SessionRepo          SessionRepository
	PermissionRepository PermissionIface
	ItemRepo             ItemRepository
	CategoryRepo         CategoryRepository
	RackRepo             RackRepository
}

func NewRepository(db database.PgxIface, log *zap.Logger) Repository {
	return Repository{
		AssignmentRepo:       NewAssignmentRepository(db, log),
		SubmissionRepo:       NewSubmissionRepo(db),
		UserRepo:             NewUserRepository(db),
		SessionRepo:          NewSessionRepository(db),
		PermissionRepository: NewPermissionRepository(db),
		ItemRepo:             NewItemRepository(db, log),
		CategoryRepo:         NewCategoryRepository(db, log),
		RackRepo:             NewRackRepository(db, log),
	}
}
