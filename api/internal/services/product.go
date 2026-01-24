package services

import (
	"context"
	"database/sql"
	"errors"
	"ukoni/internal/models"
)

type ProductService struct {
	DB                 *sql.DB
	ProductModel       *models.ProductModel
	ActivityLogService *ActivityLogService
}

// Categories

func (s *ProductService) CreateCategory(ctx context.Context, userID string, name string, parentID *string) (*models.ProductCategory, error) {
	if name == "" {
		return nil, errors.New("category name cannot be empty")
	}

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	category := &models.ProductCategory{
		Name:             name,
		ParentCategoryID: parentID,
	}

	if err := s.ProductModel.CreateCategory(ctx, tx, category); err != nil {
		return nil, err
	}

	if s.ActivityLogService != nil {
		if err := s.ActivityLogService.LogActivity(ctx, tx, nil, &userID, "product_category.created", "product_category", &category.ID, nil); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *ProductService) ListCategories(ctx context.Context) ([]*models.ProductCategory, error) {
	return s.ProductModel.ListCategories(ctx)
}

// Canonical Products

func (s *ProductService) CreateCanonicalProduct(ctx context.Context, userID string, name string, description *string, categoryID *string) (*models.CanonicalProduct, error) {
	if name == "" {
		return nil, errors.New("canonical product name cannot be empty")
	}

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	cp := &models.CanonicalProduct{
		Name:        name,
		Description: description,
		CategoryID:  categoryID,
	}

	if err := s.ProductModel.CreateCanonicalProduct(ctx, tx, cp); err != nil {
		return nil, err
	}

	if s.ActivityLogService != nil {
		if err := s.ActivityLogService.LogActivity(ctx, tx, nil, &userID, "canonical_product.created", "canonical_product", &cp.ID, nil); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return cp, nil
}

func (s *ProductService) ListCanonicalProducts(ctx context.Context) ([]*models.CanonicalProduct, error) {
	return s.ProductModel.ListCanonicalProducts(ctx)
}

// Products

func (s *ProductService) CreateProduct(ctx context.Context, userID string, canonicalProductID *string, brand *string, name string, description *string, categoryID *string) (*models.Product, error) {
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	p := &models.Product{
		CanonicalProductID: canonicalProductID,
		Brand:              brand,
		Name:               name,
		Description:        description,
		CategoryID:         categoryID,
	}

	if err := s.ProductModel.CreateProduct(ctx, tx, p); err != nil {
		return nil, err
	}

	if s.ActivityLogService != nil {
		if err := s.ActivityLogService.LogActivity(ctx, tx, nil, &userID, "product.created", "product", &p.ID, nil); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *ProductService) ListProducts(ctx context.Context) ([]*models.Product, error) {
	return s.ProductModel.ListProducts(ctx)
}

// Product Variants

func (s *ProductService) CreateProductVariant(ctx context.Context, userID string, productID string, variantName string, sku *string, unit *string) (*models.ProductVariant, error) {
	if variantName == "" {
		return nil, errors.New("variant name cannot be empty")
	}
	if productID == "" {
		return nil, errors.New("product id cannot be empty")
	}

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	pv := &models.ProductVariant{
		ProductID:   productID,
		VariantName: variantName,
		SKU:         sku,
		Unit:        unit,
	}

	if err := s.ProductModel.CreateProductVariant(ctx, tx, pv); err != nil {
		return nil, err
	}

	if s.ActivityLogService != nil {
		if err := s.ActivityLogService.LogActivity(ctx, tx, nil, &userID, "product_variant.created", "product_variant", &pv.ID, nil); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return pv, nil
}

func (s *ProductService) ListProductVariants(ctx context.Context, productID string) ([]*models.ProductVariant, error) {
	return s.ProductModel.ListProductVariants(ctx, productID)
}
