package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jshelley8117/CodeCart/internal/common"
	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/persistence"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

type OrderService struct {
	OrderPersistence persistence.OrderPersistence
	Logger           *zap.Logger
}

func NewOrderService(orderPersistence persistence.OrderPersistence, logger *zap.Logger) OrderService {
	return OrderService{
		OrderPersistence: orderPersistence,
	}
}

func (os OrderService) CreateOrder(ctx context.Context, request model.CreateOrderRequest) error {
	zLog := utils.FromContext(ctx, os.Logger).Named("order_service")
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

func (os OrderService) GetAllOrders(ctx context.Context) ([]model.Order, error) {
	zLog := utils.FromContext(ctx, os.Logger).Named("order_service")
	zLog.Debug("entered GetAllOrders")

	orderRows, err := os.OrderPersistence.FetchAllOrders(ctx)
	if err != nil {
		zLog.Error("persistence invocation failed: %w", zap.Error(err))
		return nil, fmt.Errorf(common.ERR_CLIENT_DB_RETRIEVAL_FAIL)
	}
	defer orderRows.Close()

	orders := make([]model.Order, 0)

	for orderRows.Next() {
		var order model.Order
		if err := orderRows.Scan(
			&order.Id,
			&order.CustomerId,
			&order.Status,
			&order.TotalPrice,
			&order.DeliveryAddress,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.AddressId,
			&order.OrderType,
		); err != nil {
			zLog.Error("scan operation failed", zap.Error(err))
			return nil, fmt.Errorf(common.ERR_CLIENT_DB_RETRIEVAL_FAIL)
		}
		orders = append(orders, order)
	}

	if err := orderRows.Err(); err != nil {
		zLog.Error("error occured while iterating through sql rows", zap.Error(err))
		return nil, fmt.Errorf(common.ERR_CLIENT_DB_RETRIEVAL_FAIL)
	}
	return orders, nil
}
