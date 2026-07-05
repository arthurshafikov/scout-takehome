/**
 * Scout API Types
 * Mirrors backend models from internal/core/models/model.go
 */

export interface BoundingBox {
  xMin: number
  yMin: number
  xMax: number
  yMax: number
}

export interface Prediction {
  id: string
  photoId: string
  classId: string
  confidence: number
  bbox: BoundingBox
}

export interface Photo {
  id: string
  x: number
  y: number
  h: number
  width: number
  height: number
  capturedAt: string
  originalUrl: string
  predictions: Prediction[]
}

export interface PhotoPage {
  items: Photo[]
  nextCursor?: string
}

export interface UploadLink {
  uploadUrl: string
  downloadUrl: string
}

export interface HealthCheck {
  test: boolean
}

/** Pest classes (from internal/core/constants/classes.go) */
export const PEST_CLASSES = [
  'powdery_mildew',
  'mirid',
  'whitefly_aphid',
  'miner_tuta',
  'thrips',
  'spider_mites',
] as const

export type PestClass = (typeof PEST_CLASSES)[number]

/** Query parameters (snake_case per backend requirement) */
export interface ListPhotosParams {
  cursor?: string
  limit?: number
  class_id?: PestClass
  min_confidence?: number
}
