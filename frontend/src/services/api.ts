import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'
import type { Photo, PhotoPage, ListPhotosParams, HealthCheck } from '@/types/api'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

export const apiSlice = createApi({
  reducerPath: 'api',
  baseQuery: fetchBaseQuery({
    baseUrl: `${API_URL}/api`,
    prepareHeaders: (headers) => {
      const apiKey = import.meta.env.VITE_API_KEY
      if (apiKey) {
        headers.set('Authorization', `Bearer ${apiKey}`)
      }
      return headers
    },
  }),
  tagTypes: ['Photos', 'Photo'],
  endpoints: (builder) => ({
    /**
     * List photos with cursor pagination and optional filters
     * Returns paginated results with nextCursor for infinite scroll
     * 
     * For infinite scroll: keep cursor in Redux state, update it when
     * user scrolls to bottom. Each query arg change creates new cache entry.
     */
    listPhotos: builder.query<PhotoPage, ListPhotosParams>({
      query: (params) => {
        const searchParams = new URLSearchParams()
        if (params.cursor) searchParams.append('cursor', params.cursor)
        if (params.limit) searchParams.append('limit', String(params.limit))
        if (params.class_id) searchParams.append('class_id', params.class_id)
        if (params.min_confidence !== undefined)
          searchParams.append('min_confidence', String(params.min_confidence))
        return `/photos?${searchParams.toString()}`
      },
      providesTags: ['Photos'],
    }),

    /**
     * Get single photo by ID with all predictions
     */
    getPhoto: builder.query<Photo, string>({
      query: (id) => `/photos/${id}`,
      providesTags: (result, _, id) => [{ type: 'Photo', id }],
    }),

    /**
     * Get thumbnail URL for a photo
     * Note: Returns URL string, not binary blob
     * Browser caches via img tag naturally
     */
    getThumbnailUrl: builder.query<string, string>({
      queryFn: (id) => {
        const url = `${API_URL}/api/thumbnails/${id}`
        return { data: url }
      },
    }),

    /**
     * Health check endpoint
     * Returns success: true if backend is healthy
     */
    healthCheck: builder.query<HealthCheck, void>({
      query: () => '/healthz',
    }),
  }),
})

export const {
  useListPhotosQuery,
  useGetPhotoQuery,
  useGetThumbnailUrlQuery,
  useHealthCheckQuery,
} = apiSlice
