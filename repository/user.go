package repository

import (
	"context"
	"database/sql"
	"project-app-inventory/database"
	"project-app-inventory/model"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(id int) (*model.User, error)
	FindAllStudents() ([]model.User, error)
	GetUserByID(id int) (model.User, error)
}

type userRepositoryImpl struct {
	db database.PgxIface
}

func NewUserRepository(db database.PgxIface) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(user *model.User) error {
	query := `
		INSERT INTO users (name, email, password_hash, role_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id
	`
	return r.db.QueryRow(context.Background(), query, user.Name, user.Email, user.PasswordHash, user.RoleID, user.IsActive).Scan(&user.ID)
}

func (r *userRepositoryImpl) FindByEmail(email string) (*model.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.password_hash, u.role_id, r.name as role_name, u.is_active, u.created_at, u.updated_at
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.email = $1 AND u.is_active = true
	`
	var user model.User
	err := r.db.QueryRow(context.Background(), query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.RoleID, &user.RoleName, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // user tidak ditemukan
	}

	return &user, err
}

func (r *userRepositoryImpl) FindByID(id int) (*model.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.password_hash, u.role_id, r.name as role_name, u.is_active, u.created_at, u.updated_at
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1 AND u.is_active = true
	`
	var user model.User
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.RoleID, &user.RoleName, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, err
}

func (r *userRepositoryImpl) FindAllStudents() ([]model.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.password_hash, u.role_id, r.name as role_name, u.is_active, u.created_at, u.updated_at
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE r.name = 'staff' AND u.is_active = true
	`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []model.User
	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.RoleID, &u.RoleName, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		students = append(students, u)
	}
	return students, nil
}

func (r *userRepositoryImpl) GetUserByID(id int) (model.User, error) {
	var user model.User
	query := `
		SELECT u.id, u.name, u.email, u.password_hash, u.role_id, r.name as role_name, u.is_active, u.created_at, u.updated_at
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1 AND u.is_active = true
	`

	err := r.db.QueryRow(context.Background(), query, id).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.RoleID, &user.RoleName, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}

	return user, nil
}
