export default function AppLayout() {
  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
          <h1 className="text-3xl font-bold text-gray-900">Scout</h1>
        </div>
      </header>
      <main>
        <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
          {/* Placeholder for routes */}
          <div className="text-center py-12">
            <p className="text-gray-500">Frontend scaffold initialized.</p>
          </div>
        </div>
      </main>
    </div>
  )
}
