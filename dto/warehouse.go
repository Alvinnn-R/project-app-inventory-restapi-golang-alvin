package dto

type WarehouseRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Location string `json:"location" validate:"required,min=5,max=255"`
}

type WarehouseUpdateRequest struct {
	Name     string `json:"name" validate:"omitempty,min=3,max=100"`
	Location string `json:"location" validate:"omitempty,min=5,max=255"`
}

type WarehouseResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Location  string `json:"location"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
