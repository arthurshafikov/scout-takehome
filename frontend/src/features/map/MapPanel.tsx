import { FC, useRef, useState } from 'react'
import { Stage, Layer, Circle, Text, Rect, Line } from 'react-konva'
import { useListPhotosQuery } from '@/services/api'
import { PEST_CLASS_RGB, PEST_CLASS_LABELS } from '@/utils/classColors'
import type { Photo, Prediction } from '@/types/api'

/**
 * Interactive Konva.js map showing greenhouse layout (40x40m)
 * Each photo is a circle positioned at (x, y) with highest prediction color
 * Map fetches up to 200 photos with separate pagination
 */
const MapPanel: FC = () => {
  const stageRef = useRef<Konva.Stage>(null)
  const [hoveredPhoto, setHoveredPhoto] = useState<Photo | null>(null)
  const [scale, setScale] = useState(1)

  // Fetch photos for map (independent cache, limit 200)
  const { data: mapData, isLoading } = useListPhotosQuery(
    {
      limit: 200,
    },
    {
      refetchOnMountOrArgChange: 300, // Update less frequently than gallery
    },
  )

  const photos = mapData?.items || []

  // Grid constants (greenhouse: 40x40m)
  const GREENHOUSE_WIDTH = 40
  const GREENHOUSE_HEIGHT = 40
  const PIXEL_SCALE = 40 // pixels per meter

  // Get highest confidence prediction for a photo
  const getTopPrediction = (photo: Photo): Prediction | null => {
    if (photo.predictions.length === 0) return null
    return photo.predictions.reduce((prev, current) =>
      current.confidence > prev.confidence ? current : prev,
    )
  }

  // Handle map pan and zoom
  const handleWheel = (e: Konva.KonvaEventObject<WheelEvent>) => {
    e.evt.preventDefault()

    const stage = stageRef.current
    if (!stage) return

    const oldScale = scale
    const pointer = stage.getPointerPosition()
    if (!pointer) return

    const newScale = e.evt.deltaY > 0 ? oldScale * 0.95 : oldScale * 1.05
    setScale(Math.max(0.5, Math.min(newScale, 5)))

    const dx = (pointer.x / oldScale - pointer.x / newScale) * oldScale
    const dy = (pointer.y / oldScale - pointer.y / newScale) * oldScale

    stage.position({
      x: stage.x() + dx * oldScale,
      y: stage.y() + dy * oldScale,
    })
  }

  const width = GREENHOUSE_WIDTH * PIXEL_SCALE
  const height = GREENHOUSE_HEIGHT * PIXEL_SCALE

  if (isLoading) {
    return (
      <div className="bg-white p-4 rounded-lg shadow">
        <div className="flex items-center justify-center h-96">
          <div className="text-gray-500">Loading map...</div>
        </div>
      </div>
    )
  }

  return (
    <div className="bg-white p-4 rounded-lg shadow">
      <div className="mb-4">
        <h2 className="text-lg font-bold text-gray-900 mb-2">
          Greenhouse Map
        </h2>
        <p className="text-sm text-gray-600">
          {photos.length} photos • Scroll to zoom • Drag to pan
        </p>
      </div>

      {/* Konva Stage */}
      <div className="border border-gray-200 rounded overflow-hidden bg-gray-50">
        <Stage
          ref={stageRef}
          width={Math.min(800, window.innerWidth - 64)}
          height={400}
          onWheel={handleWheel}
          draggable
          scaleX={scale}
          scaleY={scale}
        >
          <Layer>
            {/* Grid background label */}
            <Text
              x={10}
              y={10}
              text={`Greenhouse: ${GREENHOUSE_WIDTH}m × ${GREENHOUSE_HEIGHT}m`}
              fontSize={12}
              fill="#999"
            />

            {/* Border */}
            <Rect
              x={0}
              y={0}
              width={width}
              height={height}
              stroke="#999"
              strokeWidth={2}
              fill="transparent"
            />

            {/* Grid lines */}
            {[...Array(5)].map((_, i) => {
              const pos = (i + 1) * (width / 5)
              return (
                <g key={`grid-${i}`}>
                  {/* Vertical line */}
                  <Line
                    key={`v-${i}`}
                    points={[pos, 0, pos, height]}
                    stroke="#ddd"
                    strokeWidth={1}
                  />
                  {/* Horizontal line */}
                  <Line
                    key={`h-${i}`}
                    points={[0, pos, width, pos]}
                    stroke="#ddd"
                    strokeWidth={1}
                  />
                </g>
              )
            })}

            {/* Photo points */}
            {photos.map((photo) => {
              const topPred = getTopPrediction(photo)
              const color = topPred
                ? PEST_CLASS_RGB[topPred.classId]
                : [128, 128, 128]
              const isHovered = hoveredPhoto?.id === photo.id

              return (
                <Circle
                  key={photo.id}
                  x={photo.x * PIXEL_SCALE}
                  y={photo.y * PIXEL_SCALE}
                  radius={isHovered ? 8 : 4}
                  fill={`rgb(${color[0]}, ${color[1]}, ${color[2]})`}
                  opacity={0.7}
                  onMouseEnter={() => setHoveredPhoto(photo)}
                  onMouseLeave={() => setHoveredPhoto(null)}
                  stroke={isHovered ? '#000' : 'none'}
                  strokeWidth={2}
                />
              )
            })}
          </Layer>
        </Stage>
      </div>

      {/* Hovered photo info */}
      {hoveredPhoto && (
        <div className="mt-4 p-3 bg-gray-50 rounded border border-gray-200">
          <div className="text-sm">
            <div className="font-semibold text-gray-900 mb-1">
              Position: {hoveredPhoto.x.toFixed(2)}m, {hoveredPhoto.y.toFixed(2)}m
            </div>
            <div className="text-gray-600 text-xs mb-2">
              {new Date(hoveredPhoto.capturedAt).toLocaleString()}
            </div>
            {hoveredPhoto.predictions.length > 0 && (
              <div className="text-xs">
                <div className="font-medium text-gray-700 mb-1">
                  Top predictions:
                </div>
                {hoveredPhoto.predictions
                  .sort((a, b) => b.confidence - a.confidence)
                  .slice(0, 3)
                  .map((pred) => (
                    <div key={pred.id} className="text-gray-600">
                      • {PEST_CLASS_LABELS[pred.classId]}: {(pred.confidence * 100).toFixed(0)}%
                    </div>
                  ))}
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  )
}

export default MapPanel
