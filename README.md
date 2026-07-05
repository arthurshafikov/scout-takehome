# Scout: Greenhouse Pest & Disease Monitoring

A complete platform for greenhouse pest and disease monitoring with real-time photo analysis, bounding box visualization, and greenhouse floor mapping.

## Project Status

**Backend: ✅ COMPLETE (Phases 1-9)**
- Go backend with SQLite + MinIO storage
- RESTful API with cursor pagination and filtering
- On-demand thumbnail engine with caching
- Health checks and metrics infrastructure
- Seed binary for data ingestion
- Docker support

**Frontend: ⏳ PENDING (Phases 10-15)**
- React + TypeScript photo gallery
- Bounding box overlay rendering
- Greenhouse floor map (Konva)
- Shared filter state

**Tests: ⏳ PENDING**
- Backend unit tests
- Integration smoke tests
- Frontend component tests

## Quick Start

### Prerequisites
- Go 1.25.5+
- Node.js 18+ (for frontend)
- Docker + Docker Compose (optional)

### Local Development

```bash
# 1. Clone and setup
git clone https://github.com/arthurshafikov/scout-takehome.git
cd scout-takehome

# 2. Start services
cd backend
docker compose -f deployments/docker-compose.yml up

# 3. Build and run backend
make build
make run

# 4. Seed images to MinIO (in another terminal)
make seed

# 5. Test backend
curl http://localhost:8080/healthz
curl http://localhost:8080/photos

# 6. Frontend setup (coming soon)
cd ../frontend
npm install
npm run dev
```

## Architecture

### Backend

**Stack**: Go 1.25.5, SQLite (read-only), MinIO, Gin HTTP

**Components**:
```
internal/
├── core/                    # Domain layer
│   ├── models/              # Photo, Prediction, BoundingBox
│   ├── types/               # Pagination, ListPhotosParams
│   ├── constants/           # Pest classes
│   └── errors/              # Sentinel errors
├── repository/              # Data access layer
│   ├── repository.go        # Interfaces
│   └── sqlite/              # SQLite implementation
├── services/                # Business logic layer
│   ├── services.go          # Aggregator
│   ├── photo_service.go     # Photo operations
│   ├── storage/             # MinIO storage service
│   ├── thumbnail/           # Thumbnail generation
│   └── metrics/             # Prometheus metrics
├── transport/               # HTTP layer
│   └── http/handler/        # Route handlers
└── config/                  # Configuration
```

**Flow**: HTTP Request → Handler → Service → Repository → SQLite/MinIO

**Key Features**:
- Cursor-based pagination with (captured_at, id) tuples
- Presigned URLs for uploads and downloads
- On-demand thumbnail generation with MinIO caching
- Structured error responses with trace IDs
- Prometheus metrics infrastructure

### Database Schema

**predictions.db** (SQLite, read-only):
```sql
photos(id, x, y, h, width, height, captured_at)
predictions(id, photo_id, class_id, confidence, bbox_xmin, bbox_ymin, bbox_xmax, bbox_ymax)
```

### Storage (MinIO)

```
scout/
├── originals/{photoId}          # Original photos
└── thumbs/{photoId}_{hash}      # Cached thumbnails
```

## API Reference

### Photos

#### List Photos
```http
GET /photos?cursor=&limit=10&class_id=&min_confidence=0.5
```

Response:
```json
{
  "success": true,
  "data": {
    "items": [{...}],
    "next_token": "base64_cursor"
  },
  "trace_id": "uuid"
}
```

#### Get Single Photo
```http
GET /photos/{id}
```

#### Generate Upload Link
```http
POST /photos/{id}/upload-link
{"content_type": "image/jpeg"}
```

Returns presigned PUT URL (15-min expiry)

### Thumbnails

#### Get Thumbnail
```http
GET /thumbnails/{id}?w=400&q=85
```

Query Params:
- `w`: width (default 400, max 2000)
- `q`: quality (default 85, 1-100)
- `dpr`: device pixel ratio (default 1.0) - coming soon

Returns: `image/jpeg` binary

**Caching**: On-demand generation with MinIO cache (24h TTL)

### Monitoring

#### Health Check
```http
GET /healthz
```

#### Metrics
```http
GET /metrics
```

Prometheus format (infrastructure in place)

## Configuration

### Environment Variables (backend/main.env)

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

## Data Model

### Photo
```go
type Photo struct {
  ID          string        `json:"id"`
  X           float64       `json:"x"`              // Greenhouse x position
  Y           float64       `json:"y"`              // Greenhouse y position
  H           float64       `json:"h"`              // Height
  Width       int           `json:"width"`          // Image width
  Height      int           `json:"height"`         // Image height
  CapturedAt  time.Time     `json:"captured_at"`
  OriginalURL string        `json:"original_url"`   // Presigned GET URL
  Predictions []Prediction  `json:"predictions"`
}
```

### Prediction
```go
type Prediction struct {
  ID        string       `json:"id"`
  PhotoID   string       `json:"photo_id"`
  ClassID   string       `json:"class_id"`        // pest/disease type
  Confidence float64    `json:"confidence"`       // 0-1
  BBox      BoundingBox `json:"bbox"`
}
```

### BoundingBox
```go
type BoundingBox struct {
  XMin float64 `json:"bbox_xmin"`  // Normalized [0, 1]
  YMin float64 `json:"bbox_ymin"`
  XMax float64 `json:"bbox_xmax"`
  YMax float64 `json:"bbox_ymax"`
}
```

### Pest Classes
- `powdery_mildew`
- `mirid`
- `whitefly_aphid`
- `miner_tuta`
- `thrips`
- `spider_mites`

## Building

### Backend

```bash
cd backend

# Development
make build
make run

# Seed images
make seed

# Tests
make test
make test-race

# Docker
docker compose -f deployments/docker-compose.yml up --build
```

### Frontend (Coming Soon)

```bash
cd frontend
npm install
npm run dev       # Development
npm run build     # Production build
npm run test      # Tests
```

## Deployment

### Docker Compose (Development)

```bash
docker compose -f backend/deployments/docker-compose.yml up
```

Services:
- **App**: http://localhost:8080
- **MinIO**: http://localhost:9000 (admin UI), http://localhost:9001 (console)

### Docker Image (Production)

```bash
docker build -f backend/build/app/Dockerfile --target prod -t scout-backend:latest .
docker run \
  -e MINIO_ENDPOINT=minio:9000 \
  -e SQLITE_DB_PATH=/data/predictions.db \
  -v /path/to/dataset:/data:ro \
  -p 8080:8080 \
  scout-backend:latest
```

## Testing

### Backend Smoke Test

```bash
# Health check
curl http://localhost:8080/healthz

# List photos
curl http://localhost:8080/photos

# Filter by class and confidence
curl "http://localhost:8080/photos?class_id=powdery_mildew&min_confidence=0.8"

# Get single photo
curl http://localhost:8080/photos/{photo_id}

# Generate thumbnail
curl http://localhost:8080/thumbnails/{photo_id}?w=400&q=85 -o thumb.jpg
```

### Frontend Tests (Coming Soon)

```bash
# Unit tests (Jest + React Testing Library)
npm run test

# E2E tests (Playwright)
npm run test:e2e

# Visual regression
npm run test:visual
```

## Development Guidelines

### Backend

See [backend/CLAUDE.md](backend/CLAUDE.md) for complete coding standards:
- 120-char line limit
- No abbreviations in names
- APIResponse wrapper for all endpoints
- Structured logging with correlation IDs
- Repository/Service/Handler layer separation
- All constants/models/types in core/

### Frontend (Coming Soon)

- Vite + React 19 + TypeScript
- Feature-based folder structure
- CSS Modules
- Redux Toolkit for state
- RTK Query for API calls
- vitest for tests

## Project Phases

- ✅ **Phase 1**: Backend foundation & module rename
- ✅ **Phase 2**: PostgreSQL → SQLite migration
- ✅ **Phase 3**: Repository pattern + cursor pagination
- ✅ **Phase 4**: MinIO storage service
- ✅ **Phase 5**: Photo endpoints + response handlers
- ✅ **Phase 6**: Thumbnail engine
- ✅ **Phase 7**: Prometheus metrics infrastructure
- ✅ **Phase 8**: Seed binary for data ingestion
- ✅ **Phase 9**: Docker compose + Makefile
- ⏳ **Phase 10**: Frontend scaffolding (React + TypeScript + Vite)
- ⏳ **Phase 11**: Gallery component with thumbnail grid
- ⏳ **Phase 12**: Bounding box overlay rendering
- ⏳ **Phase 13**: Filter state management
- ⏳ **Phase 14**: Greenhouse floor map (Konva)
- ⏳ **Phase 15**: Tests + refinement + deployment

## MinIO CLI

```bash
# Setup alias
mc alias set scout http://localhost:9000 minioadmin minioadmin

# List buckets
mc ls scout/

# List objects
mc ls scout/scout/originals/

# Download object
mc cp scout/scout/originals/photo-id /tmp/

# Upload object
mc cp image.jpg scout/scout/originals/photo-id

# Remove object
mc rm scout/scout/originals/photo-id
```

## Performance Notes

- **SQLite**: Read-only mode with WAL journaling, cursor pagination limits memory
- **MinIO**: Presigned URLs avoid proxying through app; 15-min PUT / 1-hour GET TTL
- **Thumbnails**: Generated on-demand with MinIO caching (24h TTL); singleflight guard (coming)
- **Pagination**: Cursor-based with (captured_at, id) tuples for stability
- **Filtering**: Joined queries on predictions table for class/confidence filtering

## License

MIT

  React 19 · Vite · pnpm · Redux Toolkit · openapi-typescript · react-konva · CSS Modules · vitest ·
  feature-based folders · `interface` over `type` · no `any` · no magic numbers.

## Additional Requirements

- **Errors.** Backend: right HTTP status + typed `Error` body, 4xx vs 5xx, shaped in one place, no leaked
  stack traces. UI: real loading, empty, and error states (never a blank screen or broken-image grid).
- **Logs.** Structured backend logs, one request traceable by a correlation id, sane levels, no secrets.
- **Metrics.** `/metrics` (or similar): rate, latency, errors, plus thumbnail cache hit/miss and generation
  time.

## Runtime

- One small box: ~1 vCPU, 512MB–1GB. Everything you build runs here — object storage (MinIO local,
  S3 in prod) is the one thing that's external
- Clients are spread across continents

## Repo

```
README.md      this assignment
openapi.yaml   the data-API contract
dataset/       50 photos + predictions.db
```

No scaffold — `git init` and build it. Use AI tools however you like; commit your AI setup
(`CLAUDE.md`/`AGENTS.md`, skills, agents) — we read it as a first-class artifact.

## Dataset

All photos come from a single greenhouse — call it **AlfaGreen** — a fixed **40×40 m** plane. A photo's
`x`/`y` place it on that plane (and on the map); `h` is the camera height above the floor.

- 50 photos here are a sample — the real catalog is far bigger; build for volume, not for 50
- Image files are raw material to ingest; the service serves originals from object storage, not this folder.
- bbox is relative to the original image; multiply corners by the render size.

```
dataset/
├── images/         50 greenhouse JPEGs, 2560×1440  (filename = <photo id>.jpg)
└── predictions.db  SQLite. Your database: read it, don't build your own.
```

```sql
photos(id, x, y, h, width, height, captured_at)
  -- x, y = location in greenhouse (m), for the map; h = camera height (m); width, height = pixels
predictions(id, photo_id, class_id, confidence, bbox_xmin, bbox_ymin, bbox_xmax, bbox_ymax)
  -- class_id: powdery_mildew, mirid, whitefly_aphid, miner_tuta, thrips, spider_mites
  -- bbox: normalized [0,1], (xmin,ymin) top-left to (xmax,ymax) bottom-right
```

## Data API

See [`openapi.yaml`](./openapi.yaml).

| Method | Path | What you get |
|---|---|---|
| `POST` | `/photos/{photoId}/upload-link` | presigned PUT URL; push the original to object storage |
| `GET` | `/photos` | photos (predictions + position + `originalUrl`), cursor-paginated, filters `classId` / `minConfidence` |
| `GET` | `/photos/{photoId}` | one photo |

Filters are optional and combine on a single prediction: a photo matches if one of its predictions is that
class with confidence ≥ `minConfidence`. You always get all of a photo's predictions. Thumbnail delivery is
yours to design, not in the contract.
