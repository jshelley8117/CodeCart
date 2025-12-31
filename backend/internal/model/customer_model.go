package model

import "time"

type Customer struct {
	Id          string    `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateCustomerRequest struct {
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	PhoneNumber string `json:"phone_number"` // TODO: should validate that phone number conforms to a specific regex pattern
	Email       string `json:"email" validate:"required,email"`
}
