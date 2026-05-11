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

type InventoryPersistence struct {
	DbHandle *sql.DB
}

func NewInventoryPersistence(dbHandle *sql.DB) InventoryPersistence {
	return InventoryPersistence{
		DbHandle: dbHandle,
	}
}

func (ip InventoryPersistence) PersistCreateInventory(ctx context.Context, inventoryDomain model.Inventory) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered PersistCreateInventory")

	query := `
		INSERT INTO inventory (product_variant_id, location_id, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := ip.DbHandle.ExecContext(
		ctx,
		query,
		inventoryDomain.ProductVariantId,
		inventoryDomain.LocationId,
		inventoryDomain.Quantity,
		inventoryDomain.CreatedAt,
		inventoryDomain.UpdatedAt,
	)
	if err != nil {
		z.Error("ExecContext failed", zap.Error(err))
		return err
	}
	return nil
}

func (ip InventoryPersistence) FetchAllInventory(ctx context.Context, page, pageSize int) (*sql.Rows, int64, error) {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered FetchAllInventory")

	var total int64
	countQuery := "SELECT COUNT(*) FROM inventory"
	if err := ip.DbHandle.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		z.Error("QueryRowContext failed on the pagination count query", zap.Error(err))
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	query := `
		SELECT id, product_variant_id, location_id, quantity, created_at, updated_at
		FROM inventory
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := ip.DbHandle.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		z.Error("QueryContext failed", zap.Error(err))
		return nil, 0, err
	}
	return rows, total, nil
}

func (ip InventoryPersistence) FetchInventoryById(ctx context.Context, id int) *sql.Row {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered FetchInventoryById")

	query := `
		SELECT id, product_variant_id, location_id, quantity, created_at, updated_at
		FROM inventory
		WHERE id = $1
	`

	return ip.DbHandle.QueryRowContext(ctx, query, id)
}

func (ip InventoryPersistence) PersistUpdateInventoryById(ctx context.Context, id int, updates map[string]any) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered PersistUpdateInventoryById")

	allowedFields := map[string]bool{
		"quantity":    true,
		"location_id": true,
	}

	query := "UPDATE inventory SET "
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

	_, err := ip.DbHandle.ExecContext(ctx, query, args...)
	if err != nil {
		z.Error("ExecContext failed for PersistUpdateInventoryById", zap.Error(err))
		return err
	}
	return nil
}

func (ip InventoryPersistence) PersistDeleteInventoryById(ctx context.Context, id int) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered PersistDeleteInventoryById")

	query := `
		DELETE FROM inventory
		WHERE id = $1
	`

	if _, err := ip.DbHandle.ExecContext(ctx, query, id); err != nil {
		z.Error("ExecContext failed for PersistDeleteInventoryById", zap.Error(err))
		return err
	}
	return nil
}
