package handler

import (
	"context"

	"github.com/arthurshafikov/scout-takehome/backend/internal/services"
	"github.com/arthurshafikov/scout-takehome/backend/internal/services/metrics"
	"github.com/arthurshafikov/scout-takehome/backend/internal/services/thumbnail"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handler struct {
	ctx                context.Context
	services           *services.Services
	thumbnailGenerator *thumbnail.ThumbnailGenerator
	metrics            *metrics.Metrics
}

func NewHandler(
	ctx context.Context,
	services *services.Services,
	thumbnailGenerator *thumbnail.ThumbnailGenerator,
	m *metrics.Metrics,
) *Handler {
	return &Handler{
		ctx:                ctx,
		services:           services,
		thumbnailGenerator: thumbnailGenerator,
		metrics:            m,
	}
}

func (h *Handler) Init(e *gin.Engine) {
	h.initHealthCheck(e)

	// Photos endpoints
	photos := e.Group("/photos")
	{
		photos.GET("", h.listPhotos)
		photos.GET("/:id", h.getPhoto)
		photos.POST("/:id/upload-link", h.generateUploadLink)
	}

	// Thumbnails endpoint
	e.GET("/thumbnails/:id", h.getThumbnail)

	// Metrics endpoint - expose Prometheus metrics
	e.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

func (h *Handler) initHealthCheck(e *gin.Engine) {
	e.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
