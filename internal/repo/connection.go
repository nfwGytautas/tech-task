package repo

import (
	"sync"

	"github.com/nfwGytautas/oxylabs/internal/model"
)

type ConnectionRepo struct {
	rw          sync.RWMutex
	connections []*model.Connection
}

func (r *ConnectionRepo) GetAllConnections() []*model.Connection {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return r.connections
}

func (r *ConnectionRepo) AddConnection(connection *model.Connection) {
	r.rw.Lock()
	defer r.rw.Unlock()
	r.connections = append(r.connections, connection)
}

func (r *ConnectionRepo) RemoveConnection(id model.ConnectionID) {
	r.rw.Lock()
	defer r.rw.Unlock()

	for i, c := range r.connections {
		if c.ID == id {
			r.connections = append(r.connections[:i], r.connections[i+1:]...)
			return
		}
	}
}

func (r *ConnectionRepo) GetConnection(id model.ConnectionID) *model.Connection {
	r.rw.RLock()
	defer r.rw.RUnlock()

	for _, c := range r.connections {
		if c.ID == id {
			return c
		}
	}
	return nil
}
