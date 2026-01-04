package service

import (
	"context"
	"strings"
	"time"

	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/persistence"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

type CustomerService struct {
	CustomerPersistence persistence.CustomerPersistence
	Logger              *zap.Logger
}

func NewCustomerService(customerPersistence persistence.CustomerPersistence, logger *zap.Logger) CustomerService {
	return CustomerService{
		CustomerPersistence: customerPersistence,
		Logger:              logger,
	}
}

func (cs CustomerService) CreateCustomer(ctx context.Context, request model.CreateCustomerRequest) error {
	zLog := utils.FromContext(ctx, cs.Logger).Named("customer_service")
	zLog.Debug("entered CustomerService")

	if err := cs.CustomerPersistence.PersistCreateCustomer(ctx, model.Customer{
		FirstName:   strings.ToLower(request.FirstName),
		LastName:    strings.ToLower(request.LastName),
		PhoneNumber: request.PhoneNumber,
		Email:       strings.ToLower(request.Email),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}); err != nil {
		zLog.Error("persistence invocation failed: %w", zap.Error(err))
		return err
	}

	return nil
}
