package model

import "time"

type Rack struct {
	ID          int       `json:"id"`
	WarehouseID int       `json:"warehouse_id"`
	Code        string    `json:"code"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
