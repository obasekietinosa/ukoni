package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ukoni/internal/models"
)

func createConsumptionTestUser(router http.Handler) string {
	payload := map[string]string{
		"name":     "Consumption User",
		"email":    "consumption@example.com",
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

func createConsumptionTestInventory(router http.Handler, token string) string {
	payload := map[string]string{
		"name": "Consumption Inventory",
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

func createConsumptionTestCanonicalProduct(router http.Handler, token, inventoryID, name string) string {
	payload := map[string]string{
		"name": name,
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/inventories/"+inventoryID+"/canonical-products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)
	return response["id"].(string)
}

func createConsumptionTestProduct(router http.Handler, token, inventoryID string) string {
	payload := map[string]string{
		"brand":       "TestBrand",
		"name":        "TestProduct",
		"description": "A test product",
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/inventories/"+inventoryID+"/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)
	return response["id"].(string)
}

func createConsumptionTestVariant(router http.Handler, token, productID string) string {
	payload := map[string]string{
		"variant_name": "Variant1",
		"sku":          "SKU123",
		"unit":         "box",
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/products/"+productID+"/variants", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)
	return response["id"].(string)
}

func addInventory(router http.Handler, token, inventoryID, variantID string, quantity float64) {
	payload := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"product_variant_id": variantID,
				"quantity":           quantity,
			},
		},
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/inventories/"+inventoryID+"/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
}

func getInventoryQuantity(router http.Handler, token, inventoryID, variantID string) float64 {
	req, _ := http.NewRequest("GET", "/inventories/"+inventoryID+"/inventory-products", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var items []map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &items)

	for _, item := range items {
		if item["product_variant_id"] == variantID {
			return item["quantity"].(float64)
		}
	}
	return 0
}

func TestConsumptionEvents(t *testing.T) {
	if testDB == nil {
		t.Skip("skipping integration test: no database connection")
	}

	clearDB()
	router := setupRouter()

	token := createConsumptionTestUser(router)
	inventoryID := createConsumptionTestInventory(router, token)
	cpID := createConsumptionTestCanonicalProduct(router, token, inventoryID, "Milk")

	t.Run("Create Canonical Consumption", func(t *testing.T) {
		qty := 1.5
		unit := "L"
		note := "Cereal"
		source := "manual"
		eventReq := map[string]interface{}{
			"canonical_product_id": cpID,
			"quantity":             qty,
			"unit":                 unit,
			"note":                 note,
			"source":               source,
			"consumed_at":          time.Now().Format(time.RFC3339),
		}
		body, _ := json.Marshal(eventReq)
		req, _ := http.NewRequest("POST", "/inventories/"+inventoryID+"/consumption-events", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
		}

		var createdEvent models.ConsumptionEvent
		json.Unmarshal(w.Body.Bytes(), &createdEvent)
		if createdEvent.ID == "" {
			t.Error("expected ID to be set")
		}
		if createdEvent.CanonicalProductID == nil || *createdEvent.CanonicalProductID != cpID {
			t.Errorf("expected canonical_product_id %s, got %v", cpID, createdEvent.CanonicalProductID)
		}
	})

	t.Run("List Events", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/inventories/"+inventoryID+"/consumption-events", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", w.Code)
		}

		var events []*models.ConsumptionEvent
		json.Unmarshal(w.Body.Bytes(), &events)
		if len(events) != 1 {
			t.Errorf("expected 1 event, got %d", len(events))
		}
	})

	t.Run("Consume Variant and Reduce Inventory", func(t *testing.T) {
		productID := createConsumptionTestProduct(router, token, inventoryID)
		variantID := createConsumptionTestVariant(router, token, productID)

		// Add inventory: 10
		addInventory(router, token, inventoryID, variantID, 10)
		initialQty := getInventoryQuantity(router, token, inventoryID, variantID)
		if initialQty != 10 {
			t.Fatalf("expected initial quantity 10, got %f", initialQty)
		}

		// Consume: 3
		qty := 3.0
		eventReq := map[string]interface{}{
			"product_variant_id": variantID,
			"quantity":           qty,
			"unit":               "box",
			"consumed_at":        time.Now().Format(time.RFC3339),
		}
		body, _ := json.Marshal(eventReq)
		req, _ := http.NewRequest("POST", "/inventories/"+inventoryID+"/consumption-events", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
		}

		// Check Inventory: Should be 7
		newQty := getInventoryQuantity(router, token, inventoryID, variantID)
		if newQty != 7 {
			t.Errorf("expected quantity 7, got %f", newQty)
		}

		// Verify event stored variant ID
		var createdEvent models.ConsumptionEvent
		json.Unmarshal(w.Body.Bytes(), &createdEvent)
		if createdEvent.ProductVariantID == nil || *createdEvent.ProductVariantID != variantID {
			t.Errorf("expected product_variant_id %s, got %v", variantID, createdEvent.ProductVariantID)
		}
	})
}
