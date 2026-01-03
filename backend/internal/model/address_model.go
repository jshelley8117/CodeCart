package model

import (
	"time"
)

type Address struct {
	StreetAddress string    `json:"street_address"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	ZipCode       string    `json:"zip_code"`
	Country       string    `json:"country"`
	UserId        string    `json:"user_id"`
	Id            string    `json:"id"`
	IsDefault     bool      `json:"is_default"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateAddressRequest struct {
	UserId        string `json:"user_id" validate:"required"`
	StreetAddress string `json:"street_address" validate:"required"`
	City          string `json:"city" validate:"required"`
	State         string `json:"state" validate:"required"`
	ZipCode       string `json:"zip_code" validate:"required"`
	Country       string `json:"country" validate:"required"`
}
