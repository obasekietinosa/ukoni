#!/bin/bash
set -e

BASE_URL="http://localhost:8080"
EMAIL="test@example.com"
PASSWORD="password123"

# 1. Login
echo "Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}")

# Extract token (assuming JSON response {"token": "..."} or checking structure)
# The output for login needs to be checked. For now I'll just print it.
echo "Login Response: $LOGIN_RESPONSE"

# Simple extraction if it's {"token": "..."}
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "Failed to extract token"
    exit 1
fi

echo "Token: $TOKEN"

# 2. Create Inventory
echo "Creating Inventory..."
CREATE_RESP=$(curl -s -X POST "$BASE_URL/inventories" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Home Inventory"}')
echo "Create Response: $CREATE_RESP"

INV_ID=$(echo $CREATE_RESP | grep -o '"id":"[^"]*' | cut -d'"' -f4)
if [ -z "$INV_ID" ]; then
    echo "Failed to extract inventory ID"
    exit 1
fi
echo "Inventory ID: $INV_ID"

# 3. List Inventories
echo "Listing Inventories..."
curl -s -X GET "$BASE_URL/inventories" \
  -H "Authorization: Bearer $TOKEN"
echo ""

# 4. Get Inventory by ID
echo "Getting Inventory $INV_ID..."
curl -s -X GET "$BASE_URL/inventories/$INV_ID" \
  -H "Authorization: Bearer $TOKEN"
echo ""
