import { Component, ReactNode } from 'react'

interface ErrorBoundaryProps {
  children: ReactNode
  fallback?: ReactNode
}

interface ErrorBoundaryState {
  hasError: boolean
  error: Error | null
}

/**
 * Error Boundary component for graceful error handling
 * Catches errors in child components and displays fallback UI
 * Prevents entire app crash when a section fails
 */
export class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  constructor(props: ErrorBoundaryProps) {
    super(props)
    this.state = { hasError: false, error: null }
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { hasError: true, error }
  }

  componentDidCatch(error: Error) {
    console.error('ErrorBoundary caught error:', error)
  }

  render() {
    if (this.state.hasError) {
      return (
        this.props.fallback || (
          <div className="bg-red-50 border border-red-200 rounded-lg p-4 m-4">
            <h2 className="text-red-800 font-bold mb-2">Something went wrong</h2>
            <p className="text-red-700 text-sm">{this.state.error?.message}</p>
          </div>
        )
      )
    }

    return this.props.children
  }
}
