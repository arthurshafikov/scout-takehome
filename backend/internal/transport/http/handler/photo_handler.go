package handler

import (
	"context"
	"strconv"

	"github.com/arthurshafikov/scout-takehome/backend/internal/core/models"
	"github.com/arthurshafikov/scout-takehome/backend/internal/core/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PhotoPage struct {
	Items     []models.Photo `json:"items"`
	NextToken string         `json:"next_token,omitempty"`
}

type UploadLinkRequest struct {
	ContentType string `json:"content_type" binding:"required"`
}

func (h *Handler) listPhotos(c *gin.Context) {
	traceID := getTraceID(c)

	cursor := c.Query("cursor")
	limitStr := c.DefaultQuery("limit", "10")
	classID := c.Query("class_id")
	minConfidenceStr := c.Query("min_confidence")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(400, APIResponse[interface{}]{
			Success: false,
			Error: &ErrorBody{
				Code:    "invalid_limit",
				Message: "Invalid limit parameter",
			},
			TraceID: traceID,
		})
		return
	}

	if limit < 1 || limit > 100 {
		limit = 10
	}

	var minConfidence float64 = 0.0
	if minConfidenceStr != "" {
		conf, err := strconv.ParseFloat(minConfidenceStr, 64)
		if err != nil {
			c.JSON(400, APIResponse[interface{}]{
				Success: false,
				Error: &ErrorBody{
					Code:    "invalid_min_confidence",
					Message: "Invalid min_confidence parameter",
				},
				TraceID: traceID,
			})
			return
		}
		minConfidence = conf
	}

	ctx := context.Background()
	photos, nextToken, err := h.services.PhotoService.ListPhotos(ctx, types.ListPhotosParams{
		Cursor:        cursor,
		Limit:         limit,
		ClassID:       classID,
		MinConfidence: minConfidence,
	})
	if err != nil {
		errorResponse(c, err, traceID)
		return
	}

	// Enrich photos with original URLs
	for i := range photos {
		url, err := h.services.StorageService.GetOriginalURL(ctx, photos[i].ID)
		if err != nil {
			// Log but don't fail - continue with nil URL
			continue
		}
		photos[i].OriginalURL = url
	}

	successResponse(c, 200, PhotoPage{
		Items:     photos,
		NextToken: nextToken,
	}, traceID)
}

func (h *Handler) getPhoto(c *gin.Context) {
	traceID := getTraceID(c)
	photoID := c.Param("id")

	if photoID == "" {
		c.JSON(400, APIResponse[interface{}]{
			Success: false,
			Error: &ErrorBody{
				Code:    "missing_photo_id",
				Message: "Photo ID is required",
			},
			TraceID: traceID,
		})
		return
	}

	ctx := context.Background()
	photo, err := h.services.PhotoService.GetPhoto(ctx, photoID)
	if err != nil {
		errorResponse(c, err, traceID)
		return
	}

	// Enrich with original URL
	url, err := h.services.StorageService.GetOriginalURL(ctx, photo.ID)
	if err == nil {
		photo.OriginalURL = url
	}

	successResponse(c, 200, photo, traceID)
}

func (h *Handler) generateUploadLink(c *gin.Context) {
	traceID := getTraceID(c)
	photoID := c.Param("id")

	if photoID == "" {
		c.JSON(400, APIResponse[interface{}]{
			Success: false,
			Error: &ErrorBody{
				Code:    "missing_photo_id",
				Message: "Photo ID is required",
			},
			TraceID: traceID,
		})
		return
	}

	var req UploadLinkRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, APIResponse[interface{}]{
			Success: false,
			Error: &ErrorBody{
				Code:    "invalid_request",
				Message: "Invalid request body",
			},
			TraceID: traceID,
		})
		return
	}

	// Validate photoID is valid UUID
	if _, err := uuid.Parse(photoID); err != nil {
		c.JSON(400, APIResponse[interface{}]{
			Success: false,
			Error: &ErrorBody{
				Code:    "invalid_photo_id",
				Message: "Photo ID must be a valid UUID",
			},
			TraceID: traceID,
		})
		return
	}

	ctx := context.Background()
	uploadLink, err := h.services.StorageService.GenerateUploadLink(ctx, photoID, req.ContentType)
	if err != nil {
		errorResponse(c, err, traceID)
		return
	}

	successResponse(c, 200, uploadLink, traceID)
}
