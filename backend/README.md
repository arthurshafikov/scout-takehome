# Scout Backend

Greenhouse pest and disease monitoring platform backend written in Go.

## Overview

The Scout backend is a high-performance API server that:
- Reads pest/disease predictions from SQLite (predictions.db)
- Serves photo metadata with bounding box predictions
- Stores original photos and thumbnails in MinIO S3-compatible storage
- Generates thumbnails on-demand with caching
- Provides RESTful endpoints for frontend consumption

## Architecture

**Stack**: Go 1.25.5, fasthttp (transitioning), SQLite (read-only), MinIO, Prometheus metrics

**Core Components**:
- `internal/core/models` - Domain models (Photo, Prediction, BoundingBox)
- `internal/repository` - Data access layer (SQLite with cursor pagination)
- `internal/services` - Business logic (PhotoService, StorageService, ThumbnailGenerator)
- `internal/transport/http` - HTTP handlers and response formatting
- `internal/config` - Configuration management (YAML + env vars)

## Getting Started

### Prerequisites
- Go 1.25.5+
- MinIO server (local or containerized)
- SQLite predictions database at `../dataset/predictions.db`

### Local Development

1. **Setup environment**:
   ```bash
   cd backend
   cp main.env.example main.env
   # Edit main.env with your MinIO credentials
   ```

2. **Build**:
   ```bash
   make build
   # or: go build -o ./.bin/app cmd/app/main.go
   ```

3. **Run**:
   ```bash
   make run
   # or: ./.bin/app -cfgFolder ./configs -env ./
   ```

4. **Seed MinIO with images** (in separate terminal):
   ```bash
   make seed
   # or: ./.bin/seed -endpoint localhost:9000 -access-key minioadmin -secret-key minioadmin -bucket scout
   ```

5. **Test**:
   ```bash
   curl http://localhost:8080/healthz
   ```

### Docker Compose

```bash
docker compose -f deployments/docker-compose.yml up
```

Services:
- App: http://localhost:8080
- MinIO: http://localhost:9000 (admin), http://localhost:9001 (console)

## API Endpoints

### Health Check
```http
GET /healthz
```
Returns `{"status":"ok"}`

### List Photos
```http
GET /photos?cursor=&limit=10&class_id=&min_confidence=0.5
```

Query Parameters:
- `cursor` (string): Pagination cursor (base64-encoded JSON of {CapturedAt, ID})
- `limit` (int): Results per page (default 10, max 100)
- `class_id` (string): Filter by pest class (optional)
- `min_confidence` (float): Minimum prediction confidence (0-1, optional)

Response:
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "uuid",
        "x": 0.5,
        "y": 0.3,
        "h": 25.5,
        "width": 1920,
        "height": 1080,
        "captured_at": "2025-07-05T13:50:00Z",
        "original_url": "https://...",
        "predictions": [
          {
            "id": "uuid",
            "photo_id": "uuid",
            "class_id": "powdery_mildew",
            "confidence": 0.95,
            "bbox_xmin": 100,
            "bbox_ymin": 50,
            "bbox_xmax": 300,
            "bbox_ymax": 200
          }
        ]
      }
    ],
    "next_token": "base64_cursor"
  },
  "trace_id": "request-id"
}
```

### Get Single Photo
```http
GET /photos/{id}
```

Returns single photo with predictions and original URL.

### Generate Upload Link
```http
POST /photos/{id}/upload-link
Content-Type: application/json

{
  "content_type": "image/jpeg"
}
```

Response:
```json
{
  "success": true,
  "data": {
    "url": "https://...",
    "method": "PUT",
    "expires_at": "2025-07-05T14:05:00Z"
  },
  "trace_id": "request-id"
}
```

### Get Thumbnail
```http
GET /thumbnails/{id}?w=400&q=85
```

Query Parameters:
- `w` (int): Width in pixels (default 400, max 2000)
- `q` (int): JPEG quality 1-100 (default 85)

Returns: `image/jpeg` binary data

Caching: Thumbnails are cached in MinIO under `thumbs/{photoId}_{hash}` with 24-hour TTL.

### Metrics
```http
GET /metrics
```

Prometheus format metrics (TODO: full implementation).

## Configuration

### Environment Variables (main.env)

```env
# App
APP_ENV=local
APP_DEBUG=true
APP_PORT=8080

# Database
SQLITE_DB_PATH=../dataset/predictions.db

# API
API_KEY=scout-api-key-12345

# MinIO
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=scout
MINIO_USE_SSL=false

# Thumbnails
THUMBNAIL_CACHE_TTL_HOURS=24
```

## Pest Classes

- `powdery_mildew`
- `mirid`
- `whitefly_aphid`
- `miner_tuta`
- `thrips`
- `spider_mites`

## Error Handling

All API errors follow standard format:

```json
{
  "success": false,
  "error": {
    "code": "photo_not_found",
    "message": "Photo not found"
  },
  "trace_id": "request-id"
}
```

HTTP Status Codes:
- 200: Success
- 400: Bad request (invalid params, bad cursor)
- 401: Unauthorized
- 404: Not found (photo_not_found)
- 500: Server error

## Building

```bash
# Development build
make build

# Build seed binary
make build-seed
```

## Testing

```bash
# Unit tests
make test

# Health check
curl http://localhost:8080/healthz

# List photos
curl http://localhost:8080/photos

# Get single photo
curl http://localhost:8080/photos/{id}
```

## MinIO CLI Reference

```bash
# Connect to local MinIO
mc alias set scout http://localhost:9000 minioadmin minioadmin

# List buckets
mc ls scout/

# List objects
mc ls scout/scout/originals/
```
