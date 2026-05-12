package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jshelley8117/CodeCart/internal/common"
	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/persistence"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

type LocationService struct {
	LocationPersistence persistence.LocationPersistence
}

func NewLocationService(locationPersistence persistence.LocationPersistence) LocationService {
	return LocationService{
		LocationPersistence: locationPersistence,
	}
}

func (ls LocationService) CreateLocation(ctx context.Context, request model.CreateLocationRequest) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("entered LocationService CreateLocation")

	if err := ls.LocationPersistence.PersistCreateLocation(ctx, model.Location{
		StoreName:    request.StoreName,
		StoreAddress: request.StoreAddress,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}); err != nil {
		z.Error("persistence invocation failed", zap.Error(err))
		return fmt.Errorf(common.ERR_CLIENT_DB_PERSISTENCE_FAIL)
	}

	return nil
}

func (ls LocationService) FetchAllLocations(ctx context.Context, page, pageSize int) ([]model.Location, int64, error) {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("entered LocationService FetchAllLocations")

	locationRows, total, err := ls.LocationPersistence.FetchAllLocations(ctx, page, pageSize)
	if err != nil {
		z.Error("persistence invocation failed", zap.Error(err))
		return nil, 0, err
	}
	defer locationRows.Close()

	locations := make([]model.Location, 0)

	for locationRows.Next() {
		var location model.Location
		if err := locationRows.Scan(
			&location.ID,
			&location.StoreName,
			&location.StoreAddress,
			&location.CreatedAt,
			&location.UpdatedAt,
		); err != nil {
			z.Error("scan operation failed", zap.Error(err))
			return nil, 0, err
		}
		locations = append(locations, location)
	}

	if err := locationRows.Err(); err != nil {
		z.Error("error occurred while iterating through sql rows", zap.Error(err))
		return nil, 0, err
	}

	return locations, total, nil
}

func (ls LocationService) FetchLocationById(ctx context.Context, id int) (model.Location, error) {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("entered LocationService FetchLocationById")

	var location model.Location
	row := ls.LocationPersistence.FetchLocationById(ctx, id)

	if err := row.Scan(
		&location.ID,
		&location.StoreName,
		&location.StoreAddress,
		&location.CreatedAt,
		&location.UpdatedAt,
	); err != nil {
		z.Error("scan operation failed", zap.Error(err))
		return model.Location{}, err
	}

	return location, nil
}

func (ls LocationService) UpdateLocationById(ctx context.Context, id int, request model.UpdateLocationRequest) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("entered LocationService UpdateLocationById")

	updates := make(map[string]any)

	if request.StoreName != "" {
		updates["store_name"] = request.StoreName
	}
	if len(request.StoreAddress) > 0 {
		updates["store_address"] = request.StoreAddress
	}

	if len(updates) == 0 {
		z.Error("no updates found", zap.Int("location_id", id))
		return fmt.Errorf("no updates found")
	}

	if err := ls.LocationPersistence.PersistUpdateLocationById(ctx, id, updates); err != nil {
		z.Error("persistence invocation failed", zap.Error(err))
		return err
	}

	return nil
}

func (ls LocationService) DeleteLocationById(ctx context.Context, id int) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("entered LocationService DeleteLocationById")

	if err := ls.LocationPersistence.PersistDeleteLocationById(ctx, id); err != nil {
		z.Error("persistence invocation failed", zap.Error(err))
		return err
	}

	return nil
}
