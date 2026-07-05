import type { BoundingBox, PestClass } from '@/types/api'

/**
 * Transform normalized BBox [0,1] to pixel coordinates
 * Normalized coords: 0-1 relative to image size
 * Pixel coords: absolute pixels in image
 */
export function normalizedToPixels(
  bbox: BoundingBox,
  width: number,
  height: number,
): BoundingBox {
  return {
    xMin: bbox.xMin * width,
    yMin: bbox.yMin * height,
    xMax: bbox.xMax * width,
    yMax: bbox.yMax * height,
  }
}

/**
 * Transform pixel coordinates to normalized BBox [0,1]
 */
export function pixelsToNormalized(
  bbox: BoundingBox,
  width: number,
  height: number,
): BoundingBox {
  return {
    xMin: bbox.xMin / width,
    yMin: bbox.yMin / height,
    xMax: bbox.xMax / width,
    yMax: bbox.yMax / height,
  }
}

/**
 * Get bounding box width in pixels
 */
export function bboxWidth(bbox: BoundingBox): number {
  return bbox.xMax - bbox.xMin
}

/**
 * Get bounding box height in pixels
 */
export function bboxHeight(bbox: BoundingBox): number {
  return bbox.yMax - bbox.yMin
}

/**
 * Get center point of bounding box
 */
export function bboxCenter(bbox: BoundingBox): { x: number; y: number } {
  return {
    x: (bbox.xMin + bbox.xMax) / 2,
    y: (bbox.yMin + bbox.yMax) / 2,
  }
}
