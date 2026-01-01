package model

import "time"

type Item struct {
	ID           int       `json:"id"`
	SKU          string    `json:"sku"`
	Name         string    `json:"name"`
	CategoryID   int       `json:"category_id"`
	RackID       int       `json:"rack_id"`
	Stock        int       `json:"stock"`
	MinimumStock int       `json:"minimum_stock"`
	Price        float64   `json:"price"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
