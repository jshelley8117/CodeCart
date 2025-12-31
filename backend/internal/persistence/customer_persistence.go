package persistence

import (
	"context"
	"database/sql"

	"github.com/jshelley8117/CodeCart/internal/model"
)

type CustomerPersistence struct {
	DbHandle *sql.DB
}

func NewCustomerPersistence(dbHandle *sql.DB) CustomerPersistence {
	return CustomerPersistence{
		DbHandle: dbHandle,
	}
}

func (cp CustomerPersistence) PersistCreateCustomer(ctx context.Context, customerDomain model.Customer) error {
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
		return err
	}

	return nil
}
