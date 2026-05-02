package service

import (
	"context"

	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/persistence"
)

type InventoryService struct {
	InventoryPersistence persistence.InventoryPersistence
}

func NewInventoryService(inventoryPersistence persistence.InventoryPersistence) InventoryService {
	return InventoryService{
		InventoryPersistence: inventoryPersistence,
	}
}

func (is InventoryService) CreateInventory(ctx context.Context, request model.CreateInventoryRequest) error

func (is InventoryService) GetAllInventory(ctx context.Context) ([]model.Inventory, error)

func (is InventoryService) GetInventoryById()

func (is InventoryService) UpdateInventoryByI()

func (is InventoryService) DeleteInventoryById()
