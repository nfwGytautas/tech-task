package usecases_test

import (
	"testing"
	"time"

	"github.com/nfwGytautas/oxylabs/internal/model"
	"github.com/nfwGytautas/oxylabs/internal/usecases"
)

func TestSender(t *testing.T) {
	testConnection := model.Connection{ID: "1"}
	repo := &TestRepo{
		OnGetAllConnections: func() []*model.Connection {
			return []*model.Connection{
				&testConnection,
			}
		},
		OnRemoveConnection: func(id model.ConnectionID) {
			t.Fatalf("Expected no connection to be removed, got %s", id)
		},
	}
	connector := &TestConnector{
		OnSend: func(id model.ConnectionID, data []byte) {
			if id != "1" {
				t.Fatalf("Expected data to be sent to connection 1, got %s", id)
			}

			if string(data) != "test" {
				t.Fatalf("Expected data to be 'test', got %s", string(data))
			}
		},
		OnClose: func(id model.ConnectionID) {
			t.Fatalf("Expected no connection to be closed, got %s", id)
		},
	}

	dataQueue := make(chan model.Data)

	usecases := usecases.Usecases{
		ConnectionRepo: repo,
		DataLimit:      100,
		DataQueue:      dataQueue,
		Connector:      connector,
	}

	go usecases.SenderLoop()

	dataQueue <- model.Data{
		Sender: "2",
		Data:   []byte("test"),
	}

	// Give it time to process
	time.Sleep(10 * time.Millisecond)

	if !connector.SendCalled {
		t.Fatalf("Expected data to be sent to connection 1")
	}
}

func TestSenderExceededLimit(t *testing.T) {
	testConnection := model.Connection{ID: "1", OutgoingBytes: 99}
	repo := &TestRepo{
		OnGetAllConnections: func() []*model.Connection {
			return []*model.Connection{
				&testConnection,
			}
		},
	}
	connector := &TestConnector{
		OnSend: func(id model.ConnectionID, data []byte) {
			if id != "1" {
				t.Fatalf("Expected data to be sent to connection 1, got %s", id)
			}

			if string(data) != "Exceeded data limit" {
				t.Fatalf("Expected data to be 'Exceeded data limit', got %s", string(data))
			}
		},
		OnClose: func(id model.ConnectionID) {
			if id != "1" {
				t.Fatalf("Expected connection 1 to be closed, got %s", id)
			}
		},
	}

	dataQueue := make(chan model.Data)

	usecases := usecases.Usecases{
		ConnectionRepo: repo,
		DataLimit:      100,
		DataQueue:      dataQueue,
		Connector:      connector,
	}

	go usecases.SenderLoop()

	dataQueue <- model.Data{
		Sender: "2",
		Data:   []byte("test"),
	}

	// Give it time to process
	time.Sleep(10 * time.Millisecond)

	if !connector.SendCalled {
		t.Fatalf("Expected data to be sent to connection 1")
	}
	if !connector.CloseCalled {
		t.Fatalf("Expected connection 1 to be closed")
	}
	if !repo.RemoveCalled {
		t.Fatalf("Expected connection 1 to be removed")
	}
}

func TestSenderMultipleConnections(t *testing.T) {
	payload := []byte("test")

	testConnection1 := model.Connection{ID: "1"}
	testConnection2 := model.Connection{ID: "2"}

	repo := &TestRepo{
		OnGetAllConnections: func() []*model.Connection {
			return []*model.Connection{
				&testConnection1,
				&testConnection2,
			}
		},
		OnRemoveConnection: func(id model.ConnectionID) {
			t.Fatalf("Expected no connection to be removed, got %s", id)
		},
	}
	connector := &TestConnector{
		OnClose: func(id model.ConnectionID) {
			t.Fatalf("Expected no connection to be closed, got %s", id)
		},
	}

	dataQueue := make(chan model.Data)

	usecases := usecases.Usecases{
		ConnectionRepo: repo,
		DataLimit:      100,
		DataQueue:      dataQueue,
		Connector:      connector,
	}

	go usecases.SenderLoop()

	dataQueue <- model.Data{
		Sender: "0",
		Data:   payload,
	}

	time.Sleep(10 * time.Millisecond)

	if !connector.SendCalled {
		t.Fatalf("Expected data to be sent to connections")
	}

	if testConnection1.OutgoingBytes != len(payload) {
		t.Fatalf("Expected connection 1 outgoing bytes to be %d, got %d", len(payload), testConnection1.OutgoingBytes)
	}
	if testConnection2.OutgoingBytes != len(payload) {
		t.Fatalf("Expected connection 2 outgoing bytes to be %d, got %d", len(payload), testConnection2.OutgoingBytes)
	}
}

func TestSenderMultipleConnectionsNotSelf(t *testing.T) {
	payload := []byte("test")

	testConnection1 := model.Connection{ID: "1"}
	testConnection2 := model.Connection{ID: "2"}

	repo := &TestRepo{
		OnGetAllConnections: func() []*model.Connection {
			return []*model.Connection{
				&testConnection1,
				&testConnection2,
			}
		},
		OnRemoveConnection: func(id model.ConnectionID) {
			t.Fatalf("Expected no connection to be removed, got %s", id)
		},
	}
	connector := &TestConnector{
		OnSend: func(id model.ConnectionID, data []byte) {
			if id == "1" {
				t.Fatalf("Expected data to not be sent to connection 1")
			}
		},
		OnClose: func(id model.ConnectionID) {
			t.Fatalf("Expected no connection to be closed, got %s", id)
		},
	}

	dataQueue := make(chan model.Data)

	usecases := usecases.Usecases{
		ConnectionRepo: repo,
		DataLimit:      100,
		DataQueue:      dataQueue,
		Connector:      connector,
	}

	go usecases.SenderLoop()

	dataQueue <- model.Data{
		Sender: "1",
		Data:   payload,
	}

	time.Sleep(10 * time.Millisecond)

	if !connector.SendCalled {
		t.Fatalf("Expected data to be sent to connections")
	}

	if testConnection1.OutgoingBytes != 0 {
		t.Fatalf("Expected connection 1 outgoing bytes to be 0, got %d", testConnection1.OutgoingBytes)
	}
	if testConnection2.OutgoingBytes != len(payload) {
		t.Fatalf("Expected connection 2 outgoing bytes to be %d, got %d", len(payload), testConnection2.OutgoingBytes)
	}
}
