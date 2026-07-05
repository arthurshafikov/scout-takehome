package repository

import (
	"context"

	"github.com/arthurshafikov/scout-takehome/backend/internal/core/models"
	"github.com/arthurshafikov/scout-takehome/backend/internal/core/types"
	"github.com/arthurshafikov/scout-takehome/backend/internal/repository/pgsql"
	"gorm.io/gorm"
)

type BaseRepo interface {
	WrapInTransaction(ctx *types.Context, txFunc types.TXFunction) error
}

type Test interface {
	Create(ctx context.Context, model models.TestModel) error
}

type Repository struct {
	Test
}

func NewRepository(db *gorm.DB) *Repository {
	baseRepo := pgsql.NewBaseRepo(db)

	return &Repository{
		Test: pgsql.NewTestRepository(*baseRepo),
	}
}
