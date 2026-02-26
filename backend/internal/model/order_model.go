package model

import (
	"encoding/json"
	"time"
)

type PaymentStatus string

const (
	PaymentStatusSuccess  PaymentStatus = "SUCCESS"
	PaymentStatusPending  PaymentStatus = "PENDING"
	PaymentStatusError    PaymentStatus = "ERROR"
	PaymentStatusCanceled PaymentStatus = "CANCELED"
)

type FulfillmentStatus string

const (
	FulfillmentStatusComplete   FulfillmentStatus = "COMPLETE"
	FulfillmentStatusInProgress FulfillmentStatus = "IN_PROGRESS"
	FulfillmentStatusPending    FulfillmentStatus = "PENDING"
	FulfillmentStatusCanceled   FulfillmentStatus = "CANCELED"
)

type OrderType string

const (
	OrderTypePickup   OrderType = "PICKUP"
	OrderTypeDelivery OrderType = "DELIVERY"
)

type Order struct {
	Id                int               `json:"id"`
	CustomerId        int               `json:"customer_id"`
	PaymentStatus     PaymentStatus     `json:"payment_status"`
	FulfillmentStatus FulfillmentStatus `json:"fulfillment_status"`
	TotalPrice        float64           `json:"total_price"`
	DeliveryAddress   json.RawMessage   `json:"delivery_address"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
	AddressId         int               `json:"address_id"`
	OrderType         OrderType         `json:"order_type"`
}

type CreateOrderRequest struct {
	CustomerId      int             `json:"customer_id" validate:"required"`
	TotalPrice      float64         `json:"total_price" validate:"required"`
	DeliveryAddress json.RawMessage `json:"delivery_address"`
	OrderType       OrderType       `json:"order_type" validate:"required"`
	AddressId       int             `json:"address_id"`
}

type UpdateOrderRequest struct {
	PaymentStatus     PaymentStatus     `json:"payment_status"`
	FulfillmentStatus FulfillmentStatus `json:"Fulfillment_status"`
	TotalPrice        float64           `json:"total_price"`
	DeliveryAddress   json.RawMessage   `json:"delivery_address"`
	AddressId         int               `json:"address_id"`
	OrderType         OrderType         `json:"order_type"`
}
