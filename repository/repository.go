package repository

import (
	"project-app-inventory/database"

	"go.uber.org/zap"
)

type Repository struct {
	AssignmentRepo       AssignmentRepository
	SubmissionRepo       SubmissionRepo
	UserRepo             UserRepository
	PermissionRepository PermissionIface
}

func NewRepository(db database.PgxIface, log *zap.Logger) Repository {
	return Repository{
		AssignmentRepo:       NewAssignmentRepository(db, log),
		SubmissionRepo:       NewSubmissionRepo(db),
		UserRepo:             NewUserRepository(db),
		PermissionRepository: NewPermissionRepository(db),
	}
}
