package workerpool

import (
	"context"

	"go.uber.org/zap"
)

const (
	DefaultNumberOfWorkers int = 1001
)

type WorkerPool struct {
	taskStream      <-chan interface{}
	log             *zap.Logger
	do              func(interface{}) error
	numberOfWorkers int
}

func (w *WorkerPool) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case task := <-w.taskStream:
			go func() {
				if err := w.do(task); err != nil {

				}
			}()
		}
	}
}
