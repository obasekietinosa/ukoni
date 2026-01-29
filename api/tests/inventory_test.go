package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestUser(router http.Handler) string {
	payload := map[string]string{
		"name":     "Inventory User",
		"email":    "inventory@example.com",
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

func TestInventoryCRUD(t *testing.T) {
	clearDB()
	router := setupRouter()
	token := createTestUser(router)

	var inventoryID string

	t.Run("Create Inventory", func(t *testing.T) {
		payload := map[string]string{
			"name": "My Kitchen",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/inventories", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, "My Kitchen", response["name"])
		inventoryID = response["id"].(string)
	})

	t.Run("List Inventories", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/inventories", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var inventories []map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &inventories)

		assert.Len(t, inventories, 1)
		assert.Equal(t, "My Kitchen", inventories[0]["name"])
	})

	t.Run("Get Inventory", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/inventories/"+inventoryID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, inventoryID, response["id"])
		assert.Equal(t, "My Kitchen", response["name"])
	})
}
