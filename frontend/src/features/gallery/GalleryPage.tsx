import { FC, useEffect, useRef, useCallback, useState } from 'react'
import PhotoCard from './PhotoCard'
import PhotoModal from './PhotoModal'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { setCursor } from '@/features/filters/filtersSlice'
import { useListPhotosQuery } from '@/services/api'
import type { Photo } from '@/types/api'

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
      limit: 50,
      class_id: selectedClass || undefined,
      min_confidence: minConfidence > 0 ? minConfidence : undefined,
    },
    {
      refetchOnMountOrArgChange: 60,
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
      { threshold: 0.1 },
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

  // Calculate distance between two points (meters)
  const calculateDistance = (x1: number, y1: number, x2: number, y2: number): number => {
    return Math.sqrt((x2 - x1) ** 2 + (y2 - y1) ** 2)
  }

  // Filter photos by location proximity
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
