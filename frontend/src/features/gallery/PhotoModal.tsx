import { FC, useState } from 'react'
import BBoxOverlay from './BBoxOverlay'
import { getThumbnailUrl, getThumbnailSrcSet, MODAL_SIZES } from '@/utils/thumbnail'
import { PEST_CLASS_LABELS, PEST_CLASS_COLORS } from '@/utils/classColors'
import type { Photo } from '@/types/api'

interface PhotoModalProps {
  photo: Photo
  isOpen: boolean
  onClose: () => void
}

/**
 * Full-screen modal viewer showing photo with bbox overlays
 * Displays predictions sorted by confidence with detailed info
 * Responsive image: uses srcset/sizes for high-quality display at any resolution
 */
const PhotoModal: FC<PhotoModalProps> = ({ photo, isOpen, onClose }) => {
  const [showOverlay, setShowOverlay] = useState(true)

  if (!isOpen) return null

  const sortedPredictions = [...photo.predictions].sort(
    (a, b) => b.confidence - a.confidence,
  )

  return (
    <div className="fixed inset-0 bg-black bg-opacity-75 z-50 flex items-center justify-center p-4">
      <div className="relative w-full max-w-5xl max-h-96vh bg-white rounded-lg overflow-hidden">
        {/* Close button */}
        <button
          onClick={onClose}
          className="absolute top-4 right-4 bg-red-500 hover:bg-red-600 text-white w-8 h-8 rounded-full flex items-center justify-center z-10"
        >
          ✕
        </button>

        <div className="flex flex-col lg:flex-row h-full max-h-screen">
          {/* Image section */}
          <div className="flex-1 flex items-center justify-center bg-black min-h-96">
            <div className="relative w-full h-full flex items-center justify-center">
              <img
                srcSet={getThumbnailSrcSet(photo.id, 90)}
                sizes={MODAL_SIZES}
                src={getThumbnailUrl(photo.id, 900, 90)}
                alt={`Photo ${photo.id}`}
                className="max-w-full max-h-full object-contain"
              />
              {showOverlay && (
                <BBoxOverlay
                  photo={photo}
                  imageWidth={photo.width}
                  imageHeight={photo.height}
                  opacity={0.6}
                />
              )}
            </div>

            {/* Overlay toggle */}
            <button
              onClick={() => setShowOverlay(!showOverlay)}
              className="absolute bottom-4 left-4 bg-gray-800 hover:bg-gray-700 text-white px-3 py-2 rounded text-sm"
            >
              {showOverlay ? 'Hide' : 'Show'} Overlays
            </button>
          </div>

          {/* Info section */}
          <div className="w-full lg:w-80 overflow-y-auto p-6 bg-gray-50">
            {/* Photo metadata */}
            <div className="mb-6">
              <h2 className="text-lg font-bold text-gray-900 mb-3">
                Photo Details
              </h2>
              <div className="space-y-2 text-sm text-gray-700">
                <div>
                  <span className="font-semibold">ID:</span>
                  <div className="text-xs text-gray-500 font-mono break-all">
                    {photo.id}
                  </div>
                </div>
                <div>
                  <span className="font-semibold">Position:</span>
                  <div>
                    {photo.x.toFixed(2)}m, {photo.y.toFixed(2)}m (H:{' '}
                    {photo.h.toFixed(1)}m)
                  </div>
                </div>
                <div>
                  <span className="font-semibold">Resolution:</span>
                  <div>
                    {photo.width} × {photo.height}px
                  </div>
                </div>
                <div>
                  <span className="font-semibold">Captured:</span>
                  <div>{new Date(photo.capturedAt).toLocaleString()}</div>
                </div>
              </div>
            </div>

            {/* Predictions */}
            <div>
              <h3 className="text-lg font-bold text-gray-900 mb-3">
                Predictions ({photo.predictions.length})
              </h3>
              {photo.predictions.length === 0 ? (
                <p className="text-sm text-gray-500">No predictions found</p>
              ) : (
                <div className="space-y-2">
                  {sortedPredictions.map((pred) => (
                    <div
                      key={pred.id}
                      className={`p-3 rounded border-2 ${PEST_CLASS_COLORS[pred.classId as keyof typeof PEST_CLASS_COLORS]} border-opacity-50`}
                    >
                      <div className="flex justify-between items-start mb-1">
                        <div className="font-semibold text-gray-900">
                          {PEST_CLASS_LABELS[pred.classId as keyof typeof PEST_CLASS_LABELS]}
                        </div>
                        <div className="text-sm font-bold text-gray-700">
                          {(pred.confidence * 100).toFixed(1)}%
                        </div>
                      </div>
                      <div className="w-full bg-gray-200 rounded-full h-2">
                        <div
                          className="bg-blue-500 h-2 rounded-full"
                          style={{
                            width: `${pred.confidence * 100}%`,
                          }}
                        />
                      </div>
                      <div className="text-xs text-gray-500 mt-2">
                        BBox: ({pred.bbox.xMin.toFixed(3)}, {pred.bbox.yMin.toFixed(3)}) →
                        ({pred.bbox.xMax.toFixed(3)}, {pred.bbox.yMax.toFixed(3)})
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default PhotoModal
