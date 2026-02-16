package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlanGroupShoppingList(t *testing.T) {
	clearDB()
	router := setupRouter()
	token := createTestUser(router)

	var inventoryID string
	var parentPlanID string
	var childPlanID string
	var productID string

	// Create Inventory
	t.Run("Create Inventory", func(t *testing.T) {
		payload := map[string]string{"name": "My Inventory"}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/inventories", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)
		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		inventoryID = resp["id"].(string)
	})

	// Create Product
	t.Run("Create Product", func(t *testing.T) {
		payload := map[string]string{
			"name": "Test Product",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/inventories/"+inventoryID+"/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)
		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		productID = resp["id"].(string)
	})

	// Create Parent Plan
	t.Run("Create Parent Plan", func(t *testing.T) {
		payload := map[string]interface{}{
			"title": "Weekly Plan",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/inventories/"+inventoryID+"/plans", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)
		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		parentPlanID = resp["id"].(string)
	})

	// Create Child Plan
	t.Run("Create Child Plan", func(t *testing.T) {
		payload := map[string]interface{}{
			"title":          "Monday Dinner",
			"parent_plan_id": parentPlanID,
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/inventories/"+inventoryID+"/plans", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)
		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		childPlanID = resp["id"].(string)
	})

	// Add Item to Child Plan
	t.Run("Add Item to Child Plan", func(t *testing.T) {
		qty := 2.0
		payload := map[string]interface{}{
			"target_type": "product",
			"target_id":   productID,
			"quantity":    &qty,
			"unit":        "pcs",
			"note":        "For lunch",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/plans/"+childPlanID+"/items", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	// Create Shopping List from Group
	t.Run("Create Shopping List from Group", func(t *testing.T) {
		payload := map[string]string{
			"name": "Group Shopping List",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/plans/"+parentPlanID+"/shopping-list", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)

		var list map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &list)
		assert.Equal(t, "Group Shopping List", list["name"])
		listID := list["id"].(string)

		// Verify items in the list
		req, _ = http.NewRequest("GET", "/shopping-lists/"+listID+"/items", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr = httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		var items []map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &items)
		assert.Len(t, items, 1)
		assert.Equal(t, "product", items[0]["target_type"])
		assert.Equal(t, productID, items[0]["target_id"])
		assert.Equal(t, 2.0, items[0]["quantity"])

		// Verify link to parent plan
		req, _ = http.NewRequest("GET", "/plans/"+parentPlanID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr = httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		var plan map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &plan)
		linkedLists := plan["shopping_lists"].([]interface{})
		assert.Contains(t, linkedLists, listID)
	})
}
