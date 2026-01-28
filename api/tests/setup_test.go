package tests

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"ukoni/internal/config"
	"ukoni/internal/database"
	"ukoni/internal/server"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	testDB    *sql.DB
	cfg       *config.Config
	dbService database.Service
)

func TestMain(m *testing.M) {
	// Setup
	cfg = config.Load()
	// Override DB URL for tests if needed
	// ensure we use the local db
	// cfg.DBURL = "postgres://postgres:postgres@localhost:5432/ukoni?sslmode=disable"

	var err error
	dbService, err = database.New(cfg.DBURL)
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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	srv := server.New(cfg, dbService, logger)
	return srv.SetupRouter()
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
