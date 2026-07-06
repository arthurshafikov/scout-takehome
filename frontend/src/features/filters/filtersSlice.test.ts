import { describe, it, expect } from 'vitest'
import filtersReducer, {
  setSelectedClass,
  setMinConfidence,
  setCursor,
  resetFilters,
} from '@/features/filters/filtersSlice'
import type { FiltersState } from '@/features/filters/filtersSlice'

describe('filtersSlice', () => {
  const initialState: FiltersState = {
    selectedClass: null,
    minConfidence: 0,
    cursor: null,
    locationCenter: null,
  }

  it('has correct initial state', () => {
    const state = filtersReducer(undefined, { type: 'unknown' })
    expect(state).toEqual(initialState)
  })

  describe('setSelectedClass', () => {
    it('sets selected class and resets cursor', () => {
      const state = filtersReducer(
        { ...initialState, cursor: 'old-cursor' },
        setSelectedClass('powdery_mildew'),
      )
      expect(state.selectedClass).toBe('powdery_mildew')
      expect(state.cursor).toBeNull()
    })

    it('clears selected class when set to null', () => {
      const state = filtersReducer(
        { ...initialState, selectedClass: 'thrips' },
        setSelectedClass(null),
      )
      expect(state.selectedClass).toBeNull()
    })
  })

  describe('setMinConfidence', () => {
    it('sets confidence threshold and resets cursor', () => {
      const state = filtersReducer(
        { ...initialState, cursor: 'old-cursor' },
        setMinConfidence(0.7),
      )
      expect(state.minConfidence).toBe(0.7)
      expect(state.cursor).toBeNull()
    })

    it('clamps confidence to valid range', () => {
      const state1 = filtersReducer(initialState, setMinConfidence(0.5))
      expect(state1.minConfidence).toBe(0.5)

      const state2 = filtersReducer(state1, setMinConfidence(0))
      expect(state2.minConfidence).toBe(0)
    })
  })

  describe('setCursor', () => {
    it('sets cursor for pagination', () => {
      const state = filtersReducer(
        initialState,
        setCursor('next-page-cursor'),
      )
      expect(state.cursor).toBe('next-page-cursor')
    })

    it('can clear cursor', () => {
      const state = filtersReducer(
        { ...initialState, cursor: 'some-cursor' },
        setCursor(null),
      )
      expect(state.cursor).toBeNull()
    })
  })

  describe('resetFilters', () => {
    it('resets all filters to initial state', () => {
      const dirtyState: FiltersState = {
        selectedClass: 'thrips',
        minConfidence: 0.8,
        cursor: 'page-2',
      }
      const state = filtersReducer(dirtyState, resetFilters())
      expect(state).toEqual(initialState)
    })
  })
})
