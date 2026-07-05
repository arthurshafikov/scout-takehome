package metrics
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	ThumbnailCacheHits   prometheus.Counter
	ThumbnailCacheMisses prometheus.Counter
	ThumbnailGenTime     prometheus.Histogram
	PhotosListCount      prometheus.Counter
	PhotosGetCount       prometheus.Counter
	APIErrors            prometheus.CounterVec
}

func NewMetrics() *Metrics {
	return &Metrics{
		ThumbnailCacheHits: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "scout_thumbnail_cache_hits_total",
			Help: "Total number of thumbnail cache hits",
		}),
		ThumbnailCacheMisses: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "scout_thumbnail_cache_misses_total",
			Help: "Total number of thumbnail cache misses",
		}),
		ThumbnailGenTime: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "scout_thumbnail_generation_seconds",
			Help:    "Time spent generating thumbnails",
			Buckets: prometheus.DefBuckets,
		}),
		PhotosListCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "scout_photos_list_total",
			Help: "Total number of photos list requests",
		}),
		PhotosGetCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "scout_photos_get_total",
			Help: "Total number of photos get requests",
		}),
		APIErrors: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "scout_api_errors_total",
			Help: "Total number of API errors",
		}, []string{"endpoint", "code"}),
	}
}

func (m *Metrics) Register() error {
	prometheus.MustRegister(m.ThumbnailCacheHits)
	prometheus.MustRegister(m.ThumbnailCacheMisses)
	prometheus.MustRegister(m.ThumbnailGenTime)
	prometheus.MustRegister(m.PhotosListCount)
	prometheus.MustRegister(m.PhotosGetCount)
	prometheus.MustRegister(&m.APIErrors)
	return nil
}
