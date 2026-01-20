package service

import (
	"context"
	"fmt"
	"log"
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
	Logger             *zap.Logger
}

func NewAddressService(addressPersistence persistence.AddressPersistence) AddressService {
	return AddressService{
		AddressPersistence: addressPersistence,
	}
}

func (as AddressService) CreateAddress(ctx context.Context, request model.CreateAddressRequest) error {
	log.Println("Entered CreateAddress")
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
	zLog := as.getZLog(ctx)
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

func (as AddressService) getZLog(ctx context.Context) *zap.Logger {
	return utils.FromContext(ctx, as.Logger)
}
