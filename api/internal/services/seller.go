package services

import (
	"database/sql"
	"ukoni/internal/models"
)

type SellerService struct {
	DB          *sql.DB
	SellerModel *models.SellerModel
}

func (s *SellerService) CreateSeller(name, sellerType string) (*models.Seller, error) {
	return s.SellerModel.Create(name, sellerType)
}

func (s *SellerService) GetSeller(id string) (*models.Seller, error) {
	return s.SellerModel.Get(id)
}

func (s *SellerService) ListSellers() ([]*models.Seller, error) {
	return s.SellerModel.List()
}

func (s *SellerService) UpdateSeller(id, name, sellerType string) (*models.Seller, error) {
	return s.SellerModel.Update(id, name, sellerType)
}

func (s *SellerService) DeleteSeller(id string) error {
	return s.SellerModel.Delete(id)
}
