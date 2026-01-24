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
	authMiddleware := middleware.NewAuthMiddleware(cfg)

	inventoryModel := &models.InventoryModel{DB: dbService.GetDB()}
	membershipModel := &models.MembershipModel{DB: dbService.GetDB()}
	activityLogModel := &models.ActivityLogModel{DB: dbService.GetDB()}

	activityLogService := &services.ActivityLogService{
		Model: activityLogModel,
	}

	inventoryService := &services.InventoryService{
		DB:                 dbService.GetDB(),
		InventoryModel:     inventoryModel,
		MembershipModel:    membershipModel,
		ActivityLogService: activityLogService,
	}
	inventoryHandler := &handlers.InventoryHandler{Service: inventoryService}

	membershipService := &services.MembershipService{
		MembershipModel:    membershipModel,
		InventoryModel:     inventoryModel,
		ActivityLogService: activityLogService,
	}
	membershipHandler := &handlers.MembershipHandler{Service: membershipService}

	productModel := &models.ProductModel{DB: dbService.GetDB()}
	productService := &services.ProductService{
		DB:           dbService.GetDB(),
		ProductModel: productModel,
	}
	productHandler := &handlers.ProductHandler{Service: productService}

	canonicalProductModel := &models.CanonicalProductModel{DB: dbService.GetDB()}
	canonicalProductService := &services.CanonicalProductService{
		DB:                    dbService.GetDB(),
		CanonicalProductModel: canonicalProductModel,
	}
	canonicalProductHandler := &handlers.CanonicalProductHandler{Service: canonicalProductService}

	router := http.NewServeMux()
	router.HandleFunc("POST /signup", authHandler.Signup)
	router.HandleFunc("POST /login", authHandler.Login)

	router.HandleFunc("POST /inventories", authMiddleware.Auth(inventoryHandler.CreateInventory))
	router.HandleFunc("GET /inventories", authMiddleware.Auth(inventoryHandler.ListInventories))
	router.HandleFunc("GET /inventories/{id}", authMiddleware.Auth(inventoryHandler.GetInventory))

	router.HandleFunc("POST /inventories/{id}/invitations", authMiddleware.Auth(membershipHandler.InviteUser))
	router.HandleFunc("GET /inventories/{id}/members", authMiddleware.Auth(membershipHandler.ListMembers))
	router.HandleFunc("DELETE /inventories/{id}/members/{userId}", authMiddleware.Auth(membershipHandler.RemoveMember))
	router.HandleFunc("POST /invitations/{id}/accept", authMiddleware.Auth(membershipHandler.AcceptInvite))

	router.HandleFunc("POST /products", authMiddleware.Auth(productHandler.CreateProduct))
	router.HandleFunc("GET /products", authMiddleware.Auth(productHandler.ListProducts))
	router.HandleFunc("GET /products/{id}", authMiddleware.Auth(productHandler.GetProduct))
	router.HandleFunc("PUT /products/{id}", authMiddleware.Auth(productHandler.UpdateProduct))
	router.HandleFunc("DELETE /products/{id}", authMiddleware.Auth(productHandler.DeleteProduct))
	router.HandleFunc("POST /products/{id}/variants", authMiddleware.Auth(productHandler.CreateVariant))
	router.HandleFunc("GET /products/{id}/variants", authMiddleware.Auth(productHandler.ListVariants))

	router.HandleFunc("POST /canonical-products", authMiddleware.Auth(canonicalProductHandler.CreateCanonicalProduct))
	router.HandleFunc("GET /canonical-products", authMiddleware.Auth(canonicalProductHandler.ListCanonicalProducts))
	router.HandleFunc("GET /canonical-products/{id}", authMiddleware.Auth(canonicalProductHandler.GetCanonicalProduct))
	router.HandleFunc("PUT /canonical-products/{id}", authMiddleware.Auth(canonicalProductHandler.UpdateCanonicalProduct))
	router.HandleFunc("DELETE /canonical-products/{id}", authMiddleware.Auth(canonicalProductHandler.DeleteCanonicalProduct))

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
