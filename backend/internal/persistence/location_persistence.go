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

type LocationPersistence struct {
	DbHandle *sql.DB
}

func NewLocationPersistence(dbHandle *sql.DB) LocationPersistence {
	return LocationPersistence{
		DbHandle: dbHandle,
	}
}

func (lp LocationPersistence) PersistCreateLocation(ctx context.Context, locationDomain model.Location) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered PersistCreateLocation")

	query := `
		INSERT INTO locations (store_name, store_address, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := lp.DbHandle.ExecContext(
		ctx,
		query,
		locationDomain.StoreName,
		locationDomain.StoreAddress,
		locationDomain.CreatedAt,
		locationDomain.UpdatedAt,
	)
	if err != nil {
		z.Error("ExecContext failed for PersistCreateLocation", zap.Error(err))
		return err
	}
	return nil
}

func (lp LocationPersistence) FetchAllLocations(ctx context.Context, page, pageSize int) (*sql.Rows, int64, error) {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered FetchAllLocations")

	var total int64
	countQuery := "SELECT COUNT(*) FROM locations"
	if err := lp.DbHandle.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		z.Error("QueryRowContext failed on the pagination count query", zap.Error(err))
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	query := `
		SELECT id, store_name, store_address, created_at, updated_at
		FROM locations
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := lp.DbHandle.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		z.Error("QueryContext failed", zap.Error(err))
		return nil, 0, err
	}
	return rows, total, nil
}

func (lp LocationPersistence) FetchLocationById(ctx context.Context, id int) *sql.Row {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered FetchLocationById")

	query := `
		SELECT id, store_name, store_address, created_at, updated_at
		FROM locations
		WHERE id = $1
	`

	return lp.DbHandle.QueryRowContext(ctx, query, id)
}

func (lp LocationPersistence) PersistUpdateLocationById(ctx context.Context, locationId int, updates map[string]any) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered PersistUpdateLocationById")

	allowedFields := map[string]bool{
		"store_name":    true,
		"store_address": true,
	}

	query := `
		UPDATE locations SET
	`
	args := []any{}
	argPosition := 1

	for field, value := range updates {
		if !allowedFields[field] {
			z.Error("Attempted to update invalid field", zap.String("invalid_field", field))
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
	args = append(args, locationId)

	_, err := lp.DbHandle.ExecContext(ctx, query, args...)
	if err != nil {
		z.Error("ExecContext failed for PersistUpdateLocationById", zap.Error(err))
		return err
	}
	return nil
}

func (lp LocationPersistence) PersistDeleteLocationById(ctx context.Context, locationId int) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered PersistDeleteLocationById")

	query := `
		DELETE FROM locations
		WHERE id = $1
	`

	if _, err := lp.DbHandle.ExecContext(ctx, query, locationId); err != nil {
		z.Error("ExecContext failed for PersistDeleteLocationById", zap.Error(err))
		return err
	}
	return nil
}
