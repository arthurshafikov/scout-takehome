import { describe, it, expect } from 'vitest'
import {
  normalizedToPixels,
  pixelsToNormalized,
  bboxWidth,
  bboxHeight,
  bboxCenter,
} from '@/utils/bbox'
import type { BoundingBox } from '@/types/api'

describe('bbox utilities', () => {
  const testBbox: BoundingBox = {
    xMin: 0.1,
    yMin: 0.2,
    xMax: 0.7,
    yMax: 0.9,
  }

  describe('normalizedToPixels', () => {
    it('converts normalized bbox to pixels', () => {
      const result = normalizedToPixels(testBbox, 1000, 500)
      expect(result.xMin).toBe(100)
      expect(result.yMin).toBe(100)
      expect(result.xMax).toBe(700)
      expect(result.yMax).toBe(450)
    })

    it('handles image dimensions correctly', () => {
      const result = normalizedToPixels(testBbox, 2560, 1440)
      expect(result.xMin).toBeCloseTo(256)
      expect(result.yMin).toBeCloseTo(288)
      expect(result.xMax).toBeCloseTo(1792)
      expect(result.yMax).toBeCloseTo(1296)
    })
  })

  describe('pixelsToNormalized', () => {
    it('converts pixels back to normalized bbox', () => {
      const pixels = normalizedToPixels(testBbox, 1000, 500)
      const result = pixelsToNormalized(pixels, 1000, 500)
      expect(result.xMin).toBeCloseTo(testBbox.xMin)
      expect(result.yMin).toBeCloseTo(testBbox.yMin)
      expect(result.xMax).toBeCloseTo(testBbox.xMax)
      expect(result.yMax).toBeCloseTo(testBbox.yMax)
    })
  })

  describe('bboxWidth', () => {
    it('calculates width correctly', () => {
      const pixels = normalizedToPixels(testBbox, 1000, 500)
      const width = bboxWidth(pixels)
      expect(width).toBe(600) // (0.7 - 0.1) * 1000
    })
  })

  describe('bboxHeight', () => {
    it('calculates height correctly', () => {
      const pixels = normalizedToPixels(testBbox, 1000, 500)
      const height = bboxHeight(pixels)
      expect(height).toBe(350) // (0.9 - 0.2) * 500
    })
  })

  describe('bboxCenter', () => {
    it('calculates center point', () => {
      const pixels = normalizedToPixels(testBbox, 1000, 500)
      const center = bboxCenter(pixels)
      expect(center.x).toBe(400) // (100 + 700) / 2
      expect(center.y).toBeCloseTo(275) // (100 + 450) / 2
    })
  })
})
