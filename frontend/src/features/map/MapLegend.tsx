import { FC } from 'react'
import { PEST_CLASSES, PEST_CLASS_LABELS } from '@/types/api'
import { PEST_CLASS_COLORS } from '@/utils/classColors'

/**
 * Legend showing pest class colors and labels
 * Displayed below the map for reference
 */
const MapLegend: FC = () => {
  return (
    <div className="bg-white p-4 rounded-lg shadow">
      <h3 className="text-sm font-bold text-gray-900 mb-3">Pest Classes</h3>
      <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-3">
        {PEST_CLASSES.map((pestClass) => (
          <div key={pestClass} className="flex items-center gap-2">
            <div
              className={`w-4 h-4 rounded border border-gray-300 ${PEST_CLASS_COLORS[pestClass]}`}
            />
            <span className="text-xs text-gray-700">
              {PEST_CLASS_LABELS[pestClass]}
            </span>
          </div>
        ))}
      </div>
    </div>
  )
}

export default MapLegend
