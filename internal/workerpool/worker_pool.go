package workerpool

import (
	"context"
	"fmt"

	"github.com/starrybarry/schedule/pkg/scheduler"

	jsoniter "github.com/json-iterator/go"
	"github.com/streadway/amqp"

	"go.uber.org/zap"
)

type consumer interface {
	GetDeliveryChannel(ctx context.Context) (<-chan amqp.Delivery, <-chan error)
}

type WorkerPool struct {
	consumer consumer
	storage  storage
	log      *zap.Logger
}

func NewWorkerPool(consumer consumer, storage storage, log *zap.Logger) *WorkerPool {
	return &WorkerPool{consumer: consumer, storage: storage, log: log}
}

func (w *WorkerPool) Start(ctx context.Context) error {
	ch, chErr := w.consumer.GetDeliveryChannel(ctx)
	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-chErr:
			if err != nil {
				return fmt.Errorf("delivery channel, error: %w", err)
			}

		case msg, ok := <-ch:
			if !ok {
				w.log.Warn("delivery channel closed!")
				return nil
			}

			if err := w.Do(ctx, msg); err != nil {
				w.log.Error("worker do", zap.Error(err))
				return err
			}

			w.log.Info("Делл Сделанно!")

		}
	}
}

func (w *WorkerPool) Do(ctx context.Context, msg amqp.Delivery) error {
	var task scheduler.Task

	if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(msg.Body, &task); err != nil {
		return fmt.Errorf("unmarshal delivery body, error: %w", err)
	}

	w.log.Info("Опять работа!", zap.Any("task", task))

	task.InProgress = true

	if err := msg.Ack(false); err != nil {
		return fmt.Errorf("msg ack, error: %w", err)
	}

	task.IsFinished = true

	if err := w.storage.UpdateTask(ctx, task); err != nil {
		return fmt.Errorf("update task, error: %w", err)
	}

	return nil
}
