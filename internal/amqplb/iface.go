package amqplb

import (
	"context"

	"github.com/streadway/amqp"
)

type Consumer interface {
	DisconnectNotifier
	ReconnectNotifier
	ErrorNotifier
	CloseNotifier

	SubscribeOn(exchange string, topic string, opts *SubscribeOptions) error

	Consume(ctx context.Context) chan amqp.Delivery

	Close() error

	Error() error
}

type Publisher interface {
	DisconnectNotifier
	ReconnectNotifier

	Publish(exchange string, topic string, msg amqp.Publishing) error

	Setup(exchange string, options ...*PublisherOptions) error
}

type DisconnectNotifier interface {
	NotifyDisconnect(chan interface{}) chan interface{}
	UnsubscribeDisconnect(chan interface{})
}

type ReconnectNotifier interface {
	NotifyReconnect(chan interface{}) chan interface{}
	UnsubscribeReconnect(chan interface{})
}

type CloseNotifier interface {
	NotifyClose(chan interface{}) chan interface{}
	UnsubscribeClose(chan interface{})
}

type ErrorNotifier interface {
	NotifyError(chan error) chan error
	UnsubscribeError(chan error)
}

type Client interface {
	DisconnectNotifier
	ReconnectNotifier
	CloseNotifier
	ErrorNotifier

	Connect() error
	Consumer(queue string) (Consumer, error)
	Publisher() Publisher
	Close() error
}
