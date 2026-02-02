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

func (s *CategoryService) CreateCategory(ctx context.Context, name string, parentCategoryID *string) (*models.Category, error) {
	if name == "" {
		return nil, fmt.Errorf("%w: category name is required", ErrInvalidInput)
	}

	category := &models.Category{
		Name:             name,
		ParentCategoryID: parentCategoryID,
	}

	err := s.CategoryModel.Create(ctx, s.DB, category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) ListCategories(ctx context.Context) ([]*models.Category, error) {
	categories, err := s.CategoryModel.List(ctx)
	if err != nil {
		return nil, err
	}
	if categories == nil {
		categories = []*models.Category{}
	}
	return categories, nil
}
