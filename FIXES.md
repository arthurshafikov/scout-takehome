# Scout - Outstanding Fixes & Improvements

This document tracks gaps against the Business Requirements Document (BRD) and recommended improvements.

## Status Summary

| Category | Status | Priority |
|----------|--------|----------|
| BRD Compliance | ✅ Complete | Critical |
| Code Quality | 🟡 5 items | Medium |
| Testing | ✅ Partial | Medium |
| Performance | ✅ Optimized | - |

---

## 🔴 Critical: BRD Requirements

### 1. Metrics Endpoint (`/metrics`)

**Status:** ✅ **DONE** | **Effort:** 30 min | **Impact:** HIGH

**What's needed:**
- BRD requires: `rate, latency, errors, plus thumbnail cache hit/miss and generation time`
- Current: Stub endpoint returns placeholder JSON
- Missing metrics:
  - Request rate per endpoint
  - Latency percentiles (p50, p95, p99)
  - Error rates by status code
  - Thumbnail cache hits
  - Thumbnail cache misses
  - Thumbnail generation time

**Implementation approach:**
1. Add Prometheus client library (`github.com/prometheus/client_golang`)
2. Create metrics collectors in `backend/internal/services/metrics/` (new package)
3. Track in handlers and thumbnail service:
   - Counter: `requests_total{endpoint,method,status}`
   - Histogram: `request_duration_seconds{endpoint}`
   - Counter: `thumbnail_cache_hits_total{photo_id,size}`
   - Counter: `thumbnail_cache_misses_total{photo_id,size}`
   - Histogram: `thumbnail_generation_seconds`
4. Update handler: `backend/internal/transport/http/handler/handler.go`
   - Replace TODO comment with real metrics collection
5. Export at `GET /metrics` in Prometheus text format

**Files to modify:**
- `backend/go.mod` (add prometheus client)
- `backend/internal/services/metrics/metrics.go` (new)
- `backend/internal/transport/http/handler/handler.go` (replace getMetrics stub)
- `backend/internal/services/thumbnail/thumbnail.go` (add metric calls)
- `backend/internal/app/app.go` (initialize metrics)

**Reference log already added:**
```go
tg.logger.Debugf("Cache HIT: %s", cacheKey)
tg.logger.Debugf("Cache MISS: %s (%v)", cacheKey, err)
```

**Implementation completed:** ✅ 2026-07-05
- ✅ Added Prometheus client library (github.com/prometheus/client_golang v1.20.0)
- ✅ Created metrics service with promauto collectors in backend/internal/services/metrics/
- ✅ Implemented HTTP middleware to track all requests (endpoint, method, status, duration)
- ✅ Added thumbnail metrics recording (cache hits/misses counters, generation time histogram)
- ✅ Exposed /metrics endpoint with Prometheus text format (GET /metrics)
- ✅ Verified metrics collection: cache miss ~150ms, cache hit ~3ms, HTTP requests tracked
- ✅ All containers healthy, metrics actively being collected in production format

**Metrics verified working:**
```
scout_http_requests_total{endpoint="/thumbnails/...",method="GET",status="200"}: 2
scout_http_request_duration_seconds{endpoint="/thumbnails/..."}: histogram with buckets
scout_thumbnail_cache_hits_total: 1
scout_thumbnail_cache_misses_total: 1
scout_thumbnail_generation_seconds: 0.149592258 seconds (cache miss latency)
```

---

### 2. Backend Smoke Test

**Status:** ✅ **DONE** | **Effort:** 20 min | **Impact:** HIGH

**What's needed:**
- BRD requires: `seed/ingest-then-read backend smoke test. Runs from a clean clone.`
- Current: Seed binary exists but no automated test
- Test must: seed photos → read back via API → verify data integrity

**Test scenario:**
1. Create temporary MinIO bucket
2. Run seed with test photos
3. Call `GET /photos` and verify:
   - All photos returned
   - Each has predictions, position (x,y,h), captured_at
4. Call `GET /photos/{id}` for one photo
5. Call `GET /photos?class_id=thrips&min_confidence=0.5` and verify filtering works
6. Verify `originalUrl` is presigned URL format

**Implementation approach:**
1. Create `backend/cmd/seed/integration_test.go`
2. Use testcontainers-go for MinIO + PostgreSQL/SQLite
3. Run seed → API calls → assertions
4. Run via: `go test ./cmd/seed/... -tags=integration`

**Files to create:**
- `backend/cmd/seed/integration_test.go` (new)

**Files to modify:**
- `backend/go.mod` (add testcontainers)
- `backend/Makefile` (add test target)

**Implementation completed:** ✅ 2026-07-06
- ✅ Created integration_test.go with comprehensive validation
- ✅ Test validates: photo listing, fetching by ID, predictions, filtering by class_id and min_confidence
- ✅ Test validates: presigned URLs generated correctly for MinIO
- ✅ Added `make test-smoke` target to run the test easily
- ✅ Test PASSES with 100% success rate:
  ```
  ✓ Found 10 photos (pagination working)
  ✓ Photo structure valid: id, x, y, h, width, height, capturedAt
  ✓ Prediction valid: classId=thrips, confidence=0.73
  ✓ Single photo fetch successful
  ✓ Original URL is valid presigned URL (http://...)
  ✓ Class filter working: found 2 photos with class=thrips
  ✓ Confidence filter working: found 11 photos with confidence >= 0.63
  ```
- ✅ Validates complete seed → read → filter pipeline

---

### 3. Map Click-to-Filter (Bonus - High Polish)

**Status:** ✅ **DONE** | **Effort:** 45 min | **Impact:** MEDIUM

**What's needed:**
- BRD bonus: "Click a spot to filter the gallery to photos near it"
- Current: Map shows all photos, no click interaction
- Missing: Click detection + proximity filter + gallery sync

**Feature:**
1. Click any point on the map
2. Find photos within 5m radius
3. Update gallery filters to show only nearby photos
4. Show "X photos near this location" message

**Implementation approach:**
1. Add click handler to Konva Stage in `MapPanel.tsx`
2. Calculate clicked coordinates (convert pixels → greenhouse meters)
3. Calculate distance to each photo using Pythagorean theorem
4. Dispatch Redux action to set "location filter"
5. `GalleryPage` subscribes to location filter and applies it
6. Show applied filter in FilterPanel

**Data flow:**
```
Map click 
  → (x, y) in greenhouse meters 
  → Find photos where sqrt((photo.x - x)² + (photo.y - y)²) ≤ 5
  → Update Redux filtersSlice.locationCenter
  → GalleryPage filters display list
```

**Files to modify:**
- `frontend/src/features/map/MapPanel.tsx` (add click handler + distance calc)
- `frontend/src/features/filters/filtersSlice.ts` (add locationCenter state)
- `frontend/src/features/gallery/GalleryPage.tsx` (apply location filter)
- `frontend/src/features/filters/FilterPanel.tsx` (show active location filter)

**Implementation completed:** ✅ 2026-07-06
- ✅ Added LocationCenter interface to filtersSlice with x, y, radius fields
- ✅ Added setLocationCenter action and reducer to Redux
- ✅ Added click handler to MapPanel Stage that:
  - Detects clicks on map background (not photo circles)
  - Converts pixel coordinates to greenhouse meters accounting for zoom level
  - Clamps coordinates to greenhouse bounds (0-40m)
  - Dispatches setLocationCenter with 5m radius
- ✅ Added visual indicators in MapPanel:
  - Blue circle showing selected location
  - Blue dot at click center
  - Dynamic help text: "Click to filter by location"
  - Shows active location radius in title
- ✅ Implemented proximity filtering in GalleryPage:
  - calculateDistance() function using Pythagorean theorem
  - applyLocationFilter() to find photos within radius
  - Seamlessly integrates with class and confidence filters
- ✅ Updated FilterPanel to display location filter:
  - Shows location coordinates (e.g., "Location (20.5m, 15.3m)")
  - Quick clear button (✕) to remove location filter
  - Integrated with "Reset Filters" button
- ✅ Works alongside existing class_id and min_confidence filters
- ✅ Results update instantly when clicking map

**Feature workflow:**
1. User sees map with all 200 photos
2. User clicks any location on the map → blue circle appears
3. Gallery automatically updates to show only nearby photos (within 5m)
4. FilterPanel displays active location filter
5. User can click X to clear location filter
6. Gallery reverts to showing all photos (respecting class/confidence filters)

---

## 🟡 High Priority: Code Quality & Best Practices

### 4. Code Quality Audit

**Status:** 🟡 **IN PROGRESS** | **Effort:** 30 min | **Impact:** MEDIUM

**Checklist:**
- [ ] Replace `type` with `interface` for object types (better for extension)
- [ ] Scan for `any` types and replace with proper types
- [ ] Extract magic numbers to named constants
- [ ] Add JSDoc comments to exported functions

**Magic numbers found:**
```typescript
// Default thumbnail sizes
400, 300, 85  → Move to constants

// Aspect ratios
0.75  → Use enum or const

// Distance for map proximity
5  → Extract as NEARBY_DISTANCE_METERS

// Confidence thresholds
0, 1, 100  → Use CONFIDENCE_MIN, CONFIDENCE_MAX
```

**Files to audit:**
- `frontend/src/features/gallery/GalleryPage.tsx`
- `frontend/src/features/map/MapPanel.tsx`
- `frontend/src/utils/bbox.ts`
- `frontend/src/utils/thumbnail.ts`
- `backend/internal/services/thumbnail/thumbnail.go`

---

### 5. UI State Verification

**Status:** 🟡 **NEEDS AUDIT** | **Effort:** 30 min | **Impact:** MEDIUM

**What to check:**
- [x] Loading state (shows "Loading photos...")
- [ ] Empty state (gallery shows message when no photos match filter)
- [ ] Error state (displays error message instead of blank)
- [ ] Never shows blank screen

**Files to audit:**
- `frontend/src/features/gallery/GalleryPage.tsx` - gallery loading/empty/error
- `frontend/src/features/map/MapPanel.tsx` - map loading state
- `frontend/src/components/AppLayout.tsx` - main layout error boundary

---

## 📝 Implementation Order

### Phase 1: Critical (BRD Compliance) - 50 min total ✅ COMPLETE
```
1. ✅ Metrics endpoint (30 min) - DONE 2026-07-05
   └─ Creates observability for production
   
2. ✅ Backend smoke test (20 min) - DONE 2026-07-06
   └─ Validates data pipeline end-to-end (seed → read → filter)
```

### Phase 2: High Polish (User Experience) - 45 min total ✅ COMPLETE
```
3. ✅ Map click-to-filter (45 min) - DONE 2026-07-06
   └─ Click map → filter gallery by proximity (5m radius)
```

### Phase 3: Code Quality - 60 min total
```
4. Code audit & cleanup (30 min)
   └─ Extract constants, fix types, JSDoc
   
5. UI state verification (30 min)
   └─ Ensure no blank screens in edge cases
```

---

## Testing Checklist

After implementing each fix:

### Metrics Endpoint
```bash
# Test endpoint returns Prometheus format
curl http://localhost:8080/metrics | head -20

# Should show:
# request_duration_seconds_bucket{...}
# request_duration_seconds_sum{...}
# request_duration_seconds_count{...}
# thumbnail_cache_hits_total{...}
# thumbnail_cache_misses_total{...}
```

### Backend Smoke Test
```bash
# Run test
cd backend
go test ./cmd/seed/... -tags=integration -v

# Should output:
# ✓ Photos seeded
# ✓ Photos retrieved
# ✓ Filtering works
# ✓ Predictions correct
```

### Map Click-to-Filter
```
1. Open Map View
2. Click center of map
3. Gallery updates to show only nearby photos
4. Verify distance calculation (use console to debug)
5. Test edge cases:
   - Click corner (0,0)
   - Click far corner (40,40)
   - Click with no photos nearby
```

### Code Quality
```bash
# Frontend type check
npm run build  # Should have no TypeScript errors

# Run tests
npm test       # All tests pass
```

---

## Deployment Checklist

Before shipping to production:

- [ ] All critical fixes (Phase 1) implemented
- [ ] Metrics endpoint validated
- [ ] Smoke test passing
- [ ] Code quality audit complete
- [ ] Tests pass (npm test, go test ./...)
- [ ] Docker build succeeds (`docker compose build`)
- [ ] All containers healthy (`docker compose up`)
- [ ] Frontend loads at http://localhost:5173
- [ ] Gallery + Map + Filters all work
- [ ] Metrics accessible at http://localhost:8080/metrics

---

## Known Limitations

These are acceptable per BRD and can be addressed in Phase 2:

1. **CSS Modules vs Tailwind**
   - BRD suggested CSS Modules; we use Tailwind
   - Rationale: Tailwind is faster for prototyping, easier to maintain
   - Could switch later if needed

2. **openapi-typescript**
   - BRD suggested generating types from openapi.yaml
   - Current: Types defined manually in `frontend/src/types/api.ts`
   - Rationale: Simpler for 50-photo demo, would add for production

3. **pnpm vs npm**
   - BRD recommended pnpm
   - Current: Using npm
   - Rationale: Works identically, npm is more universal

---

## References

- **BRD:** See `BRD.md` for full requirements
- **Status:** See `IMPLEMENTATION_STATUS.md` for phase completion
- **Backend logs:** Check `docker compose logs app` for debug output
- **Metrics format:** https://prometheus.io/docs/instrumenting/exposition_formats/
