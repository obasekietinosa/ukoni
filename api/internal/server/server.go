package server

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

type Server struct {
	Config *config.Config
	DB     database.Service
	Logger *slog.Logger
}

func New(cfg *config.Config, db database.Service, logger *slog.Logger) *Server {
	return &Server{
		Config: cfg,
		DB:     db,
		Logger: logger,
	}
}

func (s *Server) SetupRouter() *http.ServeMux {
	// Initialize models
	userModel := &models.UserModel{DB: s.DB.GetDB()}
	inventoryModel := &models.InventoryModel{DB: s.DB.GetDB()}
	membershipModel := &models.MembershipModel{DB: s.DB.GetDB()}
	activityLogModel := &models.ActivityLogModel{DB: s.DB.GetDB()}
	productModel := &models.ProductModel{DB: s.DB.GetDB()}
	canonicalProductModel := &models.CanonicalProductModel{DB: s.DB.GetDB()}
	sellerModel := &models.SellerModel{DB: s.DB.GetDB()}
	outletModel := &models.OutletModel{DB: s.DB.GetDB()}
	shoppingListModel := &models.ShoppingListModel{DB: s.DB.GetDB()}
	transactionModel := &models.TransactionModel{DB: s.DB.GetDB()}
	inventoryProductModel := &models.InventoryProductModel{DB: s.DB.GetDB()}
	consumptionModel := &models.ConsumptionModel{DB: s.DB.GetDB()}

	// Initialize services
	authService := &services.AuthService{
		UserModel: userModel,
		JWTSecret: s.Config.JWTSecret,
	}

	activityLogService := &services.ActivityLogService{
		Model: activityLogModel,
	}

	inventoryService := &services.InventoryService{
		DB:                 s.DB.GetDB(),
		InventoryModel:     inventoryModel,
		MembershipModel:    membershipModel,
		ActivityLogService: activityLogService,
	}

	membershipService := &services.MembershipService{
		MembershipModel:    membershipModel,
		InventoryModel:     inventoryModel,
		ActivityLogService: activityLogService,
	}

	productService := &services.ProductService{
		DB:           s.DB.GetDB(),
		ProductModel: productModel,
	}

	canonicalProductService := &services.CanonicalProductService{
		DB:                    s.DB.GetDB(),
		CanonicalProductModel: canonicalProductModel,
	}

	sellerService := &services.SellerService{
		DB:          s.DB.GetDB(),
		SellerModel: sellerModel,
	}

	outletService := &services.OutletService{
		DB:          s.DB.GetDB(),
		OutletModel: outletModel,
	}

	shoppingListService := &services.ShoppingListService{
		ShoppingListModel:  shoppingListModel,
		InventoryModel:     inventoryModel,
		MembershipModel:    membershipModel,
		ActivityLogService: activityLogService,
	}

	inventoryProductService := &services.InventoryProductService{
		InventoryProductModel: inventoryProductModel,
		ProductModel:          productModel,
	}

	transactionService := &services.TransactionService{
		DB:                      s.DB.GetDB(),
		TransactionModel:        transactionModel,
		MembershipModel:         membershipModel,
		OutletModel:             outletModel,
		ActivityLogService:      activityLogService,
		InventoryProductService: inventoryProductService,
	}

	consumptionService := &services.ConsumptionService{
		DB:                 s.DB.GetDB(),
		ConsumptionModel:   consumptionModel,
		MembershipModel:    membershipModel,
		ActivityLogService: activityLogService,
	}

	// Initialize handlers
	authHandler := &handlers.AuthHandler{Service: authService}
	inventoryHandler := &handlers.InventoryHandler{Service: inventoryService}
	membershipHandler := &handlers.MembershipHandler{Service: membershipService}
	productHandler := &handlers.ProductHandler{Service: productService}
	canonicalProductHandler := &handlers.CanonicalProductHandler{Service: canonicalProductService}
	sellerHandler := &handlers.SellerHandler{Service: sellerService}
	outletHandler := &handlers.OutletHandler{Service: outletService}
	shoppingListHandler := &handlers.ShoppingListHandler{Service: shoppingListService}
	transactionHandler := &handlers.TransactionHandler{Service: transactionService}
	consumptionHandler := &handlers.ConsumptionHandler{Service: consumptionService}

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(s.Config)

	// Setup router
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

	router.HandleFunc("POST /sellers", authMiddleware.Auth(sellerHandler.CreateSeller))
	router.HandleFunc("GET /sellers", authMiddleware.Auth(sellerHandler.ListSellers))
	router.HandleFunc("GET /sellers/{id}", authMiddleware.Auth(sellerHandler.GetSeller))
	router.HandleFunc("PUT /sellers/{id}", authMiddleware.Auth(sellerHandler.UpdateSeller))
	router.HandleFunc("DELETE /sellers/{id}", authMiddleware.Auth(sellerHandler.DeleteSeller))

	router.HandleFunc("POST /sellers/{id}/outlets", authMiddleware.Auth(outletHandler.CreateOutlet))
	router.HandleFunc("GET /sellers/{id}/outlets", authMiddleware.Auth(outletHandler.ListOutlets))
	router.HandleFunc("GET /outlets/{id}", authMiddleware.Auth(outletHandler.GetOutlet))
	router.HandleFunc("PUT /outlets/{id}", authMiddleware.Auth(outletHandler.UpdateOutlet))
	router.HandleFunc("DELETE /outlets/{id}", authMiddleware.Auth(outletHandler.DeleteOutlet))

	router.HandleFunc("POST /inventories/{id}/shopping-lists", authMiddleware.Auth(shoppingListHandler.CreateList))
	router.HandleFunc("GET /inventories/{id}/shopping-lists", authMiddleware.Auth(shoppingListHandler.ListLists))
	router.HandleFunc("GET /shopping-lists/{id}", authMiddleware.Auth(shoppingListHandler.GetList))
	router.HandleFunc("PUT /shopping-lists/{id}", authMiddleware.Auth(shoppingListHandler.UpdateList))
	router.HandleFunc("DELETE /shopping-lists/{id}", authMiddleware.Auth(shoppingListHandler.DeleteList))
	router.HandleFunc("GET /shopping-lists/{id}/items", authMiddleware.Auth(shoppingListHandler.ListItems))
	router.HandleFunc("POST /shopping-lists/{id}/items", authMiddleware.Auth(shoppingListHandler.AddItem))
	router.HandleFunc("PUT /shopping-list-items/{itemId}", authMiddleware.Auth(shoppingListHandler.UpdateItem))
	router.HandleFunc("DELETE /shopping-list-items/{itemId}", authMiddleware.Auth(shoppingListHandler.DeleteItem))

	router.HandleFunc("POST /inventories/{id}/transactions", authMiddleware.Auth(transactionHandler.CreateTransaction))
	router.HandleFunc("GET /inventories/{id}/transactions", authMiddleware.Auth(transactionHandler.ListTransactions))
	router.HandleFunc("GET /transactions/{id}", authMiddleware.Auth(transactionHandler.GetTransaction))

	router.HandleFunc("POST /inventories/{id}/consumption-events", authMiddleware.Auth(consumptionHandler.CreateConsumptionEvent))
	router.HandleFunc("GET /inventories/{id}/consumption-events", authMiddleware.Auth(consumptionHandler.ListConsumptionEvents))

	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	return router
}

func (s *Server) Run() error {
	router := s.SetupRouter()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.Config.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown channel
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s.Logger.Info("starting server", "port", s.Config.Port, "env", s.Config.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Logger.Error("server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	<-done
	s.Logger.Info("server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		s.Logger.Error("server shutdown failed", "error", err)
		return err
	}
	s.Logger.Info("server exited properly")
	return nil
}
