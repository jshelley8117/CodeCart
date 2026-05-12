package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/persistence"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

var validPaymentStatuses = map[string]bool{
	"PENDING":  true,
	"SUCCESS":  true,
	"FAILED":   true,
	"REFUNDED": true,
	"CANCELED": true,
}

type PaymentService struct {
	PaymentPersistence persistence.PaymentPersistence
}

func NewPaymentService(paymentPersistence persistence.PaymentPersistence) PaymentService {
	return PaymentService{
		PaymentPersistence: paymentPersistence,
	}
}

func (ps PaymentService) CreatePayment(ctx context.Context, request model.CreatePaymentRequest, authId string) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered CreatePayment")

	paymentDomain := model.Payment{
		OrderId:         request.OrderId,
		PaymentMethod:   request.PaymentMethod,
		Status:          "PENDING",
		Currency:        request.Currency,
		PaymentIntentId: request.PaymentIntentId,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := ps.PaymentPersistence.PersistCreatePayment(ctx, paymentDomain, authId); err != nil {
		z.Error("persistence invocation failed", zap.Error(err))
		return err
	}
	return nil
}

func (ps PaymentService) GetPaymentById(ctx context.Context, id int, authId string) (model.Payment, error) {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered GetPaymentById")

	row := ps.PaymentPersistence.FetchPaymentById(ctx, id, authId)
	if row == nil {
		z.Warn("payment not found", zap.Int("payment_id", id))
		return model.Payment{}, nil
	}

	var payment model.Payment
	if err := row.Scan(
		&payment.Id,
		&payment.OrderId,
		&payment.PaymentMethod,
		&payment.Status,
		&payment.Currency,
		&payment.PaymentIntentId,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	); err != nil {
		z.Error("scan operation failed", zap.Error(err))
		return model.Payment{}, err
	}

	return payment, nil
}

func (ps PaymentService) UpdatePaymentById(ctx context.Context, request model.UpdatePaymentRequest, id int, authId string) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered UpdatePaymentById")

	updates := make(map[string]any)

	if request.PaymentMethod != "" {
		updates["payment_method"] = request.PaymentMethod
	}
	if request.Status != "" {
		if !validPaymentStatuses[request.Status] {
			z.Error("invalid payment status", zap.String("status", request.Status))
			return fmt.Errorf("invalid payment status: %s", request.Status)
		}
		updates["status"] = request.Status
	}
	if request.Currency != "" {
		updates["currency"] = request.Currency
	}
	if request.PaymentIntentId != "" {
		updates["payment_intent_id"] = request.PaymentIntentId
	}

	if len(updates) == 0 {
		z.Error("no updates provided", zap.Int("payment_id", id))
		return fmt.Errorf("no updates provided")
	}

	if err := ps.PaymentPersistence.PersistUpdatePaymentById(ctx, id, authId, updates); err != nil {
		z.Error("persistence invocation failed", zap.Error(err))
		return err
	}

	return nil
}

func (ps PaymentService) DeletePaymentById(ctx context.Context, id int, authId string) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered DeletePaymentById")

	if err := ps.PaymentPersistence.PersistDeletePaymentById(ctx, id, authId); err != nil {
		z.Error("persistence invocation failed", zap.Error(err))
		return err
	}

	return nil
}
