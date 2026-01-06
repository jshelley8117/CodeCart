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
		Logger:   logger.Named("customer_persistence"),
	}
}

func (cp CustomerPersistence) PersistCreateCustomer(ctx context.Context, customerDomain model.Customer) error {
	zLog := cp.getZLog(ctx)
	zLog.Debug("entered PersistCreateCustomer")
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
		zLog.Error("ExecContext failed for PersistCreateCustomer", zap.Error(err))
		return err
	}

	return nil
}

func (cp CustomerPersistence) FetchAllCustomersById(ctx context.Context) (*sql.Rows, error) {
	zLog := cp.getZLog(ctx)
	zLog.Debug("entered FetchAllCustomersById")
	query := `
		SELECT id, first_name, last_name, phone_number, email, created_at, updated_at
		FROM customers
	`

	rows, err := cp.DbHandle.QueryContext(ctx, query)
	if err != nil {
		zLog.Error("QueryContext failed for FetchAllCustomersById", zap.Error(err))
		return nil, err
	}
	return rows, nil
}

func (cp CustomerPersistence) PersistDeleteCustomerById(ctx context.Context, id int) error {
	zLog := cp.getZLog(ctx)
	zLog.Debug("entered DeleteCustomerById")
	query := `
		DELETE FROM customers
		WHERE id = $1
	`

	if _, err := cp.DbHandle.ExecContext(ctx, query, id); err != nil {
		zLog.Error("ExecContext failed for DeleteCustomerById", zap.Error(err))
		return err
	}
	return nil
}

func (cp CustomerPersistence) getZLog(ctx context.Context) *zap.Logger {
	return utils.FromContext(ctx, cp.Logger)
}
