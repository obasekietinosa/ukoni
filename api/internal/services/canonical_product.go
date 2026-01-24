package services

import (
	"context"
	"database/sql"
	"fmt"
	"ukoni/internal/models"
)

type CanonicalProductService struct {
	DB                    *sql.DB
	CanonicalProductModel *models.CanonicalProductModel
}

func (s *CanonicalProductService) CreateCanonicalProduct(ctx context.Context, name, description, categoryID string) (*models.CanonicalProduct, error) {
	if name == "" {
		return nil, fmt.Errorf("%w: product name is required", ErrInvalidInput)
	}

	product := &models.CanonicalProduct{
		Name: name,
	}
	if description != "" {
		product.Description = &description
	}
	if categoryID != "" {
		product.CategoryID = &categoryID
	}

	err := s.CanonicalProductModel.Create(ctx, s.DB, product)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *CanonicalProductService) GetCanonicalProduct(ctx context.Context, id string) (*models.CanonicalProduct, error) {
	return s.CanonicalProductModel.GetByID(ctx, id)
}

func (s *CanonicalProductService) ListCanonicalProducts(ctx context.Context, limit, offset int, search string) ([]*models.CanonicalProduct, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.CanonicalProductModel.List(ctx, limit, offset, search)
}

func (s *CanonicalProductService) UpdateCanonicalProduct(ctx context.Context, id, name, description, categoryID string) (*models.CanonicalProduct, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: product id is required", ErrInvalidInput)
	}
	if name == "" {
		return nil, fmt.Errorf("%w: product name is required", ErrInvalidInput)
	}

	product := &models.CanonicalProduct{
		ID:   id,
		Name: name,
	}
	if description != "" {
		product.Description = &description
	}
	if categoryID != "" {
		product.CategoryID = &categoryID
	}

	err := s.CanonicalProductModel.Update(ctx, s.DB, product)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return s.CanonicalProductModel.GetByID(ctx, id)
}

func (s *CanonicalProductService) DeleteCanonicalProduct(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("%w: product id is required", ErrInvalidInput)
	}
	err := s.CanonicalProductModel.Delete(ctx, s.DB, id)
	if err == sql.ErrNoRows {
		return ErrNotFound
	}
	return err
}
