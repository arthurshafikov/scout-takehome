import type { BoundingBox } from '@/types/api'

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
 * Get center point of bounding box
 */
export function bboxCenter(bbox: BoundingBox): { x: number; y: number } {
  return {
    x: (bbox.xMin + bbox.xMax) / 2,
    y: (bbox.yMin + bbox.yMax) / 2,
  }
}
