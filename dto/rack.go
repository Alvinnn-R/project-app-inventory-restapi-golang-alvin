package dto

type RackRequest struct {
	WarehouseID int    `json:"warehouse_id" validate:"required,gt=0"`
	Code        string `json:"code" validate:"required,min=2,max=50"`
	Description string `json:"description" validate:"omitempty,max=500"`
}

type RackUpdateRequest struct {
	WarehouseID int    `json:"warehouse_id" validate:"omitempty,gt=0"`
	Code        string `json:"code" validate:"omitempty,min=2,max=50"`
	Description string `json:"description" validate:"omitempty,max=500"`
}

type RackResponse struct {
	ID          int    `json:"id"`
	WarehouseID int    `json:"warehouse_id"`
	Code        string `json:"code"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
