package scheduler

import (
	"context"

	"github.com/starrybarry/schedule/pkg/scheduler"
)

type Storage interface {
	AddTask(ctx context.Context, task scheduler.Task) error
	DeleteTask(ctx context.Context, task scheduler.Task) error
}

func NewService(storage Storage) Service {
	return Service{storage: storage}
}

type Service struct {
	storage Storage
}

func (s *Service) AddTask(ctx context.Context, task scheduler.Task) error {
	return nil
}

func (s *Service) DeleteTask(ctx context.Context, task scheduler.Task) error {
	return nil
}
