package main

import (
	"log/slog"
	"os"
	"time"

	"ukoni/internal/config"
	"ukoni/internal/database"
	"ukoni/internal/models"
	"ukoni/internal/services"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	dbService, err := database.New(cfg.DBURL)
	if err != nil {
		logger.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer dbService.Close()
	logger.Info("database connected")

	userModel := &models.UserModel{DB: dbService.GetDB()}
	authService := &services.AuthService{
		UserModel: userModel,
		JWTSecret: cfg.JWTSecret,
	}

	// Seed basic user
	email := "test@example.com"
	_, err = userModel.GetByEmail(email)
	if err == nil {
		logger.Info("test user already exists", "email", email)
		return
	}

	logger.Info("creating test user", "email", email)
	start := time.Now()
	_, err = authService.Signup("Test User", email, "password123")
	if err != nil {
		logger.Error("failed to create test user", "error", err)
		os.Exit(1)
	}
	logger.Info("test user created", "duration", time.Since(start))
}
