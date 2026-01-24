package model

import "time"

type Category string

const (
	Grocery            Category = "GROCERY"
	GeneralMerchandise Category = "GENERAL_MERCHANDISE"
	Produce            Category = "PRODUCE"
	Seafood            Category = "SEAFOOD"
	MeatMarket         Category = "MEAT_MARKET"
)

type Product struct {
	Id              int       `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	UnitPrice       int       `json:"unit_price"`
	Category        Category  `json:"category"`
	Brand           string    `json:"brand"`
	IsAgeRestricted bool      `json:"is_age_restricted"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	IsActive        bool      `json:"is_active"`
}

type CreateProductRequest struct {
	Name            string   `json:"name" validate:"required"`
	Description     string   `json:"description"`
	UnitPrice       int      `json:"price" validate:"required"`
	Category        Category `json:"category" validate:"required"`
	Brand           string   `json:"brand" validate:"required"`
	IsAgeRestricted bool     `json:"is_age_restricted" validate:"required"`
}
