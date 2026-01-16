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

func (cp CustomerPersistence) FetchAllCustomers(ctx context.Context) (*sql.Rows, error) {
	zLog := cp.getZLog(ctx)
	zLog.Debug("entered FetchAllCustomersById")
	query := `
		SELECT id, first_name, last_name, phone_number, email, created_at, updated_at
		FROM customers
	`

	rows, err := cp.DbHandle.QueryContext(ctx, query)
	if err != nil {
		zLog.Error("QueryContext failed for FetchAllCustomers", zap.Error(err))
		return nil, err
	}
	return rows, nil
}

func (cp CustomerPersistence) PersistDeleteCustomerById(ctx context.Context, id int) error {
	zLog := cp.getZLog(ctx)
	zLog.Debug("entered PersistDeleteCustomerById")
	query := `
		DELETE FROM customers
		WHERE id = $1
	`

	if _, err := cp.DbHandle.ExecContext(ctx, query, id); err != nil {
		zLog.Error("ExecContext failed for PersistDeleteCustomerById", zap.Error(err))
		return err
	}
	return nil
}

func (cp CustomerPersistence) PersistUpdateCustomerById(ctx context.Context, id int, updates map[string]any) error {
	zLog := cp.getZLog(ctx)
	zLog.Debug("entered PersistUpdateCustomerById")

	allowedFields := map[string]bool{
		"first_name":   true,
		"last_name":    true,
		"email":        true,
		"phone_number": true,
	}

	query := "UPDATE customers SET "
	args := []any{}
	argPosition := 1

	// range/loop thru all updates and populate 'args' slice/array and construct the query using placeholders
	// iteration 1: args[value1] -> "UPDATE customers SET key1 = $1"
	// iteration 2: args[value1, value2] -> "UPDATE customers SET key1 = $1, key2 = $2"
	// iteration 3: args[value1, value2, value3] -> "UPDATE customers SET key1 = $1, key2 = $2, key3 = $3"
	for field, value := range updates {
		if !allowedFields[field] {
			zLog.Error("Attempted to update invalid field", zap.String("field", field))
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

	// "args..." will inject the values into the placeholders in the query
	_, err := cp.DbHandle.ExecContext(ctx, query, args...)
	if err != nil {
		zLog.Error("ExecContext failed for PersistUpdateCustomerById", zap.Error(err))
		return err
	}
	return nil
}

func (cp CustomerPersistence) getZLog(ctx context.Context) *zap.Logger {
	return utils.FromContext(ctx, cp.Logger)
}
