import { FC, useEffect, useRef, useCallback, useState } from 'react'
import PhotoCard from './PhotoCard'
import PhotoModal from './PhotoModal'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { setCursor } from '@/features/filters/filtersSlice'
import { useListPhotosQuery } from '@/services/api'
import type { Photo } from '@/types/api'
import {
  GALLERY_FETCH_LIMIT,
  GALLERY_REFETCH_INTERVAL_MS,
  INFINITE_SCROLL_THRESHOLD,
} from '@/constants/greenhouse'

/**
 * Gallery page with infinite scroll pagination
 * Fetches photos with applied filters (class_id, min_confidence)
 * Uses cursor-based pagination for efficient scrolling
 */
const GalleryPage: FC = () => {
  const [selectedPhoto, setSelectedPhoto] = useState<Photo | null>(null)
  const [isModalOpen, setIsModalOpen] = useState(false)
  const observerTarget = useRef<HTMLDivElement>(null)

  const dispatch = useAppDispatch()
  const { selectedClass, minConfidence, cursor, locationCenter } = useAppSelector(
    (state) => state.filters,
  )

  // Fetch photos with filters
  const { data, isLoading, isFetching } = useListPhotosQuery(
    {
      cursor: cursor || undefined,
      limit: GALLERY_FETCH_LIMIT,
      class_id: selectedClass || undefined,
      min_confidence: minConfidence > 0 ? minConfidence : undefined,
    },
    {
      refetchOnMountOrArgChange: GALLERY_REFETCH_INTERVAL_MS,
    },
  )

  // Intersection observer for infinite scroll
  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        if (
          entries[0].isIntersecting &&
          data?.nextCursor &&
          !isFetching
        ) {
          dispatch(setCursor(data.nextCursor))
        }
      },
      { threshold: INFINITE_SCROLL_THRESHOLD },
    )

    if (observerTarget.current) {
      observer.observe(observerTarget.current)
    }

    return () => observer.disconnect()
  }, [data?.nextCursor, isFetching, dispatch])

  const handlePhotoClick = useCallback((photo: Photo) => {
    setSelectedPhoto(photo)
    setIsModalOpen(true)
  }, [])

  const handleCloseModal = useCallback(() => {
    setIsModalOpen(false)
    setTimeout(() => setSelectedPhoto(null), 300)
  }, [])

  /**
   * Calculate Euclidean distance between two points
   * @param x1 - First point X coordinate (meters)
   * @param y1 - First point Y coordinate (meters)
   * @param x2 - Second point X coordinate (meters)
   * @param y2 - Second point Y coordinate (meters)
   * @returns Distance in meters
   */
  const calculateDistance = (x1: number, y1: number, x2: number, y2: number): number => {
    return Math.sqrt((x2 - x1) ** 2 + (y2 - y1) ** 2)
  }

  /**
   * Filter photos by location proximity radius
   * Returns all photos if no location center is set
   * @param photosToFilter - Array of photos to filter
   * @returns Photos within proximity radius of selected location
   */
  const applyLocationFilter = (photosToFilter: Photo[]): Photo[] => {
    if (!locationCenter) return photosToFilter

    return photosToFilter.filter((photo) => {
      const distance = calculateDistance(
        locationCenter.x,
        locationCenter.y,
        photo.x,
        photo.y,
      )
      return distance <= locationCenter.radius
    })
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-96">
        <div className="text-gray-500">Loading photos...</div>
      </div>
    )
  }

  const photos = applyLocationFilter(data?.items || [])

  return (
    <>
      {/* Gallery grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        {photos.map((photo) => (
          <PhotoCard
            key={photo.id}
            photo={photo}
            onClick={() => handlePhotoClick(photo)}
          />
        ))}
      </div>

      {/* Infinite scroll trigger */}
      <div
        ref={observerTarget}
        className="flex justify-center py-8"
      >
        {isFetching && (
          <div className="text-gray-500">Loading more...</div>
        )}
        {!isFetching && !data?.nextCursor && photos.length > 0 && (
          <div className="text-gray-400 text-sm">No more photos</div>
        )}
        {photos.length === 0 && (
          <div className="text-gray-400">No photos found</div>
        )}
      </div>

      {/* Modal */}
      {selectedPhoto && (
        <PhotoModal
          photo={selectedPhoto}
          isOpen={isModalOpen}
          onClose={handleCloseModal}
        />
      )}
    </>
  )
}

export default GalleryPage
