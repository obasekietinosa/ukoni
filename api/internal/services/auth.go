package services

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"log"
	"time"
	"ukoni/internal/mailer"
	"ukoni/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserModel *models.UserModel
	JWTSecret string
	Mailer    mailer.Mailer
	WebappURL string
}

func (s *AuthService) Signup(name, email, password string) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
	}

	if err := s.UserModel.Insert(user); err != nil {
		return nil, err
	}

	if s.Mailer != nil {
		go func() {
			err := s.Mailer.SendEmail(user.Email, "Welcome to Ukoni!", "Welcome to Ukoni, "+user.Name+"!")
			if err != nil {
				log.Printf("failed to send welcome email to %s: %v", user.Email, err)
			}
		}()
	}

	return user, nil
}

func (s *AuthService) RequestPasswordReset(email string) error {
	user, err := s.UserModel.GetByEmail(email)
	if err != nil {
		// Do not leak information about whether the user exists or not
		return nil
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return err
	}
	tokenString := hex.EncodeToString(tokenBytes)
	expiresAt := time.Now().Add(1 * time.Hour)

	if err := s.UserModel.CreatePasswordResetToken(user.ID, tokenString, expiresAt); err != nil {
		return err
	}

	if s.Mailer != nil {
		go func() {
			resetLink := fmt.Sprintf("%s/reset-password?token=%s", s.WebappURL, tokenString)

			tmpl, err := template.ParseFS(mailer.TemplatesFS, "templates/password-reset.html")
			if err != nil {
				log.Printf("failed to parse password reset template: %v", err)
				return
			}

			data := struct {
				UserName  string
				ResetLink string
			}{
				UserName:  user.Name,
				ResetLink: resetLink,
			}

			var body bytes.Buffer
			if err := tmpl.Execute(&body, data); err != nil {
				log.Printf("failed to execute password reset template: %v", err)
				return
			}

			err = s.Mailer.SendEmail(user.Email, "Reset Password", body.String())
			if err != nil {
				log.Printf("failed to send password reset email to %s: %v", user.Email, err)
			}
		}()
	}

	return nil
}

func (s *AuthService) ValidatePasswordResetToken(tokenString string) error {
	token, err := s.UserModel.GetPasswordResetToken(tokenString)
	if err != nil {
		return errors.New("invalid token")
	}

	if token.UsedAt != nil {
		return errors.New("token has already been used")
	}

	if time.Now().After(token.ExpiresAt) {
		return errors.New("token has expired")
	}

	return nil
}

func (s *AuthService) ResetPassword(tokenString, newPassword string) error {
	if err := s.ValidatePasswordResetToken(tokenString); err != nil {
		return err
	}

	token, err := s.UserModel.GetPasswordResetToken(tokenString)
	if err != nil {
		return errors.New("invalid token")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := s.UserModel.UpdatePassword(token.UserID, string(hash)); err != nil {
		return err
	}

	if err := s.UserModel.MarkPasswordResetTokenUsed(tokenString); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	user, err := s.UserModel.GetByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return nil, "", err
	}

	return user, tokenString, nil
}
