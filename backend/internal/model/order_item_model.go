package model

type OrderItem struct {
	Id               int     `json:"id"`
	OrderId          int     `json:"order_id"`
	ProductVariantId int     `json:"product_variant_id"`
	Quantity         int     `json:"quantity"`
	UnitPrice        float64 `json:"unit_price"`
	Discount         float64 `json:"discount"`
}

type CreateOrderItemRequest struct {
	OrderId          int     `json:"order_id" validate:"required"`
	ProductVariantId int     `json:"product_variant_id" validate:"required"`
	Quantity         int     `json:"quantity" validate:"required"`
	UnitPrice        float64 `json:"unit_price" validate:"required"`
	Discount         float64 `json:"discount" validate:"omitempty,min=0"`
}

type UpdateOrderItemRequest struct {
	ProductVariantId *int     `json:"product_variant_id"`
	Quantity         *int     `json:"quantity" validate:"omitempty,min=1"`
	UnitPrice        *float64 `json:"unit_price" validate:"omitempty,min=0"`
	Discount         *float64 `json:"discount" validate:"omitempty,min=0"`
}
