package model

import "time"

type Payment struct {
	Id              int       `json:"id"`
	OrderId         int       `json:"order_id"`
	PaymentMethod   string    `json:"payment_method"`
	Status          string    `json:"status"`
	Currency        string    `json:"currency"`
	PaymentIntentId string    `json:"payment_intent_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreatePaymentRequest struct {
	OrderId         int    `json:"order_id" validate:"required"`
	PaymentMethod   string `json:"payment_method" validate:"required"`
	Currency        string `json:"currency" validate:"required"`
	PaymentIntentId string `json:"payment_intent_id"`
}

type UpdatePaymentRequest struct {
	PaymentMethod   string `json:"payment_method,omitempty"`
	Status          string `json:"status,omitempty"`
	Currency        string `json:"currency,omitempty"`
	PaymentIntentId string `json:"payment_intent_id,omitempty"`
}
