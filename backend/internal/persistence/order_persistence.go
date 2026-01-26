package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

func (op OrderPersistence) FetchAllOrders(ctx context.Context) (*sql.Rows, error) {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("Entered FetchAllOrders")

	query := `
		SELECT id, customer_id, status, total_price, delivery_address, created_at, updated_at, address_id, order_type
		FROM orders
	`

	rows, err := op.DbHandle.QueryContext(ctx, query)
	if err != nil {
		zLog.Error("QueryContext failed for FetchAllOrders", zap.Error(err))
		return nil, err
	}
	return rows, nil
}

func (op OrderPersistence) FetchOrderById(ctx context.Context, id int) *sql.Row {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("Entered FetchOrderById")

	query := `
		SELECT id, customer_id, status, total_price, delivery_address, created_at, updated_at, address_id, order_type
		FROM customers
		WHERE id = $1
	`

	row := op.DbHandle.QueryRowContext(ctx, query, id)
	return row
}

func (op OrderPersistence) PersistUpdateOrderById(ctx context.Context, id int, updates map[string]any) error {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("Entered PersistUpdateOrderById")

	allowedFields := map[string]bool{
		"status":           true,
		"total_price":      true,
		"delivery_address": true,
		"address_id":       true,
		"order_type":       true,
	}

	query := "UPDATE customers SET"
	args := []any{}
	argPosition := 1

	for field, value := range updates {
		if !allowedFields[field] {
			zLog.Error("Attempted to update invalid field", zap.String("field", field))
			return fmt.Errorf("invalid field: %s", field)
		}

		if argPosition > 1 {
			query += ", "
		}
		query += field + " = $" + fmt.Sprintf("%d", argPosition)
		args = append(args, value)
		argPosition++
	}

	query += ", updated_at = $" + fmt.Sprintf("%d", argPosition)
	args = append(args, time.Now())
	argPosition++

	query += " WHERE id = $" + fmt.Sprintf("%d", argPosition)
	args = append(args, id)

	_, err := op.DbHandle.ExecContext(ctx, query, args...)
	if err != nil {
		zLog.Error("ExecContext failed for PersistUpdateOrderById", zap.Error(err))
		return err
	}
	return nil
}
