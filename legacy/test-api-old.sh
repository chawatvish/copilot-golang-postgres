#!/bin/bash

# Gin REST API Test Script
# This script demonstrates all the API endpoints

echo "=== Gin REST API Test Script ==="
echo ""

BASE_URL="http://localhost:8080"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}1. Testing Health Endpoint${NC}"
curl -s "$BASE_URL/health" | jq .
echo ""

echo -e "${BLUE}2. Testing Root Endpoint${NC}"
curl -s "$BASE_URL/" | jq .
echo ""

echo -e "${BLUE}3. Getting All Users${NC}"
curl -s "$BASE_URL/api/v1/users" | jq .
echo ""

echo -e "${BLUE}4. Getting User by ID (ID: 1)${NC}"
curl -s "$BASE_URL/api/v1/users/1" | jq .
echo ""

echo -e "${BLUE}5. Creating New User${NC}"
curl -s -X POST "$BASE_URL/api/v1/users" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com"
  }' | jq .
echo ""

echo -e "${BLUE}6. Getting All Users Again (to see new user)${NC}"
curl -s "$BASE_URL/api/v1/users" | jq .
echo ""

echo -e "${BLUE}7. Updating User (ID: 2)${NC}"
curl -s -X PUT "$BASE_URL/api/v1/users/2" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Smith Updated",
    "email": "jane.updated@example.com"
  }' | jq .
echo ""

echo -e "${BLUE}8. Getting Updated User (ID: 2)${NC}"
curl -s "$BASE_URL/api/v1/users/2" | jq .
echo ""

echo -e "${BLUE}9. Deleting User (ID: 3)${NC}"
curl -s -X DELETE "$BASE_URL/api/v1/users/3" | jq .
echo ""

echo -e "${BLUE}10. Final User List${NC}"
curl -s "$BASE_URL/api/v1/users" | jq .
echo ""

echo -e "${GREEN}=== Test Complete ===${NC}"
