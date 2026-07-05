import type { PestClass } from '@/types/api'

/**
 * Tailwind color classes for each pest class
 * Used for bbox overlays and legend display
 */
export const PEST_CLASS_COLORS: Record<PestClass, string> = {
  powdery_mildew: 'bg-red-200',
  mirid: 'bg-yellow-200',
  whitefly_aphid: 'bg-blue-200',
  miner_tuta: 'bg-purple-200',
  thrips: 'bg-orange-200',
  spider_mites: 'bg-pink-200',
}

/**
 * RGB values for each pest class (for Konva drawing)
 */
export const PEST_CLASS_RGB: Record<PestClass, [number, number, number]> = {
  powdery_mildew: [239, 68, 68], // red-500
  mirid: [234, 179, 8], // yellow-500
  whitefly_aphid: [59, 130, 246], // blue-500
  miner_tuta: [168, 85, 247], // purple-500
  thrips: [249, 115, 22], // orange-500
  spider_mites: [236, 72, 153], // pink-500
}

/**
 * Friendly display names for pest classes
 */
export const PEST_CLASS_LABELS: Record<PestClass, string> = {
  powdery_mildew: 'Powdery Mildew',
  mirid: 'Mirid',
  whitefly_aphid: 'Whitefly/Aphid',
  miner_tuta: 'Tuta Absoluta',
  thrips: 'Thrips',
  spider_mites: 'Spider Mites',
}

export function getColorClass(pestClass: PestClass): string {
  return PEST_CLASS_COLORS[pestClass] || 'bg-gray-200'
}

export function getColorRGB(pestClass: PestClass): [number, number, number] {
  return PEST_CLASS_RGB[pestClass] || [128, 128, 128]
}

export function getLabel(pestClass: PestClass): string {
  return PEST_CLASS_LABELS[pestClass] || pestClass
}
