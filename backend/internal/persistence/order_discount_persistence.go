package persistence

import (
	"context"
	"database/sql"

	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

type OrderDiscountPersistence struct {
	DbHandle *sql.DB
}

func NewOrderDiscountPersistence(dbHandle *sql.DB) OrderDiscountPersistence {
	return OrderDiscountPersistence{
		DbHandle: dbHandle,
	}
}

func (odp OrderDiscountPersistence) PersistCreateOrderDiscount(ctx context.Context, orderDiscountDomain model.OrderDiscount) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered PersistCreateOrderDiscount")

	query := `NewOrderDiscountPersistence
		INSERT INTO order_discounts (order_id, discount_id)
		VALUES ($1, $2)
	`

	_, err := odp.DbHandle.ExecContext(
		ctx,
		query,
		orderDiscountDomain.OrderId,
		orderDiscountDomain.DiscountId,
	)

	if err != nil {
		z.Error("ExexContext failed", zap.Error(err))
		return err
	}

	return nil
}

func (odp OrderDiscountPersistence) FetchOrderDiscountById(ctx context.Context, id int) *sql.Row {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered FetchOrderDiscountById")

	query := `
		SELECT id, order_id, discount_id
		FROM order_discounts
		WHERE id = $1
	`

	return odp.DbHandle.QueryRowContext(ctx, query, id)
}
