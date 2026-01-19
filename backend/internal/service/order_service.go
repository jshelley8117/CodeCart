package service

import (
	"context"
	"log"
	"time"

	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/persistence"
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
	log.Println("Entered CreateOrderService")
	var orderDomainModel model.Order
	if request.AddressId == 0 {
		orderDomainModel = model.Order{
			CustomerId:      request.CustomerId,
			Status:          request.Status,
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
			Status:          request.Status,
			TotalPrice:      request.TotalPrice,
			DeliveryAddress: request.DeliveryAddress,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			OrderType:       request.OrderType,
			AddressId:       request.AddressId,
		}
	}

	if err := os.OrderPersistence.PersistCreateOrder(ctx, orderDomainModel); err != nil {
		return err
	}
	return nil
}
