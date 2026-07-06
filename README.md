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
docker compose up --wait
```

The `--wait` flag automatically waits for all services to become healthy before returning.

### Step 3: Seed Initial Dataset (REQUIRED)

After services are healthy, seed the database with 50 sample photos:

```bash
cd backend && make seed && cd ../
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
