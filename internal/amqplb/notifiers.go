package amqplb

import "sync"

type disconnectNotifier struct {
	sync.Mutex
	channel             string
	disconnectListeners map[chan interface{}]interface{}
}

func (c *disconnectNotifier) NotifyDisconnect(ch chan interface{}) chan interface{} {
	c.Lock()
	c.disconnectListeners[ch] = nil
	c.Unlock()

	return ch
}

func (c *disconnectNotifier) UnsubscribeDisconnect(ch chan interface{}) {
	c.Lock()
	delete(c.disconnectListeners, ch)
	c.Unlock()
}

func (c *disconnectNotifier) notifyDisconnect() {
	c.Lock()
	for listener := range c.disconnectListeners {
		listener <- nil
	}
	c.Unlock()
}

type reconnectNotifier struct {
	sync.Mutex
	channel            string
	reconnectListeners map[chan interface{}]interface{}
}

func (c *reconnectNotifier) NotifyReconnect(ch chan interface{}) chan interface{} {
	c.Lock()
	c.reconnectListeners[ch] = nil
	c.Unlock()

	return ch
}

func (c *reconnectNotifier) UnsubscribeReconnect(ch chan interface{}) {
	c.Lock()
	delete(c.reconnectListeners, ch)
	c.Unlock()
}

func (c *reconnectNotifier) notifyReconnect() {
	c.Lock()
	for listener := range c.reconnectListeners {
		listener <- nil
	}
	c.Unlock()
}

type closeNotifier struct {
	sync.Mutex
	channel        string
	closeListeners map[chan interface{}]interface{}
}

func (c *closeNotifier) NotifyClose(ch chan interface{}) chan interface{} {
	c.Lock()
	c.closeListeners[ch] = nil
	c.Unlock()

	return ch
}

func (c *closeNotifier) UnsubscribeClose(ch chan interface{}) {
	c.Lock()
	delete(c.closeListeners, ch)
	c.Unlock()
}

func (c *closeNotifier) notifyClose() {
	c.Lock()
	for listener := range c.closeListeners {
		listener <- nil
	}
	c.Unlock()
}

type errorNotifier struct {
	sync.Mutex
	channel        string
	errorListeners map[chan error]interface{}
}

func (c *errorNotifier) NotifyError(ch chan error) chan error {
	c.Lock()
	c.errorListeners[ch] = nil
	c.Unlock()

	return ch
}

func (c *errorNotifier) UnsubscribeError(ch chan error) {
	c.Lock()
	delete(c.errorListeners, ch)
	c.Unlock()
}

func (c *errorNotifier) notifyError(err error) {
	c.Lock()
	for listener := range c.errorListeners {
		listener <- err
	}
	c.Unlock()
}
