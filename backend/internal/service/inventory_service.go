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

type InventoryService struct {
	InventoryPersistence persistence.InventoryPersistence
}

func NewInventoryService(inventoryPersistence persistence.InventoryPersistence) InventoryService {
	return InventoryService{
		InventoryPersistence: inventoryPersistence,
	}
}

func (is InventoryService) CreateInventory(ctx context.Context, request model.CreateInventoryRequest) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("entered CreateInventory")

	inventoryDomainModel := model.Inventory{
		ProductVariantId: request.ProductVariantId,
		LocationId:       request.LocationId,
		Quantity:         request.Quantity,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := is.InventoryPersistence.PersistCreateInventory(ctx, inventoryDomainModel); err != nil {
		z.Error("persistence invocation failed", zap.Error(err))
		return err
	}
	return nil
}

func (is InventoryService) GetAllInventory(ctx context.Context, page, pageSize int) ([]model.Inventory, int64, error) {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("entered GetAllInventory")

	inventoryRows, total, err := is.InventoryPersistence.FetchAllInventory(ctx, page, pageSize)
	if err != nil {
		z.Error("persistence invocation failed", zap.Error(err))
		return nil, 0, err
	}
	defer inventoryRows.Close()

	inventory := make([]model.Inventory, 0)

	for inventoryRows.Next() {
		var item model.Inventory
		if err := inventoryRows.Scan(
			&item.Id,
			&item.ProductVariantId,
			&item.LocationId,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			z.Error("scan operation failed", zap.Error(err))
			return nil, 0, err
		}
		inventory = append(inventory, item)
	}

	if err := inventoryRows.Err(); err != nil {
		z.Error("error occurred while iterating through sql rows", zap.Error(err))
		return nil, 0, err
	}

	return inventory, total, nil
}

func (is InventoryService) GetInventoryById(ctx context.Context, id int) (model.Inventory, error) {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("entered GetInventoryById")

	row := is.InventoryPersistence.FetchInventoryById(ctx, id)

	var item model.Inventory
	if err := row.Scan(
		&item.Id,
		&item.ProductVariantId,
		&item.LocationId,
		&item.Quantity,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		z.Error("scan operation failed", zap.Error(err))
		return model.Inventory{}, err
	}

	return item, nil
}

func (is InventoryService) UpdateInventoryById(ctx context.Context, id int, request model.UpdateInventoryRequest) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("entered UpdateInventoryById")

	updates := make(map[string]any)

	if request.Quantity != nil {
		updates["quantity"] = *request.Quantity
	}
	if request.LocationId != nil {
		updates["location_id"] = *request.LocationId
	}

	if len(updates) == 0 {
		z.Error("no updates found", zap.Int("inventory_id", id))
		return fmt.Errorf("no updates found")
	}

	if err := is.InventoryPersistence.PersistUpdateInventoryById(ctx, id, updates); err != nil {
		z.Error("persistence invocation failed", zap.Error(err))
		return err
	}

	return nil
}

func (is InventoryService) DeleteInventoryById(ctx context.Context, id int) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("entered DeleteInventoryById")

	if err := is.InventoryPersistence.PersistDeleteInventoryById(ctx, id); err != nil {
		z.Error("persistence invocation failed", zap.Error(err))
		return err
	}

	return nil
}
