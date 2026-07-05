package handler

import (
	"context"
	"fmt"
	"strconv"

	"github.com/arthurshafikov/boilerplate/internal/services"
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
	h.initTestRoutes(e)
}

func (h *Handler) parseIntegerFromParam(ctx *gin.Context, key string) (int, error) {
	param := ctx.Param(key)
	if param == "" {
		return 0, fmt.Errorf("the param %s is missing", key)
	}
	paramInt, err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}

	return paramInt, err
}
