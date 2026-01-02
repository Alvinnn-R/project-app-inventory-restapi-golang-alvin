package model

import "time"

type Sale struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	TotalAmount float64    `json:"total_amount"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type SaleItem struct {
	ID          int     `json:"id"`
	SaleID      int     `json:"sale_id"`
	ItemID      int     `json:"item_id"`
	Quantity    int     `json:"quantity"`
	PriceAtSale float64 `json:"price_at_sale"`
	Subtotal    float64 `json:"subtotal"`
}
