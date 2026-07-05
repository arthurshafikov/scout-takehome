package services

import (
	"context"

	"github.com/arthurshafikov/boilerplate/internal/config"
	"github.com/arthurshafikov/boilerplate/internal/core/constants/events"
	"github.com/arthurshafikov/boilerplate/internal/core/models"
	"github.com/arthurshafikov/boilerplate/internal/repository"
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
