package repository

import (
	"context"
	"database/sql"

	"github.com/arthurshafikov/scout-takehome/backend/internal/core/models"
	"github.com/arthurshafikov/scout-takehome/backend/internal/core/types"
	"github.com/arthurshafikov/scout-takehome/backend/internal/repository/sqlite"
)

type PhotoRepository interface {
	GetPhotoByID(ctx context.Context, id string) (*models.Photo, error)
	ListPhotos(ctx context.Context, params types.ListPhotosParams) ([]models.Photo, string, error)
}

type ListPhotosParams struct {
	Cursor        string
	Limit         int
	ClassID       string
	MinConfidence float64
}

type Repository struct {
	PhotoRepository
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		PhotoRepository: sqlite.NewSQLitePhotoRepository(db),
		db:              db,
	}
}
