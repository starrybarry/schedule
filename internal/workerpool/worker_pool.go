package workerpool

import (
	"context"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/starrybarry/schedule/pkg/scheduler"

	"github.com/streadway/amqp"

	"go.uber.org/zap"
)

type consumer interface {
	GetDeliveryChannel(ctx context.Context) (<-chan amqp.Delivery, error)
}

type WorkerPool struct {
	consumer consumer
	log      *zap.Logger
}

func NewWorkerPool(consumer consumer, log *zap.Logger) *WorkerPool {
	return &WorkerPool{consumer: consumer, log: log}
}

func (w *WorkerPool) Start(ctx context.Context) error {
	ch, err := w.consumer.GetDeliveryChannel(ctx)
	if err != nil {
		return fmt.Errorf("get delivery channel, error: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-ch:
			go func() {
				if err := w.Do(ctx, msg); err != nil {
					w.log.Error("worker do", zap.Error(err))
					return
				}

				w.log.Info("Делл Сделанно!")
			}()
		}
	}
}

func (w *WorkerPool) Do(_ context.Context, msg amqp.Delivery) error {
	var task scheduler.Task

	if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(msg.Body, &task); err != nil {
		return fmt.Errorf("unmarshal delivery body, error: %w", err)
	}

	w.log.Info("Опять работа!", zap.Any("task", task))

	if err := msg.Ack(false); err != nil {
		return fmt.Errorf("msg ack, error: %w", err)
	}

	return nil
}
