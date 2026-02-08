package tests

import (
	"database/sql"
	"fmt"
	"math/rand"
	"testing"

	"github.com/google/uuid"
)

func seedBenchmarkData(b *testing.B, db *sql.DB) ([]string, []string) {
	// Create a user
	userID := uuid.New().String()
	_, err := db.Exec("INSERT INTO users (id, email, name, password_hash) VALUES ($1, $2, $3, $4)", userID, "bench@example.com", "Bench User", "hash")
	if err != nil {
		b.Fatalf("Failed to create user: %v", err)
	}

	// Create an inventory
	inventoryID := uuid.New().String()
	_, err = db.Exec("INSERT INTO inventories (id, name, owner_user_id) VALUES ($1, $2, $3)", inventoryID, "Bench Inventory", userID)
	if err != nil {
		b.Fatalf("Failed to create inventory: %v", err)
	}

	// Create product variants (let's say 100 variants)
	variantIDs := make([]string, 100)

	// We need canonical product, product category (optional), and product first
	canonicalProductID := uuid.New().String()
	_, err = db.Exec("INSERT INTO canonical_products (id, name, inventory_id) VALUES ($1, $2, $3)", canonicalProductID, "Bench CP", inventoryID)
    if err != nil {
        b.Fatalf("Failed to create canonical product: %v", err)
    }

	productID := uuid.New().String()
	_, err = db.Exec("INSERT INTO products (id, canonical_product_id, name, inventory_id) VALUES ($1, $2, $3, $4)", productID, canonicalProductID, "Bench Product", inventoryID)
    if err != nil {
        b.Fatalf("Failed to create product: %v", err)
    }

	for i := 0; i < 100; i++ {
		id := uuid.New().String()
		variantIDs[i] = id
		_, err = db.Exec("INSERT INTO product_variants (id, product_id, variant_name) VALUES ($1, $2, $3)", id, productID, fmt.Sprintf("Variant %d", i))
		if err != nil {
			b.Fatalf("Failed to create variant: %v", err)
		}
	}

	// Create inventory products for benchmark
	for _, vID := range variantIDs {
		_, err = db.Exec("INSERT INTO inventory_products (inventory_id, product_variant_id, quantity) VALUES ($1, $2, $3)", inventoryID, vID, 10.0)
		if err != nil {
			b.Fatalf("Failed to create inventory product: %v", err)
		}
	}

	// Create transactions and items
	// 5000 transactions
	numTransactions := 5000
	txIDs := make([]string, numTransactions)

	// Prepare statement for faster insertion
	stmtTx, err := db.Prepare("INSERT INTO transactions (id, inventory_id, created_by_user_id) VALUES ($1, $2, $3)")
	if err != nil {
		b.Fatalf("Failed to prepare tx stmt: %v", err)
	}
	defer stmtTx.Close()

	stmtItem, err := db.Prepare("INSERT INTO transaction_items (id, transaction_id, product_variant_id, quantity, price_per_unit) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		b.Fatalf("Failed to prepare item stmt: %v", err)
	}
	defer stmtItem.Close()

	for i := 0; i < numTransactions; i++ {
		txID := uuid.New().String()
		txIDs[i] = txID
		_, err := stmtTx.Exec(txID, inventoryID, userID)
		if err != nil {
			b.Fatalf("Failed to insert transaction: %v", err)
		}

		// Add 5 items per transaction
		for j := 0; j < 5; j++ {
			itemID := uuid.New().String()
			variantID := variantIDs[rand.Intn(len(variantIDs))]
			_, err := stmtItem.Exec(itemID, txID, variantID, 1.0, 10.0)
			if err != nil {
				b.Fatalf("Failed to insert transaction item: %v", err)
			}
		}
	}

	_, err = db.Exec("ANALYZE transaction_items")
	if err != nil {
		b.Logf("Failed to analyze table: %v", err)
	}
	_, err = db.Exec("ANALYZE inventory_products")
	if err != nil {
		b.Logf("Failed to analyze table: %v", err)
	}

	return txIDs, variantIDs
}

func BenchmarkInventoryProductsByInventoryID(b *testing.B) {
	if testDB == nil {
		b.Skip("Database not initialized")
	}

	clearDB()
	b.StopTimer()
	// We need to capture the inventoryID, but seedBenchmarkData returns txIDs and variantIDs.
	// For simplicity, we'll just query the first inventory we find, as seedBenchmarkData creates only one.
	seedBenchmarkData(b, testDB)

	var inventoryID string
	err := testDB.QueryRow("SELECT id FROM inventories LIMIT 1").Scan(&inventoryID)
	if err != nil {
		b.Fatalf("Failed to get inventory ID: %v", err)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		rows, err := testDB.Query("SELECT * FROM inventory_products WHERE inventory_id = $1", inventoryID)
		if err != nil {
			b.Fatalf("Query failed: %v", err)
		}
		rows.Close()
	}
}

func BenchmarkTransactionItemsByTransactionID(b *testing.B) {
	if testDB == nil {
		b.Skip("Database not initialized")
	}

	clearDB() // Ensure clean state
	b.StopTimer()
	txIDs, _ := seedBenchmarkData(b, testDB)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		// Pick a random transaction ID
		// Use modulo to cycle through IDs if N > len(txIDs)
		// but using random might be more realistic distribution?
		// Actually sequential or cycled access is fine for benchmarking the index lookup.

		id := txIDs[i % len(txIDs)]

		rows, err := testDB.Query("SELECT * FROM transaction_items WHERE transaction_id = $1", id)
		if err != nil {
			b.Fatalf("Query failed: %v", err)
		}
		rows.Close()
	}
}

func BenchmarkTransactionItemsByVariantID(b *testing.B) {
	if testDB == nil {
		b.Skip("Database not initialized")
	}

	clearDB() // Ensure clean state
	b.StopTimer()
	_, variantIDs := seedBenchmarkData(b, testDB)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		id := variantIDs[i % len(variantIDs)]

		rows, err := testDB.Query("SELECT * FROM transaction_items WHERE product_variant_id = $1", id)
		if err != nil {
			b.Fatalf("Query failed: %v", err)
		}
		rows.Close()
	}
}
