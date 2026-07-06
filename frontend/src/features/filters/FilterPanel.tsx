import { FC } from 'react'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import {
  setSelectedClass,
  setMinConfidence,
  setLocationCenter,
  resetFilters,
} from '@/features/filters/filtersSlice'
import { PEST_CLASSES } from '@/types/api'
import { PEST_CLASS_LABELS } from '@/utils/classColors'

/**
 * Filter panel with pest class selector and confidence slider
 * Changes reset pagination to show results from the beginning
 */
const FilterPanel: FC = () => {
  const dispatch = useAppDispatch()
  const { selectedClass, minConfidence, locationCenter } = useAppSelector(
    (state) => state.filters,
  )

  const handleClassChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const value = e.target.value as any
    dispatch(setSelectedClass(value === 'none' ? null : value))
  }

  const handleConfidenceChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    dispatch(setMinConfidence(Number(e.target.value)))
  }

  const handleReset = () => {
    dispatch(resetFilters())
  }

  return (
    <div className="bg-white p-6 rounded-lg shadow mb-6">
      <h2 className="text-xl font-bold text-gray-900 mb-4">Filters</h2>

      <div className="space-y-4">
        {/* Pest class selector */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Pest Class
          </label>
          <select
            value={selectedClass || 'none'}
            onChange={handleClassChange}
            className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="none">All Classes</option>
            {PEST_CLASSES.map((pest) => (
              <option key={pest} value={pest}>
                {PEST_CLASS_LABELS[pest]}
              </option>
            ))}
          </select>
        </div>

        {/* Confidence threshold slider */}
        <div>
          <div className="flex justify-between items-center mb-2">
            <label className="block text-sm font-medium text-gray-700">
              Minimum Confidence
            </label>
            <span className="text-sm font-semibold text-blue-600">
              {(minConfidence * 100).toFixed(0)}%
            </span>
          </div>
          <input
            type="range"
            min="0"
            max="1"
            step="0.05"
            value={minConfidence}
            onChange={handleConfidenceChange}
            className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer"
          />
          <div className="flex justify-between text-xs text-gray-500 mt-1">
            <span>0%</span>
            <span>50%</span>
            <span>100%</span>
          </div>
        </div>

        {/* Reset button */}
        <button
          onClick={handleReset}
          disabled={!selectedClass && minConfidence === 0 && !locationCenter}
          className="w-full px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 disabled:opacity-50 disabled:cursor-not-allowed rounded-md transition-colors"
        >
          Reset Filters
        </button>
      </div>

      {/* Active filters display */}
      {(selectedClass || minConfidence > 0 || locationCenter) && (
        <div className="mt-4 pt-4 border-t border-gray-200">
          <div className="text-sm text-gray-600">
            <strong>Active filters:</strong>
            <div className="flex flex-wrap gap-2 mt-2">
              {selectedClass && (
                <span className="inline-block bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-xs">
                  {PEST_CLASS_LABELS[selectedClass]}
                </span>
              )}
              {minConfidence > 0 && (
                <span className="inline-block bg-green-100 text-green-800 px-3 py-1 rounded-full text-xs">
                  Confidence ≥ {(minConfidence * 100).toFixed(0)}%
                </span>
              )}
              {locationCenter && (
                <div className="flex items-center gap-2">
                  <span className="inline-block bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-xs">
                    Location ({locationCenter.x.toFixed(1)}m, {locationCenter.y.toFixed(1)}m)
                  </span>
                  <button
                    onClick={() => dispatch(setLocationCenter(null))}
                    className="inline-block text-blue-600 hover:text-blue-800 text-xs font-semibold"
                    title="Clear location filter"
                  >
                    ✕
                  </button>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export default FilterPanel
