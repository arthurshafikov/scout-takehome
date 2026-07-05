package services

import (
	"context"

	"github.com/arthurshafikov/boilerplate/internal/core/models"
)

type TestService struct {
}

func NewTestService() *TestService {
	return &TestService{}
}

func (s *TestService) Create(ctx context.Context, models models.TestModel) error {
	return nil
}
