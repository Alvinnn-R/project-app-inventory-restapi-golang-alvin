package dto

type SaleItemRequest struct {
	ItemID   int `json:"item_id" validate:"required,gt=0"`
	Quantity int `json:"quantity" validate:"required,gt=0"`
}

type SaleRequest struct {
	Items []SaleItemRequest `json:"items" validate:"required,min=1,dive"`
}

type SaleItemResponse struct {
	ID          int     `json:"id"`
	ItemID      int     `json:"item_id"`
	ItemName    string  `json:"item_name,omitempty"`
	Quantity    int     `json:"quantity"`
	PriceAtSale float64 `json:"price_at_sale"`
	Subtotal    float64 `json:"subtotal"`
}

type SaleResponse struct {
	ID          int                `json:"id"`
	UserID      int                `json:"user_id"`
	UserName    string             `json:"user_name,omitempty"`
	TotalAmount float64            `json:"total_amount"`
	Items       []SaleItemResponse `json:"items,omitempty"`
	CreatedAt   string             `json:"created_at"`
	UpdatedAt   string             `json:"updated_at"`
	DeletedAt   *string            `json:"deleted_at,omitempty"`
}
