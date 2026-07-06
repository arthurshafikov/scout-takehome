package app

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	_ "modernc.org/sqlite"

	configPkg "github.com/arthurshafikov/scout-takehome/backend/internal/config"
	"github.com/arthurshafikov/scout-takehome/backend/internal/repository"
	"github.com/arthurshafikov/scout-takehome/backend/internal/services"
	"github.com/arthurshafikov/scout-takehome/backend/internal/services/metrics"
	"github.com/arthurshafikov/scout-takehome/backend/internal/services/thumbnail"
	"github.com/arthurshafikov/scout-takehome/backend/internal/transport/http"
	"github.com/arthurshafikov/scout-takehome/backend/internal/transport/http/handler"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

var (
	envFolderPath    string
	configFolderPath string
)

func init() {
	flag.StringVar(&envFolderPath, "env", "", "Path to .env file folder")
	flag.StringVar(&configFolderPath, "cfgFolder", "", "Path to configs folder")
}

func Run() {
	flag.Parse()

	time.Local = time.UTC

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	g, ctx := errgroup.WithContext(ctx)
	defer cancel()

	config := configPkg.NewConfig(envFolderPath, configFolderPath)
	logger := logrus.New()
	logger.SetReportCaller(true)
	logger.SetOutput(os.Stdout)

	if config.Debug {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	logger.Info("Starting Scout backend...")

	defer func() {
		if err := recover(); err != nil {
			logger.Panic(err)
		}
	}()

	// Connect to SQLite
	db, err := sql.Open("sqlite", config.DBPath)
	if err != nil {
		logger.Fatalf("Failed to open SQLite database: %v", err)

		return
	}

	defer func() {
		if err := db.Close(); err != nil {
			logger.Warnf("Failed to close database: %v", err)
		}
	}()

	if err := db.PingContext(ctx); err != nil {
		logger.Fatalf("Failed to connect to SQLite database: %v", err)

		return
	}

	logger.Info("Connected to SQLite database")

	// Initialize metrics
	m := metrics.NewMetrics()

	repo := repository.NewRepository(db)

	svc, err := services.NewServices(services.Deps{
		Repository: repo,
		Logger:     logger,
		Config:     config,
	})
	if err != nil {
		logger.Fatalf("Failed to initialize services: %v", err)

		return
	}

	// Initialize thumbnail generator
	minioClient, err := svc.GetMinIOClient()
	if err != nil {
		logger.Fatalf("Failed to get MinIO client for thumbnails: %v", err)

		return
	}

	// Dataset path is the directory containing the database
	datasetPath := filepath.Dir(config.DBPath)
	thumbnailGen := thumbnail.NewThumbnailGenerator(minioClient, config.Bucket, datasetPath, m)

	h := handler.NewHandler(ctx, svc, thumbnailGen, m, config.APIKey)

	server := http.NewServer(logger, h, config.Debug, m)

	g.Go(func() error {
		defer func() {
			if err := recover(); err != nil {
				logger.Panic(err)
			}
		}()

		return server.Serve(ctx, g, config.Port)
	})

	if err := g.Wait(); err != nil {
		logger.Error(err)
	}
}
