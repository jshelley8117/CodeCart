package model

import (
	"encoding/json"
	"time"
)

type OrderStatus string

const (
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusCanceled  OrderStatus = "canceled"
)

type OrderType string

const (
	OrderTypePickup   OrderStatus = "pickup"
	OrderTypeDelivery OrderStatus = "delivery"
)

type Order struct {
	Id              int             `json:"id"`
	CustomerId      int             `json:"customer_id"`
	Status          OrderStatus     `json:"status"`
	TotalPrice      float64         `json:"total_price"`
	DeliveryAddress json.RawMessage `json:"delivery_address"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	AddressId       int             `json:"address_id"`
	OrderType       OrderType       `json:"orderType"`
}

type CreateOrderRequest struct {
	CustomerId      int             `json:"customer_id" validate:"required"`
	Status          OrderStatus     `json:"status" validate:"required"`
	TotalPrice      float64         `json:"total_price" validate:"required"`
	DeliveryAddress json.RawMessage `json:"delivery_address"`
	OrderType       OrderType       `json:"orderType" validate:"required"`
}
