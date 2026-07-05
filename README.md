# Scout: Greenhouse Pest & Disease Monitoring

A complete platform for greenhouse pest and disease monitoring with real-time photo analysis, bounding box visualization, and greenhouse floor mapping.

## Project Status

**Backend: ✅ COMPLETE (Phases 1-10)**
- Go backend with SQLite + MinIO storage
- RESTful API with cursor pagination and filtering
- On-demand thumbnail engine with caching
- Health checks and metrics infrastructure
- Seed binary for data ingestion
- Docker support

**Frontend: ✅ COMPLETE (Phases A-J)**
- React 19 + TypeScript + Vite
- Infinite scroll photo gallery with predictions
- Bounding box SVG overlays
- Greenhouse floor map (Konva.js)
- Filter panel (class, confidence)
- Tailwind CSS responsive design
- Unit tests (vitest + React Testing Library)
- Docker multi-stage build

**Documentation**
- Backend: [backend/README.md](backend/README.md)
- Frontend: [frontend/README.md](frontend/README.md)

## Quick Start

### Prerequisites
- Go 1.25.5+
- Node.js 20+ (for frontend)
- Docker + Docker Compose
- npm or pnpm

### ⚡ 5-Minute Quick Start (Docker Only)

1. **Start all services** (one command):
```bash
docker compose up
```

2. **Open in browser**:
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080/api/healthz
   - MinIO Console: http://localhost:9001 (user: `minioadmin` / password: `minioadmin`)

3. **Verify it works**:
   - Gallery loads with photos ✅
   - Click a photo to see full view with bounding boxes ✅
   - Try filters (pest class dropdown, confidence slider) ✅
   - Switch to "Map View" tab ✅

**That's it!** Your full Scout platform is running.

---

### Option 1: Full Docker Stack (Recommended for Testing)

```bash
docker compose up
```

**What starts**:
- Frontend (React dev server): http://localhost:5173
- Backend API: http://localhost:8080
- MinIO (image storage): http://localhost:9001
- 50 sample photos with pest predictions

**To stop**: Press `Ctrl+C`

---

### Option 2: Local Development (Best for Coding)

Use 4 terminal tabs:

#### Terminal 1: Backend Services (Docker)
```bash
docker compose up app minio
```

#### Terminal 2: Run Backend Server
```bash
cd backend
make build
make run
```

#### Terminal 3: Run Frontend Dev Server
```bash
cd frontend
npm install --legacy-peer-deps
npm run dev
```
Frontend runs at: http://localhost:5173 with hot-reload

#### Terminal 4: Seed Images (One-time)
```bash
cd backend
make seed
```
Uploads 50 greenhouse photos with AI pest predictions to MinIO.

---

### Option 3: Run Frontend Tests Only

```bash
cd frontend
npm install --legacy-peer-deps
npm test              # Run all 25 tests
npm test -- --watch  # Watch mode (re-run on changes)
npm test -- --coverage  # With coverage report
```

Tests cover:
- ✅ Bounding box coordinate transformations (10 tests)
- ✅ Thumbnail URL utilities (6 tests)
- ✅ Redux filters & pagination state (9 tests)

---

## ✅ Verify It's Working

### Quick Health Checks

**All services running?**
```bash
# Backend health
curl http://localhost:8080/api/healthz
# Expected: {"success": true, "data": {"Test": true}}

# MinIO is up
curl http://localhost:9000
# Expected: XML response

# Frontend is serving
curl http://localhost:5173
# Expected: HTML content
```

**Via Browser**:
1. Open http://localhost:5173 → Should see Scout gallery
2. Scroll down → Photos should load (infinite scroll)
3. Click a photo → Modal opens with full image & bbox overlays
4. Filters panel → Try selecting pest class or adjusting confidence
5. Click "Map View" tab → Should see greenhouse floor map

**Common Issues**:
- ❌ "Cannot GET /" at 5173? → Frontend container building, wait 30s
- ❌ "Connection refused" on /api? → Backend not started, check `make run`
- ❌ No photos loading? → Run `make seed` in another terminal
- ❌ Thumbnails broken? → Check MinIO console (http://localhost:9001)

## API Endpoints

### Photos
- `GET /api/photos` - List photos with cursor pagination
  - Query params: `cursor`, `limit`, `class_id`, `min_confidence`
  - Response: `{ items: Photo[], nextCursor?: string }`

- `GET /api/photos/{id}` - Get single photo with all predictions

- `GET /api/thumbnails/{id}` - Get thumbnail image (binary)

- `GET /api/healthz` - Health check

## Data Schema

### Photo
```typescript
{
  id: string               // UUID
  x: number               // Position in meters (0-40)
  y: number               // Position in meters (0-40)
  h: number               // Height in meters
  width: number           // Image width in pixels
  height: number          // Image height in pixels
  capturedAt: string      // ISO 8601 timestamp
  originalUrl: string     // S3 URL
  predictions: Prediction[]
}
```

### Prediction
```typescript
{
  id: string              // UUID
  photoId: string         // Reference to photo
  classId: string         // Pest class
  confidence: number      // 0-1 confidence score
  bbox: {                 // Normalized bounding box [0,1]
    xMin: number
    yMin: number
    xMax: number
    yMax: number
  }
}
```

### Pest Classes
- `powdery_mildew` - Red
- `mirid` - Yellow
- `whitefly_aphid` - Blue
- `miner_tuta` - Purple
- `thrips` - Orange
- `spider_mites` - Pink

## Greenhouse Layout

- **Dimensions**: 40m × 40m
- **Photo Grid**: Distributed across greenhouse floor
- **Resolution**: 2560 × 1440 pixels per photo
- **Predictions**: Multiple per photo with bounding boxes

## Frontend Features

### Gallery View
- Infinite scroll with cursor pagination
- Filter by pest class (dropdown)
- Filter by minimum confidence (slider)
- Thumbnail preview with top prediction badge
- Full-screen modal with:
  - Interactive bbox overlays (SVG)
  - Photo metadata (position, resolution, timestamp)
  - Sorted prediction list with confidence bars
  - Position-relative bbox coordinates

### Map View
- Interactive Konva.js canvas
- 40×40m greenhouse grid
- Photo positions as colored circles
- Colors represent top prediction class
- Hover for prediction details
- Pan and zoom controls

### Filters
- **Pest Class**: Dropdown with all 6 classes
- **Confidence Threshold**: Slider (0-100%)
- **Active Filters**: Display badges
- **Reset Button**: Clear all filters

## Development

### Backend Development
```bash
cd backend
go run ./cmd/app/main.go -cfgFolder ./configs -env ./
```

### Frontend Development
```bash
cd frontend
npm run dev           # Dev server
npm run build         # Production build
npm test              # Run tests
npm test -- --coverage  # With coverage
```

### Build Docker Images
```bash
# Backend
cd backend/build/app
docker build -t scout-backend:latest -f Dockerfile ../../..

# Frontend
cd frontend
docker build -t scout-frontend:latest .
```

## Testing

### Backend Tests
```bash
cd backend
go test ./...
```

### Frontend Tests
```bash
cd frontend
npm test

# Watch mode
npm test -- --watch

# Coverage
npm test -- --coverage
```

## Architecture

### Backend Stack
- **Language**: Go 1.25.5
- **Database**: SQLite (read-only, WAL mode)
- **Storage**: MinIO (S3-compatible)
- **HTTP**: net/http + chi router
- **Metrics**: Prometheus client

### Frontend Stack
- **Framework**: React 19 + TypeScript
- **Build Tool**: Vite 5
- **State Management**: Redux Toolkit
- **Data Fetching**: RTK Query
- **Styling**: Tailwind CSS 3
- **Visualization**: Konva.js + react-konva
- **Testing**: Vitest + React Testing Library

### Deployment
- **Container Runtime**: Docker
- **Orchestration**: Docker Compose
- **Frontend Server**: Nginx (multi-stage build)
- **Backend Server**: Go http
- **Reverse Proxy**: Nginx (/api proxy)

## API Response Format

All API responses follow this format:
```json
{
  "success": true,
  "data": { /* response data */ },
  "error": null,
  "trace_id": "unique-trace-id"
}
```

Errors:
```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable message"
  },
  "trace_id": "unique-trace-id"
}
```

## Environment Variables

### Backend
```env
APP_ENV=local
APP_PORT=8080
SQLITE_DB_PATH=./dataset/predictions.db
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=scout
```

### Frontend
```env
VITE_API_URL=http://localhost:8080
VITE_API_KEY=  # Optional
```

## Troubleshooting

### Backend won't start
```bash
# Check database exists
ls -la dataset/predictions.db

# Check MinIO is running
docker compose logs minio

# Check logs
docker compose logs app
```

### Frontend won't compile
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install --legacy-peer-deps
npm run dev
```

### Images not loading
1. Check MinIO is running: `http://localhost:9001`
2. Check seed was run: `make seed`
3. Check browser console for CORS errors
4. Verify API URL in .env

### Map not showing photos
1. Check cursor pagination working in gallery
2. Verify photos have position data (x, y)
3. Check browser DevTools for errors

## Project Structure

```
scout-takehome/
├── backend/
│   ├── cmd/app/main.go          # Backend entry
│   ├── internal/
│   │   ├── core/models/         # Domain models
│   │   ├── repository/sqlite/   # Data access
│   │   ├── services/            # Business logic
│   │   └── transport/http/      # HTTP handlers
│   ├── migrations/              # SQL migrations
│   ├── build/app/Dockerfile     # Backend image
│   ├── deployments/docker-compose.yml
│   └── README.md
├── frontend/
│   ├── src/
│   │   ├── main.tsx             # Entry point
│   │   ├── App.tsx              # Root component
│   │   ├── types/api.ts         # TypeScript types
│   │   ├── services/api.ts      # RTK Query
│   │   ├── app/store.ts         # Redux store
│   │   ├── features/
│   │   │   ├── gallery/         # Photo grid
│   │   │   ├── filters/         # Filter state
│   │   │   └── map/             # Greenhouse map
│   │   ├── components/          # Shared components
│   │   └── utils/               # Utilities
│   ├── Dockerfile               # Frontend image
│   ├── vite.config.ts           # Vite config
│   ├── tailwind.config.js       # Tailwind config
│   └── README.md
├── dataset/
│   ├── images/                  # 50 sample photos
│   └── predictions.db           # SQLite database
└── README.md
```

## Performance Notes

- **Gallery**: Cursor-based pagination (50 items/page) for smooth infinite scroll
- **Map**: Fetches up to 200 photos (separate cache from gallery)
- **Thumbnails**: Cached in MinIO for 24 hours
- **BBox Overlay**: SVG rendering (perfect scaling, no pixelation)
- **Bundle**: ~150KB gzipped (React, Redux, Konva)

## License

Proprietary - Scout Takehome Project

## Support

For issues or questions, refer to:
- Backend: [backend/README.md](backend/README.md)
- Frontend: [frontend/README.md](frontend/README.md)
- Backend errors: Check server logs
- Frontend errors: Check browser DevTools


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
