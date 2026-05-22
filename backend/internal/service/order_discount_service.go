package service

import (
	"context"

	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/persistence"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

type OrderDiscountService struct {
	OrderDiscountPersistence persistence.OrderDiscountPersistence
}

func NewOrderDiscountService(orderDiscountPersistence persistence.OrderDiscountPersistence) OrderDiscountService {
	return OrderDiscountService{
		OrderDiscountPersistence: orderDiscountPersistence,
	}
}

func (ods OrderDiscountService) CreateOrderDiscount(ctx context.Context, request model.OrderDiscountRequest) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("entered CreateOrderDiscount CreateLocation")

	orderDiscountDomainModel := model.OrderDiscount{
		OrderId:    request.OrderId,
		DiscountId: request.DiscountId,
	}

	if err := ods.OrderDiscountPersistence.PersistCreateOrderDiscount(ctx, orderDiscountDomainModel); err != nil {
		z.Error("persistence invocation failed", zap.Error(err))
		return err
	}

	return nil
}

func (ods OrderDiscountService) GetOrderDiscountById(ctx context.Context, id int) (model.OrderDiscount, error) {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("entered GetOrderDiscountById")

	row := ods.OrderDiscountPersistence.FetchOrderDiscountById(ctx, id)

	var orderDiscount model.OrderDiscount
	if err := row.Scan(
		&orderDiscount.Id,
		&orderDiscount.OrderId,
		&orderDiscount.DiscountId,
	); err != nil {
		z.Error("scan operation failed", zap.Error(err))
		return model.OrderDiscount{}, err
	}

	return orderDiscount, nil
}
