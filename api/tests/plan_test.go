package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlanCRUD(t *testing.T) {
	clearDB()
	router := setupRouter()
	token := createTestUser(router) // defined in inventory_test.go

	var inventoryID string
	var planID string
	var itemID string
	var shoppingListID string
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

	// Create Product (needed for plan item)
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

	// Create Shopping List (needed for linking)
	t.Run("Create Shopping List", func(t *testing.T) {
		payload := map[string]string{"name": "My List"}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/inventories/"+inventoryID+"/shopping-lists", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)
		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		shoppingListID = resp["id"].(string)
	})

	// Create Plan
	t.Run("Create Plan", func(t *testing.T) {
		payload := map[string]interface{}{
			"title":       "Weekly Meal Plan",
			"description": "Plan for the week",
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
		planID = resp["id"].(string)
		assert.Equal(t, "Weekly Meal Plan", resp["title"])
	})

	// List Plans
	t.Run("List Plans", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/inventories/"+inventoryID+"/plans", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		var plans []map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &plans)
		assert.Len(t, plans, 1)
		assert.Equal(t, planID, plans[0]["id"])
	})

	// Add Plan Item
	t.Run("Add Plan Item", func(t *testing.T) {
		qty := 2.0
		payload := map[string]interface{}{
			"target_type": "product",
			"target_id":   productID,
			"quantity":    &qty,
			"unit":        "pcs",
			"note":        "For lunch",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/plans/"+planID+"/items", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)
		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		itemID = resp["id"].(string)
		assert.Equal(t, "product", resp["target_type"])
	})

	// Link Shopping List
	t.Run("Link Shopping List", func(t *testing.T) {
		payload := map[string]string{
			"shopping_list_id": shoppingListID,
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/plans/"+planID+"/shopping-lists", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	// Get Plan Details
	t.Run("Get Plan Details", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/plans/"+planID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)

		assert.Equal(t, planID, resp["id"])
		items := resp["items"].([]interface{})
		assert.Len(t, items, 1)
		assert.Equal(t, itemID, items[0].(map[string]interface{})["id"])

		lists := resp["shopping_lists"].([]interface{})
		assert.Len(t, lists, 1)
		assert.Equal(t, shoppingListID, lists[0].(string))
	})

	// Update Plan Item
	t.Run("Update Plan Item", func(t *testing.T) {
		qty := 3.0
		payload := map[string]interface{}{
			"quantity": &qty,
			"note":     "Updated note",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("PUT", "/plan-items/"+itemID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.Equal(t, 3.0, resp["quantity"])
	})

	// Unlink Shopping List
	t.Run("Unlink Shopping List", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/plans/%s/shopping-lists/%s", planID, shoppingListID), nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNoContent, rr.Code)
	})

	// Remove Plan Item
	t.Run("Remove Plan Item", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/plan-items/"+itemID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNoContent, rr.Code)
	})

	// Create Child Plan
	t.Run("Create Child Plan", func(t *testing.T) {
		payload := map[string]interface{}{
			"title":          "Monday Dinner",
			"parent_plan_id": planID,
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
		assert.Equal(t, planID, resp["parent_plan_id"])
	})

	// Verify Child Plan in Get Plan
	t.Run("Verify Child Plan", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/plans/"+planID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)

		children := resp["children"].([]interface{})
		assert.Len(t, children, 1)
	})

	// Delete Plan
	t.Run("Delete Plan", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/plans/"+planID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNoContent, rr.Code)
	})
}
