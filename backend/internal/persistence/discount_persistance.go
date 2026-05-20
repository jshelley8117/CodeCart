package persistence

import (
	"context"
	"database/sql"

	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

type DiscountPersistence struct {
	DbHandle *sql.DB
}

func NewDiscountPersistence(dbHandle *sql.DB) DiscountPersistence {
	return DiscountPersistence{
		DbHandle: dbHandle,
	}
}

func (dp DiscountPersistence) PersistCreateDiscount(ctx context.Context, discountDomain model.Discount) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered PersistCreateDiscount")

	query := `
		INSERT INTO discounts (code, value, description, type, is_active, starts_at, expires_at, usage_limit)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := dp.DbHandle.ExecContext(
		ctx,
		query,
		discountDomain.Code,
		discountDomain.Value,
		discountDomain.Description,
		discountDomain.Type,
		discountDomain.IsActive,
		discountDomain.StartsAt,
		discountDomain.ExpiresAt,
		discountDomain.UsageLimit,
	)
	if err != nil {
		z.Error("ExecContext failed", zap.Error(err))
		return err
	}
	return nil
}
