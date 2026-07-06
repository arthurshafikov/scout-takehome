# Scout Backend — Coding Conventions & Setup

This document outlines all conventions and project-specific setup for the Scout backend.

## General Conventions

### API Response Structure
Always use the `APIResponse` wrapper for all endpoints:
```go
type ServerResponse struct {
    Success bool      `json:"success"`
    Data    any       `json:"data,omitempty"`
    Error   *APIError `json:"error,omitempty"`
    Meta    any       `json:"meta,omitempty"`
}

type APIError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}
```

### Error Handling
- Return a custom `APIError` struct with a code and message
- Sentinel errors live in `internal/core/errors/errors.go`
- Errors are always wrapped with `fmt.Errorf("context: %w", err)` to preserve the call stack
- Never leak stack traces in API responses
- HTTP status codes are chosen based on error type: map all sentinels to correct 4xx or 5xx in one place (see handler error mapping)

### Logging
- Use `logrus` library for all logging
- Write informative logs for errors
- Include correlation ID (`X-Request-ID` header) in all log entries for request traceability
- Use structured logging with `logrus.WithFields()`
- Log levels: `Debug` for development info, `Info` for important events, `Warn` for recoverable issues, `Error` for failures, `Panic` only for unrecoverable states
- Never log secrets (API keys, credentials, passwords, etc.)

### Testing
- Follow AAA pattern: **Arrange**, **Act**, **Assert**
- Use mocks from `testify/mock`, not hand-written mocks
- Mocks must be in a `mocks/` folder inside the package being tested (e.g., `internal/services/mocks/`)
- Only implement methods actually used in the test — do not implement every method of the interface
- Mocks are generated with `mockgen` (see Makefile)

### Code Organization

#### Constants
All constants live in `internal/core/constants/`:
- `constants.go` — general constants
- `classes.go` — pest class enum values (powdery_mildew, mirid, whitefly_aphid, miner_tuta, thrips, spider_mites)
- `events/` — event type constants (if needed)

#### Errors
All errors live in `internal/core/errors/errors.go`:
- Sentinel errors (e.g., `ErrPhotoNotFound`, `ErrBadCursor`)
- Error mapping helpers

#### Helpers
All project-wide helpers live in `internal/core/helpers/`:
- Generic utility functions
- **BBox coordinate transformation** — convert normalized `[0,1]` boxes to pixel coords at any rendered size/DPR

#### Models
All domain models live in `internal/core/models/`:
- `photo.go` — Photo struct
- `prediction.go` — Prediction struct
- `bounding_box.go` — BoundingBox struct

#### Types
Additional types unrelated to models live in `internal/core/types/`:
- `server.go` — ServerResponse, APIError
- `cursor.go` — pagination cursor type

### Naming & Style

- **No abbreviations** (use `account`, not `acc`; `householdAccount` over `hs`), except universally accepted short forms (`id`, `url`, `db`, `ctx`, `err`, etc.)
- **Readability over brevity**
- **IDs are always UUIDs**, never integers
- Prefer **DTOs as function arguments** over multiple parameters
- **Interface** over `type` (use `interface` keyword for declarations)
- **No magic numbers** — define them as constants
- **No `any` type** unless unavoidable

### Formatting & Linting

- **Line length must never exceed 120 characters** — format code for readability
- **Always add a blank line before `return`** to visually separate it
- Follow `.golangci.yml` settings for linting and formatting
- Run `make lint` before committing

### Repository Layer

The repository layer must be as simple as possible:
- **Only SQLite queries and mapping DB models to application models**
- **No business logic** — all business logic belongs in services
- Direct `database/sql` queries (no GORM)
- Row mapping from SQLite to application models

### Handler Layer

The handler layer must be as simple as possible:
- **No business logic** — all business logic belongs in services
- Parse incoming data (query params, body, headers)
- Call services
- Form responses using `ServerResponse` wrapper
- Error mapping

### HTTP Framework

- Use Gin for the HTTP server
- Handler functions receive `*fasthttp.RequestCtx`
- Middleware: correlation ID, API key auth, structured logging

---

## Scout-Specific Setup

### Database: SQLite

**Location**: `predictions.db` (read-only)
- Path configured in `main.env` as `SQLITE_DB_PATH`
- Schema:
  ```sql
  photos(id TEXT PRIMARY KEY, x REAL, y REAL, h REAL, width INTEGER, height INTEGER, captured_at TEXT)
  predictions(id TEXT PRIMARY KEY, photo_id TEXT, class_id TEXT, confidence REAL, bbox_xmin REAL, bbox_ymin REAL, bbox_xmax REAL, bbox_ymax REAL)
  ```
- Opened in read-only mode via `sqlite://path?mode=ro&_journal_mode=WAL`
- No migrations — database is pre-populated and immutable

### Object Storage: MinIO

**Configuration** (in `main.env`):
- `MINIO_ENDPOINT`: MinIO server endpoint (e.g., `localhost:9000`)
- `MINIO_ACCESS_KEY`: Access key
- `MINIO_SECRET_KEY`: Secret key
- `MINIO_BUCKET`: Bucket name (e.g., `scout`)

**Usage**:
- Original photos uploaded to `originals/{photoId}` via presigned PUT links
- Thumbnails cached in `thumbnails/{photoId}_{hash}` (write-once)
- `GetOriginalURL()` returns presigned GET URL (1 hour TTL) or public URL

### API Authentication

**Static API Key**:
- Configured in `main.env` as `API_KEY`
- All requests require `X-API-Key` header matching the configured key
- Return `401 AuthenticationRequired` if missing or invalid

### Thumbnail Engine

**Endpoint**: `GET /thumbnails/{photoId}?w=<width>&dpr=<dpr>&q=<quality>`
- `w` (width): 1–2560 pixels
- `dpr` (device pixel ratio): 1.0–3.0 (default 1.0)
- `q` (quality): 1–100 JPEG quality (default 80)

**Caching**:
- Cache key: `SHA256(photoId + w + dpr + q)`
- Stored in MinIO under `thumbnails/` prefix
- Per-request singleflight prevents duplicate concurrent generation
- Metrics: cache hit/miss rate and generation time

**Generation**:
1. Check MinIO `thumbs/` for cached thumbnail → stream if found
2. If miss: fetch original from MinIO
3. Resize using `github.com/disintegration/imaging`
4. Store in MinIO under `thumbs/`
5. Stream to client

### Seed Client

**Binary**: `cmd/seed/main.go`
**Purpose**: Upload all 50 photos from `dataset/images/` to MinIO

**Configuration** (via environment variables):
- `API_URL`: Backend URL (e.g., `http://localhost:8080`)
- `API_KEY`: API key matching backend config
- `IMAGES_DIR`: Path to images folder (e.g., `../dataset/images`)
- `DB_PATH`: Path to predictions.db (e.g., `../dataset/predictions.db`)

**Behavior**:
- Reads all photos from `predictions.db`
- For each photo:
  1. Check if already in MinIO (idempotent)
  2. If missing, call `POST /photos/{id}/upload-link` to get presigned URL
  3. PUT the image file from `dataset/images/{photoId}.jpg`
- Print progress to stdout

**Idempotency**: Safe to run multiple times — skips photos already uploaded

---

## Core Dependencies

### Go Module

**Module path**: `github.com/arthurshafikov/scout-takehome/backend`

**Key packages**:
- `modernc.org/sqlite` — SQLite driver (pure Go, no CGO)
- `github.com/minio/minio-go/v7` — MinIO client
- `github.com/valyala/fasthttp` — HTTP server
- `github.com/disintegration/imaging` — Image resizing
- `github.com/prometheus/client_golang` — Prometheus metrics
- `github.com/sirupsen/logrus` — Structured logging
- `github.com/google/uuid` — UUID generation
- `golang.org/x/sync/singleflight` — Concurrency guard
- `github.com/stretchr/testify` — Testing utilities

---

## Configuration

**File**: `main.env`

**Variables**:
```env
APP_ENV=local                    # local, testing, or production
APP_DEBUG=true
APP_PORT=8080

SQLITE_DB_PATH=../dataset/predictions.db
API_KEY=your-api-key-here

MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=scout
MINIO_USE_SSL=false

THUMBNAIL_CACHE_TTL_HOURS=24
```

---

## Build & Run

### Local Development

```bash
cd backend
make run          # Build and run
make test         # Run tests
make lint         # Run linter
make mocks        # Generate mocks
```

### Docker

```bash
docker compose up --build    # Start all services (backend, frontend, MinIO)
docker compose run --rm seed  # Run seed client
```

---

## Notes for Future Implementation

- **Do not add PostgreSQL/GORM back** — SQLite is read-only for predictions, MinIO handles originals
- All error responses must include `request_id` (from correlation ID middleware)
- Thumbnail generation is the expensive operation — cache aggressively
- BBox coordinate transformation is critical: test at multiple sizes/DPRs
