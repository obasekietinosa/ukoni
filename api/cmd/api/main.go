package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ukoni/internal/config"
	"ukoni/internal/database"
	"ukoni/internal/handlers"
	"ukoni/internal/middleware"
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
	authHandler := &handlers.AuthHandler{Service: authService}

	inventoryModel := &models.InventoryModel{DB: dbService.GetDB()}
	inventoryService := &services.InventoryService{InventoryModel: inventoryModel}
	inventoryHandler := &handlers.InventoryHandler{Service: inventoryService}

	router := http.NewServeMux()
	router.HandleFunc("POST /signup", authHandler.Signup)
	router.HandleFunc("POST /login", authHandler.Login)

	router.HandleFunc("POST /inventories", middleware.Auth(inventoryHandler.CreateInventory))
	router.HandleFunc("GET /inventories", middleware.Auth(inventoryHandler.ListInventories))
	router.HandleFunc("GET /inventories/{id}", middleware.Auth(inventoryHandler.GetInventory))

	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown channel
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("starting server", "port", cfg.Port, "env", cfg.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	<-done
	logger.Info("server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed", "error", err)
		os.Exit(1)
	}
	logger.Info("server exited properly")
}
