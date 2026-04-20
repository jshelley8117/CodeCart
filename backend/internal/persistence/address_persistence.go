package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

type AddressPersistence struct {
	DbHandle *sql.DB
}

func NewAddressPersistence(dbHandle *sql.DB) AddressPersistence {
	return AddressPersistence{
		DbHandle: dbHandle,
	}
}

func (ap AddressPersistence) PersistCreateAddress(ctx context.Context, addressDomain model.Address) error {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("Entered PersistCreateAddress")
	query := `
		INSERT INTO addresses (id, user_id, street_address, city, state, zip_code, country, is_default, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := ap.DbHandle.ExecContext(
		ctx,
		query,
		addressDomain.Id,
		addressDomain.UserId,
		addressDomain.StreetAddress,
		addressDomain.City,
		addressDomain.State,
		addressDomain.ZipCode,
		addressDomain.Country,
		addressDomain.IsDefault,
		addressDomain.CreatedAt,
		addressDomain.UpdatedAt,
	)
	if err != nil {
		log.Printf("Error in PersistCreateAddress, %s", err)
		return err
	}
	return nil
}

func (ap AddressPersistence) FetchAllAddresses(ctx context.Context) (*sql.Rows, error) {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("Entered FetchAllAddresses")

	query := `
		SELECT street_address, city, state, zip_code, country, user_id, id, is_default, created_at, updated_at
		FROM addresses
	`

	rows, err := ap.DbHandle.QueryContext(ctx, query)
	if err != nil {
		zLog.Error("QueryContext failed for FetchAllAddresses", zap.Error(err))
		return nil, err
	}
	return rows, nil
}

func (ap AddressPersistence) FetchAddressById(ctx context.Context, id int) *sql.Row {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered FetchAddressById")

	query := `
		SELECT street_address, city, state, zip_code, country, user_id, id, is_default, created_at, updated_at
		FROM addresses
		WHERE id = $1
	`

	return ap.DbHandle.QueryRowContext(ctx, query, id)
}

func (ap AddressPersistence) PersistUpdateAddressById(ctx context.Context, id int, updates map[string]any) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered PersistUpdateAddressById")

	allowedFields := map[string]bool{
		"street_address": true,
		"city":           true,
		"state":          true,
		"zip_code":       true,
		"country":        true,
		"is_default":     true,
	}

	query := `
		UPDATE addresses SET
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
	args = append(args, id)

	if _, err := ap.DbHandle.ExecContext(ctx, query, args...); err != nil {
		z.Error("ExecContext failed for PersistUpdateAddressById", zap.Error(err))
		return err
	}
	return nil
}

func (ap AddressPersistence) PersistDeleteAddressById(ctx context.Context, id int) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered PersistDeleteAddressById")

	query := `
		DELETE FROM addresses
		WHERE id = $1
	`

	if _, err := ap.DbHandle.ExecContext(ctx, query, id); err != nil {
		z.Error("ExecContext failed for PersistDeleteAddressById", zap.Error(err))
		return err
	}

	return nil
}
