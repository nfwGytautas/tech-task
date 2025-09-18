package usecases

import (
	"sync"

	"github.com/nfwGytautas/oxylabs/internal/model"
)

type Connector interface {
	Send(model.ConnectionID, []byte)
	Close(model.ConnectionID)
}

type Repo interface {
	GetAllConnections() []*model.Connection
	AddConnection(connection *model.Connection)
	RemoveConnection(id model.ConnectionID)
	GetConnection(id model.ConnectionID) *model.Connection
}

type Usecases struct {
	rw sync.RWMutex

	ConnectionRepo Repo
	DataLimit      int
	DataQueue      chan model.Data
	Connector      Connector
}
