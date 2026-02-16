package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/persistence"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

type ProductService struct {
	ProductPersistence persistence.ProductPersistence
}

func NewProductService(productPersistence persistence.ProductPersistence) ProductService {
	return ProductService{
		ProductPersistence: productPersistence,
	}
}

func (ps ProductService) CreateProduct(ctx context.Context, request model.CreateProductRequest) error {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("entered ServiceCreateProduct")

	productDomainModel := model.Product{
		Name:            request.Name,
		Description:     request.Description,
		UnitPrice:       request.UnitPrice,
		Category:        request.Category,
		Brand:           request.Brand,
		IsAgeRestricted: *request.IsAgeRestricted,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		IsActive:        true,
	}

	err := ps.ProductPersistence.PersistCreateProduct(ctx, productDomainModel)
	if err != nil {
		zLog.Error("persistance invocation failed", zap.Error(err))
		return err
	}

	return nil
}

func (ps ProductService) FetchAllProducts(ctx context.Context, page, pageSize int) ([]model.Product, int64, error) {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("entered ServiceFetchAllProducts")

	productRows, total, err := ps.ProductPersistence.FetchAllProducts(ctx, page, pageSize)
	if err != nil {
		zLog.Error("persistence invocation failed", zap.Error(err))
		return nil, 0, err
	}
	defer productRows.Close()

	products := make([]model.Product, 0)

	for productRows.Next() {
		var product model.Product
		if err := productRows.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&product.UnitPrice,
			&product.Category,
			&product.Brand,
			&product.IsAgeRestricted,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			zLog.Error("scan operation failed", zap.Error(err))
			return nil, 0, err
		}
		products = append(products, product)
	}

	if err := productRows.Err(); err != nil {
		zLog.Error("error occurred while iterating through sql rows", zap.Error(err))
		return nil, 0, err
	}

	return products, total, nil
}

func (ps ProductService) ServiceFetchProductById(ctx context.Context, id int) (model.Product, error) {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("entered ServiceFetchProductById")

	var product model.Product
	row := ps.ProductPersistence.FetchProductById(ctx, id)

	if err := row.Scan(
		&product.Id,
		&product.Name,
		&product.Description,
		&product.UnitPrice,
		&product.Category,
		&product.Brand,
		&product.IsAgeRestricted,
		&product.CreatedAt,
		&product.UpdatedAt,
	); err != nil {
		zLog.Error("scan operation failed", zap.Error(err))
		return model.Product{}, err
	}

	return product, nil
}

func (ps ProductService) UpdateProductById(ctx context.Context, id int, request model.UpdateProductRequest) error {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("entered ServiceUpdateProductById")

	updates := make(map[string]any)

	if request.Name != "" {
		updates["name"] = request.Name
	}
	if request.Description != "" {
		updates["description"] = request.Description
	}
	if request.UnitPrice != 0 {
		updates["unit_price"] = request.UnitPrice
	}
	if request.Category != "" {
		updates["category"] = request.Category
	}
	if request.Brand != "" {
		updates["brand"] = request.Brand
	}
	if request.IsAgeRestricted {
		updates["age_restricted"] = request.IsAgeRestricted
	}

	if len(updates) == 0 {
		zLog.Error("no updates found", zap.Int("product_id", id))
		return fmt.Errorf("no updates found")
	}

	if err := ps.ProductPersistence.PersistUpdateProductById(ctx, id, updates); err != nil {
		zLog.Error("persistence invocation failed", zap.Error(err))
		return err
	}

	return nil
}

func (ps ProductService) FetchAllProductVariantsByProductId(ctx context.Context, productId, page, pageSize int) ([]model.ProductVariant, int64, error) {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("entered ServiceFetchAllProductVariantsByProductId")

	variantRows, total, err := ps.ProductPersistence.FetchAllProductVariantsByProductId(ctx, productId, page, pageSize)
	if err != nil {
		zLog.Error("persistence invocation failed", zap.Error(err))
		return nil, 0, err
	}
	defer variantRows.Close()

	variants := make([]model.ProductVariant, 0)

	for variantRows.Next() {
		var variant model.ProductVariant
		var sku, size, flavor, imagePath sql.NullString
		if err := variantRows.Scan(
			&variant.Id,
			&sku,
			&size,
			&flavor,
			&variant.IsActive,
			&variant.CreatedAt,
			&variant.UpdatedAt,
			&imagePath,
			&variant.ProductId,
		); err != nil {
			zLog.Error("scan operation failed", zap.Error(err))
			return nil, 0, err
		}
		variant.Sku = sku.String
		variant.Size = size.String
		variant.Flavor = flavor.String
		variant.ImagePath = imagePath.String
		variants = append(variants, variant)
	}

	if err := variantRows.Err(); err != nil {
		zLog.Error("error occurred while iterating through sql rows", zap.Error(err))
		return nil, 0, err
	}

	return variants, total, nil
}

func (ps ProductService) UpdateProductVariantById(ctx context.Context, id int, request model.UpdateProductVariant) error {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("entered UpdateProductVariantById")

	updates := make(map[string]any)

	if request.Size != "" {
		updates["size"] = request.Size
	}
	if request.Flavor != "" {
		updates["flavor"] = request.Flavor
	}
	if request.ImagePath != "" {
		updates["image_path"] = request.ImagePath
	}
	if request.IsActive {
		updates["is_active"] = request.IsActive
	}

	if len(updates) == 0 {
		zLog.Error("no updates found", zap.Int("variant_id", id))
		return fmt.Errorf("no updates found")
	}

	if err := ps.ProductPersistence.PersistUpdateProductVariantById(ctx, id, updates); err != nil {
		zLog.Error("persistence invocation failed", zap.Error(err))
		return err
	}

	return nil
}

func (ps ProductService) DeleteProductById(ctx context.Context, id int) error {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("entered DeleteProductById")

	if err := ps.ProductPersistence.PersistDeleteVariantsByProductId(ctx, id); err != nil {
		zLog.Error("persistence invocation failed", zap.Error(err))
		return err
	}

	if err := ps.ProductPersistence.PersistDeleteProductById(ctx, id); err != nil {
		zLog.Error("persistence invocation failed", zap.Error(err))
		return err
	}

	return nil
}

func (ps ProductService) DeleteProductVariantById(ctx context.Context, id int) error {
	zLog := utils.FromContext(ctx, zap.NewNop())
	zLog.Debug("entered DeleteProductVariantById")

	if err := ps.ProductPersistence.PersistDeleteProductVariantById(ctx, id); err != nil {
		zLog.Error("persistence invocation failed", zap.Error(err))
		return err
	}

	return nil
}
