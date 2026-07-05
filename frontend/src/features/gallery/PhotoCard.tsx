import { FC } from 'react'
import { getThumbnailUrl } from '@/utils/thumbnail'
import { PEST_CLASS_LABELS } from '@/utils/classColors'
import type { Photo } from '@/types/api'

interface PhotoCardProps {
  photo: Photo
  onClick: () => void
}

/**
 * Grid card displaying photo thumbnail with prediction summary
 */
const PhotoCard: FC<PhotoCardProps> = ({ photo, onClick }) => {
  const highestConfidence =
    photo.predictions.length > 0
      ? Math.max(...photo.predictions.map((p) => p.confidence))
      : 0

  const topPrediction = photo.predictions.find(
    (p) => p.confidence === highestConfidence,
  )

  return (
    <div
      onClick={onClick}
      className="cursor-pointer bg-white rounded-lg shadow hover:shadow-lg transition-shadow overflow-hidden"
    >
      {/* Image container */}
      <div className="relative w-full aspect-square bg-gray-100 overflow-hidden">
        <img
          src={getThumbnailUrl(photo.id)}
          alt={`Photo ${photo.id}`}
          className="w-full h-full object-cover"
        />

        {/* Position indicator */}
        <div className="absolute top-2 right-2 bg-black bg-opacity-50 text-white text-xs px-2 py-1 rounded">
          {photo.x.toFixed(1)}, {photo.y.toFixed(1)}m
        </div>

        {/* Prediction badge */}
        {topPrediction && (
          <div className="absolute bottom-2 left-2 bg-white bg-opacity-90 rounded px-2 py-1 text-xs">
            <div className="font-semibold text-gray-900">
              {PEST_CLASS_LABELS[topPrediction.classId as keyof typeof PEST_CLASS_LABELS]}
            </div>
            <div className="text-gray-600">
              {(topPrediction.confidence * 100).toFixed(0)}%
            </div>
          </div>
        )}
      </div>

      {/* Info section */}
      <div className="p-3">
        <div className="text-xs text-gray-500 mb-2">
          {new Date(photo.capturedAt).toLocaleString()}
        </div>
        <div className="text-xs font-medium text-gray-700">
          {photo.predictions.length} prediction
          {photo.predictions.length !== 1 ? 's' : ''}
        </div>
      </div>
    </div>
  )
}

export default PhotoCard
