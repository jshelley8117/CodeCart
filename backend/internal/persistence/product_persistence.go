package persistence

import (
	"context"
	"database/sql"

	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

type ProductPersistence struct {
	DbHandle *sql.DB
}

func NewProductPersistence(dbHandle *sql.DB) ProductPersistence {
	return ProductPersistence{
		DbHandle: dbHandle,
	}
}

func (pp ProductPersistence) PersistCreateProduct(ctx context.Context, productDomain model.Product) error {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("Entered PersistCreateProduct")

	query := `
		INSERT INTO products (name, description, unit_price, category, brand, is_age_restricted, created_at, updated_at, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := pp.DbHandle.ExecContext(
		ctx,
		query,
		productDomain.Name,
		productDomain.Description,
		productDomain.UnitPrice,
		productDomain.Category,
		productDomain.Brand,
		productDomain.IsAgeRestricted,
		productDomain.CreatedAt,
		productDomain.UpdatedAt,
		productDomain.IsActive,
	)
	if err != nil {
		zLog.Error("ExecContext failed", zap.Error(err))
		return err
	}
	return nil
}
