package model

import (
	"time"
)

type User struct {
	Id         string    `json:"id"`
	Email      string    `json:"email" validate:"required,email"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	IsActive   bool      `json:"is_active"`
	CustomerId string    `json:"customer_id"`
	AuthId     string    `json:"auth_id"`
}

type CreateUserRequest struct {
	Id       string `json:"id"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	IsActive  bool   `json:"is_active"`
}
