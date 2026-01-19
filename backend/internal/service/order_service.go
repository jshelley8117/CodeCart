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

// CreateOrder creates an order in the database
//
// request - The create order request
//
// Returns an error if the persistence layer fails to create the order
func (os OrderService) CreateOrder(ctx context.Context, request model.CreateOrderRequest) error {
	log.Println("Entered CreateOrderService")
	orderDomainModel := model.Order{
		CustomerId:      request.CustomerId,
		Status:          request.Status,
		TotalPrice:      request.TotalPrice,
		DeliveryAddress: request.DeliveryAddress,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		OrderType:       request.OrderType,
	}

	log.Println("Trying to save Order with AddressId: ", orderDomainModel.CustomerId)

	if err := os.OrderPersistence.PersistCreateOrder(ctx, orderDomainModel); err != nil {
		return err
	}
	return nil
}
