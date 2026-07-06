import { describe, it, expect } from 'vitest'
import {
  getThumbnailUrl,
} from '@/utils/thumbnail'

describe('thumbnail utilities', () => {
  describe('getThumbnailUrl', () => {
    it('constructs thumbnail URL correctly with defaults', () => {
      const photoId = 'test-photo-123'
      const url = getThumbnailUrl(photoId)
      expect(url).toBe('/api/thumbnails/test-photo-123?w=600&q=85')
    })

    it('constructs thumbnail URL with custom width and quality', () => {
      const photoId = 'test-photo-456'
      const url = getThumbnailUrl(photoId, 900, 90)
      expect(url).toBe('/api/thumbnails/test-photo-456?w=900&q=90')
    })

    it('handles special characters in photoId', () => {
      const photoId = 'photo-with-dashes-and-uuid-format'
      const url = getThumbnailUrl(photoId, 300)
      expect(url).toContain(photoId)
      expect(url).toContain('w=300')
    })
  })
})
