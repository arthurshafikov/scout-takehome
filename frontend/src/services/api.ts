import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'
import type { Photo, PhotoPage, ListPhotosParams, HealthCheck } from '@/types/api'

// Always use relative path for API calls - this ensures they work through nginx proxy
// In production (Docker): localhost:5173 -> nginx -> /api -> backend at http://app:8080
// In development: localhost:5173 -> Vite proxy or direct to http://localhost:8080
const API_BASE_URL = '/api'

export const apiSlice = createApi({
  reducerPath: 'api',
  baseQuery: fetchBaseQuery({
    baseUrl: API_BASE_URL,
    prepareHeaders: (headers) => {
      const apiKey = import.meta.env.VITE_API_KEY
      if (apiKey) {
        headers.set('X-API-Key', apiKey)
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
      queryFn: async (params, _api, _options, baseQuery) => {
        const searchParams = new URLSearchParams()
        if (params.cursor) searchParams.append('cursor', params.cursor)
        if (params.limit) searchParams.append('limit', String(params.limit))
        if (params.class_id) searchParams.append('class_id', params.class_id)
        if (params.min_confidence !== undefined)
          searchParams.append('min_confidence', String(params.min_confidence))
        
        const result = await baseQuery(`/photos?${searchParams.toString()}`)
        
        if (result.error) {
          return { error: result.error }
        }
        
        // Transform response: backend returns { success, data: { items, next_token } }
        // We need to return { items, nextCursor }
        const response = result.data as { success: boolean; data: { items: Photo[]; next_token?: string } }
        if (response.data?.items) {
          return {
            data: {
              items: response.data.items,
              nextCursor: response.data.next_token,
            }
          }
        }
        
        return { data: { items: [] } }
      },
      providesTags: ['Photos'],
    }),

    /**
     * Get single photo by ID with all predictions
     */
    getPhoto: builder.query<Photo, string>({
      query: (id) => `/photos/${id}`,
      providesTags: (_, __, id) => [{ type: 'Photo', id }],
    }),

    /**
     * Get thumbnail URL for a photo
     * Note: Returns URL string, not binary blob
     * Browser caches via img tag naturally
     */
    getThumbnailUrl: builder.query<string, string>({
      queryFn: (id) => {
        const url = `${API_BASE_URL}/thumbnails/${id}`
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
