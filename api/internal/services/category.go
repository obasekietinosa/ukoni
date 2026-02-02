package services

import (
	"context"
	"database/sql"
	"fmt"
	"ukoni/internal/models"
)

type CategoryService struct {
	DB            *sql.DB
	CategoryModel *models.CategoryModel
}

func (s *CategoryService) CreateCategory(ctx context.Context, inventoryID, name string, parentCategoryID *string) (*models.Category, error) {
	if inventoryID == "" {
		return nil, fmt.Errorf("%w: inventory id is required", ErrInvalidInput)
	}
	if name == "" {
		return nil, fmt.Errorf("%w: category name is required", ErrInvalidInput)
	}

	category := &models.Category{
		InventoryID:      inventoryID,
		Name:             name,
		ParentCategoryID: parentCategoryID,
	}

	err := s.CategoryModel.Create(ctx, s.DB, category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) ListCategories(ctx context.Context, inventoryID string) ([]*models.Category, error) {
	if inventoryID == "" {
		return nil, fmt.Errorf("%w: inventory id is required", ErrInvalidInput)
	}
	categories, err := s.CategoryModel.List(ctx, inventoryID)
	if err != nil {
		return nil, err
	}
	if categories == nil {
		categories = []*models.Category{}
	}
	return categories, nil
}
