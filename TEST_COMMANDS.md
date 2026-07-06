# Scout API - Error Test Command Reference

Quick reference guide with all curl commands to test each error scenario.

## Quick Test (All Scenarios)

```bash
cd /home/arthur/Projects/scout-takehome
python3 test-errors.py
```

Expected output:
```
Passed: 15
Failed: 0
✓ All tests passed!
```

---

## Manual Test Commands by Scenario

### Authentication Errors (401)

#### Scenario 1: Missing X-API-Key Header
```bash
curl -v http://localhost:8080/photos
# Expected: 401 Unauthorized
# Error: "Missing X-API-Key header"
```

#### Scenario 2: Invalid API Key
```bash
curl -v -H "X-API-Key: wrong-api-key" http://localhost:8080/photos
# Expected: 401 Unauthorized
# Error: "Invalid API key"
```

---

### GET /photos - Photo Listing

#### Scenario 3: Valid Request
```bash
curl -s -H "X-API-Key: scout-api-key-12345" \
  "http://localhost:8080/photos?limit=5" | jq '.data | keys'
# Expected: 200 OK
# Response: {"items": [...], "next_cursor": "..."}
```

#### Scenario 4: Invalid Limit (Non-Numeric)
```bash
curl -s -H "X-API-Key: scout-api-key-12345" \
  "http://localhost:8080/photos?limit=invalid" | jq '.'
# Expected: 400 Bad Request
# Error Code: "invalid_limit"
```

#### Scenario 5: Invalid Min Confidence (Non-Numeric)
```bash
curl -s -H "X-API-Key: scout-api-key-12345" \
  "http://localhost:8080/photos?min_confidence=invalid" | jq '.'
# Expected: 400 Bad Request
# Error Code: "invalid_min_confidence"
```

#### Scenario 6: Valid Request with Filters
```bash
curl -s -H "X-API-Key: scout-api-key-12345" \
  "http://localhost:8080/photos?class_id=powdery_mildew&min_confidence=0.5" | jq '.data.items | length'
# Expected: 200 OK
# Response: Array of filtered photos
```

---

### GET /photos/:id - Photo Detail

#### Scenario 7: Valid Photo ID
```bash
curl -s -H "X-API-Key: scout-api-key-12345" \
  "http://localhost:8080/photos/14068d2d-34ea-442e-8c5d-a0b1f771e0fc" | jq '.data.id'
# Expected: 200 OK
# Response: Photo object with details
```

#### Scenario 8: Invalid Photo ID (Not Found)
```bash
curl -s -H "X-API-Key: scout-api-key-12345" \
  "http://localhost:8080/photos/invalid-photo-id-xyz123" | jq '.error'
# Expected: 404 Not Found
# Error Code: "photo_not_found"
```

---

### GET /thumbnails/:id - Thumbnail Generation

#### Scenario 9: Valid Thumbnail Request
```bash
curl -s -H "X-API-Key: scout-api-key-12345" \
  "http://localhost:8080/thumbnails/14068d2d-34ea-442e-8c5d-a0b1f771e0fc?w=600&q=85" \
  --output test-thumbnail.jpg && file test-thumbnail.jpg
# Expected: 200 OK with JPEG image
```

#### Scenario 10: Invalid Photo ID for Thumbnail
```bash
curl -s -H "X-API-Key: scout-api-key-12345" \
  "http://localhost:8080/thumbnails/invalid-photo-id-xyz123?w=600&q=85" | jq '.error'
# Expected: 404 Not Found
# Error Code: "photo_not_found"
```

#### Scenario 11: Invalid Width Parameter (Non-Numeric)
```bash
curl -s -H "X-API-Key: scout-api-key-12345" \
  "http://localhost:8080/thumbnails/14068d2d-34ea-442e-8c5d-a0b1f771e0fc?w=invalid&q=85" \
  --output thumb-invalid-width.jpg && file thumb-invalid-width.jpg
# Expected: 200 OK with image (defaults to w=400)
```

#### Scenario 12: Width Out of Range (> 2000)
```bash
curl -s -H "X-API-Key: scout-api-key-12345" \
  "http://localhost:8080/thumbnails/14068d2d-34ea-442e-8c5d-a0b1f771e0fc?w=5000&q=85" \
  --output thumb-large-width.jpg && file thumb-large-width.jpg
# Expected: 200 OK with image (defaults to w=400)
```

#### Scenario 13: Invalid Quality Parameter (Non-Numeric)
```bash
curl -s -H "X-API-Key: scout-api-key-12345" \
  "http://localhost:8080/thumbnails/14068d2d-34ea-442e-8c5d-a0b1f771e0fc?q=invalid" \
  --output thumb-invalid-quality.jpg && file thumb-invalid-quality.jpg
# Expected: 200 OK with image (defaults to q=85)
```

#### Scenario 14: Quality Out of Range (> 100)
```bash
curl -s -H "X-API-Key: scout-api-key-12345" \
  "http://localhost:8080/thumbnails/14068d2d-34ea-442e-8c5d-a0b1f771e0fc?q=150" \
  --output thumb-high-quality.jpg && file thumb-high-quality.jpg
# Expected: 200 OK with image (defaults to q=85)
```

---

### Health Check (Public Endpoint)

#### Scenario 15: Health Check (No Auth Required)
```bash
curl -s http://localhost:8080/healthz | jq '.'
# Expected: 200 OK
# Response: {"status": "ok"}
```

---

## Testing Workflow

### Step 1: Get Valid Photo ID
```bash
PHOTO_ID=$(curl -s -H "X-API-Key: scout-api-key-12345" \
  "http://localhost:8080/photos?limit=1" | jq -r '.data.items[0].id')
echo "Using photo ID: $PHOTO_ID"
```

### Step 2: Test Each Error Scenario
```bash
# Test authentication errors
echo "=== Testing Authentication ==="
curl -H "X-API-Key: scout-api-key-12345" \
  "http://localhost:8080/photos?limit=5" | jq '.success'

# Test validation errors
echo "=== Testing Invalid Parameters ==="
curl -H "X-API-Key: scout-api-key-12345" \
  "http://localhost:8080/photos?limit=invalid" | jq '.error.code'

# Test not found errors
echo "=== Testing Not Found ==="
curl -H "X-API-Key: scout-api-key-12345" \
  "http://localhost:8080/photos/invalid-id" | jq '.error.code'
```

### Step 3: Verify Error Response Format
```bash
curl -s -H "X-API-Key: scout-api-key-12345" \
  "http://localhost:8080/photos?limit=invalid" | jq '{
    success: .success,
    error_code: .error.code,
    error_message: .error.message,
    has_trace_id: (.trace_id != null)
  }'

# Expected output:
# {
#   "success": false,
#   "error_code": "invalid_limit",
#   "error_message": "Invalid limit parameter",
#   "has_trace_id": true
# }
```

---

## Error Patterns Summary

| HTTP Status | Pattern | Example |
|---|---|---|
| **200** | Request processed successfully | Valid photo list, valid thumbnail |
| **400** | Client error - bad request | Invalid parameter format or value |
| **401** | Client error - unauthorized | Missing or invalid API key |
| **404** | Client error - not found | Photo ID doesn't exist |
| **500** | Server error (should not occur) | Unexpected internal error |

---

## Test Result Verification

### Check Status Code
```bash
curl -s -o /dev/null -w "%{http_code}" \
  -H "X-API-Key: scout-api-key-12345" \
  "http://localhost:8080/photos?limit=invalid"
# Should output: 400
```

### Check Error Code
```bash
curl -s -H "X-API-Key: scout-api-key-12345" \
  "http://localhost:8080/photos?limit=invalid" | \
  jq -r '.error.code'
# Should output: invalid_limit
```

### Check Trace ID Exists
```bash
curl -s -H "X-API-Key: scout-api-key-12345" \
  "http://localhost:8080/photos?limit=invalid" | \
  jq -r '.trace_id'
# Should output: UUID-like string
```

---

## Batch Testing Script

Save as `batch-test.sh`:

```bash
#!/bin/bash

API_KEY="scout-api-key-12345"
VALID_ID="14068d2d-34ea-442e-8c5d-a0b1f771e0fc"
INVALID_ID="invalid-id"

test_endpoint() {
  local desc=$1
  local url=$2
  local expected_code=$3
  
  code=$(curl -s -o /dev/null -w "%{http_code}" -H "X-API-Key: $API_KEY" "$url")
  if [ "$code" = "$expected_code" ]; then
    echo "✓ $desc: $code"
  else
    echo "✗ $desc: Expected $expected_code, got $code"
  fi
}

echo "Testing error scenarios..."
test_endpoint "Valid request" "http://localhost:8080/photos?limit=5" "200"
test_endpoint "Invalid limit" "http://localhost:8080/photos?limit=invalid" "400"
test_endpoint "Invalid confidence" "http://localhost:8080/photos?min_confidence=invalid" "400"
test_endpoint "Valid photo" "http://localhost:8080/photos/$VALID_ID" "200"
test_endpoint "Invalid photo" "http://localhost:8080/photos/$INVALID_ID" "404"
test_endpoint "Valid thumbnail" "http://localhost:8080/thumbnails/$VALID_ID?w=600" "200"
test_endpoint "Invalid photo thumbnail" "http://localhost:8080/thumbnails/$INVALID_ID?w=600" "404"
test_endpoint "Missing auth" "http://localhost:8080/photos" "401"
test_endpoint "Health check" "http://localhost:8080/healthz" "200"
```

Run with:
```bash
chmod +x batch-test.sh
./batch-test.sh
```

---

## Docker Service Verification

### Check Backend is Running
```bash
docker compose ps | grep scout-backend
# Should show "Up" status
```

### Check Backend Logs for Errors
```bash
docker compose logs app --tail 20
# Look for error messages or "HTTP server started"
```

### Verify Database Connectivity
```bash
docker compose exec app curl -s http://localhost:8080/healthz | jq '.'
# Should return: {"status": "ok"}
```

---

## Troubleshooting

### Backend Not Responding (Connection Refused)
```bash
# Check if services are running
docker compose ps

# Start services if needed
docker compose up -d

# Wait for backend to be healthy
sleep 5
docker compose ps
```

### Getting 500 Errors Instead of Expected Status
```bash
# Check backend logs
docker compose logs app --tail 50

# Verify error is properly mapped in response.go
grep -A 2 "errors.Is" backend/internal/transport/http/handler/response.go
```

### API Key Not Working
```bash
# Verify API key in backend
grep "apiKey" backend/internal/transport/http/handler/handler.go

# Check .env file
cat backend/main.env | grep -i api

# Verify middleware is applied to routes
grep -B 2 "authMiddleware" backend/internal/transport/http/handler/handler.go
```
