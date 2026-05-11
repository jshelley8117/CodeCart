package model

import (
	"time"
)

type User struct {
	Id         int       `json:"id"`
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	IsActive   bool      `json:"is_active"`
	CustomerId *int      `json:"customer_id"`
	AuthId     string    `json:"auth_id"`
	Role       *string   `json:"role"`
}

type CreateUserRequest struct {
	Email      string  `json:"email" validate:"required,email"`
	CustomerId *int    `json:"customer_id"`
	AuthId     string  `json:"auth_id" validate:"required"`
	Role       *string `json:"role"`
}

type UserResponse struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	IsActive  bool   `json:"is_active"`
	Role      string `json:"role"`
}
