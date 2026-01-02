package dto

type UserRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	RoleID   int    `json:"role_id" validate:"required,gt=0"`
	IsActive bool   `json:"is_active"`
}

type UserUpdateRequest struct {
	Name     string `json:"name" validate:"omitempty,min=3,max=100"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,min=6"`
	RoleID   int    `json:"role_id" validate:"omitempty,gt=0"`
	IsActive *bool  `json:"is_active" validate:"omitempty"`
}

type UserResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	RoleID    int    `json:"role_id"`
	RoleName  string `json:"role_name,omitempty"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
