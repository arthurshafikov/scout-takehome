# Scout: Greenhouse Pest & Disease Monitoring

A complete platform for greenhouse pest/disease monitoring with real-time photo analysis, bounding box visualization, and greenhouse floor mapping.

## ✅ Project Status

**Backend**: Complete (Go 1.25.5, SQLite, MinIO, Prometheus)  
**Frontend**: Complete (React 18, TypeScript, Vite, Redux Toolkit)  
**Tests**: All passing (7/7 backend smoke tests, frontend unit tests)  
**Build**: Zero errors (golangci-lint clean)

## 🚀 Quick Start (5 Minutes)

### Prerequisites
- Docker & Docker Compose only (everything runs in containers)

### Start Everything
```bash
docker compose up
```

### Access the App
- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080/api/thumbnails/{photoId}?w=600&q=85
- **MinIO Console**: http://localhost:9001 (minioadmin / minioadmin)
- **Prometheus Metrics**: http://localhost:8080/metrics

### Verify It Works
```bash
# Backend health
curl http://localhost:8080/api/healthz

# List photos (should return 10 photos with predictions)
curl http://localhost:8080/api/photos

# Filter by pest class
curl 'http://localhost:8080/api/photos?class_id=thrips'

# Filter by confidence
curl 'http://localhost:8080/api/photos?min_confidence=0.8'

# Get single photo
curl http://localhost:8080/api/photos/{photo_id}

# Get thumbnail
curl http://localhost:8080/api/thumbnails/{photo_id}?w=400 -o thumb.jpg
```

## 🔧 Setup from Scratch (New Clone)

Automated setup script handles everything:

```bash
# After cloning
./setup.sh
```

This script:
- ✓ Creates `backend/main.env` from template
- ✓ Creates `frontend/.env` from template  
- ✓ Validates all required files exist

Then start the application:
```bash
docker compose up
```

**Access URLs** (after ~30 seconds):
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080/api
- MinIO Console: http://localhost:9001

**Verify it works:**
```bash
curl http://localhost:8080/api/healthz
```

Optional: Re-seed dataset
```bash
cd backend
make seed
```

## 📊 Seed Initial Dataset

The database comes pre-seeded with 50 photos. If you need to re-seed:

```bash
cd backend
make seed
```

This uploads all images from `dataset/images/` to MinIO bucket `scout`.

## 🏗 Architecture

### Backend Stack
- **Language**: Go 1.25.5
- **Database**: SQLite (read-only, WAL mode)
- **Storage**: MinIO (S3-compatible)
- **HTTP**: Gin framework
- **Metrics**: Prometheus (custom scout_* metrics)

**Layer Design**:
```
HTTP Request
  ↓
Handler (validation, response formatting)
  ↓
Service (business logic)
  ↓
Repository (data access)
  ↓
SQLite/MinIO Storage
```

### Frontend Stack
- **Framework**: React 18 + TypeScript 5.4
- **Build Tool**: Vite 6.0
- **State Management**: Redux Toolkit + RTK Query
- **Styling**: Tailwind CSS 3.4
- **Canvas**: Konva.js (greenhouse map)
- **Tests**: Vitest + React Testing Library

### Data Flow
```
Gallery/Map Component
  ↓
Redux State (filters, pagination)
  ↓
RTK Query (API calls with caching)
  ↓
Backend API
  ↓
SQLite + MinIO
```

## 📡 API Endpoints

All endpoints under `/api/`:

```
GET  /healthz                          Health check
GET  /photos                           List photos (cursor paginated)
     ├─ ?cursor=...                    Pagination cursor
     ├─ ?limit=10                      Results per page
     ├─ ?class_id=thrips               Filter by pest class
     └─ ?min_confidence=0.8            Min confidence (0-1)
GET  /photos/{id}                      Get single photo
POST /photos/{id}/upload-link          Generate upload link
GET  /thumbnails/{id}?w=400&q=85       Get thumbnail
GET  /metrics                          Prometheus metrics
```

**Response Format**:
```json
{
  "success": true,
  "data": { /* response payload */ },
  "trace_id": "unique-id"
}
```

Errors return `success: false` with typed error codes and HTTP status codes.

## ⚙️ Configuration

### Backend (backend/main.env)
```env
APP_ENV=local
APP_PORT=8080
SQLITE_DB_PATH=../dataset/predictions.db
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=scout
```

### Frontend (frontend/.env)
```env
VITE_API_URL=http://localhost:8080
```

## 🐛 Troubleshooting

| Issue | Solution |
|-------|----------|
| "Cannot GET /" at :5173 | Frontend building, wait 30s |
| Connection refused on /api | Backend not started, check logs |
| No photos loading | Run `make seed` in backend dir |
| Images broken | Check MinIO console (http://localhost:9001) |
| Lint errors | Run `make lint` in backend |

## 📚 Detailed Documentation

- **Backend**: [backend/README.md](backend/README.md) - Setup, API reference, deployment
- **Frontend**: [frontend/README.md](frontend/README.md) - Architecture, phases, testing
- **API Contract**: [openapi.yaml](openapi.yaml)
