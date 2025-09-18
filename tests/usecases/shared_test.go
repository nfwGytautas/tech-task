package usecases_test

import (
	"github.com/nfwGytautas/oxylabs/internal/model"
)

// In reality this would be some actual mock from a library

type TestConnector struct {
	SendCalled  bool
	CloseCalled bool

	OnSend  func(id model.ConnectionID, data []byte)
	OnClose func(id model.ConnectionID)
}

type TestRepo struct {
	AddCalled    bool
	RemoveCalled bool

	OnGetAllConnections func() []*model.Connection
	OnAddConnection     func(connection *model.Connection)
	OnRemoveConnection  func(id model.ConnectionID)
	OnGetConnection     func(id model.ConnectionID) *model.Connection
}

func NewTestDataQueue(f func(data model.Data)) chan model.Data {
	t := make(chan model.Data)

	go func() {
		for data := range t {
			f(data)
		}
	}()

	return t
}

func (c *TestConnector) Send(id model.ConnectionID, data []byte) {
	if c.OnSend != nil {
		c.OnSend(id, data)
	}

	c.SendCalled = true
}

func (c *TestConnector) Close(id model.ConnectionID) {
	if c.OnClose != nil {
		c.OnClose(id)
	}

	c.CloseCalled = true
}

func (r *TestRepo) GetAllConnections() []*model.Connection {
	return r.OnGetAllConnections()
}

func (r *TestRepo) AddConnection(connection *model.Connection) {
	if r.OnAddConnection != nil {
		r.OnAddConnection(connection)
	}

	r.AddCalled = true
}

func (r *TestRepo) RemoveConnection(id model.ConnectionID) {
	if r.OnRemoveConnection != nil {
		r.OnRemoveConnection(id)
	}

	r.RemoveCalled = true
}

func (r *TestRepo) GetConnection(id model.ConnectionID) *model.Connection {
	return r.OnGetConnection(id)
}
