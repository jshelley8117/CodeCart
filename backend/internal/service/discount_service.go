package service

import (
	"context"
	"strings"

	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/persistence"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

type DiscountService struct {
	DiscountPersistence persistence.DiscountPersistence
}

func NewDiscountService(discountPersistence persistence.DiscountPersistence) DiscountService {
	return DiscountService{
		DiscountPersistence: discountPersistence,
	}
}

func (ds DiscountService) CreateDiscount(ctx context.Context, request model.CreateDiscount) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("entered CreateDiscount")

	discountDomainModel := model.Discount{
		Code:        strings.ToUpper(request.Code),
		Value:       request.Value,
		Description: request.Description,
		Type:        request.Type,
		IsActive:    request.IsActive,
		StartsAt:    request.StartsAt,
		ExpiresAt:   request.ExpiresAt,
		UsageLimit:  request.UsageLimit,
	}

	if err := ds.DiscountPersistence.PersistCreateDiscount(ctx, discountDomainModel); err != nil {
		z.Error("persistence invocation failed", zap.Error(err))
		return err
	}

	return nil
}
