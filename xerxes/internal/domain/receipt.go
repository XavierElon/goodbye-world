package domain

import (
	"time"
)

// ReceiptItem represents a single item on a receipt
type ReceiptItem struct {
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Total    float64 `json:"total"`
}

// Receipt represents a complete receipt
type Receipt struct {
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	StoreID   string            `json:"store_id"`
	Items     []ReceiptItem     `json:"items"`
	Subtotal  float64           `json:"subtotal"`
	Tax       float64           `json:"tax"`
	Total     float64           `json:"total"`
	CreatedAt time.Time         `json:"created_at"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// ReceiptCreation represents the data needed to create a receipt
type ReceiptCreation struct {
	PhoneNumber string            `json:"phone_number" validate:"required"`
	StoreID     string            `json:"store_id" validate:"required"`
	Items       []ReceiptItem     `json:"items" validate:"required"`
	Subtotal    float64           `json:"subtotal" validate:"required"`
	Tax         float64           `json:"tax" validate:"required"`
	Total       float64           `json:"total" validate:"required"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// Store represents a store where receipts are created
type Store struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}