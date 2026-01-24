package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ukoni/internal/models"
)

var (
	ErrInvalidInput = errors.New("invalid input")
)

type ProductService struct {
	DB           *sql.DB
	ProductModel *models.ProductModel
}

func (s *ProductService) CreateProduct(ctx context.Context, brand, name, description, categoryID string) (*models.Product, error) {
	if name == "" {
		return nil, fmt.Errorf("%w: product name is required", ErrInvalidInput)
	}

	product := &models.Product{
		Name: name,
	}
	if brand != "" {
		product.Brand = &brand
	}
	if description != "" {
		product.Description = &description
	}
	if categoryID != "" {
		product.CategoryID = &categoryID
	}

	err := s.ProductModel.Create(ctx, s.DB, product)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id string) (*models.Product, error) {
	return s.ProductModel.GetByID(ctx, id)
}

func (s *ProductService) ListProducts(ctx context.Context, limit, offset int, search string) ([]*models.Product, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.ProductModel.List(ctx, limit, offset, search)
}

func (s *ProductService) CreateVariant(ctx context.Context, productID, variantName, sku, unit string) (*models.ProductVariant, error) {
	if productID == "" {
		return nil, fmt.Errorf("%w: product id is required", ErrInvalidInput)
	}
	if variantName == "" {
		return nil, fmt.Errorf("%w: variant name is required", ErrInvalidInput)
	}

	variant := &models.ProductVariant{
		ProductID:   productID,
		VariantName: variantName,
	}
	if sku != "" {
		variant.SKU = &sku
	}
	if unit != "" {
		variant.Unit = &unit
	}

	err := s.ProductModel.CreateVariant(ctx, s.DB, variant)
	if err != nil {
		return nil, err
	}
	return variant, nil
}

func (s *ProductService) ListVariants(ctx context.Context, productID string) ([]*models.ProductVariant, error) {
	return s.ProductModel.ListVariants(ctx, productID)
}

func (s *ProductService) UpdateProduct(ctx context.Context, id, brand, name, description, categoryID string) (*models.Product, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: product id is required", ErrInvalidInput)
	}
	if name == "" {
		return nil, fmt.Errorf("%w: product name is required", ErrInvalidInput)
	}

	product := &models.Product{
		ID:   id,
		Name: name,
	}
	if brand != "" {
		product.Brand = &brand
	}
	if description != "" {
		product.Description = &description
	}
	if categoryID != "" {
		product.CategoryID = &categoryID
	}

	err := s.ProductModel.Update(ctx, s.DB, product)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return s.ProductModel.GetByID(ctx, id)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("%w: product id is required", ErrInvalidInput)
	}
	err := s.ProductModel.Delete(ctx, s.DB, id)
	if err == sql.ErrNoRows {
		return ErrNotFound
	}
	return err
}
