package app

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "modernc.org/sqlite"

	configPkg "github.com/arthurshafikov/scout-takehome/backend/internal/config"
	"github.com/arthurshafikov/scout-takehome/backend/internal/repository"
	"github.com/arthurshafikov/scout-takehome/backend/internal/services"
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

	if config.App.Debug {
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
	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?mode=ro&_journal_mode=WAL", config.SQLiteConfig.DBPath))
	if err != nil {
		logger.Fatalf("Failed to open SQLite database: %v", err)

		return
	}

	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		logger.Fatalf("Failed to connect to SQLite database: %v", err)

		return
	}

	logger.Info("Connected to SQLite database")

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
	minioClient, err := svc.StorageService.GetMinIOClient()
	if err != nil {
		logger.Fatalf("Failed to get MinIO client for thumbnails: %v", err)

		return
	}

	thumbnailGen := thumbnail.NewThumbnailGenerator(minioClient, config.MinIOConfig.Bucket)

	h := handler.NewHandler(ctx, svc, thumbnailGen)

	server := http.NewServer(logger, h, config.App.Debug)

	g.Go(func() error {
		defer func() {
			if err := recover(); err != nil {
				logger.Panic(err)
			}
		}()

		return server.Serve(ctx, g, config.App.Port)
	})

	if err := g.Wait(); err != nil {
		logger.Error(err)
	}
}
