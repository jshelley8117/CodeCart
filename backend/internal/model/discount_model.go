package model

import "time"

type DiscountType string

const (
	DiscountTypePercentage DiscountType = "PERCENTAGE"
	DiscountTypeFixed      DiscountType = "FIXED"
	PaymentStatusBogo      DiscountType = "BOGO"
)

type Discount struct {
	Id         string       `json:"id" `
	Code       string       `json:"code"`
	Value      string       `json:"value"`
	Describe   string       `json:"describe"`
	Type       DiscountType `json:"type"`
	IsActive   bool         `json:"is_active"`
	StartsAt   time.Time    `json:"starts_at"`
	ExpiresAt  time.Time    `json:"expires_at"`
	UsageLimit int          `json:"usage_limit"`
	TimesUsed  int          `json:"times_used"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

type CreateDiscount struct {
	Code       string    `json:"code" validate:"required"`
	Value      string    `json:"value" validate:"required"`
	Describe   string    `json:"describe" validate:"required"`
	Type       string    `json:"type" validate:"required"`
	IsActive   bool      `json:"is_active" validate:"required"`
	StartsAt   time.Time `json:"starts_at" validate:"required"`
	ExpiresAt  time.Time `json:"expires_at" validate:"required"`
	UsageLimit int       `json:"usage_limit" validate:"required"`
}

type UpdateDiscount struct {
	Code       string    `json:"code"`
	Value      string    `json:"value"`
	Describe   string    `json:"describe"`
	Type       string    `json:"type"`
	IsActive   bool      `json:"is_active"`
	StartsAt   time.Time `json:"starts_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	UsageLimit int       `json:"usage_limit"`
}
