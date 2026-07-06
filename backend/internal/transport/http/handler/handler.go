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
	apiKey             string
}

func NewHandler(
	ctx context.Context,
	services *services.Services,
	thumbnailGenerator *thumbnail.ThumbnailGenerator,
	m *metrics.Metrics,
	apiKey string,
) *Handler {
	return &Handler{
		ctx:                ctx,
		services:           services,
		thumbnailGenerator: thumbnailGenerator,
		metrics:            m,
		apiKey:             apiKey,
	}
}

func (h *Handler) Init(e *gin.Engine) {
	h.initHealthCheck(e)

	// API key middleware for protected routes
	authMiddleware := h.apiKeyMiddleware()

	// Photos endpoints
	photos := e.Group("/photos")
	photos.Use(authMiddleware)
	{
		photos.GET("", h.listPhotos)
		photos.GET("/:id", h.getPhoto)
		photos.POST("/:id/upload-link", h.generateUploadLink)
	}

	// Thumbnails endpoint (protected)
	e.GET("/thumbnails/:id", authMiddleware, h.getThumbnail)

	// Metrics endpoint - expose Prometheus metrics
	e.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

func (h *Handler) apiKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			traceID := getTraceID(c)
			c.JSON(401, APIResponse[interface{}]{
				Success: false,
				Error: &ErrorBody{
					Code:    "unauthorized",
					Message: "Missing X-API-Key header",
				},
				TraceID: traceID,
			})
			c.Abort()
			return
		}

		if apiKey != h.apiKey {
			traceID := getTraceID(c)
			c.JSON(401, APIResponse[interface{}]{
				Success: false,
				Error: &ErrorBody{
					Code:    "unauthorized",
					Message: "Invalid API key",
				},
				TraceID: traceID,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (h *Handler) initHealthCheck(e *gin.Engine) {
	e.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
