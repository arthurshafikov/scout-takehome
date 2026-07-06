#!/usr/bin/env python3
"""
Scout API Error Handling Test Suite
Tests all error scenarios and validates HTTP status codes
"""

import requests
import json
import sys
from typing import Tuple, Optional

BASE_URL = "http://localhost:8080"
API_KEY = "scout-api-key-12345"
WRONG_KEY = "wrong-api-key"
VALID_PHOTO_ID = "14068d2d-34ea-442e-8c5d-a0b1f771e0fc"
INVALID_PHOTO_ID = "invalid-photo-id-xyz123"

# Colors for output
GREEN = '\033[92m'
RED = '\033[91m'
YELLOW = '\033[93m'
RESET = '\033[0m'

passed = 0
failed = 0

def test_endpoint(
    name: str,
    method: str,
    endpoint: str,
    expected_status: int,
    api_key: Optional[str] = None,
    expect_json: bool = True
) -> bool:
    """Test an endpoint and verify the status code"""
    global passed, failed
    
    print(f"{YELLOW}Testing: {name}{RESET}")
    
    url = f"{BASE_URL}{endpoint}"
    headers = {}
    if api_key is not None:
        headers["X-API-Key"] = api_key
    
    try:
        if method == "GET":
            response = requests.get(url, headers=headers, timeout=5)
        elif method == "POST":
            response = requests.post(url, headers=headers, timeout=5)
        else:
            raise ValueError(f"Unsupported method: {method}")
        
        if response.status_code == expected_status:
            try:
                body = response.json() if expect_json else response.text[:100]
            except:
                body = response.text[:100]
            
            print(f"{GREEN}✓ PASS{RESET} - Status: {response.status_code}")
            
            # Pretty print JSON if it's an error response
            if expect_json and "error" in str(body):
                try:
                    error_msg = json.dumps(response.json(), indent=2)
                    print(f"  Response:\n{json.dumps(response.json(), indent=2)}")
                except:
                    print(f"  Response: {body}")
            else:
                print(f"  Response: {str(body)[:150]}")
            passed += 1
            return True
        else:
            print(f"{RED}✗ FAIL{RESET} - Expected: {expected_status}, Got: {response.status_code}")
            try:
                print(f"  Response:\n{json.dumps(response.json(), indent=2)}")
            except:
                print(f"  Response: {response.text[:200]}")
            failed += 1
            return False
    except requests.exceptions.RequestException as e:
        print(f"{RED}✗ ERROR{RESET} - {e}")
        failed += 1
        return False
    except Exception as e:
        print(f"{RED}✗ ERROR{RESET} - {e}")
        failed += 1
        return False
    finally:
        print()

def main():
    print(f"{YELLOW}{'='*50}{RESET}")
    print(f"{YELLOW}Scout API - Error Handling Test Suite{RESET}")
    print(f"{YELLOW}{'='*50}{RESET}")
    print()
    
    # ============================================
    # Authentication Tests
    # ============================================
    print(f"{YELLOW}--- Authentication Errors ---{RESET}")
    print()
    
    test_endpoint(
        "Missing X-API-Key header",
        "GET",
        "/photos",
        401,
        api_key=None
    )
    
    test_endpoint(
        "Invalid API Key",
        "GET",
        "/photos",
        401,
        api_key=WRONG_KEY
    )
    
    # ============================================
    # GET /photos Tests
    # ============================================
    print(f"{YELLOW}--- GET /photos - Photo Listing ---{RESET}")
    print()
    
    test_endpoint(
        "Valid request with auth",
        "GET",
        "/photos?limit=5",
        200,
        api_key=API_KEY
    )
    
    test_endpoint(
        "Invalid limit parameter (non-numeric)",
        "GET",
        "/photos?limit=invalid",
        400,
        api_key=API_KEY
    )
    
    test_endpoint(
        "Invalid min_confidence parameter (non-numeric)",
        "GET",
        "/photos?min_confidence=invalid",
        400,
        api_key=API_KEY
    )
    
    test_endpoint(
        "Valid request with filters",
        "GET",
        "/photos?class_id=powdery_mildew&min_confidence=0.5",
        200,
        api_key=API_KEY
    )
    
    # ============================================
    # GET /photos/:id Tests
    # ============================================
    print(f"{YELLOW}--- GET /photos/:id - Photo Detail ---{RESET}")
    print()
    
    test_endpoint(
        "Valid photo ID",
        "GET",
        f"/photos/{VALID_PHOTO_ID}",
        200,
        api_key=API_KEY
    )
    
    test_endpoint(
        "Invalid photo ID (not found)",
        "GET",
        f"/photos/{INVALID_PHOTO_ID}",
        404,
        api_key=API_KEY
    )
    
    # ============================================
    # GET /thumbnails/:id Tests
    # ============================================
    print(f"{YELLOW}--- GET /thumbnails/:id - Thumbnail Generation ---{RESET}")
    print()
    
    test_endpoint(
        "Valid photo thumbnail",
        "GET",
        f"/thumbnails/{VALID_PHOTO_ID}?w=600&q=85",
        200,
        api_key=API_KEY,
        expect_json=False
    )
    
    test_endpoint(
        "Invalid photo ID for thumbnail",
        "GET",
        f"/thumbnails/{INVALID_PHOTO_ID}?w=600&q=85",
        404,
        api_key=API_KEY
    )
    
    test_endpoint(
        "Thumbnail with invalid width (defaults to 400)",
        "GET",
        f"/thumbnails/{VALID_PHOTO_ID}?w=invalid&q=85",
        200,
        api_key=API_KEY,
        expect_json=False
    )
    
    test_endpoint(
        "Thumbnail with width > 2000 (defaults to 400)",
        "GET",
        f"/thumbnails/{VALID_PHOTO_ID}?w=5000&q=85",
        200,
        api_key=API_KEY,
        expect_json=False
    )
    
    test_endpoint(
        "Thumbnail with invalid quality (defaults to 85)",
        "GET",
        f"/thumbnails/{VALID_PHOTO_ID}?q=invalid",
        200,
        api_key=API_KEY,
        expect_json=False
    )
    
    test_endpoint(
        "Thumbnail with quality > 100 (defaults to 85)",
        "GET",
        f"/thumbnails/{VALID_PHOTO_ID}?q=150",
        200,
        api_key=API_KEY,
        expect_json=False
    )
    
    # ============================================
    # Health Check (Public Endpoint)
    # ============================================
    print(f"{YELLOW}--- Health Check (Public Endpoint) ---{RESET}")
    print()
    
    test_endpoint(
        "Health check without auth",
        "GET",
        "/healthz",
        200,
        api_key=None
    )
    
    # ============================================
    # Summary
    # ============================================
    total = passed + failed
    print(f"{YELLOW}{'='*50}{RESET}")
    print(f"{YELLOW}Test Summary{RESET}")
    print(f"{YELLOW}{'='*50}{RESET}")
    print(f"{GREEN}Passed: {passed}{RESET}")
    print(f"{RED}Failed: {failed}{RESET}")
    print(f"{YELLOW}Total: {total}{RESET}")
    print()
    
    if failed == 0:
        print(f"{GREEN}✓ All tests passed!{RESET}")
        return 0
    else:
        print(f"{RED}✗ Some tests failed{RESET}")
        return 1

if __name__ == "__main__":
    sys.exit(main())
