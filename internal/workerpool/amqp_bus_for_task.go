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
	publisher    amqplb.Publisher
	consumer     amqplb.Consumer
	log          *zap.Logger
}

func NewBus(exchangeName string, publisher amqplb.Publisher, consumer amqplb.Consumer, log *zap.Logger) *Bus {
	return &Bus{exchangeName: exchangeName, publisher: publisher, consumer: consumer, log: log}
}

func (b *Bus) PublishTask(task scheduler.Task) error {
	body, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(task)
	if err != nil {
		return fmt.Errorf("marshal task:error: %w", err)
	}

	return b.publisher.Publish(b.exchangeName, "tasks",
		amqp.Publishing{
			ContentType: "application/json",
			Timestamp:   time.Now(),
			Body:        body,
		})
}

func (b *Bus) GetDeliveryChannel(ctx context.Context) (<-chan amqp.Delivery, <-chan error) {
	chErr := b.consumer.NotifyError(make(chan error))
	ch := b.consumer.Consume(ctx)
	return ch, chErr
}
