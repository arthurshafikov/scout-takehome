/**
 * Thumbnail URL utilities
 * Constructs URLs for thumbnail images from the backend
 */

/**
 * Get thumbnail URL for a photo
 * Format: GET /api/thumbnails/{photoId}
 * Returns binary image data (JPEG or PNG)
 * Uses relative path to work through nginx proxy in Docker
 */
export function getThumbnailUrl(photoId: string): string {
  return `/api/thumbnails/${photoId}`
}
