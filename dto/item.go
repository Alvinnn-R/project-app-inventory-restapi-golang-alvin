package dto

type ItemRequest struct {
	SKU          string  `json:"sku" validate:"required,min=3,max=50"`
	Name         string  `json:"name" validate:"required,min=3,max=150"`
	CategoryID   int     `json:"category_id" validate:"required,gt=0"`
	RackID       int     `json:"rack_id" validate:"required,gt=0"`
	Stock        int     `json:"stock" validate:"required,gte=0"`
	MinimumStock int     `json:"minimum_stock" validate:"required,gte=0"`
	Price        float64 `json:"price" validate:"required,gt=0"`
}

type ItemUpdateRequest struct {
	SKU          string  `json:"sku" validate:"omitempty,min=3,max=50"`
	Name         string  `json:"name" validate:"omitempty,min=3,max=150"`
	CategoryID   int     `json:"category_id" validate:"omitempty,gt=0"`
	RackID       int     `json:"rack_id" validate:"omitempty,gt=0"`
	Stock        int     `json:"stock" validate:"omitempty,gte=0"`
	MinimumStock int     `json:"minimum_stock" validate:"omitempty,gte=0"`
	Price        float64 `json:"price" validate:"omitempty,gt=0"`
}

type ItemResponse struct {
	ID           int     `json:"id"`
	SKU          string  `json:"sku"`
	Name         string  `json:"name"`
	CategoryID   int     `json:"category_id"`
	RackID       int     `json:"rack_id"`
	Stock        int     `json:"stock"`
	MinimumStock int     `json:"minimum_stock"`
	Price        float64 `json:"price"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}
