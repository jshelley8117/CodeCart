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

func (pp ProductPersistence) FetchAllProducts(ctx context.Context, page, pageSize int) (*sql.Rows, int64, error) {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("Entered FetchAllProducts")

	var total int64
	countQuery := "SELECT COUNT(*) FROM products"
	if err := pp.DbHandle.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		zLog.Error("QueryRowContext failed on the pagination count query", zap.Error(err))
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	query := `
		SELECT id, name, description, unit_price, category, brand, is_age_restricted, created_at, updated_at, is_active
		FROM products
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := pp.DbHandle.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		zLog.Error("QueryContext failed", zap.Error(err))
		return nil, 0, err
	}
	return rows, total, nil
}

func (pp ProductPersistence) FetchAllProductVariantsByProductId(ctx context.Context, productId int) (*sql.Rows, error) {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("Entered FetchAllProductVariantsByProductId")

	query := `
		SELECT id, sku, size, flavor, is_active, created_at, updated_at, image_path, product_id
		FROM product_variants
		WHERE id = $1
	`

	rows, err := pp.DbHandle.QueryContext(ctx, query, productId)
	if err != nil {
		zLog.Error("QueryContext failed", zap.Error(err))
		return nil, err
	}
	return rows, nil
}
