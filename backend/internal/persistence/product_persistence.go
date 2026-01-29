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

func (pp ProductPersistence) FetchProductById(ctx context.Context, id int) *sql.Row {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("Entered FetchProductById")

	query := `
		SELECT id, name, description, unit_price, category, brand, is_age_restricted, created_at, updated_at, is_active
		FROM products
		WHERE id = $1
	`

	return pp.DbHandle.QueryRowContext(ctx, query, id)
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

func (pp ProductPersistence) PersistUpdateProductById(ctx context.Context, productId int, updates map[string]any) error {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("Entered PersistUpdateProductById")

	allowedFields := map[string]bool{
		"name":              true,
		"description":       true,
		"unit_price":        true,
		"category":          true,
		"brand":             true,
		"is_age_restricted": true,
		"is_active":         true,
	}

	query := `
		UPDATE products SET
	`
	args := []any{}
	argPosition := 1

	for field, value := range updates {
		if !allowedFields[field] {
			zLog.Error("Attempted to update invalid field", zap.String("invalid_field", field))
			return fmt.Errorf("invalid field: %v", field)
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
	args = append(args, productId)

	_, err := pp.DbHandle.ExecContext(ctx, query, args...)
	if err != nil {
		zLog.Error("ExecContext failed for PersistUpdateProductById", zap.Error(err))
		return err
	}
	return nil
}

func (pp ProductPersistence) PersistUpdateProductVariantById(ctx context.Context, variantId int, updates map[string]any) error {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("Entered PersistUpdateProductVariantById")

	allowedFields := map[string]bool{
		"size":       true,
		"flavor":     true,
		"is_active":  true,
		"image_path": true,
	}

	query := `
		UPDATE product_variants SET
	`
	args := []any{}
	argPosition := 1

	for field, value := range updates {
		if !allowedFields[field] {
			zLog.Error("Attempted to update invalid field", zap.String("invalid_field", field))
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
	args = append(args, variantId)

	_, err := pp.DbHandle.ExecContext(ctx, query, args...)
	if err != nil {
		zLog.Error("ExecContext failed for PersistUpdateProductVariantById", zap.Error(err))
		return err
	}
	return nil
}

func (pp ProductPersistence) PersistDeleteProductById(ctx context.Context, productId int) error {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("Entered PersistDeleteProductById")

	query := `
		DELETE FROM products
		WHERE id = $1
	`

	if _, err := pp.DbHandle.ExecContext(ctx, query, productId); err != nil {
		zLog.Error("ExecContext failed for PersistDeleteProductById", zap.Error(err))
		return err
	}
	return nil
}

func (pp ProductPersistence) PersistDeleteProductVariantById(ctx context.Context, variantId int) error {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("Entered PersistDeleteProductVariantById")

	query := `
		DELETE FROM product_variants
		WHERE id = $1
	`

	if _, err := pp.DbHandle.ExecContext(ctx, query, variantId); err != nil {
		zLog.Error("ExecContext failed for PersistDeleteProductVariantById")
		return err
	}
	return nil
}
