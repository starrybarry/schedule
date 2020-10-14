package scheduler

import (
	"context"

	"github.com/starrybarry/schedule/pkg/scheduler"
	"go.uber.org/zap"
)

type TaskStorage struct {
	log *zap.Logger
}

func (ts *TaskStorage) AddTask(ctx context.Context, task scheduler.Task) error {
	return nil
}

func (ts *TaskStorage) DeleteTask(ctx context.Context, task scheduler.Task) error {
	return nil
}

func (ts *TaskStorage) GetTask(ctx context.Context) (scheduler.Task, error) {
	return scheduler.Task{}, nil
}
