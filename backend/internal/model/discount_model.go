package model

import "time"

type DiscountType string

const (
	DiscountTypePercentage DiscountType = "PERCENTAGE"
	DiscountTypeFixed      DiscountType = "FIXED"
	PaymentStatusBogo      DiscountType = "BOGO"
)

type Discount struct {
	Id          int          `json:"id"`
	Code        string       `json:"code"`
	Value       float64      `json:"value"`
	Description string       `json:"description"`
	Type        DiscountType `json:"type"`
	IsActive    bool         `json:"is_active"`
	StartsAt    time.Time    `json:"starts_at"`
	ExpiresAt   time.Time    `json:"expires_at"`
	UsageLimit  int          `json:"usage_limit"`
	TimesUsed   int          `json:"times_used"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type CreateDiscount struct {
	Code        string       `json:"code" validate:"required"`
	Value       float64      `json:"value" validate:"required"`
	Description string       `json:"description" validate:"required"`
	Type        DiscountType `json:"type" validate:"required"`
	IsActive    bool         `json:"is_active"`
	StartsAt    time.Time    `json:"starts_at"`
	ExpiresAt   time.Time    `json:"expires_at"`
	UsageLimit  int          `json:"usage_limit"`
}

type UpdateDiscountRequest struct {
	Code        *string       `json:"code"`
	Value       *float64      `json:"value"`
	Description *string       `json:"description"`
	Type        *DiscountType `json:"type"`
	IsActive    *bool         `json:"is_active"`
	StartsAt    *time.Time    `json:"starts_at"`
	ExpiresAt   *time.Time    `json:"expires_at"`
	UsageLimit  *int          `json:"usage_limit"`
}
