package workerpool

import (
	"context"
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/streadway/amqp"

	"github.com/starrybarry/schedule/internal/amqplb"
	"github.com/starrybarry/schedule/pkg/scheduler"
	"go.uber.org/zap"
)

type Bus struct {
	exchangeName string
	client       amqplb.Client
	log          *zap.Logger
}

func NewBus(exchangeName string, client amqplb.Client, log *zap.Logger) *Bus {
	return &Bus{exchangeName: exchangeName, client: client, log: log}
}

func (b *Bus) PublishTask(task scheduler.Task) error {
	body, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(task)
	if err != nil {
		return fmt.Errorf("marshal task:error: %w", err)
	}

	return b.client.Publisher().Publish(b.exchangeName, "tasks", amqp.Publishing{
		ContentType: "application/json",
		Timestamp:   time.Now(),
		Body:        body,
	})
}

func (b *Bus) GetDeliveryChannel(ctx context.Context) (<-chan amqp.Delivery, error) {
	cons, err := b.client.Consumer("task")
	if err != nil {
		return nil, fmt.Errorf("get consumer, error: %w", err)
	}

	return cons.Consume(ctx), nil
}
