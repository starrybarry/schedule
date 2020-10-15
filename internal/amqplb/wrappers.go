package amqplb

import (
	"context"

	"github.com/streadway/amqp"
)

func AsyncConsume(ctx context.Context, consumer Consumer, handler func(msg amqp.Delivery) error) {
	var msg amqp.Delivery
	for msg = range consumer.Consume(ctx) {
		go func(msg amqp.Delivery) {
			if err := handler(msg); err != nil {
				msg.Nack(false, true)

				return
			}

			msg.Ack(false)
		}(msg)
	}
}

func SyncConsume(ctx context.Context, consumer Consumer, handler func(msg amqp.Delivery) error) error {
	errCh := make(chan error, 1)

	consumer.NotifyError(errCh)
	defer consumer.UnsubscribeError(errCh)

	var msg amqp.Delivery
	var open bool
	ch := consumer.Consume(ctx)

	for {
		select {
		case msg, open = <-ch:
			if open == false {
				return nil
			}

			if err := handler(msg); err != nil {
				msg.Nack(false, true)

				continue
			}

			if err := msg.Ack(false); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return err
		}
	}
}
