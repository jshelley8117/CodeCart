package service

import (
	"context"
	"time"

	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/persistence"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

type OrderService struct {
	OrderPersistence persistence.OrderPersistence
}

func NewOrderService(orderPersistence persistence.OrderPersistence) OrderService {
	return OrderService{
		OrderPersistence: orderPersistence,
	}
}

func (os OrderService) CreateOrder(ctx context.Context, request model.CreateOrderRequest) error {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("entered OrderService")

	var orderDomainModel model.Order
	if request.AddressId == 0 {
		orderDomainModel = model.Order{
			CustomerId:      request.CustomerId,
			Status:          "PENDING",
			TotalPrice:      request.TotalPrice,
			DeliveryAddress: request.DeliveryAddress,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			OrderType:       request.OrderType,
			AddressId:       -1,
		}
	} else {
		orderDomainModel = model.Order{
			CustomerId:      request.CustomerId,
			Status:          "PENDING",
			TotalPrice:      request.TotalPrice,
			DeliveryAddress: request.DeliveryAddress,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			OrderType:       request.OrderType,
			AddressId:       request.AddressId,
		}
	}

	if err := os.OrderPersistence.PersistCreateOrder(ctx, orderDomainModel); err != nil {
		zLog.Error("persistence invocation failed: %w", zap.Error(err))
		return err
	}
	return nil
}
