package room_tests

import "sync"

type Message []byte

type DummyDriver struct {
	sync.RWMutex

	received chan []byte

	closed bool

	sent []Message
}

func (d *DummyDriver) Receive() []byte {
	return <-d.received
}

func (d *DummyDriver) Send(data []byte) {
	d.Lock()
	defer d.Unlock()
	d.sent = append(d.sent, data)
}

func (d *DummyDriver) Close() {
	d.closed = true
}

func NewDummyDriver() *DummyDriver {
	return &DummyDriver{
		received: make(chan []byte),
		closed:   false,
		sent:     []Message{},
	}
}
