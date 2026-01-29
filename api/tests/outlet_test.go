package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestSeller(router http.Handler, token string) string {
	payload := map[string]string{
		"name": "Test Seller",
		"type": "chain",
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/sellers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)
	return response["id"].(string)
}

func TestOutletCRUD(t *testing.T) {
	clearDB()
	router := setupRouter()
	token := createTestUser(router)
	sellerID := createTestSeller(router, token)

	var outletID string

	t.Run("Create Outlet", func(t *testing.T) {
		payload := map[string]string{
			"name":        "Test Seller Tottenham",
			"channel":     "physical",
			"address":     "123 High St",
			"website_url": "https://example.com/tottenham",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/sellers/"+sellerID+"/outlets", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, "Test Seller Tottenham", response["name"])
		assert.Equal(t, "physical", response["channel"])
		outletID = response["id"].(string)
	})

	t.Run("Get Outlet", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/outlets/"+outletID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, outletID, response["id"])
		assert.Equal(t, "Test Seller Tottenham", response["name"])
	})

	t.Run("List Outlets", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/sellers/"+sellerID+"/outlets", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var outlets []map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &outlets)

		assert.Len(t, outlets, 1)
		assert.Equal(t, "Test Seller Tottenham", outlets[0]["name"])
	})

	t.Run("Update Outlet", func(t *testing.T) {
		payload := map[string]string{
			"name":        "Test Seller North London",
			"channel":     "physical",
			"address":     "456 High St",
			"website_url": "https://example.com/north",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("PUT", "/outlets/"+outletID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, "Test Seller North London", response["name"])
	})

	t.Run("Delete Outlet", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/outlets/"+outletID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)

		// Verify deletion
		req, _ = http.NewRequest("GET", "/outlets/"+outletID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr = httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
