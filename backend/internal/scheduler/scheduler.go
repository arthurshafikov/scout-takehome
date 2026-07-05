package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/arthurshafikov/scout-takehome/backend/internal/config"
	"github.com/arthurshafikov/scout-takehome/backend/internal/scheduler/jobs"
	"github.com/arthurshafikov/scout-takehome/backend/internal/services"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type Deps struct {
	Logger   Logger
	Services *services.Services
	Config   *config.Config
}

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Panic(args ...interface{})
}

type Scheduler struct {
	ctx context.Context

	scheduler gocron.Scheduler

	logger   Logger
	services *services.Services
	config   *config.Config
}

func NewScheduler(ctx context.Context, deps Deps) (*Scheduler, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	return &Scheduler{
		ctx:       ctx,
		scheduler: scheduler,

		logger:   deps.Logger,
		services: deps.Services,
		config:   deps.Config,
	}, nil
}

func (s *Scheduler) Start() error {
	s.initJobs()

	s.scheduler.Start()
	s.logger.Info("Scheduler started...")

	<-s.ctx.Done()
	s.logger.Info("Shutting down the scheduler...")
	err := s.scheduler.Shutdown()
	s.logger.Info("Scheduler has been stopped...")

	return err
}

func (s *Scheduler) initJobs() {
	testJob := jobs.NewTestJob(
		s.ctx,
		s.services,
	)
	if _, err := s.scheduler.NewJob(
		gocron.DurationJob(time.Hour*24*365),
		gocron.NewTask(
			func() {
				testJob.Handle()
			},
		),
		gocron.WithEventListeners(
			gocron.AfterJobRunsWithPanic(func(_ uuid.UUID, _ string, recoverData any) {
				s.logger.Error(fmt.Errorf("testJob panic: %s", recoverData))
			}),
			gocron.AfterJobRunsWithError(func(_ uuid.UUID, _ string, err error) {
				s.logger.Error(fmt.Errorf("testJob error: %w", err))
			}),
		),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	); err != nil {
		s.logger.Error(err)
	}
}
