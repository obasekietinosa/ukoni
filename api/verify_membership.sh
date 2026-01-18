#!/bin/bash
set -e

BASE_URL="http://localhost:8080"
EMAIL_A="userA_$(date +%s)@test.com"
EMAIL_B="userB_$(date +%s)@test.com"
PASSWORD="password123"

echo "=== Registering User A: $EMAIL_A ==="
TOKEN_A=$(curl -s -X POST $BASE_URL/signup \
  -H "Content-Type: application/json" \
  -d "{\"name\": \"User A\", \"email\": \"$EMAIL_A\", \"password\": \"$PASSWORD\"}" | jq -r .token)

if [ "$TOKEN_A" == "null" ]; then
    echo "Signup A failed, logging in..."
    TOKEN_A=$(curl -s -X POST $BASE_URL/login \
      -H "Content-Type: application/json" \
      -d "{\"email\": \"$EMAIL_A\", \"password\": \"$PASSWORD\"}" | jq -r .token)
fi
echo "Token A: ${TOKEN_A:0:10}..."

echo "=== Registering User B: $EMAIL_B ==="
TOKEN_B=$(curl -s -X POST $BASE_URL/signup \
  -H "Content-Type: application/json" \
  -d "{\"name\": \"User B\", \"email\": \"$EMAIL_B\", \"password\": \"$PASSWORD\"}" | jq -r .token)

if [ "$TOKEN_B" == "null" ]; then
    echo "Signup B failed, logging in..."
    TOKEN_B=$(curl -s -X POST $BASE_URL/login \
      -H "Content-Type: application/json" \
      -d "{\"email\": \"$EMAIL_B\", \"password\": \"$PASSWORD\"}" | jq -r .token)
fi
echo "Token B: ${TOKEN_B:0:10}..."

echo "=== User A Creating Inventory ==="
INVENTORY_ID=$(curl -s -X POST $BASE_URL/inventories \
  -H "Authorization: Bearer $TOKEN_A" \
  -H "Content-Type: application/json" \
  -d '{"name": "Home Inventory"}' | jq -r .id)
echo "Inventory ID: $INVENTORY_ID"

echo "=== User A Inviting User B ==="
INVITE_RESP=$(curl -s -X POST $BASE_URL/inventories/$INVENTORY_ID/invitations \
  -H "Authorization: Bearer $TOKEN_A" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL_B\", \"role\": \"editor\"}")
INVITE_ID=$(echo $INVITE_RESP | jq -r .id)
echo "Invite ID: $INVITE_ID"

if [ "$INVITE_ID" == "null" ]; then
    echo "Failed to create invitation"
    echo $INVITE_RESP
    exit 1
fi

echo "=== User B Accepting Invite ==="
ACCEPT_RESP=$(curl -s -X POST $BASE_URL/invitations/$INVITE_ID/accept \
  -H "Authorization: Bearer $TOKEN_B" \
  -H "Content-Type: application/json" \
  -d '{}')
echo "Accept Resp: $ACCEPT_RESP"

echo "=== User A Listing Members ==="
MEMBERS=$(curl -s -X GET $BASE_URL/inventories/$INVENTORY_ID/members \
  -H "Authorization: Bearer $TOKEN_A")
echo "Members: $MEMBERS"

COUNT=$(echo $MEMBERS | jq '. | length')
if [ "$COUNT" -ge 1 ]; then
    echo "SUCCESS: Found $COUNT members"
else
    echo "FAILURE: Expected members, found none"
    exit 1
fi
