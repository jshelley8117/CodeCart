package persistence

import (
	"context"
	"database/sql"

	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

type OrderPersistence struct {
	DbHandle *sql.DB
}

func NewOrderPersistence(dbHandle *sql.DB) OrderPersistence {
	return OrderPersistence{
		DbHandle: dbHandle,
	}
}

func (op OrderPersistence) PersistCreateOrder(ctx context.Context, orderDomain model.Order) error {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("Entered PersistCreateOrder")

	query := `
		INSERT INTO orders (customer_id, status, total_price, delivery_address, created_at, updated_at, address_id, orderType)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := op.DbHandle.ExecContext(
		ctx,
		query,
		orderDomain.CustomerId,
		orderDomain.Status,
		orderDomain.TotalPrice,
		orderDomain.DeliveryAddress,
		orderDomain.CreatedAt,
		orderDomain.UpdatedAt,
		orderDomain.AddressId,
		orderDomain.OrderType,
	)
	if err != nil {
		zLog.Error("ExecContext failed", zap.Error(err))
		return err
	}
	return nil
}
