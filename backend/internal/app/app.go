package app

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	configPkg "github.com/arthurshafikov/scout-takehome/backend/internal/config"
	"github.com/arthurshafikov/scout-takehome/backend/internal/events"
	"github.com/arthurshafikov/scout-takehome/backend/internal/repository"
	"github.com/arthurshafikov/scout-takehome/backend/internal/repository/pgsql"
	"github.com/arthurshafikov/scout-takehome/backend/internal/scheduler"
	"github.com/arthurshafikov/scout-takehome/backend/internal/services"
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
	logrus.SetLevel(logrus.DebugLevel)
	logger.SetReportCaller(true)
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
	logger.Info("Starting the app...")

	defer func() {
		if err := recover(); err != nil {
			logger.Panic(err)
		}
	}()

	eventsHandler := events.NewHandler(ctx, g, logger, config.App.MaxEventHandlerGoroutines)

	db := pgsql.ConnectToDatabase(ctx, &config.DBConfig, config.App.Debug)
	repository := repository.NewRepository(db)

	services := services.NewServices(services.Deps{
		Repository:    repository,
		Logger:        logger,
		Config:        config,
		EventsHandler: eventsHandler,
	})

	eventsHandler.InitEvents(&events.Deps{
		Services: services,
		Config:   config,
	})

	scheduler, err := scheduler.NewScheduler(ctx, scheduler.Deps{
		Logger:   logger,
		Services: services,
		Config:   config,
	})
	if err != nil {
		logger.Error(err)

		return
	}
	g.Go(func() error {
		defer func() {
			if err := recover(); err != nil {
				logger.Panic(err)
			}
		}()

		return scheduler.Start()
	})

	handler := handler.NewHandler(ctx, services)

	server := http.NewServer(logger, handler, config.App.Debug)

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
