package service

import (
	"context"
	"fmt"
	"strconv"
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
		Logger:              logger.Named("customer_service"),
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
		zLog.Error("persistence invocation failed", zap.Error(err))
		return err
	}

	return nil
}

func (cs CustomerService) GetAllCustomers(ctx context.Context) ([]model.Customer, error) {
	zLog := cs.getZLog(ctx)
	zLog.Debug("entered GetAllCustomers")

	customerRows, err := cs.CustomerPersistence.FetchAllCustomersById(ctx)
	if err != nil {
		zLog.Error("persistence invocation failed", zap.Error(err))
		return nil, err
	}
	defer customerRows.Close()

	customers := make([]model.Customer, 0)

	for customerRows.Next() {
		var cust model.Customer
		if err := customerRows.Scan(
			&cust.Id,
			&cust.FirstName,
			&cust.LastName,
			&cust.PhoneNumber,
			&cust.Email,
			&cust.CreatedAt,
			&cust.UpdatedAt,
		); err != nil {
			zLog.Error("scan operation failed", zap.Error(err))
			return nil, err
		}
		customers = append(customers, cust)
	}

	if err := customerRows.Err(); err != nil {
		zLog.Error("error occured while iterating through sql rows", zap.Error(err))
		return nil, err
	}

	return customers, nil
}

func (cs CustomerService) DeleteCustomerById(ctx context.Context, id int) error {
	zLog := cs.getZLog(ctx)
	zLog.Debug("entered DeleteCustomerById")

	if err := cs.CustomerPersistence.PersistDeleteCustomerById(ctx, id); err != nil {
		zLog.Error("persistence invocation failed", zap.Error(err))
		return err
	}
	return nil
}

func (cs CustomerService) UpdateCustomerById(ctx context.Context, request model.UpdateCustomerRequest, id int) error {
	zLog := cs.getZLog(ctx)
	zLog.Debug("entered UpdateCustomerById")

	updates := make(map[string]any)

	if request.FirstName != "" {
		updates["first_name"] = strings.ToLower(request.FirstName)
	}
	if request.LastName != "" {
		updates["last_name"] = strings.ToLower(request.LastName)
	}
	if request.PhoneNumber != nil && *request.PhoneNumber != "" {
		updates["phone_number"] = strings.ToLower(*request.PhoneNumber)
	}
	if request.Email != nil && *request.Email != "" {
		updates["email"] = strings.ToLower(*request.Email)
	}

	if len(updates) == 0 {
		zLog.Error("No updates found", zap.String("customer_id", strconv.Itoa(id)))
		return fmt.Errorf("No updates found")
	}

	if err := cs.CustomerPersistence.PersistUpdateCustomerById(ctx, id, updates); err != nil {
		zLog.Error("persistence invocation failed", zap.Error(err))
		return err
	}

	return nil
}

func (cs CustomerService) getZLog(ctx context.Context) *zap.Logger {
	return utils.FromContext(ctx, cs.Logger)
}
