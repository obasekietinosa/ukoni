package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategoryCRUD(t *testing.T) {
	clearDB()
	router := setupRouter()
	token := createTestUser(router)
	inventoryID := createProductTestInventory(router, token)

	var categoryID string

	t.Run("Create Category", func(t *testing.T) {
		payload := map[string]string{
			"name": "Dairy",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/inventories/"+inventoryID+"/categories", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.NotEmpty(t, response["id"])
		assert.Equal(t, "Dairy", response["name"])
		assert.Equal(t, inventoryID, response["inventory_id"])
		categoryID = response["id"].(string)
	})

	t.Run("Create Sub-Category", func(t *testing.T) {
		payload := map[string]interface{}{
			"name":               "Milk",
			"parent_category_id": categoryID,
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/inventories/"+inventoryID+"/categories", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, "Milk", response["name"])
		assert.Equal(t, categoryID, response["parent_category_id"])
	})

	t.Run("List Categories", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/inventories/"+inventoryID+"/categories", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var categories []map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &categories)

		assert.Len(t, categories, 2)
		// Sorted by name, so Dairy then Milk
		assert.Equal(t, "Dairy", categories[0]["name"])
		assert.Equal(t, "Milk", categories[1]["name"])
	})

	t.Run("Create Category without Auth", func(t *testing.T) {
		payload := map[string]string{
			"name": "Produce",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/inventories/"+inventoryID+"/categories", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		// No auth header
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Access Other Inventory Categories Forbidden", func(t *testing.T) {
		// To be safe, let's manually register a second user
		payload := map[string]string{
			"name":     "Other User",
			"email":    "other@example.com",
			"password": "password123",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		otherToken := resp["token"].(string)

		req, _ = http.NewRequest("GET", "/inventories/"+inventoryID+"/categories", nil)
		req.Header.Set("Authorization", "Bearer "+otherToken)
		rr = httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code)
	})
}
