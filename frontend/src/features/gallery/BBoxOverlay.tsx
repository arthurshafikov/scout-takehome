import { FC } from 'react'
import { normalizedToPixels } from '@/utils/bbox'
import { PEST_CLASS_RGB, PEST_CLASS_LABELS } from '@/utils/classColors'
import type { Photo } from '@/types/api'

interface BBoxOverlayProps {
  photo: Photo
  imageWidth: number
  imageHeight: number
  opacity?: number
}

/**
 * SVG overlay rendering bounding boxes for predictions
 * Uses normalized [0,1] coordinates and scales to image dimensions
 * No pixel-rounding issues - native SVG scaling
 */
const BBoxOverlay: FC<BBoxOverlayProps> = ({
  photo,
  imageWidth,
  imageHeight,
  opacity = 0.5,
}) => {
  return (
    <svg
      className="absolute inset-0 w-full h-full pointer-events-none"
      viewBox={`0 0 ${imageWidth} ${imageHeight}`}
      preserveAspectRatio="xMidYMid meet"
    >
      {/* Bounding boxes */}
      {photo.predictions.map((prediction) => {
        const bbox = normalizedToPixels(
          prediction.bbox,
          imageWidth,
          imageHeight,
        )
        const [r, g, b] = PEST_CLASS_RGB[prediction.classId]
        const color = `rgba(${r}, ${g}, ${b}, ${opacity})`
        const borderColor = `rgb(${r}, ${g}, ${b})`

        return (
          <g key={prediction.id}>
            {/* Box fill */}
            <rect
              x={bbox.xMin}
              y={bbox.yMin}
              width={bbox.xMax - bbox.xMin}
              height={bbox.yMax - bbox.yMin}
              fill={color}
            />
            {/* Box border */}
            <rect
              x={bbox.xMin}
              y={bbox.yMin}
              width={bbox.xMax - bbox.xMin}
              height={bbox.yMax - bbox.yMin}
              fill="none"
              stroke={borderColor}
              strokeWidth="2"
            />
            {/* Label background */}
            <rect
              x={bbox.xMin}
              y={Math.max(0, bbox.yMin - 24)}
              width={120}
              height={22}
              fill={borderColor}
            />
            {/* Label text */}
            <text
              x={bbox.xMin + 2}
              y={Math.max(16, bbox.yMin - 6)}
              fill="white"
              fontSize="12"
              fontWeight="bold"
              fontFamily="sans-serif"
            >
              {PEST_CLASS_LABELS[prediction.classId]} (
              {(prediction.confidence * 100).toFixed(0)}%)
            </text>
          </g>
        )
      })}
    </svg>
  )
}

export default BBoxOverlay
