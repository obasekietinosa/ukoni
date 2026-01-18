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

	inventoryService := &services.InventoryService{
		DB:              testDB,
		InventoryModel:  inventoryModel,
		MembershipModel: membershipModel,
	}
	inventoryHandler := &handlers.InventoryHandler{Service: inventoryService}

	membershipService := &services.MembershipService{
		MembershipModel: membershipModel,
		InventoryModel:  inventoryModel,
	}
	membershipHandler := &handlers.MembershipHandler{Service: membershipService}

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
