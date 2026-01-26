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
		Logger:           logger,
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
		zLog.Error("persistence invocation failed", zap.Error(err))
		return err
	}
	return nil
}

func (os OrderService) GetAllOrders(ctx context.Context) ([]model.Order, error) {
	zLog := utils.FromContext(ctx, os.Logger).Named("order_service")
	zLog.Debug("entered GetAllOrders")

	orderRows, err := os.OrderPersistence.FetchAllOrders(ctx)
	if err != nil {
		zLog.Error("persistence invocation failed", zap.Error(err))
		return nil, err
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
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := orderRows.Err(); err != nil {
		zLog.Error("error occured while iterating through sql rows", zap.Error(err))
		return nil, err
	}
	return orders, nil
}

func (os OrderService) FetchOrderById(ctx context.Context, id int) (model.Order, error) {
	zLog := utils.FromContext(ctx, os.Logger).Named("order_service")
	zLog.Debug("entered FetchOrderById")

	orderRow := os.OrderPersistence.FetchOrderById(ctx, id)
	if orderRow == nil {
		zLog.Error("order not found", zap.Int("order_id", id))
		return model.Order{}, fmt.Errorf(common.ERR_CLIENT_DB_PERSISTENCE_FAIL)
	}

	var order model.Order
	if err := orderRow.Scan(
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
		return model.Order{}, err
	}

	return order, nil
}

func (os OrderService) UpdateOrderById(ctx context.Context, request model.UpdateOrderRequest, id int) error {
	zLog := utils.FromContext(ctx, os.Logger).Named("order_service")
	zLog.Debug("entered UpdateOrderById")

	updates := make(map[string]any)

	if request.Status != "" {
		if !validateStatus(request.Status) {
			zLog.Error("invalid status", zap.Int("order_id", id))
			return fmt.Errorf(common.ERR_CLIENT_DB_PERSISTENCE_FAIL)
		}
		updates["status"] = request.Status
	}
	if request.TotalPrice != 0 {
		updates["total_price"] = request.TotalPrice
	}
	if request.DeliveryAddress != nil {
		updates["delivery_address"] = request.DeliveryAddress
	}
	if request.AddressId != 0 {
		updates["address_id"] = request.AddressId
	}
	if request.OrderType != "" {
		if !validateType(request.OrderType) {
			zLog.Error("invalid order type", zap.Int("order_id", id))
			return fmt.Errorf(common.ERR_CLIENT_DB_PERSISTENCE_FAIL)
		}
		updates["order_type"] = request.OrderType
	}

	if len(updates) == 0 {
		zLog.Error("No updates found", zap.Int("order_id", id))
		return fmt.Errorf(common.ERR_CLIENT_DB_PERSISTENCE_FAIL)
	}

	if err := os.OrderPersistence.PersistUpdateOrderById(ctx, id, updates); err != nil {
		zLog.Error("persistence invocation failed", zap.Error(err))
		return err
	}

	return nil
}

func validateStatus(status model.OrderStatus) bool {
	return status == model.OrderStatus("PENDING") || status == model.OrderStatus("DELIVERED") || status == model.OrderStatus("CANCELLED")
}

func validateType(orderType model.OrderType) bool {
	return orderType == model.OrderType("DELIVERY") || orderType == model.OrderType("PICKUP")
}
