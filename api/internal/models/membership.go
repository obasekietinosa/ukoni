package models

import (
	"context"
	"database/sql"
	"time"
)

type InventoryMembership struct {
	ID          string     `json:"id"`
	InventoryID string     `json:"inventory_id"`
	UserID      string     `json:"user_id"`
	Role        string     `json:"role"`
	InvitedAt   time.Time  `json:"invited_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type Invitation struct {
	ID              string     `json:"id"`
	InventoryID     string     `json:"inventory_id"`
	Email           string     `json:"email"`
	Role            string     `json:"role"`
	InvitedByUserID string     `json:"invited_by_user_id"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	AcceptedAt      *time.Time `json:"accepted_at,omitempty"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
}

type MembershipModel struct {
	DB *sql.DB
}

func (m *MembershipModel) CreateInvitation(invitation *Invitation) error {
	query := `
		INSERT INTO invitations (inventory_id, email, role, invited_by_user_id, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	return m.DB.QueryRowContext(context.Background(), query,
		invitation.InventoryID, invitation.Email, invitation.Role, invitation.InvitedByUserID, invitation.Status,
	).Scan(&invitation.ID, &invitation.CreatedAt)
}

func (m *MembershipModel) GetInvitationByID(id string) (*Invitation, error) {
	query := `
		SELECT id, inventory_id, email, role, invited_by_user_id, status, created_at, accepted_at, expires_at
		FROM invitations
		WHERE id = $1
	`
	var i Invitation
	err := m.DB.QueryRowContext(context.Background(), query, id).Scan(
		&i.ID, &i.InventoryID, &i.Email, &i.Role, &i.InvitedByUserID, &i.Status, &i.CreatedAt, &i.AcceptedAt, &i.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func (m *MembershipModel) AcceptInvitation(inviteID, userID string, now time.Time) error {
	// Start a transaction
	tx, err := m.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Get the invitation to verify and get details
	var i Invitation
	queryInvite := `
		SELECT inventory_id, role, status
		FROM invitations
		WHERE id = $1 FOR UPDATE
	`
	err = tx.QueryRowContext(context.Background(), queryInvite, inviteID).Scan(&i.InventoryID, &i.Role, &i.Status)
	if err != nil {
		return err // Handle not found
	}

	if i.Status != "pending" {
		// Could return custom error here
		return sql.ErrNoRows // Or specific error "invitation not pending"
	}

	// 2. Update invitation status
	updateInvite := `
		UPDATE invitations
		SET status = 'accepted', accepted_at = $1
		WHERE id = $2
	`
	_, err = tx.ExecContext(context.Background(), updateInvite, now, inviteID)
	if err != nil {
		return err
	}

	// 3. Create membership
	createMember := `
		INSERT INTO inventory_memberships (inventory_id, user_id, role)
		VALUES ($1, $2, $3)
	`
	_, err = tx.ExecContext(context.Background(), createMember, i.InventoryID, userID, i.Role)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (m *MembershipModel) ListMembers(inventoryID string) ([]*InventoryMembership, error) {
	query := `
		SELECT id, inventory_id, user_id, role, invited_at, deleted_at
		FROM inventory_memberships
		WHERE inventory_id = $1 AND deleted_at IS NULL
	`
	rows, err := m.DB.QueryContext(context.Background(), query, inventoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*InventoryMembership
	for rows.Next() {
		var member InventoryMembership
		if err := rows.Scan(&member.ID, &member.InventoryID, &member.UserID, &member.Role, &member.InvitedAt, &member.DeletedAt); err != nil {
			return nil, err
		}
		members = append(members, &member)
	}
	return members, rows.Err()
}

func (m *MembershipModel) RemoveMember(inventoryID, userID string) error {
	query := `
		UPDATE inventory_memberships
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE inventory_id = $1 AND user_id = $2
	`
	_, err := m.DB.ExecContext(context.Background(), query, inventoryID, userID)
	return err
}

func (m *MembershipModel) AddMember(inventoryID, userID, role string) error {
	query := `
		INSERT INTO inventory_memberships (inventory_id, user_id, role)
		VALUES ($1, $2, $3)
	`
	_, err := m.DB.ExecContext(context.Background(), query, inventoryID, userID, role)
	return err
}

// CheckMembership checks if a user is a member of an inventory (active)
func (m *MembershipModel) GetMembership(inventoryID, userID string) (*InventoryMembership, error) {
	query := `
		SELECT id, inventory_id, user_id, role, invited_at, deleted_at
		FROM inventory_memberships
		WHERE inventory_id = $1 AND user_id = $2 AND deleted_at IS NULL
	`
	var member InventoryMembership
	err := m.DB.QueryRowContext(context.Background(), query, inventoryID, userID).Scan(
		&member.ID, &member.InventoryID, &member.UserID, &member.Role, &member.InvitedAt, &member.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &member, nil
}
