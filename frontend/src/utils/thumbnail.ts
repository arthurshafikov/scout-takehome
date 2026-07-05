/**
 * Thumbnail URL utilities
 * Constructs URLs for thumbnail images from the backend
 */

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

/**
 * Get thumbnail URL for a photo
 * Format: GET /api/thumbnails/{photoId}
 * Returns binary image data (JPEG or PNG)
 */
export function getThumbnailUrl(photoId: string): string {
  return `${API_URL}/api/thumbnails/${photoId}`
}

/**
 * Get original photo URL
 * Format: GET /api/photos/{photoId}/image
 * Can be served directly or pre-signed from storage
 */
export function getPhotoUrl(photoId: string): string {
  return `${API_URL}/api/photos/${photoId}/image`
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
