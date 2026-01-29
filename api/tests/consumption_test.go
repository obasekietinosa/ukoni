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

func TestConsumptionEvents(t *testing.T) {
	if testDB == nil {
		t.Skip("skipping integration test: no database connection")
	}

	clearDB()
	router := setupRouter()

	token := createConsumptionTestUser(router)
	inventoryID := createConsumptionTestInventory(router, token)
	cpID := createConsumptionTestCanonicalProduct(router, token, inventoryID, "Milk")

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

	req, _ = http.NewRequest("GET", "/inventories/"+inventoryID+"/consumption-events", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}

	var events []*models.ConsumptionEvent
	json.Unmarshal(w.Body.Bytes(), &events)
	if len(events) != 1 {
		t.Errorf("expected 1 event, got %d", len(events))
	}
}
