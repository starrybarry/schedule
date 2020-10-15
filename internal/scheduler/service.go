package scheduler

import (
	"context"
	"fmt"

	"github.com/starrybarry/schedule/pkg/scheduler"
)

type storage interface {
	AddTask(ctx context.Context, task scheduler.Task) error
	DeleteTask(ctx context.Context, task scheduler.Task) error
}

func NewService(storage storage) *Service {
	return &Service{storage: storage}
}

type Service struct {
	storage storage
}

func (s *Service) AddTask(ctx context.Context, task scheduler.Task) error {
	if err := s.storage.AddTask(ctx, task); err != nil {
		return fmt.Errorf("storage add task, error: %w", err)
	}

	return nil
}

func (s *Service) DeleteTask(ctx context.Context, task scheduler.Task) error {
	if err := s.storage.DeleteTask(ctx, task); err != nil {
		return fmt.Errorf("storage delete task, error: %w", err)
	}

	return nil
}
