package amqplb

import (
	"sync"
	"time"

	"github.com/streadway/amqp"
)

const maxReconnectAttempts = 5

func NewClient(uri string, heartbeat time.Duration) *client {
	return &client{
		disconnectNotifier: &disconnectNotifier{channel: "client", disconnectListeners: make(map[chan interface{}]interface{}, 0)},
		reconnectNotifier:  &reconnectNotifier{channel: "client", reconnectListeners: make(map[chan interface{}]interface{}, 0)},
		closeNotifier:      &closeNotifier{channel: "client", closeListeners: make(map[chan interface{}]interface{}, 0)},
		errorNotifier:      &errorNotifier{channel: "client", errorListeners: make(map[chan error]interface{}, 0)},

		uri:       uri,
		heartbeat: heartbeat,
	}
}

type client struct {
	*disconnectNotifier
	*reconnectNotifier
	*closeNotifier
	*errorNotifier

	uri       string
	heartbeat time.Duration

	conn           *amqp.Connection
	isConnected    bool
	reconnectingMu sync.RWMutex

	channelConfigurers []func(c *amqp.Channel)
}

func (c *client) ConfigureChannel(conf func(c *amqp.Channel)) {
	c.channelConfigurers = append(c.channelConfigurers, conf)
}

func (c *client) ready() chan interface{} {
	c.reconnectingMu.RLock()
	defer c.reconnectingMu.RUnlock()

	ready := make(chan interface{}, 1)

	if c.isConnected == true && c.conn != nil && c.conn.IsClosed() == false {
		ready <- true
		close(ready)
		return ready
	}

	c.NotifyReconnect(ready)

	return ready
}

func (c *client) Connect() error {
	c.reconnectingMu.Lock()
	defer c.reconnectingMu.Unlock()

	return c.connect()
}

func (c *client) listenConnectionError(conn *amqp.Connection) {
	err := <-conn.NotifyClose(make(chan *amqp.Error, 1))
	if err == nil {
		return
	}

	c.reconnectingMu.Lock()
	defer c.reconnectingMu.Unlock()

	c.notifyDisconnect()

	c.isConnected = false

	reconnectAttempts := 0

	for {
		if reconnectAttempts >= maxReconnectAttempts {
			c.isConnected = false

			return
		}

		time.Sleep(time.Second * time.Duration(reconnectAttempts))

		if err := c.reconnect(); err == nil {
			break
		}

		reconnectAttempts++
	}
}

func (c *client) connect() error {
	conn, err := amqp.DialConfig(c.uri, amqp.Config{
		Heartbeat: c.heartbeat,
	})
	if err != nil {
		return err
	}

	go c.listenConnectionError(conn)

	c.conn = conn
	c.isConnected = true

	return nil
}

func (c *client) disconnect() error {
	if c.conn == nil || c.conn.IsClosed() {
		return nil
	}

	if err := c.conn.Close(); err != nil {
		return err
	}

	c.isConnected = false

	c.notifyDisconnect()

	return nil
}

func (c *client) reconnect() error {
	if err := c.disconnect(); err != nil {
		return err
	}

	if err := c.connect(); err != nil {
		return err
	}

	c.notifyReconnect()

	return nil
}

func (c *client) Consumer(queue string) (Consumer, error) {
	<-c.ready()

	return newConsumer(c, queue, c.channelConfigurers), nil
}

func (c *client) Publisher() Publisher {
	<-c.ready()

	return newPublisher(c)
}

func (c *client) Close() error {
	if err := c.disconnect(); err != nil {
		return err
	}

	c.notifyClose()

	return nil
}
