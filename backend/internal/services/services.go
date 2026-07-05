package services

import (
	"github.com/arthurshafikov/scout-takehome/backend/internal/config"
	"github.com/arthurshafikov/scout-takehome/backend/internal/core/models"
	"github.com/arthurshafikov/scout-takehome/backend/internal/repository"
)

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

type PhotoService interface {
	GetPhoto(id string) (*models.Photo, error)
}

type Services struct {
	PhotoService
}

type Deps struct {
	Repository repository.Repository
	Logger     Logger
	Config     *config.Config
}

func NewServices(deps Deps) *Services {
	return &Services{
		PhotoService: NewPhotoService(deps.Repository, deps.Logger),
	}
}
