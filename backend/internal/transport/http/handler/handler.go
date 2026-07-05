package handler

import (
	"context"

	"github.com/arthurshafikov/scout-takehome/backend/internal/services"
	"github.com/arthurshafikov/scout-takehome/backend/internal/services/thumbnail"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	ctx                context.Context
	services           *services.Services
	thumbnailGenerator *thumbnail.ThumbnailGenerator
}

func NewHandler(
	ctx context.Context,
	services *services.Services,
	thumbnailGenerator *thumbnail.ThumbnailGenerator,
) *Handler {
	return &Handler{
		ctx:                ctx,
		services:           services,
		thumbnailGenerator: thumbnailGenerator,
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

	// Metrics endpoint
	e.GET("/metrics", h.getMetrics)
}

func (h *Handler) getMetrics(c *gin.Context) {
	// TODO: Implement metrics endpoint using Prometheus
	c.JSON(200, gin.H{
		"status": "ok",
		"message": "Metrics endpoint - TODO: implement Prometheus metrics",
	})
}

func (h *Handler) initHealthCheck(e *gin.Engine) {
	e.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
