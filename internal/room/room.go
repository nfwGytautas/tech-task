package room

import (
	"log"
	"sync"
)

type ClientDriver interface {
	Receive() []byte
	Send(data []byte)
	Close()
}

type client struct {
	id            int
	driver        ClientDriver
	room          *Room
	incomingCount int
	outgoingCount int
}

type Room struct {
	m sync.RWMutex

	clients   []*client
	dataLimit int

	nextId int
}

func NewRoom(dataLimit int) *Room {
	return &Room{
		m:         sync.RWMutex{},
		clients:   []*client{},
		dataLimit: dataLimit,
		nextId:    0,
	}
}

func (r *Room) AddNewClient(driver ClientDriver) {
	client := &client{
		driver:        driver,
		room:          r,
		incomingCount: 0,
		outgoingCount: 0,
		id:            r.nextId,
	}

	r.nextId++

	go client.read()

	r.m.Lock()
	defer r.m.Unlock()
	r.clients = append(r.clients, client)
}

func (c *client) read() {
	for {
		data := c.driver.Receive()
		c.incomingCount += len(data)

		log.Printf("Received %d bytes from client %d", c.incomingCount, c.id)

		// Broadcast
		go c.room.broadcast(c.id, data)

		if c.incomingCount > c.room.dataLimit {
			c.onLimitExceeded()
			return
		}
	}
}

func (c *client) write(data []byte) {
	c.outgoingCount += len(data)
	c.driver.Send(data)

	log.Printf("Sent %d bytes to client %d", c.outgoingCount, c.id)

	if c.outgoingCount > c.room.dataLimit {
		c.onLimitExceeded()
		return
	}
}

func (c *client) onLimitExceeded() {
	log.Printf("Client %d exceeded data limit: %d incoming, %d outgoing", c.id, c.incomingCount, c.outgoingCount)
	c.driver.Send([]byte("Exceeded data limit"))
	c.driver.Close()
	c.room.removeClient(c.id)
}

func (r *Room) broadcast(id int, data []byte) {
	r.m.RLock()
	defer r.m.RUnlock()

	for _, client := range r.clients {
		// Don't broadcast to the client that sent the data
		if client.id == id {
			continue
		}

		go client.write(data)
	}
}

func (r *Room) removeClient(id int) {
	r.m.Lock()
	defer r.m.Unlock()

	for i, client := range r.clients {
		if client.id == id {
			r.clients = append(r.clients[:i], r.clients[i+1:]...)
			return
		}
	}
}
