package service

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/persistence"
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
	log.Println("Entered CreateAddress")
	addressDomainModel := model.Address{
		Id:            "1234",
		StreetAddress: strings.ToLower(request.StreetAddress),
		City:          strings.ToLower(request.City),
		State:         strings.ToLower(request.State),
		ZipCode:       strings.ToLower(request.ZipCode),
		Country:       strings.ToLower(request.Country),
		UserId:        request.UserId,
		IsDefault:     true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := as.AddressPersistence.PersistCreateAddress(ctx, addressDomainModel); err != nil {
		return err
	}

	return nil
}
