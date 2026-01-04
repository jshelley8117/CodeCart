package persistence

import (
	"context"
	"database/sql"

	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

type CustomerPersistence struct {
	DbHandle *sql.DB
	Logger   *zap.Logger
}

func NewCustomerPersistence(dbHandle *sql.DB, logger *zap.Logger) CustomerPersistence {
	return CustomerPersistence{
		DbHandle: dbHandle,
		Logger:   logger,
	}
}

func (cp CustomerPersistence) PersistCreateCustomer(ctx context.Context, customerDomain model.Customer) error {
	zLog := utils.FromContext(ctx, cp.Logger).Named("customer_persistence")
	zLog.Debug("entered CustomerPersistence")
	query := `
		INSERT INTO customers (first_name, last_name, phone_number, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := cp.DbHandle.ExecContext(
		ctx,
		query,
		customerDomain.FirstName,
		customerDomain.LastName,
		customerDomain.PhoneNumber,
		customerDomain.Email,
		customerDomain.CreatedAt,
		customerDomain.UpdatedAt,
	)
	if err != nil {
		zLog.Error("ExecContext failed: %w", zap.Error(err))
		return err
	}

	return nil
}
