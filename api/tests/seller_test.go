package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSellerCRUD(t *testing.T) {
	clearDB()
	router := setupRouter()
	token := createTestUser(router)

	var sellerID string

	t.Run("Create Seller", func(t *testing.T) {
		payload := map[string]string{
			"name": "Lidl",
			"type": "chain",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/sellers", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, "Lidl", response["name"])
		assert.Equal(t, "chain", response["type"])
		sellerID = response["id"].(string)
	})

	t.Run("Get Seller", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/sellers/"+sellerID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, sellerID, response["id"])
		assert.Equal(t, "Lidl", response["name"])
	})

	t.Run("List Sellers", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/sellers", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var sellers []map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &sellers)

		assert.Len(t, sellers, 1)
		assert.Equal(t, "Lidl", sellers[0]["name"])
	})

	t.Run("Update Seller", func(t *testing.T) {
		payload := map[string]string{
			"name": "Lidl UK",
			"type": "chain",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("PUT", "/sellers/"+sellerID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, "Lidl UK", response["name"])
	})

	t.Run("Delete Seller", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/sellers/"+sellerID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)

		// Verify deletion
		req, _ = http.NewRequest("GET", "/sellers/"+sellerID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr = httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code) // Should fail to find
	})
}
