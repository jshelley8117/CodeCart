package model

import "time"

type Customer struct {
	Id          int       `json:"id"`
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
	PhoneNumber string `json:"phone_number" validate:"e164"`
	Email       string `json:"email" validate:"required,email"`
}

// pointers are used in some fields here because the zero value for strings is "", meaning that after unmarshaling, the empty/omitted fields will have
// a value of "", and that would fail the fields that have a special struct tag validation (e.g. email, e164, etc...)
type UpdateCustomerRequest struct {
	FirstName   string  `json:"first_name,omitempty"`
	LastName    string  `json:"last_name,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty" validate:"omitempty,e164"`
	Email       *string `json:"email,omitempty" validate:"omitempty,email"`
}
