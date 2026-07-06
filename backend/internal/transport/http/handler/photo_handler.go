package handler

import (
	"context"
	"strconv"

	"github.com/arthurshafikov/scout-takehome/backend/internal/core/models"
	"github.com/arthurshafikov/scout-takehome/backend/internal/core/types"
	"github.com/arthurshafikov/scout-takehome/backend/internal/services/thumbnail"
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

	var minConfidence = 0.0
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
	photos, nextToken, err := h.services.ListPhotos(ctx, types.ListPhotosParams{
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
		url, err := h.services.GetOriginalURL(ctx, photos[i].ID)
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
	photo, err := h.services.GetPhoto(ctx, photoID)
	if err != nil {
		errorResponse(c, err, traceID)
		return
	}

	// Enrich with original URL
	url, err := h.services.GetOriginalURL(ctx, photo.ID)
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
	uploadLink, err := h.services.GenerateUploadLink(ctx, photoID, req.ContentType)
	if err != nil {
		errorResponse(c, err, traceID)
		return
	}

	successResponse(c, 200, uploadLink, traceID)
}

func (h *Handler) getThumbnail(c *gin.Context) {
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

	// Parse query parameters
	widthStr := c.DefaultQuery("w", "400")
	qualityStr := c.DefaultQuery("q", "85")

	width, err := strconv.Atoi(widthStr)
	if err != nil || width < 1 || width > 2000 {
		width = 400
	}

	quality, err := strconv.Atoi(qualityStr)
	if err != nil || quality < 1 || quality > 100 {
		quality = 85
	}

	// Generate thumbnail
	ctx := context.Background()
	opts := thumbnail.ThumbnailOptions{
		Width:   width,
		Height:  int(float64(width) * 0.75), // 4:3 aspect ratio
		Quality: quality,
	}

	thumbnailData, err := h.thumbnailGenerator.Generate(ctx, photoID, opts)
	if err != nil {
		errorResponse(c, err, traceID)
		return
	}

	// Set response headers and stream the image
	c.Data(200, "image/jpeg", thumbnailData)
}
