package model

import "time"

type Category string

const (
	Produce         Category = "PRODUCE"
	Meat            Category = "MEAT"
	Seafood         Category = "SEAFOOD"
	Bakery          Category = "BAKERY"
	Dairy           Category = "DAIRY"
	Deli            Category = "DELI"
	Pantry          Category = "PANTRY"
	Frozen          Category = "FROZEN"
	Beverages       Category = "BEVERAGES"
	HomeEssentials  Category = "HOME_ESSENTIALS"
	HealthAndBeauty Category = "HEALTH_AND_BEAUTY"
	Baby            Category = "BABY"
)

type Product struct {
	Id              int              `json:"id"`
	Name            string           `json:"name"`
	Description     string           `json:"description"`
	UnitPrice       float64          `json:"unit_price"`
	Category        Category         `json:"category"`
	Brand           string           `json:"brand"`
	IsAgeRestricted bool             `json:"is_age_restricted"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	IsActive        bool             `json:"is_active"`
	ProductVariants []ProductVariant `json:"product_variants"`
}

type CreateProductRequest struct {
	Name            string   `json:"name" validate:"required"`
	Description     string   `json:"description"`
	UnitPrice       float64  `json:"unit_price" validate:"required"`
	Category        Category `json:"category" validate:"required"`
	Brand           string   `json:"brand" validate:"required"`
	IsAgeRestricted *bool    `json:"is_age_restricted" validate:"required"`
}

type UpdateProductRequest struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	UnitPrice       float64  `json:"unit_price"`
	Category        Category `json:"category"`
	Brand           string   `json:"brand"`
	IsAgeRestricted bool     `json:"is_age_restricted"`
	IsActive        bool     `json:"is_active"`
}

type ProductVariant struct {
	Id        int       `json:"id"`
	Sku       string    `json:"sku"`
	Size      string    `json:"size"`
	Flavor    string    `json:"flavor"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ImagePath string    `json:"image_path"`
	ProductId int       `json:"product_id"`
}

type CreateProductVariant struct {
	Size      string `json:"size" validate:"required"`
	Flavor    string `json:"flavor" validate:"required"`
	ProductId int    `json:"product_id" validate:"required"`
	ImagePath string `json:"image_path" validate:"required"`
}

type UpdateProductVariant struct {
	Size      string `json:"size"`
	Flavor    string `json:"flavor"`
	IsActive  bool   `json:"is_active"`
	ImagePath string `json:"image_path"`
}
