/**
 * Thumbnail URL utilities
 * Constructs URLs for thumbnail images from the backend with responsive sizing
 * 
 * Backend supports:
 * - w: width in pixels (e.g., 300, 600, 900, 1200)
 * - q: quality 0-100 (default 85)
 */

/**
 * Get thumbnail URL for a photo at specific size and quality
 * Format: GET /api/thumbnails/{photoId}?w={width}&q={quality}
 * Returns binary image data (JPEG or PNG)
 * 
 * @param photoId - Photo UUID
 * @param width - Desired width in pixels (default: 600)
 * @param quality - JPEG quality 0-100 (default: 85)
 * @returns URL string for the thumbnail
 */
export function getThumbnailUrl(
  photoId: string,
  width: number = 600,
  quality: number = 85
): string {
  return `/api/thumbnails/${photoId}?w=${width}&q=${quality}`
}

/**
 * Get srcset string for responsive images
 * Provides multiple resolution options for different screen sizes
 * 
 * @param photoId - Photo UUID
 * @param quality - JPEG quality (default: 85)
 * @returns srcset string with 1x, 2x options for standard width
 */
export function getThumbnailSrcSet(
  photoId: string,
  quality: number = 85
): string {
  return `
    ${getThumbnailUrl(photoId, 300, quality)} 300w,
    ${getThumbnailUrl(photoId, 600, quality)} 600w,
    ${getThumbnailUrl(photoId, 900, quality)} 900w,
    ${getThumbnailUrl(photoId, 1200, quality)} 1200w
  `.trim()
}

/**
 * Standard sizes attribute for responsive images in gallery grid
 * Adapts image size based on viewport and layout
 */
export const GALLERY_SIZES =
  '(max-width: 640px) 100vw, (max-width: 1024px) 50vw, (max-width: 1536px) 33vw, 25vw'

/**
 * Sizes for modal/lightbox view (larger images)
 */
export const MODAL_SIZES = '(max-width: 1024px) 100vw, 90vw'

