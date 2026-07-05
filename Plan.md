# Scout — Implementation Plan

## Phase 1: Backend Foundation

### 1.1 — Repo & module rename
- Rename Go module from `github.com/arthurshafikov/boilerplate` → `github.com/arthurshafikov/scout-takehome/backend` across all files
- Update `APP_NAME` in Makefile to `scout`

### 1.2 — Strip PostgreSQL, add SQLite + MinIO
- Remove GORM, pgsql package, goose migrations, PostgreSQL docker-compose service
- Add `modernc.org/sqlite` (pure Go, no CGO) to `go.mod`
- Add `github.com/minio/minio-go/v7` for object storage
- Add `github.com/valyala/fasthttp` (replace Gin's `net/http`)
- Update `config.go`: remove `DBConfig`, add `SQLiteConfig` (path to `predictions.db`), `MinIOConfig` (endpoint, bucket, key, secret), `APIKeyConfig`, `ThumbnailConfig`

### 1.3 — Migrate HTTP from Gin → fasthttp
- Rewrite `transport/http/server.go` using `fasthttp.Server`
- Rewrite handler wiring to use `*fasthttp.RequestCtx` throughout
- Keep `ServerResponse`/`APIError` types in `internal/core/types/server.go`
- Add correlation ID middleware (UUID per request, returned as `X-Request-ID`)
- Add API key auth middleware (static key from config, `X-API-Key` header)
- Add structured request logging middleware (logrus, includes correlation id, method, path, status, latency)

### 1.4 — Core types/errors/constants
- Add `internal/core/constants/classes.go` — pest class enum values
- Add error sentinels in `internal/core/errors/errors.go`: `ErrPhotoNotFound`, `ErrBadCursor`, `ErrInvalidParam`
- Add `internal/core/types/cursor.go` — opaque cursor (base64-encoded JSON `{CapturedAt, ID}`)

---

## Phase 2: Data Layer (SQLite Repository)

### 2.1 — Models (`internal/core/models/`)
- `photo.go`: `Photo{ID, X, Y, H, Width, Height, CapturedAt, OriginalURL, Predictions}`
- `prediction.go`: `Prediction{ID, PhotoID, ClassID, Confidence, BBox}`
- `bounding_box.go`: `BoundingBox{XMin, YMin, XMax, YMax}`

### 2.2 — Repository interface (`internal/repository/repository.go`)
- `PhotoRepository` interface:
  - `GetPhotoByID(ctx, id) (*models.Photo, error)`
  - `ListPhotos(ctx, ListPhotosParams) ([]models.Photo, string, error)` — returns items + next cursor
- Remove GORM `BaseRepo`, `WrapInTransaction` (not needed for SQLite read-only)

### 2.3 — SQLite implementation (`internal/repository/sqlite/`)
- `sqlite.go`: open SQLite DB with read-only mode
- `photo_repository.go`: implement `PhotoRepository` — direct `database/sql` queries, map rows to models
- Cursor pagination: `WHERE (captured_at, id) < (?, ?)` ordered by `(captured_at DESC, id)`
- Filter logic: `classId` + `minConfidence` → same prediction must match both (JOIN on predictions)

---

## Phase 3: Object Storage (MinIO)

### 3.1 — MinIO service (`internal/services/storage/`)
- `GenerateUploadLink(ctx, photoID, contentType) (*UploadLink, error)` — presigned PUT URL (15 min TTL)
- `GetOriginalURL(ctx, photoID) (string, error)` — presigned GET URL (1 hour TTL, or public URL)
- `ObjectExists(ctx, photoID) (bool, error)`

---

## Phase 4: Thumbnail Engine

### 4.1 — Design (`internal/services/thumbnail/`)
- Endpoint: `GET /thumbnails/{photoID}?w=<width>&dpr=<dpr>&q=<quality>`
  - `w`: pixel width (1–2560)
  - `dpr`: device pixel ratio (1.0–3.0, default 1)
  - `q`: JPEG quality (1–100, default 80)
- Server fetches original from MinIO, resizes with `github.com/disintegration/imaging`
- Cache key: `SHA256(photoID + w + dpr + q)` → stored back in MinIO under `thumbs/` prefix (write-once)
- On request: check MinIO `thumbs/` first → cache hit streams directly; miss → generate → store → stream
- Cache hit/miss and generation time tracked via Prometheus metrics

### 4.2 — Concurrency guard
- Per-key singleflight (`golang.org/x/sync/singleflight`) to avoid generating the same thumbnail twice under concurrent requests

---

## Phase 5: Services & Handlers

### 5.1 — Photo service (`internal/services/photo/`)
- `GetPhoto(ctx, GetPhotoDTO) (*models.Photo, error)` — fetches from repo, enriches with `OriginalURL` from MinIO
- `ListPhotos(ctx, ListPhotosDTO) ([]models.Photo, string, error)`
- `CreateUploadLink(ctx, CreateUploadLinkDTO) (*UploadLink, error)`

### 5.2 — HTTP handlers (`internal/transport/http/handler/`)
- `GET /photos` → parse cursor, limit, classId, minConfidence → call service → `200 PhotoPage`
- `GET /photos/{photoId}` → call service → `200 Photo` or `404`
- `POST /photos/{photoId}/upload-link` → parse body → call service → `200 UploadLink`
- `GET /thumbnails/{photoId}` → parse params → call thumbnail service → stream JPEG
- `GET /metrics` → Prometheus handler
- Error mapping in one place: sentinel errors → HTTP status + typed error body with `request_id`

---

## Phase 6: Seed Client

### 6.1 — Seed binary (`cmd/seed/main.go`)
- Reads all photos from `predictions.db`
- For each photo, calls `POST /photos/{id}/upload-link`, then PUTs the file from `dataset/images/`
- Idempotent: skips photos already in MinIO (checks object existence before uploading)
- Configurable via env: `API_URL`, `API_KEY`, `IMAGES_DIR`, `DB_PATH`
- Prints progress to stdout

---

## Phase 7: Frontend

### 7.1 — Project scaffold (`frontend/`)
- Vite + React 19 + TypeScript, pnpm
- Redux Toolkit (global state: selected classId filter, confidence filter, map-click point)
- `openapi-typescript` to generate types from `openapi.yaml`
- CSS Modules, vitest + React Testing Library
- Feature-based folder structure:
  ```
  src/
    features/
      gallery/       # photo grid + photo modal
      map/           # greenhouse Konva canvas
      filters/       # shared filter bar
    shared/
      api/           # RTK Query endpoints
      components/    # reusable UI
      hooks/
  ```

### 7.2 — Gallery feature
- Paginated grid, `IntersectionObserver` for infinite scroll
- Thumbnails via `<img srcset>` calling `/thumbnails/{id}?w=X&dpr=Y`
- BBox overlay: `<canvas>` drawn over each thumbnail; normalized coords × rendered size
- Click opens modal with full-size photo (`w=2560`) + all prediction boxes drawn
- Filter bar: class dropdown + confidence slider → updates Redux → re-fetches

### 7.3 — Greenhouse map feature
- `react-konva` canvas mapping the 40×40 m greenhouse floor plan
- Each photo = marker at `(x, y)` colored by dominant prediction class
- Zoom + pan via Konva stage
- Click on map marker → dispatches `setMapCenter({x, y, radius: 5})` to Redux → gallery filters photos within that radius
- Shared class filter drives both views

### 7.4 — Error/loading/empty states
- Loading skeleton for gallery grid
- Error boundary with retry button
- Empty state component for filtered results with zero matches

---

## Phase 8: Tests

### 8.1 — Backend tests
- `internal/core/helpers/bbox_test.go` — bbox coordinate transform (normalized → pixel coords at any size/DPR)
- `internal/services/thumbnail/thumbnail_test.go` — parse/validate thumbnail params + cache key math
- `internal/repository/sqlite/photo_repository_test.go` — ingest `predictions.db` then read (smoke)
- Handler tests for photo endpoints using testify/mock mocks for services

### 8.2 — Frontend tests
- `features/gallery/BBoxOverlay.test.tsx` — bbox transform correctness at various sizes and DPRs
- `features/filters/filtersSlice.test.ts` — reducer unit test
- `features/map/mapSlice.test.ts` — "near" radius filter logic

---

## Phase 9: Dockerization

### 9.1 — Backend Dockerfile (`backend/build/app/Dockerfile`)
- Multi-stage: `golang:1.25.5` build → `alpine` runtime
- CGO disabled (using `modernc.org/sqlite`)
- Copies `configs/`, `main.env`, mounts `predictions.db` at runtime

### 9.2 — Frontend Dockerfile (`frontend/Dockerfile`)
- `node:22-alpine` build (pnpm install + vite build) → `nginx:alpine` serve
- nginx proxies `/api` and `/thumbnails` to backend

### 9.3 — docker-compose.yml (root level)
- Services: `backend`, `frontend`, `minio`, `minio-init` (creates bucket on startup via `mc`)
- `backend` mounts `dataset/` read-only
- All wired with healthchecks + `depends_on`

### 9.4 — Root README.md
- Prerequisites (Docker + Docker Compose)
- One-command startup: `docker compose up --build`
- Seed step: `docker compose run --rm seed`
- Ports: backend `8080`, frontend `3000`, MinIO console `9001`
- Environment variable reference

---

## Phase 10: CLAUDE.md

Create `backend/CLAUDE.md` with all coding conventions from `prompt.txt` plus Scout-specific notes:
- SQLite path config
- MinIO config
- Thumbnail endpoint design rationale
- Seed binary location
- No PostgreSQL/GORM — do not add them back

---

## Implementation Order

| # | Task |
|---|------|
| 1 | CLAUDE.md |
| 2 | Module rename + remove PostgreSQL boilerplate |
| 3 | Config, SQLite repo, models |
| 4 | MinIO storage service |
| 5 | Photo service + handlers (data API) |
| 6 | Thumbnail engine |
| 7 | Metrics endpoint |
| 8 | Seed binary |
| 9 | Backend tests |
| 10 | Frontend scaffold + RTK Query + generated types |
| 11 | Gallery feature + BBox overlay |
| 12 | Greenhouse map feature |
| 13 | Frontend tests |
| 14 | Docker-compose + Dockerfiles |
| 15 | Root README |
