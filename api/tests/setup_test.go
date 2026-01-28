package tests

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"ukoni/internal/config"
	"ukoni/internal/database"
	"ukoni/internal/handlers"
	"ukoni/internal/middleware"
	"ukoni/internal/models"
	"ukoni/internal/services"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	testDB *sql.DB
	cfg    *config.Config
)

func TestMain(m *testing.M) {
	// Setup
	cfg = config.Load()
	// Override DB URL for tests if needed
	// ensure we use the local db
	// cfg.DBURL = "postgres://postgres:postgres@localhost:5432/ukoni?sslmode=disable"

	dbService, err := database.New(cfg.DBURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	testDB = dbService.GetDB()

	// Run tests
	code := m.Run()

	// Teardown
	dbService.Close()
	os.Exit(code)
}

func setupRouter() *http.ServeMux {
	userModel := &models.UserModel{DB: testDB}
	authService := &services.AuthService{
		UserModel: userModel,
		JWTSecret: cfg.JWTSecret,
	}
	authHandler := &handlers.AuthHandler{Service: authService}
	authMiddleware := middleware.NewAuthMiddleware(cfg)

	inventoryModel := &models.InventoryModel{DB: testDB}
	membershipModel := &models.MembershipModel{DB: testDB}
	activityLogModel := &models.ActivityLogModel{DB: testDB}

	activityLogService := &services.ActivityLogService{
		Model: activityLogModel,
	}

	inventoryService := &services.InventoryService{
		DB:                 testDB,
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

	productModel := &models.ProductModel{DB: testDB}
	productService := &services.ProductService{
		DB:           testDB,
		ProductModel: productModel,
	}
	productHandler := &handlers.ProductHandler{Service: productService}

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

	// router.HandleFunc("GET /products/{id}/variants", authMiddleware.Auth(productHandler.ListVariants))

	canonicalProductModel := &models.CanonicalProductModel{DB: testDB}
	canonicalProductService := &services.CanonicalProductService{
		DB:                    testDB,
		CanonicalProductModel: canonicalProductModel,
	}
	canonicalProductHandler := &handlers.CanonicalProductHandler{Service: canonicalProductService}

	router.HandleFunc("POST /canonical-products", authMiddleware.Auth(canonicalProductHandler.CreateCanonicalProduct))
	router.HandleFunc("GET /canonical-products", authMiddleware.Auth(canonicalProductHandler.ListCanonicalProducts))
	router.HandleFunc("GET /canonical-products/{id}", authMiddleware.Auth(canonicalProductHandler.GetCanonicalProduct))
	router.HandleFunc("PUT /canonical-products/{id}", authMiddleware.Auth(canonicalProductHandler.UpdateCanonicalProduct))
	router.HandleFunc("DELETE /canonical-products/{id}", authMiddleware.Auth(canonicalProductHandler.DeleteCanonicalProduct))

	sellerModel := &models.SellerModel{DB: testDB}
	sellerService := &services.SellerService{
		DB:          testDB,
		SellerModel: sellerModel,
	}
	sellerHandler := &handlers.SellerHandler{Service: sellerService}

	outletModel := &models.OutletModel{DB: testDB}
	outletService := &services.OutletService{
		DB:          testDB,
		OutletModel: outletModel,
	}
	outletHandler := &handlers.OutletHandler{Service: outletService}

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

	shoppingListModel := &models.ShoppingListModel{DB: testDB}
	shoppingListService := &services.ShoppingListService{
		ShoppingListModel:  shoppingListModel,
		InventoryModel:     inventoryModel,
		MembershipModel:    membershipModel,
		ActivityLogService: activityLogService,
	}
	shoppingListHandler := &handlers.ShoppingListHandler{Service: shoppingListService}

	router.HandleFunc("POST /inventories/{id}/shopping-lists", authMiddleware.Auth(shoppingListHandler.CreateList))
	router.HandleFunc("GET /inventories/{id}/shopping-lists", authMiddleware.Auth(shoppingListHandler.ListLists))
	router.HandleFunc("GET /shopping-lists/{id}", authMiddleware.Auth(shoppingListHandler.GetList))
	router.HandleFunc("PUT /shopping-lists/{id}", authMiddleware.Auth(shoppingListHandler.UpdateList))
	router.HandleFunc("DELETE /shopping-lists/{id}", authMiddleware.Auth(shoppingListHandler.DeleteList))
	router.HandleFunc("GET /shopping-lists/{id}/items", authMiddleware.Auth(shoppingListHandler.ListItems))
	router.HandleFunc("POST /shopping-lists/{id}/items", authMiddleware.Auth(shoppingListHandler.AddItem))
	router.HandleFunc("PUT /shopping-list-items/{itemId}", authMiddleware.Auth(shoppingListHandler.UpdateItem))
	router.HandleFunc("DELETE /shopping-list-items/{itemId}", authMiddleware.Auth(shoppingListHandler.DeleteItem))

	transactionModel := &models.TransactionModel{DB: testDB}
	inventoryProductModel := &models.InventoryProductModel{DB: testDB}
	inventoryProductService := &services.InventoryProductService{
		InventoryProductModel: inventoryProductModel,
		ProductModel:          productModel,
	}
	transactionService := &services.TransactionService{
		DB:                      testDB,
		TransactionModel:        transactionModel,
		MembershipModel:         membershipModel,
		OutletModel:             outletModel,
		ActivityLogService:      activityLogService,
		InventoryProductService: inventoryProductService,
	}
	transactionHandler := &handlers.TransactionHandler{Service: transactionService}

	router.HandleFunc("POST /inventories/{id}/transactions", authMiddleware.Auth(transactionHandler.CreateTransaction))
	router.HandleFunc("GET /inventories/{id}/transactions", authMiddleware.Auth(transactionHandler.ListTransactions))
	router.HandleFunc("GET /transactions/{id}", authMiddleware.Auth(transactionHandler.GetTransaction))

	consumptionModel := &models.ConsumptionModel{DB: testDB}
	consumptionService := &services.ConsumptionService{
		DB:                 testDB,
		ConsumptionModel:   consumptionModel,
		MembershipModel:    membershipModel,
		ActivityLogService: activityLogService,
	}
	consumptionHandler := &handlers.ConsumptionHandler{Service: consumptionService}

	router.HandleFunc("POST /inventories/{id}/consumption-events", authMiddleware.Auth(consumptionHandler.CreateConsumptionEvent))
	router.HandleFunc("GET /inventories/{id}/consumption-events", authMiddleware.Auth(consumptionHandler.ListConsumptionEvents))

	return router
}

func clearDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tables := []string{
		"shopping_list_items",
		"shopping_lists",
		"activity_logs",
		"consumption_events",
		"transaction_items",
		"transactions",
		"inventory_products",
		"outlets",
		"sellers",
		"product_variants",
		"products",
		"product_categories",
		"canonical_products",
		"invitations",
		"inventory_memberships",
		"inventories",
		"users",
	}

	for _, table := range tables {
		_, err := testDB.ExecContext(ctx, "DELETE FROM "+table)
		if err != nil {
			log.Printf("failed to clear table %s: %v", table, err)
		}
	}
}
