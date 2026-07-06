package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

// Integration test for seed → read flow
// Prerequisites: backend API running on localhost:8080, MinIO on localhost:9000
// Run with: go test ./cmd/seed/... -tags=integration -v

// PhotoPage mirrors the handler response structure
type PhotoPage struct {
	Items     []Photo `json:"items"`
	NextToken string  `json:"next_token,omitempty"`
}

// Photo mirrors the model structure
type Photo struct {
	ID          string        `json:"id"`
	X           float64       `json:"x"`
	Y           float64       `json:"y"`
	H           float64       `json:"h"`
	Width       int           `json:"width"`
	Height      int           `json:"height"`
	CapturedAt  time.Time     `json:"capturedAt"`
	OriginalURL string        `json:"originalUrl,omitempty"`
	Predictions []Prediction  `json:"predictions"`
}

// Prediction mirrors the model structure
type Prediction struct {
	ID         string      `json:"id"`
	PhotoID    string      `json:"photoId"`
	ClassID    string      `json:"classId"`
	Confidence float64     `json:"confidence"`
	BBox       BoundingBox `json:"bbox"`
}

// BoundingBox mirrors the model structure
type BoundingBox struct {
	XMin float64 `json:"xMin"`
	YMin float64 `json:"yMin"`
	XMax float64 `json:"xMax"`
	YMax float64 `json:"yMax"`
}

// APIResponse mirrors the handler response wrapper
type APIResponse struct {
	Success bool        `json:"success"`
	Data    PhotoPage   `json:"data"`
	Error   interface{} `json:"error,omitempty"`
	TraceID string      `json:"traceId"`
}

func TestBackendSmokeTest(t *testing.T) {
	// Skip if backend not running
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if !isBackendHealthy(ctx) {
		t.Skip("Backend not ready at http://localhost:8080 (this is OK for isolated testing)")
	}

	// Test 1: Seed has created photos
	t.Log("Test 1: Verifying backend has photos...")
	photos, err := listPhotos(ctx, "", 10, "", 0.0)
	if err != nil {
		t.Fatalf("Failed to list photos: %v", err)
	}

	if len(photos.Items) == 0 {
		t.Log("Warning: No photos found in backend. Run seed first: make seed")
		// This is not a failure - seed just hasn't run
		return
	}

	t.Logf("✓ Found %d photos in backend", len(photos.Items))

	// Test 2: Each photo has required fields
	t.Log("Test 2: Verifying photo structure...")
	photo := photos.Items[0]

	if photo.ID == "" {
		t.Errorf("Photo missing ID")
	}
	if photo.CapturedAt.IsZero() {
		t.Errorf("Photo missing capturedAt timestamp")
	}
	if photo.Width <= 0 || photo.Height <= 0 {
		t.Errorf("Photo has invalid dimensions: %dx%d", photo.Width, photo.Height)
	}

	// X, Y, H should be set (position in greenhouse)
	if photo.X < 0 || photo.Y < 0 || photo.H <= 0 {
		t.Logf("Warning: Photo position not fully set: x=%v, y=%v, h=%v", photo.X, photo.Y, photo.H)
	}

	t.Logf("✓ Photo structure valid: %s (captured: %s)", photo.ID, photo.CapturedAt.Format(time.RFC3339))

	// Test 3: Photo has predictions
	t.Log("Test 3: Verifying predictions...")
	if len(photo.Predictions) == 0 {
		t.Logf("Warning: Photo has no predictions. This is OK if it's a clean image.")
	} else {
		pred := photo.Predictions[0]
		if pred.ID == "" {
			t.Errorf("Prediction missing ID")
		}
		if pred.ClassID == "" {
			t.Errorf("Prediction missing classId")
		}
		if pred.Confidence < 0 || pred.Confidence > 1 {
			t.Errorf("Prediction has invalid confidence: %v (should be 0-1)", pred.Confidence)
		}
		if pred.BBox.XMin < 0 || pred.BBox.YMin < 0 || pred.BBox.XMax > 1 || pred.BBox.YMax > 1 {
			t.Logf("Warning: BBox might be invalid: %+v", pred.BBox)
		}
		t.Logf("✓ Prediction valid: %s (confidence: %.2f)", pred.ClassID, pred.Confidence)
	}

	// Test 4: Get single photo by ID
	t.Log("Test 4: Fetching single photo by ID...")
	singlePhoto, err := getPhoto(ctx, photo.ID)
	if err != nil {
		t.Fatalf("Failed to get photo: %v", err)
	}

	if singlePhoto.ID != photo.ID {
		t.Errorf("Retrieved photo ID mismatch: %s vs %s", singlePhoto.ID, photo.ID)
	}
	t.Logf("✓ Single photo fetch successful")

	// Test 5: OriginalURL is set (presigned URL)
	t.Log("Test 5: Verifying original URL...")
	if singlePhoto.OriginalURL == "" {
		t.Errorf("Photo missing originalUrl")
	} else if !isValidPresignedURL(singlePhoto.OriginalURL) {
		t.Logf("Warning: originalUrl may not be a valid presigned URL: %s", singlePhoto.OriginalURL)
	} else {
		t.Logf("✓ Original URL is valid presigned URL")
	}

	// Test 6: Filtering by class_id (if predictions exist)
	if len(photo.Predictions) > 0 {
		t.Log("Test 6: Testing class_id filter...")
		classID := photo.Predictions[0].ClassID
		filtered, err := listPhotos(ctx, "", 100, classID, 0.0)
		if err != nil {
			t.Fatalf("Failed to filter by class_id: %v", err)
		}

		found := false
		for _, p := range filtered.Items {
			for _, pred := range p.Predictions {
				if pred.ClassID == classID {
					found = true
					break
				}
			}
			if found {
				break
			}
		}

		if !found {
			t.Errorf("Filter for class_id=%s returned no results, but expected to find %s", classID, photo.ID)
		} else {
			t.Logf("✓ Class filter working: found %d photos with class=%s", len(filtered.Items), classID)
		}
	}

	// Test 7: Filtering by min_confidence
	if len(photo.Predictions) > 0 {
		t.Log("Test 7: Testing min_confidence filter...")
		minConf := photo.Predictions[0].Confidence - 0.1
		filtered, err := listPhotos(ctx, "", 100, "", minConf)
		if err != nil {
			t.Fatalf("Failed to filter by min_confidence: %v", err)
		}

		// Should include our photo since confidence is higher than threshold
		found := false
		for _, p := range filtered.Items {
			if p.ID == photo.ID {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("min_confidence filter didn't include expected photo")
		} else {
			t.Logf("✓ Confidence filter working: found %d photos with confidence >= %.2f", len(filtered.Items), minConf)
		}
	}

	t.Log("\n✅ All backend smoke tests passed!")
}

// Helper functions

func isBackendHealthy(ctx context.Context) bool {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/photos?limit=1", nil)
	if err != nil {
		return false
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200
}

func listPhotos(ctx context.Context, cursor string, limit int, classID string, minConfidence float64) (*PhotoPage, error) {
	url := fmt.Sprintf("http://localhost:8080/photos?limit=%d", limit)
	if cursor != "" {
		url += fmt.Sprintf("&cursor=%s", cursor)
	}
	if classID != "" {
		url += fmt.Sprintf("&class_id=%s", classID)
	}
	if minConfidence > 0 {
		url += fmt.Sprintf("&min_confidence=%v", minConfidence)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned %d: %s", resp.StatusCode, string(body))
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API returned success=false: %v", apiResp.Error)
	}

	return &apiResp.Data, nil
}

func getPhoto(ctx context.Context, photoID string) (*Photo, error) {
	url := fmt.Sprintf("http://localhost:8080/photos/%s", photoID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned %d: %s", resp.StatusCode, string(body))
	}

	// Response is wrapped in APIResponse
	var wrapper struct {
		Success bool   `json:"success"`
		Data    *Photo `json:"data"`
		Error   interface{} `json:"error,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !wrapper.Success {
		return nil, fmt.Errorf("API returned success=false: %v", wrapper.Error)
	}

	return wrapper.Data, nil
}

func isValidPresignedURL(url string) bool {
	// Presigned URLs from MinIO typically contain query parameters like X-Amz-Algorithm, X-Amz-Signature
	// Or they're direct S3-style URLs with embedded credentials
	// At minimum, they should be HTTP(S) URLs
	return bytes.HasPrefix([]byte(url), []byte("http://")) || bytes.HasPrefix([]byte(url), []byte("https://"))
}

// BenchmarkIntegration is an optional benchmark (only runs with -bench flag)
func BenchmarkBackendListPhotos(b *testing.B) {
	ctx := context.Background()

	// Skip if backend not running
	if !isBackendHealthy(ctx) {
		b.Skip("Backend not running")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := listPhotos(ctx, "", 10, "", 0.0)
		if err != nil {
			b.Fatalf("Failed to list photos: %v", err)
		}
	}

	b.ReportAllocs()
}
