package services

import (
	"context"
	"fmt"

	"github.com/arthurshafikov/scout-takehome/backend/internal/core/models"
	"github.com/arthurshafikov/scout-takehome/backend/internal/core/types"
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

func (s *PhotoServiceImpl) GetPhoto(ctx context.Context, id string) (*models.Photo, error) {
	photo, err := s.repo.GetPhotoByID(ctx, id)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to get photo %s: %v", id, err))
		return nil, fmt.Errorf("get photo: %w", err)
	}

	return photo, nil
}

func (s *PhotoServiceImpl) ListPhotos(
	ctx context.Context,
	params types.ListPhotosParams,
) ([]models.Photo, string, error) {
	photos, nextToken, err := s.repo.ListPhotos(ctx, params)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to list photos: %v", err))
		return nil, "", fmt.Errorf("list photos: %w", err)
	}

	return photos, nextToken, nil
}
