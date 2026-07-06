# Scout API - Error Scenarios & Testing Guide

This document lists all error scenarios and their expected HTTP status codes, along with test results.

## Test Results Summary

✅ **All 15 error handling tests passed successfully**

```
Passed: 15
Failed: 0
Total: 15
```

All error scenarios are properly handled with appropriate HTTP status codes and error response formats.

## Authentication Errors

### 1. Missing X-API-Key Header
- **Endpoint:** Any protected endpoint
- **Expected Status:** `401 Unauthorized`
- **Error Code:** `unauthorized`
- **Error Message:** `Missing X-API-Key header`
- **Status:** ✅ PASS
- **Test:**
  ```bash
  curl http://localhost:8080/photos
  ```
- **Response:**
  ```json
  {
    "success": false,
    "error": {
      "code": "unauthorized",
      "message": "Missing X-API-Key header"
    }
  }
  ```

### 2. Invalid/Wrong API Key
- **Endpoint:** Any protected endpoint  
- **Expected Status:** `401 Unauthorized`
- **Error Code:** `unauthorized`
- **Error Message:** `Invalid API key`
- **Status:** ✅ PASS
- **Test:**
  ```bash
  curl -H "X-API-Key: wrong-key" http://localhost:8080/photos
  ```
- **Response:**
  ```json
  {
    "success": false,
    "error": {
      "code": "unauthorized",
      "message": "Invalid API key"
    }
  }
  ```

---

## GET /photos - Photo Listing

### 3. Valid Request with Authentication
- **Expected Status:** `200 OK`
- **Status:** ✅ PASS
- **Test:**
  ```bash
  curl -H "X-API-Key: scout-api-key-12345" \
    "http://localhost:8080/photos?limit=5"
  ```

### 4. Invalid Limit Parameter (Not Numeric)
- **Expected Status:** `400 Bad Request`
- **Error Code:** `invalid_limit`
- **Error Message:** `Invalid limit parameter`
- **Status:** ✅ PASS
- **Test:**
  ```bash
  curl -H "X-API-Key: scout-api-key-12345" \
    "http://localhost:8080/photos?limit=invalid"
  ```
- **Response:**
  ```json
  {
    "success": false,
    "error": {
      "code": "invalid_limit",
      "message": "Invalid limit parameter"
    }
  }
  ```

### 5. Invalid Min Confidence Parameter (Not Numeric)
- **Expected Status:** `400 Bad Request`
- **Error Code:** `invalid_min_confidence`
- **Error Message:** `Invalid min_confidence parameter`
- **Status:** ✅ PASS
- **Test:**
  ```bash
  curl -H "X-API-Key: scout-api-key-12345" \
    "http://localhost:8080/photos?min_confidence=invalid"
  ```
- **Response:**
  ```json
  {
    "success": false,
    "error": {
      "code": "invalid_min_confidence",
      "message": "Invalid min_confidence parameter"
    }
  }
  ```

### 6. Valid Request with Filters
- **Expected Status:** `200 OK`
- **Status:** ✅ PASS
- **Test:**
  ```bash
  curl -H "X-API-Key: scout-api-key-12345" \
    "http://localhost:8080/photos?class_id=powdery_mildew&min_confidence=0.5"
  ```
- **Note:** Invalid class IDs are silently accepted, returning empty results

---

## GET /photos/:id - Photo Detail

### 7. Valid Photo ID
- **Expected Status:** `200 OK`
- **Status:** ✅ PASS
- **Test:**
  ```bash
  curl -H "X-API-Key: scout-api-key-12345" \
    "http://localhost:8080/photos/{valid_photo_id}"
  ```

### 8. Invalid Photo ID (Not Found)
- **Expected Status:** `404 Not Found`
- **Error Code:** `photo_not_found`
- **Error Message:** `Photo not found`
- **Status:** ✅ PASS (Fixed - was returning 500)
- **Test:**
  ```bash
  curl -H "X-API-Key: scout-api-key-12345" \
    "http://localhost:8080/photos/invalid-photo-id"
  ```
- **Response:**
  ```json
  {
    "success": false,
    "error": {
      "code": "photo_not_found",
      "message": "Photo not found"
    }
  }
  ```

---

## GET /thumbnails/:id - Thumbnail Generation

### 9. Valid Photo Thumbnail
- **Expected Status:** `200 OK`
- **Content-Type:** `image/jpeg`
- **Status:** ✅ PASS
- **Test:**
  ```bash
  curl -H "X-API-Key: scout-api-key-12345" \
    "http://localhost:8080/thumbnails/{valid_photo_id}?w=600&q=85" \
    -o thumbnail.jpg
  ```

### 10. Invalid Photo ID for Thumbnail
- **Expected Status:** `404 Not Found`
- **Error Code:** `photo_not_found`
- **Error Message:** `Photo not found`
- **Status:** ✅ PASS
- **Test:**
  ```bash
  curl -H "X-API-Key: scout-api-key-12345" \
    "http://localhost:8080/thumbnails/invalid-photo-id?w=600&q=85"
  ```
- **Response:**
  ```json
  {
    "success": false,
    "error": {
      "code": "photo_not_found",
      "message": "Photo not found"
    }
  }
  ```

### 11. Invalid Width Parameter (Not Numeric)
- **Expected Status:** `200 OK` (defaults to 400px)
- **Status:** ✅ PASS
- **Test:**
  ```bash
  curl -H "X-API-Key: scout-api-key-12345" \
    "http://localhost:8080/thumbnails/{valid_photo_id}?w=invalid"
  ```
- **Note:** Invalid width values are silently clamped to default

### 12. Width Parameter Out of Range (> 2000)
- **Expected Status:** `200 OK` (defaults to 400px)
- **Status:** ✅ PASS
- **Test:**
  ```bash
  curl -H "X-API-Key: scout-api-key-12345" \
    "http://localhost:8080/thumbnails/{valid_photo_id}?w=5000"
  ```
- **Note:** Out-of-range values are silently clamped to default

### 13. Invalid Quality Parameter (Not Numeric)
- **Expected Status:** `200 OK` (defaults to 85)
- **Status:** ✅ PASS
- **Test:**
  ```bash
  curl -H "X-API-Key: scout-api-key-12345" \
    "http://localhost:8080/thumbnails/{valid_photo_id}?q=invalid"
  ```
- **Note:** Invalid quality values are silently clamped to default

### 14. Quality Parameter Out of Range (> 100)
- **Expected Status:** `200 OK` (defaults to 85)
- **Status:** ✅ PASS
- **Test:**
  ```bash
  curl -H "X-API-Key: scout-api-key-12345" \
    "http://localhost:8080/thumbnails/{valid_photo_id}?q=150"
  ```
- **Note:** Out-of-range values are silently clamped to default

---

## Health Check (Public Endpoint)

### 15. Health Check Without Authentication
- **Expected Status:** `200 OK`
- **Status:** ✅ PASS
- **Test:**
  ```bash
  curl http://localhost:8080/healthz
  ```
- **Response:**
  ```json
  {
    "status": "ok"
  }
  ```

---

## Error Response Format

All error responses follow this standardized format:

```json
{
  "success": false,
  "error": {
    "code": "error_code",
    "message": "Human-readable error message"
  },
  "trace_id": "correlation-id-for-debugging"
}
```

---

## Success Response Format

All success responses follow this standardized format:

```json
{
  "success": true,
  "data": { ... },
  "trace_id": "correlation-id-for-debugging"
}
```

---

## HTTP Status Code Summary

| Status | Count | Scenarios |
|--------|-------|-----------|
| **200** | 6 | Valid requests |
| **400** | 3 | Invalid parameter values (non-numeric, required fields) |
| **401** | 2 | Missing or invalid authentication |
| **404** | 2 | Resource not found (photo doesn't exist) |
| **500** | 0 | None - all errors properly handled |

---

## Error Codes Reference

| Error Code | HTTP Status | Description |
|-----------|------------|-------------|
| `unauthorized` | 401 | Missing or invalid API key |
| `invalid_limit` | 400 | Limit parameter is not numeric or out of range |
| `invalid_min_confidence` | 400 | Confidence parameter is not numeric |
| `photo_not_found` | 404 | Requested photo doesn't exist in database |
| `internal_error` | 500 | Unexpected server error (should not occur) |

---

## Error Handling Architecture

### Authentication Layer (Middleware)
- Checks for X-API-Key header presence
- Validates API key against configured value
- Returns 401 immediately if auth fails
- Prevents access to all protected endpoints

### Request Validation (Handlers)
- Parses and validates query parameters
- Returns 400 for invalid formats (non-numeric, missing required fields)
- Silently clamps out-of-range values to defaults
- Prevents invalid requests from reaching service layer

### Business Logic Layer (Services)
- Wraps repository errors with context
- Preserves error types through layers
- Enables proper error mapping in response layer

### Response Layer (Handlers)
- Maps domain errors to HTTP status codes
- Converts internal errors to user-friendly messages
- Adds trace IDs for debugging and correlation
- Never leaks stack traces or internal details

---

## Error Handling Improvements Made

### Issue Fixed: Invalid Photo ID Returning 500 Instead of 404
- **Root Cause:** GET /photos/:id handler didn't properly handle "photo not found" errors from repository
- **Fix Applied:** Repository now returns sentinel error `ErrPhotoNotFound` when photo doesn't exist
- **Result:** Invalid photo IDs now return proper 404 with error response

### File Modified:
- `backend/internal/repository/sqlite/photo_repository.go` - Updated GetPhotoByID to return ErrPhotoNotFound

---

## Testing

### Running Error Test Suite

To run the comprehensive error handling test suite:

```bash
cd /home/arthur/Projects/scout-takehome
python3 test-errors.py
```

The test suite covers all 15 error scenarios and validates:
- Correct HTTP status codes
- Proper error response format
- Accurate error messages and codes
- All authentication paths
- Parameter validation
- Resource not found scenarios

### Expected Output

```
==================================================
Scout API - Error Handling Test Suite
==================================================
...
==================================================
Test Summary
==================================================
Passed: 15
Failed: 0
Total: 15

✓ All tests passed!
```

---

## Notes on Error Handling Strategy

1. **Fail-Safe Design:** Errors default to 500 if not explicitly mapped (defensive programming)
2. **Parameter Validation:** Invalid parameters are clamped to safe defaults rather than failing, improving UX
3. **No Information Leakage:** Error messages never reveal internal implementation details
4. **Traceability:** All errors include trace_id for server-side correlation and debugging
5. **Consistent Format:** All errors return same JSON structure for easy client-side handling
6. **Layered Handling:** Each layer (middleware → handler → service → repository) handles appropriate errors

