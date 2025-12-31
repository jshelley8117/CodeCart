package service

import (
	"context"
	"strings"
	"time"

	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/persistence"
)

type CustomerService struct {
	CustomerPersistence persistence.CustomerPersistence
}

func NewCustomerService(customerPersistence persistence.CustomerPersistence) CustomerService {
	return CustomerService{
		CustomerPersistence: customerPersistence,
	}
}

func (cs CustomerService) CreateCustomer(ctx context.Context, request model.CreateCustomerRequest) error {

	if err := cs.CustomerPersistence.PersistCreateCustomer(ctx, model.Customer{
		FirstName:   strings.ToLower(request.FirstName),
		LastName:    strings.ToLower(request.LastName),
		PhoneNumber: request.PhoneNumber,
		Email:       strings.ToLower(request.Email),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}); err != nil {
		return err
	}

	return nil
}
