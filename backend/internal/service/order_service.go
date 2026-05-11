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
}

func NewOrderService(orderPersistence persistence.OrderPersistence) OrderService {
	return OrderService{
		OrderPersistence: orderPersistence,
	}
}

func (os OrderService) CreateOrder(ctx context.Context, request model.CreateOrderRequest) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("entered OrderService")

	var orderDomainModel model.Order
	if request.AddressId == 0 {
		orderDomainModel = model.Order{
			CustomerId:        request.CustomerId,
			PaymentStatus:     model.PaymentStatusPending,
			FulfillmentStatus: model.FulfillmentStatusPending,
			TotalPrice:        request.TotalPrice,
			DeliveryAddress:   request.DeliveryAddress,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			OrderType:         request.OrderType,
			AddressId:         -1,
		}
	} else {
		orderDomainModel = model.Order{
			CustomerId:        request.CustomerId,
			PaymentStatus:     model.PaymentStatusPending,
			FulfillmentStatus: model.FulfillmentStatusPending,
			TotalPrice:        request.TotalPrice,
			DeliveryAddress:   request.DeliveryAddress,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			OrderType:         request.OrderType,
			AddressId:         request.AddressId,
		}
	}

	if err := os.OrderPersistence.PersistCreateOrder(ctx, orderDomainModel); err != nil {
		z.Error("persistence invocation failed", zap.Error(err))
		return err
	}
	return nil
}

func (os OrderService) GetAllOrders(ctx context.Context) ([]model.Order, error) {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("entered GetAllOrders")

	orderRows, err := os.OrderPersistence.FetchAllOrders(ctx)
	if err != nil {
		z.Error("persistence invocation failed", zap.Error(err))
		return nil, err
	}
	defer orderRows.Close()

	orders := make([]model.Order, 0)

	for orderRows.Next() {
		var order model.Order
		if err := orderRows.Scan(
			&order.Id,
			&order.CustomerId,
			&order.PaymentStatus,
			&order.FulfillmentStatus,
			&order.TotalPrice,
			&order.DeliveryAddress,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.AddressId,
			&order.OrderType,
		); err != nil {
			z.Error("scan operation failed", zap.Error(err))
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := orderRows.Err(); err != nil {
		z.Error("error occured while iterating through sql rows", zap.Error(err))
		return nil, err
	}
	return orders, nil
}

func (os OrderService) FetchOrderById(ctx context.Context, id int) (model.Order, error) {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("entered FetchOrderById")

	orderRow := os.OrderPersistence.FetchOrderById(ctx, id)
	if orderRow == nil {
		z.Warn("order not found", zap.Int("order_id", id))
		return model.Order{}, nil
	}

	var order model.Order
	if err := orderRow.Scan(
		&order.Id,
		&order.CustomerId,
		&order.PaymentStatus,
		&order.FulfillmentStatus,
		&order.TotalPrice,
		&order.DeliveryAddress,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.AddressId,
		&order.OrderType,
	); err != nil {
		z.Error("scan operation failed", zap.Error(err))
		return model.Order{}, err
	}

	return order, nil
}

func (os OrderService) UpdateOrderById(ctx context.Context, request model.UpdateOrderRequest, id int) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("entered UpdateOrderById")

	updates := make(map[string]any)

	if request.PaymentStatus != "" {
		if !validatePaymentStatus(request.PaymentStatus) {
			z.Error("invalid status", zap.Int("order_id", id))
			return fmt.Errorf("invalid status: %s", request.PaymentStatus)
		}
		updates["payment_status"] = request.PaymentStatus
	}
	if request.FulfillmentStatus != "" {
		if !validateFulfillmentStatus(request.FulfillmentStatus) {
			z.Error("invalid status", zap.Int("order_id", id))
			return fmt.Errorf("invalid status: %s", request.FulfillmentStatus)
		}
		updates["Fulfillment_status"] = request.FulfillmentStatus

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
			z.Error("invalid order type", zap.Int("order_id", id))
			return fmt.Errorf(common.ERR_CLIENT_DB_PERSISTENCE_FAIL)
		}
		updates["order_type"] = request.OrderType
	}

	if len(updates) == 0 {
		z.Error("No updates found", zap.Int("order_id", id))
		return fmt.Errorf("no updates found")
	}

	if err := os.OrderPersistence.PersistUpdateOrderById(ctx, id, updates); err != nil {
		z.Error("persistence invocation failed", zap.Error(err))
		return err
	}

	return nil
}

func validatePaymentStatus(status model.PaymentStatus) bool {
	return status == model.PaymentStatus("SUCCESS") || status == model.PaymentStatus("PENDING") || status == model.PaymentStatus("ERROR") || status == model.PaymentStatus("CANCELED")

}

func validateFulfillmentStatus(status model.FulfillmentStatus) bool {
	return status == model.FulfillmentStatus("COMPLETE") || status == model.FulfillmentStatus("IN_PROGRESS") || status == model.FulfillmentStatus("PENDING") || status == model.FulfillmentStatus("CANCELED")
}

func validateType(orderType model.OrderType) bool {
	return orderType == model.OrderType("DELIVERY") || orderType == model.OrderType("PICKUP")
}
