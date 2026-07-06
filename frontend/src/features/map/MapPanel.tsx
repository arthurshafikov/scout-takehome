import { FC, useRef, useState } from 'react'
import { Stage, Layer, Circle, Text, Rect, Line } from 'react-konva'
import type Konva from 'konva'
import { useDispatch, useSelector } from 'react-redux'
import { useListPhotosQuery } from '@/services/api'
import { setLocationCenter } from '@/features/filters/filtersSlice'
import type { RootState } from '@/app/store'
import { PEST_CLASS_RGB, PEST_CLASS_LABELS } from '@/utils/classColors'
import type { Photo, Prediction } from '@/types/api'
import {
  GREENHOUSE_WIDTH,
  GREENHOUSE_HEIGHT,
  PIXEL_SCALE,
  PROXIMITY_RADIUS,
  MIN_ZOOM,
  MAX_ZOOM,
  ZOOM_IN_FACTOR,
  ZOOM_OUT_FACTOR,
  MAP_FETCH_LIMIT,
  MAP_REFETCH_INTERVAL_MS,
} from '@/constants/greenhouse'

/**
 * Interactive Konva.js map showing greenhouse layout (40x40m)
 * Each photo is a circle positioned at (x, y) with highest prediction color
 * Map fetches up to 200 photos with separate pagination
 */
const MapPanel: FC = () => {
  const stageRef = useRef<Konva.Stage>(null)
  const dispatch = useDispatch()
  const locationCenter = useSelector((state: RootState) => state.filters.locationCenter)
  const [hoveredPhoto, setHoveredPhoto] = useState<Photo | null>(null)
  const [scale, setScale] = useState(0.5) // Default zoom out to fit entire map

  // Fetch photos for map (independent cache, up to MAP_FETCH_LIMIT)
  const { data: mapData, isLoading } = useListPhotosQuery(
    {
      limit: MAP_FETCH_LIMIT,
    },
    {
      refetchOnMountOrArgChange: MAP_REFETCH_INTERVAL_MS,
    },
  )

  const photos = mapData?.items || []

  /**
   * Get highest confidence prediction for a photo
   * @param photo - Photo object with predictions array
   * @returns Prediction with highest confidence, or null if no predictions
   */
  const getTopPrediction = (photo: Photo): Prediction | null => {
    if (photo.predictions.length === 0) return null
    return photo.predictions.reduce((prev, current) =>
      current.confidence > prev.confidence ? current : prev,
    )
  }

  /**
   * Handle mouse wheel zoom
   * Zooms in/out around pointer position with MIN_ZOOM to MAX_ZOOM constraints
   */
  const handleWheel = (e: Konva.KonvaEventObject<WheelEvent>) => {
    e.evt.preventDefault()

    const stage = stageRef.current
    if (!stage) return

    const oldScale = scale
    const pointer = stage.getPointerPosition()
    if (!pointer) return

    const newScale = e.evt.deltaY > 0 ? oldScale * ZOOM_OUT_FACTOR : oldScale * ZOOM_IN_FACTOR
    setScale(Math.max(MIN_ZOOM, Math.min(newScale, MAX_ZOOM)))

    const dx = (pointer.x / oldScale - pointer.x / newScale) * oldScale
    const dy = (pointer.y / oldScale - pointer.y / newScale) * oldScale

    stage.position({
      x: stage.x() + dx * oldScale,
      y: stage.y() + dy * oldScale,
    })
  }

  /**
   * Handle click on map to set location-based filter
   * Converts pixel coordinates to greenhouse meters, accounting for zoom
   * Dispatches setLocationCenter with PROXIMITY_RADIUS for filtering
   * Ignores clicks on photo circles (only responds to stage background)
   */
  const handleStageClick = (e: Konva.KonvaEventObject<MouseEvent>) => {
    // Only handle clicks on the stage background, not on photo circles
    if (e.target !== stageRef.current) {
      return
    }

    const stage = stageRef.current
    if (!stage) return

    const pointer = stage.getPointerPosition()
    if (!pointer) return

    // Convert pixel coordinates to greenhouse meters
    // Divide by PIXEL_SCALE and by current scale to account for zoom
    const x = pointer.x / PIXEL_SCALE / scale
    const y = pointer.y / PIXEL_SCALE / scale

    // Clamp to greenhouse bounds
    const clampedX = Math.max(0, Math.min(x, GREENHOUSE_WIDTH))
    const clampedY = Math.max(0, Math.min(y, GREENHOUSE_HEIGHT))

    // Set location center for proximity filter
    dispatch(
      setLocationCenter({
        x: clampedX,
        y: clampedY,
        radius: PROXIMITY_RADIUS,
      }),
    )
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
          {photos.length} photos • Click to filter by location • Scroll to zoom • Drag to pan
          {locationCenter && ` • Location filter active (${locationCenter.radius}m radius)`}
        </p>
      </div>

      {/* Konva Stage */}
      <div className="border border-gray-200 rounded overflow-hidden bg-gray-50">
        <Stage
          ref={stageRef}
          width={Math.min(800, window.innerWidth - 64)}
          height={400}
          onWheel={handleWheel}
          onClick={handleStageClick}
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
            {Array.from({ length: 5 }).map((_, i) => {
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

            {/* Proximity circle (if location is selected) */}
            {locationCenter && (
              <>
                <Circle
                  x={locationCenter.x * PIXEL_SCALE}
                  y={locationCenter.y * PIXEL_SCALE}
                  radius={locationCenter.radius * PIXEL_SCALE}
                  fill="rgba(59, 130, 246, 0.1)"
                  stroke="#3b82f6"
                  strokeWidth={2}
                />
                <Circle
                  x={locationCenter.x * PIXEL_SCALE}
                  y={locationCenter.y * PIXEL_SCALE}
                  radius={6}
                  fill="#3b82f6"
                  stroke="#fff"
                  strokeWidth={2}
                />
              </>
            )}

            {/* Photo points */}
            {photos.map((photo) => {
              const topPred = getTopPrediction(photo)
              const color = topPred
                ? PEST_CLASS_RGB[topPred.classId as keyof typeof PEST_CLASS_RGB]
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
                {[...hoveredPhoto.predictions]
                  .sort((a, b) => b.confidence - a.confidence)
                  .slice(0, 3)
                  .map((pred) => (
                    <div key={pred.id} className="text-gray-600">
                      • {PEST_CLASS_LABELS[pred.classId as keyof typeof PEST_CLASS_LABELS]}: {(pred.confidence * 100).toFixed(0)}%
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
