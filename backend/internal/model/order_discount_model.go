package model

type OrderDiscount struct {
	Id         int `json:"id"`
	OrderId    int `json:"order_id"`
	DiscountId int `json:"disocunt_id"`
}

type OrderDiscountRequest struct {
	OrderId    int `json:"order_id" validate:"required"`
	DiscountId int `json:"disocunt_id" validate:"required"`
}
