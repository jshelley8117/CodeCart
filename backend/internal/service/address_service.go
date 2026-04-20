package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jshelley8117/CodeCart/internal/common"
	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/persistence"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

type AddressService struct {
	AddressPersistence persistence.AddressPersistence
}

func NewAddressService(addressPersistence persistence.AddressPersistence) AddressService {
	return AddressService{
		AddressPersistence: addressPersistence,
	}
}

func (as AddressService) CreateAddress(ctx context.Context, request model.CreateAddressRequest) error {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("entered CreateAddress")

	addressDomainModel := model.Address{
		StreetAddress: strings.ToLower(request.StreetAddress),
		City:          strings.ToLower(request.City),
		State:         strings.ToLower(request.State),
		ZipCode:       strings.ToLower(request.ZipCode),
		Country:       strings.ToLower(request.Country),
		UserId:        request.UserId,
		IsDefault:     false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := as.AddressPersistence.PersistCreateAddress(ctx, addressDomainModel); err != nil {
		return err
	}

	return nil
}

func (as AddressService) GetAllAddresses(ctx context.Context) ([]model.Address, error) {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("Entered GetAllAddresses")

	addressRows, err := as.AddressPersistence.FetchAllAddresses(ctx)
	if err != nil {
		zLog.Error("persistence invocation failed", zap.Error(err))
		return nil, fmt.Errorf(common.ERR_CLIENT_DB_RETRIEVAL_FAIL)
	}
	defer addressRows.Close()

	addresses := make([]model.Address, 0)

	for addressRows.Next() {
		var addr model.Address
		if err := addressRows.Scan(
			&addr.StreetAddress,
			&addr.City,
			&addr.State,
			&addr.ZipCode,
			&addr.Country,
			&addr.UserId,
			&addr.Id,
			&addr.IsDefault,
			&addr.CreatedAt,
			&addr.UpdatedAt,
		); err != nil {
			zLog.Error("scan operation failed", zap.Error(err))
			return nil, fmt.Errorf(common.ERR_CLIENT_DB_RETRIEVAL_FAIL)
		}
		addresses = append(addresses, addr)
	}

	if err := addressRows.Err(); err != nil {
		zLog.Error("error occured while iterating through sql rows", zap.Error(err))
		return nil, fmt.Errorf(common.ERR_CLIENT_DB_RETRIEVAL_FAIL)
	}

	return addresses, nil
}

func (as AddressService) GetAddressById(ctx context.Context, id int) (model.Address, error) {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("Entered GetAddressById")

	addressRow := as.AddressPersistence.FetchAddressById(ctx, id)
	if addressRow == nil {
		zLog.Warn("order not found", zap.Int("order_id", id))
		return model.Address{}, nil
	}

	var address model.Address
	if err := addressRow.Scan(
		&address.StreetAddress,
		&address.City,
		&address.State,
		&address.ZipCode,
		&address.Country,
		&address.UserId,
		&address.Id,
		&address.IsDefault,
		&address.CreatedAt,
		&address.UpdatedAt,
	); err != nil {
		zLog.Error("scan operation failed", zap.Error(err))
		return model.Address{}, err
	}

	return address, nil
}

func (as AddressService) UpdateAddressById(ctx context.Context, request model.UpdateAddressRequest, id int) error {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("Entered UpdateAddressById")

	updates := make(map[string]any)

	if request.StreetAddress != "" {
		updates["street_address"] = request.StreetAddress
	}
	if request.City != "" {
		updates["city"] = request.City
	}
	if request.State != "" {
		updates["state"] = request.State
	}
	if request.ZipCode != "" {
		updates["zip_code"] = request.ZipCode
	}
	if request.Country != "" {
		updates["country"] = request.Country
	}
	if request.IsDefault != nil {
		updates["is_default"] = request.IsDefault
	}

	if len(updates) == 0 {
		zLog.Error("No updates found", zap.Int("address_id", id))
		return fmt.Errorf("no updates found")
	}

	if err := as.AddressPersistence.PersistUpdateAddressById(ctx, id, updates); err != nil {
		zLog.Error("persistence invocation failed", zap.Error(err))
		return err
	}

	return nil
}

func (as AddressService) DeleteAddressById(ctx context.Context, id int) error {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("Entered DeleteAddressById")

	if err := as.AddressPersistence.PersistDeleteAddressById(ctx, id); err != nil {
		zLog.Error("persistence invocation failed", zap.Error(err))
		return fmt.Errorf(common.ERR_CLIENT_DB_DELETE_FAIL)
	}
	return nil
}
