import { useState } from 'react'
import GalleryPage from '@/features/gallery/GalleryPage'
import FilterPanel from '@/features/filters/FilterPanel'
import MapPanel from '@/features/map/MapPanel'
import MapLegend from '@/features/map/MapLegend'

type TabType = 'gallery' | 'map'

export default function AppLayout() {
  const [activeTab, setActiveTab] = useState<TabType>('gallery')

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow sticky top-0 z-40">
        <div className="max-w-7xl mx-auto py-4 px-4 sm:px-6 lg:px-8">
          <h1 className="text-3xl font-bold text-gray-900">Scout</h1>
          <p className="text-sm text-gray-600 mt-1">
            Greenhouse Pest & Disease Monitoring
          </p>
        </div>
      </header>

      {/* Main content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Tab navigation */}
        <div className="flex gap-4 mb-6 border-b border-gray-200">
          <button
            onClick={() => setActiveTab('gallery')}
            className={`px-4 py-2 font-medium text-sm border-b-2 transition-colors ${
              activeTab === 'gallery'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-600 hover:text-gray-900'
            }`}
          >
            Gallery
          </button>
          <button
            onClick={() => setActiveTab('map')}
            className={`px-4 py-2 font-medium text-sm border-b-2 transition-colors ${
              activeTab === 'map'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-600 hover:text-gray-900'
            }`}
          >
            Map View
          </button>
        </div>

        {/* Gallery view */}
        {activeTab === 'gallery' && (
          <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
            {/* Filter sidebar */}
            <div className="lg:col-span-1">
              <div className="sticky top-24">
                <FilterPanel />
              </div>
            </div>

            {/* Gallery main content */}
            <div className="lg:col-span-3">
              <GalleryPage />
            </div>
          </div>
        )}

        {/* Map view */}
        {activeTab === 'map' && (
          <div className="space-y-6">
            <MapPanel />
            <MapLegend />
          </div>
        )}
      </main>
    </div>
  )
}
