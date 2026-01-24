package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Seller struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Type      string     `json:"type"` // 'chain', 'independent', 'online'
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type SellerModel struct {
	DB *sql.DB
}

func (m *SellerModel) Create(name, sellerType string) (*Seller, error) {
	seller := &Seller{
		Name: name,
		Type: sellerType,
	}

	query := `
		INSERT INTO sellers (name, type)
		VALUES ($1, $2)
		RETURNING id, created_at`

	err := m.DB.QueryRow(query, seller.Name, seller.Type).Scan(&seller.ID, &seller.CreatedAt)
	if err != nil {
		return nil, err
	}

	return seller, nil
}

func (m *SellerModel) Get(id string) (*Seller, error) {
	query := `
		SELECT id, name, type, created_at, deleted_at
		FROM sellers
		WHERE id = $1 AND deleted_at IS NULL`

	var seller Seller
	err := m.DB.QueryRow(query, id).Scan(
		&seller.ID,
		&seller.Name,
		&seller.Type,
		&seller.CreatedAt,
		&seller.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &seller, nil
}

func (m *SellerModel) List() ([]*Seller, error) {
	query := `
		SELECT id, name, type, created_at, deleted_at
		FROM sellers
		WHERE deleted_at IS NULL
		ORDER BY name ASC`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sellers []*Seller
	for rows.Next() {
		var s Seller
		err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.Type,
			&s.CreatedAt,
			&s.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		sellers = append(sellers, &s)
	}

	return sellers, nil
}

func (m *SellerModel) Update(id, name, sellerType string) (*Seller, error) {
	query := `
		UPDATE sellers
		SET name = $1, type = $2
		WHERE id = $3 AND deleted_at IS NULL
		RETURNING id, name, type, created_at, deleted_at`

	var seller Seller
	err := m.DB.QueryRow(query, name, sellerType, id).Scan(
		&seller.ID,
		&seller.Name,
		&seller.Type,
		&seller.CreatedAt,
		&seller.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &seller, nil
}

func (m *SellerModel) Delete(id string) error {
	query := `
		UPDATE sellers
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
