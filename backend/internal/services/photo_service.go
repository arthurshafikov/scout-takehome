package services

import (
	"context"
	"fmt"

	"github.com/arthurshafikov/scout-takehome/backend/internal/core/models"
	"github.com/arthurshafikov/scout-takehome/backend/internal/repository"
)

type PhotoServiceImpl struct {
	repo   repository.PhotoRepository
	logger Logger
}

func NewPhotoService(repo *repository.Repository, logger Logger) PhotoService {
	return &PhotoServiceImpl{
		repo:   repo,
		logger: logger,
	}
}

func (s *PhotoServiceImpl) GetPhoto(id string) (*models.Photo, error) {
	ctx := context.Background()
	photo, err := s.repo.GetPhotoByID(ctx, id)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to get photo %s: %v", id, err))
		return nil, fmt.Errorf("get photo: %w", err)
	}

	return photo, nil
}
