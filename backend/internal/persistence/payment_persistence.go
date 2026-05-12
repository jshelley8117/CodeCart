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

type PaymentPersistence struct {
	DbHandle *sql.DB
}

func NewPaymentPersistence(dbHandle *sql.DB) PaymentPersistence {
	return PaymentPersistence{
		DbHandle: dbHandle,
	}
}

func (pp PaymentPersistence) PersistCreatePayment(ctx context.Context, paymentDomain model.Payment, authId string) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered PersistCreatePayment")

	query := `
		INSERT INTO payments (order_id, payment_method, status, currency, payment_intent_id, created_at, updated_at)
		SELECT o.id, $2, $3, $4, $5, $6, $7
		FROM orders o
		JOIN users u ON o.customer_id = u.customer_id
		WHERE o.id = $1 AND u.auth_id = $8
	`

	_, err := pp.DbHandle.ExecContext(
		ctx,
		query,
		paymentDomain.OrderId,
		paymentDomain.PaymentMethod,
		paymentDomain.Status,
		paymentDomain.Currency,
		paymentDomain.PaymentIntentId,
		paymentDomain.CreatedAt,
		paymentDomain.UpdatedAt,
		authId,
	)
	if err != nil {
		z.Error("ExecContext failed for PersistCreatePayment", zap.Error(err))
		return err
	}
	return nil
}

func (pp PaymentPersistence) FetchPaymentById(ctx context.Context, id int, authId string) *sql.Row {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered FetchPaymentById")

	query := `
		SELECT p.id, p.order_id, p.payment_method, p.status, p.currency, p.payment_intent_id, p.created_at, p.updated_at
		FROM payments p
		JOIN orders o ON p.order_id = o.id
		JOIN users u ON o.customer_id = u.customer_id
		WHERE p.id = $1 AND u.auth_id = $2
	`

	return pp.DbHandle.QueryRowContext(ctx, query, id, authId)
}

func (pp PaymentPersistence) PersistUpdatePaymentById(ctx context.Context, id int, authId string, updates map[string]any) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered PersistUpdatePaymentById")

	allowedFields := map[string]bool{
		"payment_method":    true,
		"status":            true,
		"currency":          true,
		"payment_intent_id": true,
	}

	query := "UPDATE payments SET"
	args := []any{}
	argPosition := 1

	for field, value := range updates {
		if !allowedFields[field] {
			z.Error("Attempted to update invalid field", zap.String("field", field))
			return fmt.Errorf("invalid field: %s", field)
		}

		if argPosition > 1 {
			query += ","
		}
		query += " " + field + " = $" + fmt.Sprintf("%d", argPosition)
		args = append(args, value)
		argPosition++
	}

	query += ", updated_at = $" + fmt.Sprintf("%d", argPosition)
	args = append(args, time.Now())
	argPosition++

	query += " WHERE id = $" + fmt.Sprintf("%d", argPosition)
	args = append(args, id)
	argPosition++

	query += " AND order_id IN (SELECT o.id FROM orders o JOIN users u ON o.customer_id = u.customer_id WHERE u.auth_id = $" + fmt.Sprintf("%d", argPosition) + ")"
	args = append(args, authId)

	_, err := pp.DbHandle.ExecContext(ctx, query, args...)
	if err != nil {
		z.Error("ExecContext failed for PersistUpdatePaymentById", zap.Error(err))
		return err
	}
	return nil
}

func (pp PaymentPersistence) PersistDeletePaymentById(ctx context.Context, id int, authId string) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered PersistDeletePaymentById")

	query := `
		DELETE FROM payments
		WHERE id = $1
		AND order_id IN (
			SELECT o.id FROM orders o
			JOIN users u ON o.customer_id = u.customer_id
			WHERE u.auth_id = $2
		)
	`

	if _, err := pp.DbHandle.ExecContext(ctx, query, id, authId); err != nil {
		z.Error("ExecContext failed for PersistDeletePaymentById", zap.Error(err))
		return err
	}

	return nil
}
