package persistence

import (
	"context"
	"database/sql"
	"log"

	"github.com/jshelley8117/CodeCart/internal/model"
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
	log.Println("Entered PersistCreateAddress")
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
