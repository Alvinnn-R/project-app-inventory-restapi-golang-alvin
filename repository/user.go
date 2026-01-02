package repository

import (
	"context"
	"errors"
	"project-app-inventory/database"
	"project-app-inventory/model"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(id int) (*model.User, error)
	FindAll(page, limit int) ([]model.User, int, error)
	Update(id int, data *model.User) error
	Delete(id int) error
	FindAllStudents() ([]model.User, error)
	GetUserByID(id int) (model.User, error)
}

type userRepositoryImpl struct {
	db     database.PgxIface
	Logger *zap.Logger
}

func NewUserRepository(db database.PgxIface) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(user *model.User) error {
	query := `
		INSERT INTO users (name, email, password_hash, role_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(context.Background(), query,
		user.Name, user.Email, user.PasswordHash, user.RoleID, user.IsActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil && r.Logger != nil {
		r.Logger.Error("error creating user", zap.Error(err))
	}
	return err
}

func (r *userRepositoryImpl) FindByEmail(email string) (*model.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.password_hash, u.role_id, r.name as role_name, u.is_active, u.created_at, u.updated_at
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.email = $1
	`
	var user model.User
	err := r.db.QueryRow(context.Background(), query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.RoleID, &user.RoleName, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil && r.Logger != nil {
		r.Logger.Error("error finding user by email", zap.Error(err))
		return nil, err
	}

	return &user, nil
}

func (r *userRepositoryImpl) FindByID(id int) (*model.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.password_hash, u.role_id, r.name as role_name, u.is_active, u.created_at, u.updated_at
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1
	`
	var user model.User
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.RoleID, &user.RoleName, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil && r.Logger != nil {
		r.Logger.Error("error finding user by id", zap.Error(err))
		return nil, err
	}

	return &user, nil
}

func (r *userRepositoryImpl) FindAll(page, limit int) ([]model.User, int, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM users`
	err := r.db.QueryRow(context.Background(), countQuery).Scan(&total)
	if err != nil {
		if r.Logger != nil {
			r.Logger.Error("error counting users", zap.Error(err))
		}
		return nil, 0, err
	}

	// Get data with pagination
	query := `
		SELECT u.id, u.name, u.email, u.password_hash, u.role_id, r.name as role_name, u.is_active, u.created_at, u.updated_at
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		ORDER BY u.name ASC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		if r.Logger != nil {
			r.Logger.Error("error querying users", zap.Error(err))
		}
		return nil, 0, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.PasswordHash,
			&user.RoleID, &user.RoleName, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			if r.Logger != nil {
				r.Logger.Error("error scanning user", zap.Error(err))
			}
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}

func (r *userRepositoryImpl) Update(id int, data *model.User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, password_hash = $3, role_id = $4, is_active = $5, updated_at = NOW()
		WHERE id = $6
	`
	result, err := r.db.Exec(context.Background(), query,
		data.Name, data.Email, data.PasswordHash, data.RoleID, data.IsActive, id,
	)
	if err != nil {
		if r.Logger != nil {
			r.Logger.Error("error updating user", zap.Error(err))
		}
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *userRepositoryImpl) Delete(id int) error {
	query := `
		DELETE FROM users 
		WHERE id = $1
	`
	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		if r.Logger != nil {
			r.Logger.Error("error deleting user", zap.Error(err))
		}
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
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
