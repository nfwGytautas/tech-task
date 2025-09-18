package usecases_test

import (
	"testing"
	"time"

	"github.com/nfwGytautas/oxylabs/internal/model"
	"github.com/nfwGytautas/oxylabs/internal/usecases"
)

func TestDataReceivedNonExisting(t *testing.T) {
	repo := &TestRepo{
		OnGetConnection: func(id model.ConnectionID) *model.Connection {
			if id != "1" {
				t.Fatalf("Expected connection 1, got %s", id)
			}

			return nil
		},
		OnRemoveConnection: func(id model.ConnectionID) {
			t.Fatalf("Expected no connection to be removed, got %s", id)
		},
	}
	connector := &TestConnector{
		OnSend: func(id model.ConnectionID, data []byte) {
			t.Fatalf("Expected no data to be sent, got %s", id)
		},
		OnClose: func(id model.ConnectionID) {
			t.Fatalf("Expected no connection to be closed, got %s", id)
		},
	}

	usecases := usecases.Usecases{
		ConnectionRepo: repo,
		DataLimit:      100,
		DataQueue:      make(chan model.Data),
		Connector:      connector,
	}

	usecases.OnDataReceived("1", []byte("test"))
}

func TestDataReceivedExisting(t *testing.T) {
	data := []byte("test")

	testConnection := model.Connection{ID: "1"}
	repo := &TestRepo{
		OnGetConnection: func(id model.ConnectionID) *model.Connection {
			if id != "1" {
				t.Fatalf("Expected connection 1, got %s", id)
			}

			return &testConnection
		},
		OnRemoveConnection: func(id model.ConnectionID) {
			t.Fatalf("Expected no connection to be removed, got %s", id)
		},
	}
	connector := &TestConnector{
		OnSend: func(id model.ConnectionID, data []byte) {
			t.Fatalf("Expected no data to be sent, got %s", id)
		},
		OnClose: func(id model.ConnectionID) {
			t.Fatalf("Expected no connection to be closed, got %s", id)
		},
	}
	dataQueue := NewTestDataQueue(func(data model.Data) {
		if data.Sender != "1" {
			t.Fatalf("Expected data from connection 1, got %s", data.Sender)
		}
		if string(data.Data) != "test" {
			t.Fatalf("Expected data to be test, got %s", string(data.Data))
		}
	})

	usecases := usecases.Usecases{
		ConnectionRepo: repo,
		DataLimit:      100,
		DataQueue:      dataQueue,
		Connector:      connector,
	}

	usecases.OnDataReceived("1", data)

	// Sleep to let the queue channel process
	time.Sleep(10 * time.Millisecond)

	if testConnection.IncomingBytes != len(data) {
		t.Fatalf("Expected incoming bytes to be %d, got %d", len(data), testConnection.IncomingBytes)
	}
}

func TestDataReceivedExceededLimit(t *testing.T) {
	data := []byte("test")

	testConnection := model.Connection{ID: "1"}

	repo := &TestRepo{
		OnGetConnection: func(id model.ConnectionID) *model.Connection {
			if id != "1" {
				t.Fatalf("Expected connection 1, got %s", id)
			}

			return &testConnection
		},
		OnRemoveConnection: func(id model.ConnectionID) {
			if id != "1" {
				t.Fatalf("Expected connection 1 to be removed, got %s", id)
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
	dataQueue := NewTestDataQueue(func(data model.Data) {
		t.Fatalf("Expected no data to be enqueued, got %s", data.Sender)
	})

	usecases := usecases.Usecases{
		ConnectionRepo: repo,
		DataLimit:      100,
		DataQueue:      dataQueue,
		Connector:      connector,
	}

	testConnection.IncomingBytes = 99

	usecases.OnDataReceived("1", data)

	// Sleep to let the queue channel process
	time.Sleep(10 * time.Millisecond)

	if !connector.SendCalled {
		t.Fatalf("Expected data to be sent to connection 1")
	}
	if !repo.RemoveCalled {
		t.Fatalf("Expected connection 1 to be removed")
	}
	if !connector.CloseCalled {
		t.Fatalf("Expected connection 1 to be closed")
	}
}
