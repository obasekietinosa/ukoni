package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Outlet struct {
	ID         uuid.UUID  `json:"id"`
	SellerID   uuid.UUID  `json:"seller_id"`
	Name       string     `json:"name"`
	Channel    string     `json:"channel"` // 'physical', 'online'
	Address    string     `json:"address,omitempty"`
	WebsiteURL string     `json:"website_url,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

type OutletModel struct {
	DB *sql.DB
}

func (m *OutletModel) Create(sellerID, name, channel, address, websiteURL string) (*Outlet, error) {
	outlet := &Outlet{
		Name:       name,
		Channel:    channel,
		Address:    address,
		WebsiteURL: websiteURL,
	}

	var err error
	outlet.SellerID, err = uuid.Parse(sellerID)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO outlets (seller_id, name, channel, address, website_url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	err = m.DB.QueryRow(query, outlet.SellerID, outlet.Name, outlet.Channel, outlet.Address, outlet.WebsiteURL).Scan(&outlet.ID, &outlet.CreatedAt)
	if err != nil {
		return nil, err
	}

	return outlet, nil
}

func (m *OutletModel) Get(id string) (*Outlet, error) {
	query := `
		SELECT id, seller_id, name, channel, address, website_url, created_at, deleted_at
		FROM outlets
		WHERE id = $1 AND deleted_at IS NULL`

	var o Outlet
	var address, websiteURL sql.NullString

	err := m.DB.QueryRow(query, id).Scan(
		&o.ID,
		&o.SellerID,
		&o.Name,
		&o.Channel,
		&address,
		&websiteURL,
		&o.CreatedAt,
		&o.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	if address.Valid {
		o.Address = address.String
	}
	if websiteURL.Valid {
		o.WebsiteURL = websiteURL.String
	}

	return &o, nil
}

func (m *OutletModel) ListBySeller(sellerID string) ([]*Outlet, error) {
	query := `
		SELECT id, seller_id, name, channel, address, website_url, created_at, deleted_at
		FROM outlets
		WHERE seller_id = $1 AND deleted_at IS NULL
		ORDER BY name ASC`

	rows, err := m.DB.Query(query, sellerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var outlets []*Outlet
	for rows.Next() {
		var o Outlet
		var address, websiteURL sql.NullString
		err := rows.Scan(
			&o.ID,
			&o.SellerID,
			&o.Name,
			&o.Channel,
			&address,
			&websiteURL,
			&o.CreatedAt,
			&o.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		if address.Valid {
			o.Address = address.String
		}
		if websiteURL.Valid {
			o.WebsiteURL = websiteURL.String
		}
		outlets = append(outlets, &o)
	}

	return outlets, nil
}

func (m *OutletModel) Update(id, name, channel, address, websiteURL string) (*Outlet, error) {
	query := `
		UPDATE outlets
		SET name = $1, channel = $2, address = $3, website_url = $4
		WHERE id = $5 AND deleted_at IS NULL
		RETURNING id, seller_id, created_at, deleted_at`

	var o Outlet
	o.Name = name
	o.Channel = channel
	o.Address = address
	o.WebsiteURL = websiteURL

	err := m.DB.QueryRow(query, name, channel, address, websiteURL, id).Scan(
		&o.ID,
		&o.SellerID,
		&o.CreatedAt,
		&o.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &o, nil
}

func (m *OutletModel) Delete(id string) error {
	query := `
		UPDATE outlets
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
