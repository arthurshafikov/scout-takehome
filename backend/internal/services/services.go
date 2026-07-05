package services

import (
	"context"

	"github.com/arthurshafikov/scout-takehome/backend/internal/config"
	"github.com/arthurshafikov/scout-takehome/backend/internal/core/constants/events"
	"github.com/arthurshafikov/scout-takehome/backend/internal/core/models"
	"github.com/arthurshafikov/scout-takehome/backend/internal/repository"
)

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

type EventsHandler interface {
	Dispatch(eventName events.Event, params ...any)
}

type Test interface {
	Create(ctx context.Context, models models.TestModel) error
}

type Services struct {
	Test
}

type Deps struct {
	Repository *repository.Repository
	Logger
	Config        *config.Config
	EventsHandler EventsHandler
}

func NewServices(deps Deps) *Services {

	return &Services{
		Test: NewTestService(),
	}
}
