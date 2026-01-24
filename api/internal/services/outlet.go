package services

import (
	"database/sql"
	"ukoni/internal/models"
)

type OutletService struct {
	DB          *sql.DB
	OutletModel *models.OutletModel
}

func (s *OutletService) CreateOutlet(sellerID, name, channel, address, websiteURL string) (*models.Outlet, error) {
	return s.OutletModel.Create(sellerID, name, channel, address, websiteURL)
}

func (s *OutletService) GetOutlet(id string) (*models.Outlet, error) {
	return s.OutletModel.Get(id)
}

func (s *OutletService) ListOutlets(sellerID string) ([]*models.Outlet, error) {
	return s.OutletModel.ListBySeller(sellerID)
}

func (s *OutletService) UpdateOutlet(id, name, channel, address, websiteURL string) (*models.Outlet, error) {
	return s.OutletModel.Update(id, name, channel, address, websiteURL)
}

func (s *OutletService) DeleteOutlet(id string) error {
	return s.OutletModel.Delete(id)
}
