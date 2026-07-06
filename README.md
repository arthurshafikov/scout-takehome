# Scout: Greenhouse Pest & Disease Monitoring

A complete platform for greenhouse pest/disease monitoring with real-time photo analysis, bounding box visualization, and greenhouse floor mapping.

## ✅ Project Status

**Backend**: Complete (Go 1.25.5, SQLite, MinIO, Prometheus)  
**Frontend**: Complete (React 18, TypeScript, Vite, Redux Toolkit)  
**Tests**: All passing (7/7 backend smoke tests, frontend unit tests)  
**Build**: Zero errors (golangci-lint clean)

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
