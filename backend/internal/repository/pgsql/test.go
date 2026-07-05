package pgsql

import (
	"context"

	"github.com/arthurshafikov/boilerplate/internal/core/models"
)

type Test struct {
	BaseRepo
}

func NewTestRepository(baseRepo BaseRepo) *Test {
	return &Test{
		BaseRepo: baseRepo,
	}
}

func (r *Test) Create(ctx context.Context, model models.TestModel) error {
	return r.getDBInstance(ctx).
		Table("models").
		Create(&model).Error
}
