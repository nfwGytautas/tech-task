package usecases

import (
	"sync"

	"github.com/nfwGytautas/oxylabs/internal/model"
	"github.com/nfwGytautas/oxylabs/internal/repo"
)

type Connector interface {
	Send(model.ConnectionID, []byte)
	Close(model.ConnectionID)
}

type Usecases struct {
	rw sync.RWMutex

	ConnectionRepo *repo.ConnectionRepo
	DataLimit      int
	DataQueue      chan model.Data
	Connector      Connector
}
