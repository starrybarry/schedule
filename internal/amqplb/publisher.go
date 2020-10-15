package amqplb

import (
	"fmt"
	"sync"

	"github.com/streadway/amqp"
)

func newPublisher(c *client) *publisher {
	channel := fmt.Sprintf("publisher")

	p := &publisher{
		disconnectNotifier: &disconnectNotifier{channel: channel, disconnectListeners: make(map[chan interface{}]interface{}, 0)},
		reconnectNotifier:  &reconnectNotifier{channel: channel, reconnectListeners: make(map[chan interface{}]interface{}, 0)},
		closeNotifier:      &closeNotifier{channel: channel, closeListeners: make(map[chan interface{}]interface{}, 0)},

		client: c,
	}

	p.watchClientState()

	return p
}

type publisher struct {
	*disconnectNotifier
	*reconnectNotifier
	*closeNotifier

	client  *client
	options *PublisherOptions
	ch      *amqp.Channel
	sync.Mutex
}

func (p *publisher) watchClientState() {
	dcCh := make(chan interface{}, 1)
	rcCh := make(chan interface{}, 1)
	closeCh := make(chan interface{}, 1)

	p.client.NotifyDisconnect(dcCh)
	p.client.NotifyReconnect(rcCh)
	p.client.NotifyClose(closeCh)

	go func() {
		for {
			select {
			case <-dcCh:
				p.Lock()
				p.ch = nil
				p.Unlock()
				p.notifyDisconnect()
			case <-rcCh:
				p.notifyReconnect()
			case <-closeCh:
				p.notifyClose()

				p.client.UnsubscribeDisconnect(dcCh)
				p.client.UnsubscribeReconnect(rcCh)
				p.client.UnsubscribeClose(closeCh)
				return
			}
		}
	}()
}

func (p *publisher) channel() (*amqp.Channel, error) {
	p.Lock()
	defer p.Unlock()
	<-p.client.ready()

	if p.ch != nil {
		return p.ch, nil
	}

	ch, err := p.client.conn.Channel()

	if err == nil {
		chClosed := make(chan *amqp.Error, 1)

		go func() {
			if err := <-chClosed; err != nil {
				p.Lock()
				p.ch = nil
				p.Unlock()
			}
		}()

		ch.NotifyClose(chClosed)

		p.ch = ch

		return p.ch, nil
	}

	return ch, err
}

func (p *publisher) Setup(exchange string, options ...*PublisherOptions) error {
	ch, err := p.channel()
	if err != nil {
		return err
	}

	var o *PublisherOptions
	if len(options) == 0 {
		o = NewDefaultPublisherOptions()
	} else {
		o = options[0]
	}

	return ch.ExchangeDeclare(exchange, string(o.Type), o.Durable, o.AutoDelete, false, false, nil)
}

func (p *publisher) Publish(exchange string, topic string, msg amqp.Publishing) error {
	ch, err := p.channel()
	if err != nil {
		return err
	}

	return ch.Publish(exchange, topic, false, false, msg)
}
