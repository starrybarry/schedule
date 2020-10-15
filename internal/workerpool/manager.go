package workerpool

import (
	"context"
	"fmt"
	"time"

	"github.com/starrybarry/schedule/pkg/scheduler"
	"go.uber.org/zap"
)

type storage interface {
	GetTasks(ctx context.Context) ([]scheduler.Task, error)
}

type publisher interface {
	PublishTask(task scheduler.Task) error
}

type Manager struct {
	pollingTimer time.Duration
	storage      storage
	publisher    publisher
	log          *zap.Logger
}

func NewManager(pollingTimer time.Duration, storage storage, publisher publisher, log *zap.Logger) *Manager {
	return &Manager{pollingTimer: pollingTimer, storage: storage, publisher: publisher, log: log}
}

func (m *Manager) PollingTaskAndPublish(ctx context.Context) error {
	t := time.NewTimer(0)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-t.C:
			t.Reset(m.pollingTimer)

			tasks, err := m.storage.GetTasks(ctx)
			if err != nil {
				return fmt.Errorf("get task, error: %w", err)
			}

			for _, task := range tasks {
				if err := m.publisher.PublishTask(task); err != nil {
					return fmt.Errorf("publish task, error: %w", err)
				}
			}

			m.log.Info("tasks published!", zap.Any("tasks", tasks))
		}
	}
}
