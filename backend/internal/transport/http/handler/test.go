package handler

import (
	"net/http"

	"github.com/arthurshafikov/boilerplate/internal/core/types"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initTestRoutes(e *gin.Engine) {
	test := e.Group("/test")
	{
		test.GET("", h.test)
	}
}

func (h *Handler) test(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, types.ServerResponse{
		Success: true,
		Data: struct {
			Test bool
		}{
			Test: true,
		},
		Error: nil,
		Meta:  nil,
	})
}
