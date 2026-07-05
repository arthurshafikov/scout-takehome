package events

import (
	"context"

	"github.com/arthurshafikov/scout-takehome/backend/internal/config"
	"github.com/arthurshafikov/scout-takehome/backend/internal/core/constants/events"
	"github.com/arthurshafikov/scout-takehome/backend/internal/events/listeners"
	"github.com/arthurshafikov/scout-takehome/backend/internal/services"
	"golang.org/x/sync/errgroup"
)

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

type RegisteredEventsMap map[string][]listeners.Listener

type Handler struct {
	ctx    context.Context
	group  *errgroup.Group
	logger Logger

	registeredEvents *RegisteredEventsMap
	sem              chan struct{}
}

type Deps struct {
	Services *services.Services
	Config   *config.Config
}

func NewHandler(
	ctx context.Context,
	group *errgroup.Group,
	logger Logger,
	maxGoroutines int,
) *Handler {
	return &Handler{
		ctx:              ctx,
		group:            group,
		logger:           logger,
		registeredEvents: &RegisteredEventsMap{},
		sem:              make(chan struct{}, maxGoroutines),
	}
}

func (h *Handler) InitEvents(deps *Deps) {
	// h.subscribe(events.TestEvent, listeners.NewListener(
	// 	deps.Services.Service,
	// ))
}

func (h *Handler) Dispatch(eventName events.Event, params ...any) {
	listeners, ok := (*h.registeredEvents)[eventName.ToString()]
	if !ok {
		h.logger.Error("no listeners registered for eventName " + eventName.ToString())
		return
	}

	h.sem <- struct{}{} // blocks in case of an exceeded limit
	h.group.Go(func() error {
		defer func() { <-h.sem }() // frees up a slot
		for _, listener := range listeners {
			if err := listener.Handle(h.ctx, params...); err != nil {
				h.logger.Error(err)
			}
		}
		return nil
	})
}

//nolint:unused
func (h *Handler) subscribe(eventName events.Event, listener listeners.Listener) {
	(*h.registeredEvents)[eventName.ToString()] = append(
		(*h.registeredEvents)[eventName.ToString()],
		listener,
	)
}
