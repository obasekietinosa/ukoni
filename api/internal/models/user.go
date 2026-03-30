package models

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	ID           string     `json:"id"`
	Email        string     `json:"email"`
	Name         string     `json:"name"`
	PasswordHash string     `json:"-"`
	CreatedAt    time.Time  `json:"created_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

type PasswordResetToken struct {
	Token     string     `json:"token"`
	UserID    string     `json:"user_id"`
	ExpiresAt time.Time  `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (email, name, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	args := []interface{}{user.Email, user.Name, user.PasswordHash}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt)
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	query := `
		SELECT id, email, name, password_hash, created_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL`

	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) UpdatePassword(userID, passwordHash string) error {
	query := `
		UPDATE users
		SET password_hash = $1
		WHERE id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, passwordHash, userID)
	return err
}

func (m *UserModel) CreatePasswordResetToken(userID, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO password_reset_tokens (token, user_id, expires_at)
		VALUES ($1, $2, $3)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, token, userID, expiresAt)
	return err
}

func (m *UserModel) GetPasswordResetToken(token string) (*PasswordResetToken, error) {
	query := `
		SELECT token, user_id, expires_at, used_at, created_at
		FROM password_reset_tokens
		WHERE token = $1`

	var t PasswordResetToken
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, token).Scan(
		&t.Token,
		&t.UserID,
		&t.ExpiresAt,
		&t.UsedAt,
		&t.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (m *UserModel) MarkPasswordResetTokenUsed(token string) error {
	query := `
		UPDATE password_reset_tokens
		SET used_at = CURRENT_TIMESTAMP
		WHERE token = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, token)
	return err
}
