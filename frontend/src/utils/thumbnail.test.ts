import { describe, it, expect, beforeEach } from 'vitest'
import {
  getThumbnailUrl,
  getPhotoUrl,
  isValidPhotoUrl,
} from '@/utils/thumbnail'

describe('thumbnail utilities', () => {
  beforeEach(() => {
    // Mock import.meta.env
    ;(import.meta as any).env.VITE_API_URL = 'http://localhost:8080'
  })

  describe('getThumbnailUrl', () => {
    it('constructs thumbnail URL correctly', () => {
      const photoId = 'test-photo-123'
      const url = getThumbnailUrl(photoId)
      expect(url).toBe('http://localhost:8080/api/thumbnails/test-photo-123')
    })

    it('handles special characters in photoId', () => {
      const photoId = 'photo-with-dashes-and-numbers-123'
      const url = getThumbnailUrl(photoId)
      expect(url).toContain(photoId)
    })
  })

  describe('getPhotoUrl', () => {
    it('constructs photo URL correctly', () => {
      const photoId = 'test-photo-123'
      const url = getPhotoUrl(photoId)
      expect(url).toBe('http://localhost:8080/api/photos/test-photo-123/image')
    })
  })

  describe('isValidPhotoUrl', () => {
    it('validates correct URLs', () => {
      expect(isValidPhotoUrl('http://localhost:8080/api/thumbnails/123')).toBe(true)
      expect(isValidPhotoUrl('https://example.com/image.jpg')).toBe(true)
    })

    it('rejects invalid URLs', () => {
      expect(isValidPhotoUrl('not-a-url')).toBe(false)
      expect(isValidPhotoUrl('')).toBe(false)
      expect(isValidPhotoUrl('ht!tp://invalid')).toBe(false)
    })
  })
})
