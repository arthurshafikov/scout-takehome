# Scout: Greenhouse Pest & Disease Monitoring

A complete platform for greenhouse pest/disease monitoring with real-time photo analysis, bounding box visualization, and greenhouse floor mapping.

## ✅ Project Status

**Backend**: Complete (Go 1.25.5, SQLite, MinIO, Prometheus)  
**Frontend**: Complete (React 18, TypeScript, Vite, Redux Toolkit)  
**Tests**: All passing (7/7 backend smoke tests, frontend unit tests)  
**Build**: Zero errors (golangci-lint clean)

## 🔧 Setup from Scratch (New Clone)

### Step 1: Configure Environment

Run the automated setup script:

```bash
# After cloning
./setup.sh
```

This script:
- ✓ Creates `backend/main.env` from template
- ✓ Creates `frontend/.env` from template  
- ✓ Validates all required files exist

### Step 2: Start Services

```bash
docker compose up
```

Wait for containers to be healthy (~30 seconds):
- Backend: `scout-backend` (healthy)
- Frontend: `scout-frontend` (healthy)  
- MinIO: `scout-minio` (running)

### Step 3: Seed Initial Dataset (REQUIRED)

After services are healthy, seed the database with 50 sample photos:

```bash
cd backend
make seed
```

This uploads images from `dataset/images/` to MinIO and creates database records. Without this step, the gallery will be empty.

### Step 4: Access the Application

**Access URLs**:
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080/api
- MinIO Console: http://localhost:9001 (credentials: `minioadmin` / `minioadmin`)

**Verify it works:**
```bash
curl -H 'X-API-Key: scout-api-key-12345' http://localhost:8080/photos
```

You should see a JSON response with photo data.

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
