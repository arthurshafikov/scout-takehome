package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	// Thumbnail metrics
	ThumbnailCacheHits      prometheus.Counter
	ThumbnailCacheMisses    prometheus.Counter
	ThumbnailGenerationTime prometheus.Histogram

	// HTTP metrics
	HTTPRequestsTotal   prometheus.CounterVec
	HTTPRequestDuration prometheus.HistogramVec
	HTTPErrorsTotal     prometheus.CounterVec
}

func NewMetrics() *Metrics {
	return &Metrics{
		// Thumbnail cache metrics
		ThumbnailCacheHits: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "scout_thumbnail_cache_hits_total",
				Help: "Total number of thumbnail cache hits",
			},
		),
		ThumbnailCacheMisses: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "scout_thumbnail_cache_misses_total",
				Help: "Total number of thumbnail cache misses",
			},
		),
		ThumbnailGenerationTime: promauto.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "scout_thumbnail_generation_seconds",
				Help:    "Time spent generating thumbnails in seconds",
				Buckets: []float64{0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
			},
		),

		// HTTP request metrics
		HTTPRequestsTotal: *promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "scout_http_requests_total",
				Help: "Total HTTP requests by endpoint, method, and status",
			},
			[]string{"endpoint", "method", "status"},
		),
		HTTPRequestDuration: *promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "scout_http_request_duration_seconds",
				Help:    "HTTP request duration in seconds by endpoint",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"endpoint", "method"},
		),
		HTTPErrorsTotal: *promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "scout_http_errors_total",
				Help: "Total HTTP errors by endpoint and status code",
			},
			[]string{"endpoint", "status"},
		),
	}
}

// RecordThumbnailCacheHit records a thumbnail cache hit
func (m *Metrics) RecordThumbnailCacheHit() {
	m.ThumbnailCacheHits.Inc()
}

// RecordThumbnailCacheMiss records a thumbnail cache miss
func (m *Metrics) RecordThumbnailCacheMiss() {
	m.ThumbnailCacheMisses.Inc()
}

// RecordThumbnailGenerationTime records thumbnail generation time
func (m *Metrics) RecordThumbnailGenerationTime(duration float64) {
	m.ThumbnailGenerationTime.Observe(duration)
}

// RecordHTTPRequest records HTTP request metrics
func (m *Metrics) RecordHTTPRequest(endpoint, method, status string, duration float64) {
	m.HTTPRequestsTotal.WithLabelValues(endpoint, method, status).Inc()
	m.HTTPRequestDuration.WithLabelValues(endpoint, method).Observe(duration)
}

// RecordHTTPError records an HTTP error
func (m *Metrics) RecordHTTPError(endpoint, status string) {
	m.HTTPErrorsTotal.WithLabelValues(endpoint, status).Inc()
}
