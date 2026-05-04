package model

import "time"

type Inventory struct {
	Id               int       `json:"id"`
	ProductVariantId int       `json:"product_variant_id"`
	LocationId       int       `json:"location_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	Quantity         int       `json:"quantity"`
}

type CreateInventoryRequest struct {
	ProductVariantId int `json:"product_variant_id" validate:"required"`
	LocationId       int `json:"location_id" validate:"required"`
	Quantity         int `json:"quantity" validate:"required"`
}

type UpdateInventoryRequest struct {
	Quantity   *int `json:"quantity"`
	LocationId *int `json:"location_id"`
}
