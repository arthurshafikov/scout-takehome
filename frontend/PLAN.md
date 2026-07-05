# Plan: Scout Frontend

## TL;DR

React 19 + TypeScript + Vite frontend in `./frontend/`. Tailwind for styling, Redux Toolkit + RTK Query for state and data fetching, SVG overlays for bounding boxes, react-konva for the 40×40m greenhouse map. Minimal tests: bbox coordinate transform, thumbnail param validation, and one filter reducer test.

---

## Steps

### Phase A — Scaffold (independent)
1. Init Vite project: `pnpm create vite frontend --template react-ts`
2. Install dependencies: `tailwindcss`, `@reduxjs/toolkit`, `react-redux`, `react-konva`, `konva`, `vitest`, `@testing-library/react`, `@testing-library/jest-dom`
3. Configure Tailwind (`tailwind.config.js`, `postcss.config.js`, import in `index.css`)
4. Configure Vite proxy: `/api/*` → `http://localhost:8080` (removes need for CORS header in dev)
5. Create `CLAUDE.md` in `frontend/` with conventions

### Phase B — Types & API Client (depends on A)
6. Hand-write TypeScript types from OpenAPI — `Photo`, `Prediction`, `BoundingBox`, `PhotoPage`, `UploadLink` → `src/types/api.ts`
7. RTK Query API slice (`src/services/api.ts`): `listPhotos`, `getPhoto`, inject `X-API-Key` from `VITE_API_KEY`. Separate `getThumbnailUrl()` helper (direct URL string, no auth needed).

### Phase C — Redux Store (depends on A)
8. Store (`src/app/store.ts`) + filter slice (`src/features/filters/filtersSlice.ts`): `{ classId, minConfidence, mapCenter }`, actions `setClassId`, `setMinConfidence`, `setMapCenter`, `clearMapCenter`

### Phase D — Gallery Feature (depends on B, C)
9. `GalleryPage`: cursor pagination ("Load more")
10. `PhotoCard`: `<img>` with `srcset` at 1×/2× DPR + SVG bbox overlay, opens `PhotoModal`
11. `BBoxOverlay` (SVG): `px = xMin * renderedWidth`, color per class (6-color constant map)
12. `PhotoModal`: full-size via `originalUrl` + all predictions

### Phase E — Filter Panel (depends on C) — *parallel with D*
13. Class dropdown + confidence slider → dispatch to filter slice

### Phase F — Greenhouse Map (depends on B, C) — *parallel with D*
14. `MapPanel` (react-konva): Circle per photo at `(x,y)` on 40×40m Stage, color by dominant class, fade non-matching, hover tooltip, click → `setMapCenter`
15. `MapLegend`; map data from separate `listPhotos({ limit: 200 })` call

### Phase G — Layout (depends on D, E, F)
16. Desktop: sidebar(filters) + center(gallery) + right(map); mobile: tabs

### Phase H — Tests (depends on D, E)
17. `src/utils/bbox.test.ts` — `transformBBox(bbox, w, h)` pixel transforms + edge cases
18. `src/utils/thumbnail.test.ts` — URL construction, width/quality bounds
19. `src/features/filters/filtersSlice.test.ts` — reducer state transitions

### Phase I — Docker (depends on all)
20. `frontend/Dockerfile`: Node build → nginx with SPA fallback + `/api/*` proxy to `backend:8080`
21. Add `frontend` service to `backend/deployments/docker-compose.yml`

### Phase J — Env & Docs
22. `frontend/.env.example`; update root `README.md`

---

## Relevant Files

- `openapi.yaml` — type source of truth
- `backend/internal/core/models/model.go` — Go models to mirror in TS
- `backend/internal/transport/http/handler/response.go` — response shape `{ success, data, error, trace_id }`
- `backend/internal/transport/http/handler/photo_handler.go` — query params: `cursor`, `limit`, `class_id`, `min_confidence`
- `backend/internal/core/constants/classes.go` — 6 pest classes
- `backend/deployments/docker-compose.yml` — add `frontend` service

---

## Verification

1. `cd frontend && pnpm run dev` — gallery loads, bboxes render correctly
2. `pnpm test` — all 3 test files pass
3. Class filter → gallery refetches, map fades non-matching dots
4. Map click → `mapCenter` in Redux → gallery shows nearby photos
5. DevTools: `?w=400` at 1× DPR, `?w=800` at 2× DPR in srcset
6. `docker compose up` → frontend `:80`, backend `:8080`, MinIO `:9000`

---

## Decisions

- **State**: RTK + RTK Query — shared filter state, automatic cache invalidation, cursor pagination
- **BBox rendering**: SVG overlay — perfect scaling, no pixel-rounding, accessible
- **Thumbnail**: Direct URL string, not RTK Query — browser cache handles blobs natively
- **Map data**: Separate `listPhotos({ limit: 200 })` — full map always visible regardless of gallery scroll
- **Tests**: vitest — zero extra config, same toolchain as Vite
- **Excluded**: Upload UI, auth UI (handled by seed binary / static key in env)

---

## Open Questions

1. **Query param casing**: Backend uses `class_id`/`min_confidence` (snake_case); OpenAPI uses camelCase. Must match *actual* backend during Phase B.
2. **Production CORS**: nginx proxies `/api/*` → `backend:8080` — no cross-origin requests.
3. **Map density**: 50-photo gallery page would give sparse map — map independently fetches up to 200 photos.
