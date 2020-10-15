package amqplb

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/streadway/amqp"
)

func newConsumer(c *client, queue string, channelConfigurers []func(c *amqp.Channel)) *consumer {
	channel := fmt.Sprintf("consumer-%s", queue)

	con := &consumer{
		disconnectNotifier: &disconnectNotifier{channel: channel, disconnectListeners: make(map[chan interface{}]interface{}, 0)},
		reconnectNotifier:  &reconnectNotifier{channel: channel, reconnectListeners: make(map[chan interface{}]interface{}, 0)},
		closeNotifier:      &closeNotifier{channel: channel, closeListeners: make(map[chan interface{}]interface{}, 0)},
		errorNotifier:      &errorNotifier{channel: channel, errorListeners: make(map[chan error]interface{}, 0)},

		client:             c,
		channelConfigurers: channelConfigurers,
		queue:              queue,
		close:              make(chan interface{}),
	}

	con.watchClientState()

	return con
}

type consumer struct {
	*disconnectNotifier
	*reconnectNotifier
	*closeNotifier
	*errorNotifier
	sync.Mutex

	client             *client
	ch                 *amqp.Channel
	channelConfigurers []func(c *amqp.Channel)

	queue    string
	close    chan interface{}
	isClosed bool

	errMu sync.RWMutex
	error error
}

func (c *consumer) watchClientState() {
	dcCh := make(chan interface{}, 1)
	rcCh := make(chan interface{}, 1)
	closeCh := make(chan interface{}, 1)

	c.client.NotifyDisconnect(dcCh)
	c.client.NotifyReconnect(rcCh)
	c.client.NotifyClose(closeCh)

	go func() {
		defer c.client.UnsubscribeDisconnect(dcCh)
		defer c.client.UnsubscribeReconnect(rcCh)
		defer c.client.UnsubscribeClose(closeCh)

		for {
			select {
			case <-dcCh:
				c.Lock()
				c.ch = nil
				c.Unlock()
				c.notifyDisconnect()
			case <-rcCh:
				c.notifyReconnect()
			case <-c.close:
				c.Lock()

				if c.isClosed {
					c.Unlock()
					return
				}

				c.isClosed = true
				c.Unlock()

				return
			case <-closeCh:
				c.Lock()
				if c.isClosed {
					c.Unlock()
					continue
				}

				c.isClosed = true
				close(c.close)
				c.Unlock()
				return
			}
		}
	}()
}

func (c *consumer) channel() (*amqp.Channel, error) {
	c.Lock()
	defer c.Unlock()

	if c.isClosed {
		return nil, errors.New("consumer is closed")
	}

	<-c.client.ready()

	if c.ch != nil {
		return c.ch, nil
	}

	ch, err := c.client.conn.Channel()

	if err != nil {
		return nil, err
	}

	_ = ch.Qos(3, 0, false)

	for _, configure := range c.channelConfigurers {
		configure(ch)
	}

	chClosed := make(chan *amqp.Error, 1)
	started := make(chan struct{})

	go func() {
		started <- struct{}{}
		close(started)

		if err := <-chClosed; err != nil {
			c.Lock()
			defer c.Unlock()

			c.ch = nil
			c.setError(err)
		}

		return
	}()

	<-started

	ch.NotifyClose(chClosed)

	c.ch = ch
	c.setError(nil) // connected, so we can drop error
	return c.ch, nil
}

func (c *consumer) SubscribeOn(exchange string, topic string, opts *SubscribeOptions) error {
	if opts == nil {
		opts = NewDefaultSubscribeOptions()
	}

	ch, err := c.channel()
	if err != nil {
		c.setError(err)
		c.notifyError(err)
		return err
	}

	queue, err := ch.QueueDeclare(c.queue, opts.Durable, opts.AutoDelete, opts.Exclusive, false, opts.Args)
	if err != nil {
		return err
	}

	return ch.QueueBind(queue.Name, topic, exchange, false, nil)
}

func (c *consumer) Close() error {
	c.Lock()
	defer c.Unlock()

	if c.isClosed {
		return nil
	}

	if c.ch == nil {
		c.isClosed = true
		close(c.close)
		c.notifyClose()
		return nil
	}

	if err := c.ch.Close(); err != nil {
		if err != amqp.ErrClosed || c.isClosed == true {
			return err
		}
	}

	c.isClosed = true
	close(c.close)
	c.notifyClose()

	return nil
}

func (c *consumer) Consume(ctx context.Context) chan amqp.Delivery {
	ch := make(chan amqp.Delivery)

	if c.isClosed {
		c.notifyError(errors.New("trying to consume from closed consumer"))
		close(ch)
		return ch
	}

	<-c.client.ready()

	go func() {
		dCh := make(chan interface{}, 1)
		c.NotifyDisconnect(dCh)

		rcCh := make(chan interface{}, 1)
		c.NotifyReconnect(rcCh)

		closeCh := make(chan interface{}, 1)
		c.client.NotifyClose(closeCh)

		defer close(ch)
		defer c.client.UnsubscribeClose(closeCh)
		defer c.notifyClose()

		defer c.UnsubscribeDisconnect(dCh)
		defer c.UnsubscribeReconnect(rcCh)

		defer func() {
			if err := c.Error(); err != nil {
				c.notifyError(err)
			}
		}()

		channel, err := c.channel()
		if err != nil {
			c.setError(err)
			return
		}

		deliveryCh, err := channel.Consume(c.queue, "", false, false, false, false, nil)
		if err != nil {
			c.setError(err)
			return
		}

		var msg amqp.Delivery
		var open bool

		for {
			select {
			case <-dCh:
				select {
				case <-rcCh:
					channel, err = c.channel()
					if err != nil {
						c.setError(err)
						return
					}

					deliveryCh, err = channel.Consume(c.queue, "", false, false, true, false, nil)
					if err != nil {
						c.setError(err)
						return
					}
				case <-c.close:
					return
				case <-ctx.Done():
					return
				}
			case msg, open = <-deliveryCh:
				if open == false {
					deliveryCh = nil
					continue
				}

				ch <- msg
			case <-c.close:
				return
			case <-ctx.Done():
				return
			}
		}
	}()

	return ch
}

func (c *consumer) setError(err error) {
	c.errMu.Lock()
	c.error = err
	c.errMu.Unlock()
}

func (c *consumer) Error() error {
	c.errMu.RLock()
	defer c.errMu.RUnlock()
	return c.error
}
