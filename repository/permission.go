package repository

import (
	"context"
	"project-app-inventory/database"
)

type PermissionIface interface {
	Allowed(userID int, code string) (bool, error)
}

type permissionRepository struct {
	db database.PgxIface
}

func NewPermissionRepository(db database.PgxIface) *permissionRepository {
	return &permissionRepository{db: db}
}

func (permissionRepository *permissionRepository) Allowed(userID int, code string) (bool, error) {
	const qAllowed = `
	WITH perm AS (
	SELECT id FROM permissions WHERE code = $2
	)
	SELECT
	CASE
		WHEN EXISTS (
		SELECT 1 FROM user_permissions up, perm
		WHERE up.user_id = $1 AND up.permission_id = perm.id AND up.effect='deny'
		) THEN FALSE
		WHEN EXISTS (
		SELECT 1 FROM user_permissions up, perm
		WHERE up.user_id = $1 AND up.permission_id = perm.id AND up.effect='allow'
		) THEN TRUE
		WHEN EXISTS (
		SELECT 1
		FROM users u
		JOIN role_permissions rp ON rp.role_id = u.role_id
		JOIN perm ON perm.id = rp.permission_id
		WHERE u.id = $1
		) THEN TRUE
		ELSE FALSE
	END AS allowed;
	`
	var allowed bool
	err := permissionRepository.db.QueryRow(context.Background(), qAllowed, userID, code).Scan(&allowed)
	return allowed, err
}
