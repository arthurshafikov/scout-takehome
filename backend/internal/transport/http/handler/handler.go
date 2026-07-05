package handler

import (
	"context"

	"github.com/arthurshafikov/scout-takehome/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	ctx      context.Context
	services *services.Services
}

func NewHandler(
	ctx context.Context,
	services *services.Services,
) *Handler {
	return &Handler{
		ctx:      ctx,
		services: services,
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
}

func (h *Handler) initHealthCheck(e *gin.Engine) {
	e.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
