package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createShoppingListTestUser(router http.Handler, email string) string {
	payload := map[string]string{
		"name":     "Shopping User",
		"email":    email,
		"password": "password123",
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)
	return response["token"].(string)
}

func createTestInventory(router http.Handler, token string) string {
	payload := map[string]string{
		"name": "Shopping List Inventory",
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/inventories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)
	return response["id"].(string)
}

func TestShoppingListCRUD(t *testing.T) {
	clearDB()
	router := setupRouter()
	token := createShoppingListTestUser(router, "shopping@example.com")
	inventoryID := createTestInventory(router, token)

	var listID string
	var itemID string
	var canonicalProductID string
	var brandProductID string
	var variantID string

	// Create a canonical product first for testing item addition
	t.Run("Setup Product", func(t *testing.T) {
		payload := map[string]string{
			"name": "Test Canonical Product",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/inventories/"+inventoryID+"/canonical-products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)
		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)
		canonicalProductID = response["id"].(string)
	})

	t.Run("Setup Brand and Variant", func(t *testing.T) {
		// Create Brand Product linked to Canonical Product
		productPayload := map[string]string{
			"name":                 "Test Brand Product",
			"brand":                "Test Brand",
			"canonical_product_id": canonicalProductID,
		}
		body, _ := json.Marshal(productPayload)
		req, _ := http.NewRequest("POST", "/inventories/"+inventoryID+"/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)
		var prodResponse map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &prodResponse)
		brandProductID = prodResponse["id"].(string)

		// Create Variant
		variantPayload := map[string]interface{}{
			"variant_name": "Test Variant 1L",
			"size":         1.0,
			"unit":         "L",
		}
		body, _ = json.Marshal(variantPayload)
		req, _ = http.NewRequest("POST", "/products/"+brandProductID+"/variants", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr = httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)
		var varResponse map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &varResponse)
		variantID = varResponse["id"].(string)
	})

	t.Run("Create Shopping List", func(t *testing.T) {
		payload := map[string]string{
			"name": "Weekly Groceries",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/inventories/"+inventoryID+"/shopping-lists", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, "Weekly Groceries", response["name"])
		assert.Equal(t, inventoryID, response["inventory_id"])
		listID = response["id"].(string)
	})

	t.Run("List Shopping Lists", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/inventories/"+inventoryID+"/shopping-lists", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var lists []map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &lists)

		assert.Len(t, lists, 1)
		assert.Equal(t, "Weekly Groceries", lists[0]["name"])
	})

	t.Run("Get Shopping List", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/shopping-lists/"+listID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, listID, response["id"])
	})

	t.Run("Update Shopping List", func(t *testing.T) {
		payload := map[string]string{
			"name": "Monthly Groceries",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("PUT", "/shopping-lists/"+listID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, "Monthly Groceries", response["name"])
	})

	t.Run("Add Item to List (Canonical)", func(t *testing.T) {
		notes := "Buy 2"
		payload := map[string]interface{}{
			"target_type": "canonical_product",
			"target_id":   canonicalProductID,
			"notes":       notes,
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/shopping-lists/"+listID+"/items", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, listID, response["shopping_list_id"])
		assert.Equal(t, canonicalProductID, response["target_id"])
		assert.Equal(t, "canonical_product", response["target_type"])
		assert.Equal(t, notes, response["notes"])
		itemID = response["id"].(string)
	})

	t.Run("Add Item to List (Variant)", func(t *testing.T) {
		notes := "Buy Specific Brand"
		quantity := 2.5
		payload := map[string]interface{}{
			"target_type": "product_variant",
			"target_id":   variantID,
			"notes":       notes,
			"quantity":    quantity,
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/shopping-lists/"+listID+"/items", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, quantity, response["quantity"])
	})

	t.Run("List Items", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/shopping-lists/"+listID+"/items", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Logf("List Items failed with status %d: %s", rr.Code, rr.Body.String())
		}
		assert.Equal(t, http.StatusOK, rr.Code)

		var items []map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &items)

		if len(items) == 0 {
			t.Log("No items returned")
			t.FailNow()
		}

		assert.Len(t, items, 2)

		// Find the canonical item
		var canonicalItem, variantItem map[string]interface{}
		for _, item := range items {
			if item["target_type"] == "canonical_product" {
				canonicalItem = item
			} else if item["target_type"] == "product_variant" {
				variantItem = item
			}
		}

		// Verify canonical item
		assert.NotNil(t, canonicalItem)
		assert.Equal(t, "Buy 2", canonicalItem["notes"])
		cp := canonicalItem["canonical_product"].(map[string]interface{})
		assert.Equal(t, "Test Canonical Product", cp["name"])

		// Verify variant item
		assert.NotNil(t, variantItem)
		assert.Equal(t, "Buy Specific Brand", variantItem["notes"])
		assert.Equal(t, 2.5, variantItem["quantity"])
		pv := variantItem["product_variant"].(map[string]interface{})
		assert.Equal(t, "Test Variant 1L", pv["variant_name"])

		// Verify expanded Product details for variant
		assert.NotNil(t, variantItem["product"], "Product details should be present")
		if variantItem["product"] != nil {
			prod := variantItem["product"].(map[string]interface{})
			assert.Equal(t, "Test Brand Product", prod["name"])
			assert.Equal(t, "Test Brand", prod["brand"])
		}
	})

	t.Run("Update Item", func(t *testing.T) {
		newNotes := "Buy 5"
		newQuantity := 5.0
		payload := map[string]interface{}{
			"notes":    newNotes,
			"quantity": newQuantity,
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("PUT", "/shopping-list-items/"+itemID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, newNotes, response["notes"])
		assert.Equal(t, newQuantity, response["quantity"])
	})

	t.Run("Delete Item", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/shopping-list-items/"+itemID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Verify Item Deleted", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/shopping-lists/"+listID+"/items", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var items []map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &items)

		// Expect 1 item remaining (the variant item)
		assert.Len(t, items, 1)
		assert.Equal(t, "product_variant", items[0]["target_type"])
	})

	t.Run("Delete List", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/shopping-lists/"+listID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Verify List Deleted", func(t *testing.T) {
		// Get returns 404 or 500? Model returns error if not found?
		// Handler returns 500 if error (needs better error handling, but for now expect non-success)
		// Or GetList might return nil/err. Model GetList returns err if not found (sql.ErrNoRows).
		// Service returns err.
		req, _ := http.NewRequest("GET", "/shopping-lists/"+listID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		// Expecting 500 Internal Server Error because error handling in handler is generic
		// "sql: no rows in result set"
		// Ideally should be 404.
		assert.NotEqual(t, http.StatusOK, rr.Code)
	})
}
