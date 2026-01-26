package model

import (
	"encoding/json"
	"time"
)

type OrderStatus string

const (
	OrderStatusConfirmed OrderStatus = "CONFIRMED"
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusCanceled  OrderStatus = "CANCELED"
)

type OrderType string

const (
	OrderTypePickup   OrderStatus = "PICKUP"
	OrderTypeDelivery OrderStatus = "DELIVERY"
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
	OrderType       OrderType       `json:"order_type"`
}

type CreateOrderRequest struct {
	CustomerId      int             `json:"customer_id" validate:"required"`
	TotalPrice      float64         `json:"total_price" validate:"required"`
	DeliveryAddress json.RawMessage `json:"delivery_address"`
	OrderType       OrderType       `json:"order_type" validate:"required"`
	AddressId       int             `json:"address_id"`
}

type UpdateOrderRequest struct {
	Status          OrderStatus     `json:"status"`
	TotalPrice      float64         `json:"total_price"`
	DeliveryAddress json.RawMessage `json:"delivery_address"`
	AddressId       int             `json:"address_id"`
	OrderType       OrderType       `json:"order_type"`
}
