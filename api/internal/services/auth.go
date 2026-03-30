package services

import (
	"errors"
	"fmt"
	"log"
	"time"
	"ukoni/internal/mailer"
	"ukoni/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserModel   *models.UserModel
	JWTSecret   string
	Mailer      mailer.Mailer
	FrontendURL string
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"type": "reset",
		"exp":  time.Now().Add(1 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return err
	}

	if s.Mailer != nil {
		go func() {
			resetLink := fmt.Sprintf("%s/reset-password?token=%s", s.FrontendURL, tokenString)
			body := fmt.Sprintf("Click the following link to reset your password: <a href=\"%s\">%s</a>", resetLink, resetLink)
			err := s.Mailer.SendEmail(user.Email, "Reset Password", body)
			if err != nil {
				log.Printf("failed to send password reset email to %s: %v", user.Email, err)
			}
		}()
	}

	return nil
}

func (s *AuthService) ResetPassword(tokenString, newPassword string) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.JWTSecret), nil
	})

	if err != nil {
		return errors.New("invalid or expired token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["type"] != "reset" {
			return errors.New("invalid token type")
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			return errors.New("invalid token claims")
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		if err := s.UserModel.UpdatePassword(userID, string(hash)); err != nil {
			return err
		}

		return nil
	}

	return errors.New("invalid token")
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
