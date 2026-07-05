package jobs

import (
	"context"

	"github.com/arthurshafikov/scout-takehome/backend/internal/services"
)

type TestJob struct {
	ctx context.Context

	services *services.Services
}

func NewTestJob(
	ctx context.Context,
	services *services.Services,
) *TestJob {
	return &TestJob{
		ctx:      ctx,
		services: services,
	}
}

func (j *TestJob) Handle() {
	// if err := j.services.Service.TestJob(j.ctx); err != nil {
	// 	logrus.Error(fmt.Errorf("TestJob err: %w", err))
	// }
}
