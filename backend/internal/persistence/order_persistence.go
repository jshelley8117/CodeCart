package persistence

import (
	"context"
	"database/sql"
	"log"

	"github.com/jshelley8117/CodeCart/internal/model"
)

type OrderPersistence struct {
	DbHandle *sql.DB
}

func NewOrderPersistence(dbHandle *sql.DB) OrderPersistence {
	return OrderPersistence{
		DbHandle: dbHandle,
	}
}

func (op OrderPersistence) PersistCreateOrder(ctx context.Context, orderDomain model.Order) error {
	log.Println("Entereed PersistCreateOrder")

	query := `
		INSERT INTO orders (customer_id, status, total_price, delivery_address, created_at, updated_at, address_id, "orderType")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	log.Printf("DEBUG: Attempting to save Order with AddressId: %v", orderDomain.AddressId)
	_, err := op.DbHandle.ExecContext(
		ctx,
		query,
		orderDomain.CustomerId,
		orderDomain.Status,
		orderDomain.TotalPrice,
		orderDomain.DeliveryAddress,
		orderDomain.CreatedAt,
		orderDomain.UpdatedAt,
		orderDomain.AddressId,
		orderDomain.OrderType,
	)
	if err != nil {
		log.Printf("Error in PersistCreateOrder, %s", err)
		return err
	}
	return nil
}
