/**
 * Image loading utility with API key authentication
 * Fetches images with X-API-Key header since img tags cannot send custom headers
 */

const imageCache = new Map<string, string>()
const loadingPromises = new Map<string, Promise<string>>()

/**
 * Load image with API key header and return blob URL
 * Uses caching to avoid duplicate requests for same image
 */
export async function loadImageWithAuth(
  url: string,
  apiKey: string,
): Promise<string> {
  // Return cached blob URL if available
  if (imageCache.has(url)) {
    return imageCache.get(url)!
  }

  // Return existing promise if request is in progress
  if (loadingPromises.has(url)) {
    return loadingPromises.get(url)!
  }

  // Fetch image with API key header
  const promise = (async () => {
    try {
      const response = await fetch(url, {
        headers: {
          'X-API-Key': apiKey,
        },
      })

      if (!response.ok) {
        throw new Error(`Failed to load image: ${response.status}`)
      }

      const blob = await response.blob()
      const blobUrl = URL.createObjectURL(blob)

      // Cache the blob URL
      imageCache.set(url, blobUrl)
      loadingPromises.delete(url)

      return blobUrl
    } catch (error) {
      loadingPromises.delete(url)
      throw error
    }
  })()

  loadingPromises.set(url, promise)
  return promise
}

/**
 * Clear image cache to free memory
 */
export function clearImageCache() {
  imageCache.forEach((blobUrl) => {
    URL.revokeObjectURL(blobUrl)
  })
  imageCache.clear()
}
