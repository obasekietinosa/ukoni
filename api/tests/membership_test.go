package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createSecondUser(router *http.ServeMux) (string, string) {
	payload := map[string]string{
		"name":     "Second User",
		"email":    "second@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)
	return response["token"].(string), response["user"].(map[string]interface{})["id"].(string)
}

func TestMembership(t *testing.T) {
	clearDB()
	router := setupRouter()

	// Owner creates an inventory
	ownerToken := createTestUser(router)

	payload := map[string]string{"name": "Shared Kitchen"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/inventories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ownerToken)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var invResp map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &invResp)
	inventoryID := invResp["id"].(string)

	// Second user to be invited
	invitedToken, invitedUserID := createSecondUser(router)
	var inviteID string
	var inviteToken string

	t.Run("Invite User", func(t *testing.T) {
		payload := map[string]string{
			"email": "second@example.com",
			"role":  "editor",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/inventories/"+inventoryID+"/invitations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)

		inviteID = response["id"].(string)
		inviteToken = response["token"].(string)
		assert.Equal(t, "second@example.com", response["email"])
		assert.Equal(t, "pending", response["status"])
		assert.NotEmpty(t, inviteToken)
	})

	t.Run("Accept Invitation", func(t *testing.T) {
		payload := map[string]string{
			"token": inviteToken,
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/invitations/"+inviteID+"/accept", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+invitedToken)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, "accepted", response["status"])
	})

	t.Run("List Members", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/inventories/"+inventoryID+"/members", nil)
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var members []map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &members)

		// Owner and Invited User
		assert.Len(t, members, 2)
	})

	t.Run("Remove Member", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/inventories/"+inventoryID+"/members/"+invitedUserID, nil)
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)

		// Verify member is gone
		reqList, _ := http.NewRequest("GET", "/inventories/"+inventoryID+"/members", nil)
		reqList.Header.Set("Authorization", "Bearer "+ownerToken)
		rrList := httptest.NewRecorder()
		router.ServeHTTP(rrList, reqList)

		var members []map[string]interface{}
		json.Unmarshal(rrList.Body.Bytes(), &members)
		assert.Len(t, members, 1) // Only owner left
	})
}
