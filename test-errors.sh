#!/bin/bash

# Scout API Error Handling Test Suite
# Tests all error scenarios and validates responses

set -e

BASE_URL="http://localhost:8080"
API_KEY="scout-api-key-12345"
WRONG_KEY="wrong-api-key"
VALID_PHOTO_ID="14068d2d-34ea-442e-8c5d-a0b1f771e0fc"
INVALID_PHOTO_ID="invalid-photo-id-xyz123"

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counter
PASSED=0
FAILED=0

# Test function
test_endpoint() {
    local name=$1
    local method=$2
    local endpoint=$3
    local expected_status=$4
    local auth_key=$5
    
    echo -e "${YELLOW}Testing: $name${NC}"
    
    if [ "$auth_key" = "NOAUTH" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$BASE_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" -H "X-API-Key: $auth_key" "$BASE_URL$endpoint")
    fi
    
    # Split response and status code
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" = "$expected_status" ]; then
        echo -e "${GREEN}✓ PASS${NC} - Status: $http_code"
        echo "  Response: $(echo "$body" | jq -c '.error // .data // .' 2>/dev/null || echo "$body" | head -c 100)"
        ((PASSED++))
    else
        echo -e "${RED}✗ FAIL${NC} - Expected: $expected_status, Got: $http_code"
        echo "  Response: $body"
        ((FAILED++))
    fi
    echo ""
}

echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}Scout API - Error Handling Test Suite${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

# ============================================
# Authentication Tests
# ============================================
echo -e "${YELLOW}--- Authentication Errors ---${NC}"
echo ""

test_endpoint \
    "Missing X-API-Key header" \
    "GET" \
    "/photos" \
    "401" \
    "NOAUTH"

test_endpoint \
    "Invalid API Key" \
    "GET" \
    "/photos" \
    "401" \
    "$WRONG_KEY"

# ============================================
# GET /photos Tests
# ============================================
echo -e "${YELLOW}--- GET /photos - Photo Listing ---${NC}"
echo ""

test_endpoint \
    "Valid request with auth" \
    "GET" \
    "/photos?limit=5" \
    "200" \
    "$API_KEY"

test_endpoint \
    "Invalid limit parameter (non-numeric)" \
    "GET" \
    "/photos?limit=invalid" \
    "400" \
    "$API_KEY"

test_endpoint \
    "Invalid min_confidence parameter (non-numeric)" \
    "GET" \
    "/photos?min_confidence=invalid" \
    "400" \
    "$API_KEY"

test_endpoint \
    "Valid with valid filters" \
    "GET" \
    "/photos?class_id=powdery_mildew&min_confidence=0.5" \
    "200" \
    "$API_KEY"

# ============================================
# GET /photos/:id Tests
# ============================================
echo -e "${YELLOW}--- GET /photos/:id - Photo Detail ---${NC}"
echo ""

test_endpoint \
    "Valid photo ID" \
    "GET" \
    "/photos/$VALID_PHOTO_ID" \
    "200" \
    "$API_KEY"

test_endpoint \
    "Invalid photo ID (not found)" \
    "GET" \
    "/photos/$INVALID_PHOTO_ID" \
    "404" \
    "$API_KEY"

# ============================================
# GET /thumbnails/:id Tests
# ============================================
echo -e "${YELLOW}--- GET /thumbnails/:id - Thumbnail Generation ---${NC}"
echo ""

test_endpoint \
    "Valid photo thumbnail" \
    "GET" \
    "/thumbnails/$VALID_PHOTO_ID?w=600&q=85" \
    "200" \
    "$API_KEY"

test_endpoint \
    "Invalid photo ID for thumbnail" \
    "GET" \
    "/thumbnails/$INVALID_PHOTO_ID?w=600&q=85" \
    "404" \
    "$API_KEY"

test_endpoint \
    "Thumbnail with invalid width (defaults to 400)" \
    "GET" \
    "/thumbnails/$VALID_PHOTO_ID?w=invalid&q=85" \
    "200" \
    "$API_KEY"

test_endpoint \
    "Thumbnail with width > 2000 (defaults to 400)" \
    "GET" \
    "/thumbnails/$VALID_PHOTO_ID?w=5000&q=85" \
    "200" \
    "$API_KEY"

test_endpoint \
    "Thumbnail with invalid quality (defaults to 85)" \
    "GET" \
    "/thumbnails/$VALID_PHOTO_ID?q=invalid" \
    "200" \
    "$API_KEY"

test_endpoint \
    "Thumbnail with quality > 100 (defaults to 85)" \
    "GET" \
    "/thumbnails/$VALID_PHOTO_ID?q=150" \
    "200" \
    "$API_KEY"

# ============================================
# Health Check (No Auth Required)
# ============================================
echo -e "${YELLOW}--- Health Check (Public Endpoint) ---${NC}"
echo ""

test_endpoint \
    "Health check without auth" \
    "GET" \
    "/healthz" \
    "200" \
    "NOAUTH"

# ============================================
# Summary
# ============================================
echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}Test Summary${NC}"
echo -e "${YELLOW}========================================${NC}"
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${RED}Failed: $FAILED${NC}"
echo -e "${YELLOW}Total: $((PASSED + FAILED))${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}✗ Some tests failed${NC}"
    exit 1
fi
