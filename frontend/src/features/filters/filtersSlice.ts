import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import type { PestClass } from '@/types/api'

export interface FiltersState {
  selectedClass: PestClass | null
  minConfidence: number
  cursor: string | null
}

const initialState: FiltersState = {
  selectedClass: null,
  minConfidence: 0,
  cursor: null,
}

const filtersSlice = createSlice({
  name: 'filters',
  initialState,
  reducers: {
    setSelectedClass: (state, action: PayloadAction<PestClass | null>) => {
      state.selectedClass = action.payload
      // Reset cursor when filter changes to start from beginning
      state.cursor = null
    },
    setMinConfidence: (state, action: PayloadAction<number>) => {
      state.minConfidence = action.payload
      // Reset cursor when filter changes
      state.cursor = null
    },
    setCursor: (state, action: PayloadAction<string | null>) => {
      state.cursor = action.payload
    },
    resetFilters: (state) => {
      state.selectedClass = null
      state.minConfidence = 0
      state.cursor = null
    },
  },
})

export const { setSelectedClass, setMinConfidence, setCursor, resetFilters } =
  filtersSlice.actions
export default filtersSlice.reducer
