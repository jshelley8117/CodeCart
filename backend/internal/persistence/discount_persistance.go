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

func (dp DiscountPersistence) FetchDiscountById(ctx context.Context, id int) *sql.Row {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered FetchDiscountById")

	query := `
		SELECT id, code, value, description, type, is_active, starts_at, expires_at, usage_limit, times_used, created_at, updated_at
		FROM discounts
		WHERE id = $1
	`

	return dp.DbHandle.QueryRowContext(ctx, query, id)
}

func (dp DiscountPersistence) PersistUpdateDiscountById(ctx context.Context, id int, updates map[string]any) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered PersistUpdateDiscountById")

	allowedFields := map[string]bool{
		"code":        true,
		"value":       true,
		"description": true,
		"type":        true,
		"is_active":   true,
		"starts_at":   true,
		"expires_at":  true,
		"usage_limit": true,
	}

	query := "UPDATE discounts SET "
	args := []any{}
	argPosition := 1

	for field, value := range updates {
		if !allowedFields[field] {
			z.Error("Attempted to update invalid field", zap.String("invalid_field", field))
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

	_, err := dp.DbHandle.ExecContext(ctx, query, args...)
	if err != nil {
		z.Error("ExecContext failed for PersistUpdateDiscountById", zap.Error(err))
		return err
	}
	return nil
}

func (dp DiscountPersistence) PersistDeleteDiscountById(ctx context.Context, id int) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered PersistDeleteDiscountById")

	query := `
		DELETE FROM discounts
		WHERE id = $1
	`

	if _, err := dp.DbHandle.ExecContext(ctx, query, id); err != nil {
		z.Error("ExecContext failed for PersistDeleteDiscountById", zap.Error(err))
		return err
	}
	return nil
}
