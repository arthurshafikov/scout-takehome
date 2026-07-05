package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/arthurshafikov/scout-takehome/backend/internal/services/metrics"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type Handler interface {
	Init(*gin.Engine)
}

type Logger interface {
	Error(msg ...any)
	Info(msg ...any)
}

type Server struct {
	logger  Logger
	httpSrv *http.Server
	handler Handler
	Engine  *gin.Engine
	metrics *metrics.Metrics
}

func NewServer(logger Logger, handler Handler, debug bool, m *metrics.Metrics) *Server {
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	return &Server{
		logger:  logger,
		handler: handler,
		Engine:  gin.Default(),
		metrics: m,
	}
}

func (s *Server) Serve(ctx context.Context, g *errgroup.Group, port string) error {
	// Add middlewares
	s.Engine.Use(corsMiddleware())
	if s.metrics != nil {
		s.Engine.Use(s.metricsMiddleware())
	}

	s.handler.Init(s.Engine)

	s.httpSrv = &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: s.Engine,
	}

	s.logger.Info(fmt.Sprintf("Starting an HTTP server on port %s...", port))

	g.Go(func() error {
		<-ctx.Done()
		return s.shutdown()
	})

	if err := s.httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("could not start listener: %w", err)
	}

	return nil
}

func (s *Server) shutdown() error {
	s.logger.Info("Shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.httpSrv.Shutdown(ctx)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// metricsMiddleware records HTTP request metrics
func (s *Server) metricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		endpoint := c.Request.URL.Path
		method := c.Request.Method
		status := strconv.Itoa(c.Writer.Status())
		duration := time.Since(start).Seconds()

		// Record HTTP request metrics
		s.metrics.RecordHTTPRequest(endpoint, method, status, duration)

		// Record errors separately
		if c.Writer.Status() >= 400 {
			s.metrics.RecordHTTPError(endpoint, status)
		}
	}
}
