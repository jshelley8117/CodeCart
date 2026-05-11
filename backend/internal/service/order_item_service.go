package service

import (
	"context"
	"fmt"

	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/persistence"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

type OrderItemService struct {
	OrderItemPersistence persistence.OrderItemPersistence
}

func NewOrderItemService(orderItemPersistence persistence.OrderItemPersistence) OrderItemService {
	return OrderItemService{
		OrderItemPersistence: orderItemPersistence,
	}
}

func (ois OrderItemService) CreateOrderItem(ctx context.Context, request model.CreateOrderItemRequest, orderId int) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("entered CreateOrderItems")

	orderItemDomainModel := model.OrderItem{
		OrderId:          orderId,
		ProductVariantId: request.ProductVariantId,
		Quantity:         request.Quantity,
		UnitPrice:        request.UnitPrice,
		Discount:         request.Discount,
	}

	if err := ois.OrderItemPersistence.PersistCreateItemOrder(ctx, orderItemDomainModel); err != nil {
		z.Error("persistence invocation failed", zap.Error(err))
		return err
	}

	return nil
}

func (ois OrderItemService) GetAllOrderItems(ctx context.Context, orderId int) ([]model.OrderItem, error) {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("entered GetAllOrderItems")

	rows, err := ois.OrderItemPersistence.FetchAllOrderItemsByOrderId(ctx, orderId)
	if err != nil {
		z.Error("persistence invocation failed", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	items := make([]model.OrderItem, 0)

	for rows.Next() {
		var item model.OrderItem
		if err := rows.Scan(
			&item.Id,
			&item.OrderId,
			&item.ProductVariantId,
			&item.Quantity,
			&item.UnitPrice,
			&item.Discount,
		); err != nil {
			z.Error("scan operation failed", zap.Error(err))
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		z.Error("error occurred while iterating through sql rows", zap.Error(err))
		return nil, err
	}

	return items, nil
}

func (ois OrderItemService) UpdateOrderItemById(ctx context.Context, request model.UpdateOrderItemRequest, orderId int, id int) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("entered UpdateOrderItemsById")

	updates := make(map[string]any)

	if request.ProductVariantId != nil {
		updates["product_variant_id"] = *request.ProductVariantId
	}
	if request.Quantity != nil {
		updates["quantity"] = *request.Quantity
	}
	if request.UnitPrice != nil {
		updates["unit_price"] = *request.UnitPrice
	}
	if request.Discount != nil {
		updates["discount"] = *request.Discount
	}

	if len(updates) == 0 {
		z.Error("no updates found", zap.Int("order_item_id", id))
		return fmt.Errorf("no updates found")
	}

	if err := ois.OrderItemPersistence.PersistUpdateOrderItemById(ctx, orderId, id, updates); err != nil {
		z.Error("persistence invocation failed", zap.Error(err))
		return err
	}

	return nil
}

func (ois OrderItemService) DeleteOrderItemById(ctx context.Context, orderId int, id int) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("entered DeleteOrderItemsById")

	if err := ois.OrderItemPersistence.PersistDeleteOrderItemById(ctx, orderId, id); err != nil {
		z.Error("persistence invocation failed", zap.Error(err))
		return err
	}

	return nil
}
