/**
 * Greenhouse Layout Constants
 * All measurements in meters or pixels as specified
 */

// Greenhouse dimensions (meters)
export const GREENHOUSE_WIDTH = 40
export const GREENHOUSE_HEIGHT = 40

// Map visualization (pixels per meter at 1x zoom)
export const PIXEL_SCALE = 40

// Proximity filter default radius (meters)
export const PROXIMITY_RADIUS = 5

// Zoom constraints
export const MIN_ZOOM = 0.5
export const MAX_ZOOM = 5
export const ZOOM_IN_FACTOR = 1.05
export const ZOOM_OUT_FACTOR = 0.95

// API pagination limits
export const MAP_FETCH_LIMIT = 200
export const GALLERY_FETCH_LIMIT = 50
export const GALLERY_REFETCH_INTERVAL_MS = 60 * 1000 // 60 seconds
export const MAP_REFETCH_INTERVAL_MS = 300 * 1000 // 300 seconds

// Intersection observer threshold for infinite scroll
export const INFINITE_SCROLL_THRESHOLD = 0.1

// Photo card dimensions (for visual feedback)
export const PHOTO_CIRCLE_RADIUS_PX = 8
