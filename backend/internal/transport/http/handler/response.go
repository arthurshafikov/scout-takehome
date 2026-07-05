package handler

import (
	"errors"
	"net/http"

	apierr "github.com/arthurshafikov/scout-takehome/backend/internal/core/errors"
	"github.com/gin-gonic/gin"
)

type APIResponse[T any] struct {
	Success bool        `json:"success"`
	Data    T           `json:"data,omitempty"`
	Error   *ErrorBody  `json:"error,omitempty"`
	TraceID string      `json:"trace_id,omitempty"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func successResponse[T any](c *gin.Context, status int, data T, traceID string) {
	c.JSON(status, APIResponse[T]{
		Success: true,
		Data:    data,
		TraceID: traceID,
	})
}

func errorResponse(c *gin.Context, err error, traceID string) {
	status := http.StatusInternalServerError
	code := "internal_error"
	message := "Internal server error"

	if errors.Is(err, apierr.ErrPhotoNotFound) {
		status = http.StatusNotFound
		code = "photo_not_found"
		message = "Photo not found"
	} else if errors.Is(err, apierr.ErrBadCursor) {
		status = http.StatusBadRequest
		code = "bad_cursor"
		message = "Invalid cursor"
	} else if errors.Is(err, apierr.ErrInvalidParam) {
		status = http.StatusBadRequest
		code = "invalid_param"
		message = "Invalid parameter"
	} else if errors.Is(err, apierr.ErrUnauthorized) {
		status = http.StatusUnauthorized
		code = "unauthorized"
		message = "Unauthorized"
	}

	c.JSON(status, APIResponse[interface{}]{
		Success: false,
		Error: &ErrorBody{
			Code:    code,
			Message: message,
		},
		TraceID: traceID,
	})
}

func getTraceID(c *gin.Context) string {
	traceID := c.GetString("trace_id")
	if traceID == "" {
		traceID = c.Request.Header.Get("X-Trace-ID")
	}
	return traceID
}
