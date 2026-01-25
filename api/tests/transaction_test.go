package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createTransactionTestUser(router *http.ServeMux, email string) string {
	payload := map[string]string{
		"name":     "Transaction User",
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

func createTransactionTestInventory(router *http.ServeMux, token string) string {
	payload := map[string]string{
		"name": "Transaction Inventory",
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

func createTestVariant(t *testing.T, router *http.ServeMux, token string) string {
	// 1. Create Canonical Product
	cpPayload := map[string]string{"name": "Generic Milk"}
	cpBody, _ := json.Marshal(cpPayload)
	req, _ := http.NewRequest("POST", "/canonical-products", bytes.NewBuffer(cpBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)
	var cpResp map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &cpResp)
	canonicalProductID := cpResp["id"].(string)

	// 2. Create Product
	pPayload := map[string]string{
		"canonical_product_id": canonicalProductID,
		"name":                 "Milk",
		"brand":                "Sainsbury's",
	}
	pBody, _ := json.Marshal(pPayload)
	req, _ = http.NewRequest("POST", "/products", bytes.NewBuffer(pBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)
	var pResp map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &pResp)
	productID := pResp["id"].(string)

	// 3. Create Variant
	vPayload := map[string]interface{}{
		"variant_name": "2 Pints",
		"sku":          "123456",
		"size":         2.0,
		"unit":         "pints",
	}
	vBody, _ := json.Marshal(vPayload)
	req, _ = http.NewRequest("POST", "/products/"+productID+"/variants", bytes.NewBuffer(vBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)
	var vResp map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &vResp)
	return vResp["id"].(string)
}

func TestTransactionCRUD(t *testing.T) {
	clearDB()
	router := setupRouter()
	token := createTransactionTestUser(router, "transaction@example.com")
	inventoryID := createTransactionTestInventory(router, token)
	variantID := createTestVariant(t, router, token)

	var transactionID string

	t.Run("Create Transaction", func(t *testing.T) {
		transactionDate := time.Now().Format(time.RFC3339)
		price := 1.50
		quantity := 2.0

		payload := map[string]interface{}{
			"transaction_date": transactionDate,
			"items": []map[string]interface{}{
				{
					"product_variant_id": variantID,
					"quantity":           quantity,
					"price_per_unit":     price,
				},
			},
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/inventories/"+inventoryID+"/transactions", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.NotEmpty(t, response["id"])
		assert.Equal(t, inventoryID, response["inventory_id"])
		assert.Equal(t, 3.0, response["total_amount"]) // 2.0 * 1.50 = 3.0

		transactionID = response["id"].(string)
	})

	t.Run("List Transactions", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/inventories/"+inventoryID+"/transactions", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var list []map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &list)

		assert.Len(t, list, 1)
		assert.Equal(t, transactionID, list[0]["id"])
	})

	t.Run("Get Transaction", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/transactions/"+transactionID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, transactionID, response["id"])

		items := response["items"].([]interface{})
		assert.Len(t, items, 1)
		item := items[0].(map[string]interface{})
		assert.Equal(t, variantID, item["product_variant_id"])
		assert.Equal(t, 2.0, item["quantity"])
	})
}
