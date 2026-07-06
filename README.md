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
- **Backend API**: http://localhost:8080/api
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

## 📁 Project Structure

```
scout-takehome/
├── backend/              Go API server
│   ├── cmd/app/          Entry point
│   ├── internal/         Core logic (models, repo, services, handlers)
│   ├── configs/          Configuration templates
│   ├── deployments/      Docker Compose
│   └── README.md         Full backend guide
├── frontend/             React web app
│   ├── src/
│   │   ├── features/     Gallery, filters, map
│   │   ├── services/     RTK Query API
│   │   ├── app/          Redux store
│   │   └── utils/        Helpers & constants
│   └── README.md         Full frontend guide
├── dataset/
│   ├── images/           50 sample photos
│   └── predictions.db    SQLite database
├── docker-compose.yml    Local dev environment
└── openapi.yaml          API specification
```

## 🔑 Key Features

- **Cursor Pagination**: Efficient scrolling through large photo datasets
- **Bounding Box Overlay**: SVG rendering with normalized [0,1] coordinates
- **Greenhouse Map**: Interactive Konva.js canvas showing photo positions
- **Smart Filtering**: By pest class and confidence threshold
- **Presigned URLs**: Secure photo uploads/downloads to MinIO
- **Thumbnail Caching**: On-demand generation with 24h cache (1-3ms hits vs 150ms generation)
- **Prometheus Metrics**: Production-ready observability
- **Error Boundaries**: React error catching prevents full app crash

## 🚢 Deployment

For production deployment, see detailed guides:
- Backend: [backend/README.md](backend/README.md#production-deployment-checklist)
- Frontend: [frontend/README.md](frontend/README.md#build--deployment)

## 📝 License

Proprietary - Scout Takehome Project
