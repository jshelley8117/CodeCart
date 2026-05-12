package model

import (
	"encoding/json"
	"time"
)

type Location struct {
	ID           int             `json:"id"`
	StoreName    string          `json:"store_name"`
	StoreAddress json.RawMessage `json:"store_address"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type CreateLocationRequest struct {
	StoreName    string          `json:"store_name" validate:"required"`
	StoreAddress json.RawMessage `json:"store_address" validate:"required"`
}

type UpdateLocationRequest struct {
	StoreName    string          `json:"store_name"`
	StoreAddress json.RawMessage `json:"store_address"`
}
