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

/**
 * Get original photo URL
 * Format: GET /api/photos/{photoId}/image
 * Can be served directly or pre-signed from storage
 * Uses relative path to work through nginx proxy in Docker
 */
export function getPhotoUrl(photoId: string): string {
  return `/api/photos/${photoId}/image`
}

/**
 * Check if URL is valid
 */
export function isValidPhotoUrl(url: string): boolean {
  try {
    new URL(url)
    return true
  } catch {
    return false
  }
}
