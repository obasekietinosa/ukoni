package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createProductTestInventory(router *http.ServeMux, token string) string {
	payload := map[string]string{
		"name": "Product Test Inventory",
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

func TestProductCRUD(t *testing.T) {
	clearDB()
	router := setupRouter()
	token := createTestUser(router)
	inventoryID := createProductTestInventory(router, token)

	var productID string

	t.Run("Create Product", func(t *testing.T) {
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

		assert.Equal(t, http.StatusCreated, rr.Code)

		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.NotEmpty(t, response["id"])
		assert.Equal(t, "TestBrand", response["brand"])
		assert.Equal(t, "TestProduct", response["name"])
		productID = response["id"].(string)
	})

	t.Run("Get Product", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/products/"+productID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, productID, response["id"])
		assert.Equal(t, "TestProduct", response["name"])
	})

	t.Run("List Products", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/inventories/"+inventoryID+"/products?search=Test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var products []map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &products)

		assert.Len(t, products, 1)
		assert.Equal(t, "TestProduct", products[0]["name"])
	})

	t.Run("Update Product", func(t *testing.T) {
		payload := map[string]string{
			"brand":       "UpdatedBrand",
			"name":        "UpdatedProduct",
			"description": "Updated description",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("PUT", "/products/"+productID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, "UpdatedBrand", response["brand"])
		assert.Equal(t, "UpdatedProduct", response["name"])
	})

	t.Run("Create Variant", func(t *testing.T) {
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

		assert.Equal(t, http.StatusCreated, rr.Code)

		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.NotEmpty(t, response["id"])
		assert.Equal(t, "Variant1", response["variant_name"])
		assert.Equal(t, "SKU123", response["sku"])
	})

	t.Run("List Variants", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/products/"+productID+"/variants", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var variants []map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &variants)

		assert.Len(t, variants, 1)
		assert.Equal(t, "Variant1", variants[0]["variant_name"])
	})

	t.Run("Delete Product", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/products/"+productID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
	})

	t.Run("Get Deleted Product", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/products/"+productID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}
