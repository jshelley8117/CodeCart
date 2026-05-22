package service

import (
	"context"
	"fmt"
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

func (ds DiscountService) GetDiscountByID(ctx context.Context, id int) (model.Discount, error) {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entererd GetDiscountById")

	row := ds.DiscountPersistence.FetchDiscountById(ctx, id)

	var discount model.Discount
	if err := row.Scan(
		&discount.Id,
		&discount.Code,
		&discount.Value,
		&discount.Description,
		&discount.Type,
		&discount.IsActive,
		&discount.StartsAt,
		&discount.ExpiresAt,
		&discount.UsageLimit,
		&discount.TimesUsed,
		&discount.CreatedAt,
		&discount.UpdatedAt,
	); err != nil {
		z.Error("scan operation failed", zap.Error(err))
		return model.Discount{}, err
	}

	return discount, nil
}

func (ds DiscountService) UpdateDiscountById(ctx context.Context, id int, request model.UpdateDiscountRequest) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("Entered UpdateDiscountById")

	updates := make(map[string]any)

	if request.Code != nil {
		updates["code"] = *request.Code
	}
	if request.Value != nil {
		updates["value"] = *request.Value
	}
	if request.Description != nil {
		updates["description"] = *request.Description
	}
	if request.Type != nil {
		updates["type"] = *request.Type
	}
	if request.IsActive != nil {
		updates["is_active"] = *request.IsActive
	}
	if request.StartsAt != nil {
		updates["starts_at"] = *request.StartsAt
	}
	if request.ExpiresAt != nil {
		updates["expires_at"] = *request.ExpiresAt
	}
	if request.UsageLimit != nil {
		updates["usage_limit"] = *request.UsageLimit
	}

	if len(updates) == 0 {
		z.Error("no updates found", zap.Int("id", id))
		return fmt.Errorf("no updates found")
	}

	if err := ds.DiscountPersistence.PersistUpdateDiscountById(ctx, id, updates); err != nil {
		z.Error("persistence invocation failed", zap.Error(err))
		return err
	}

	return nil
}

func (ds DiscountService) DeleteDiscountbyID(ctx context.Context, id int) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("entered DeleteDiscount")

	if err := ds.DiscountPersistence.PersistDeleteDiscountById(ctx, id); err != nil {
		z.Error("persistence invocation failed", zap.Error(err))
		return err
	}

	return nil
}
