package usecases_test

import (
	"testing"

	"github.com/nfwGytautas/oxylabs/internal/model"
	"github.com/nfwGytautas/oxylabs/internal/usecases"
)

func TestDisconnectNonExisting(t *testing.T) {
	repo := &TestRepo{
		OnGetConnection: func(id model.ConnectionID) *model.Connection {
			return nil
		},
		OnRemoveConnection: func(id model.ConnectionID) {
			t.Fatalf("Expected no connection to be removed, got %s", id)
		},
	}
	connector := &TestConnector{}

	usecases := usecases.Usecases{
		ConnectionRepo: repo,
		DataLimit:      100,
		DataQueue:      make(chan model.Data),
		Connector:      connector,
	}

	usecases.Disconnect("1")
}

func TestDisconnectExisting(t *testing.T) {
	repo := &TestRepo{
		OnGetConnection: func(id model.ConnectionID) *model.Connection {
			return &model.Connection{ID: id}
		},
		OnRemoveConnection: func(id model.ConnectionID) {
			if id != "1" {
				t.Fatalf("Expected connection 1 to be removed, got %s", id)
			}
		},
	}
	connector := &TestConnector{}

	usecases := usecases.Usecases{
		ConnectionRepo: repo,
		DataLimit:      100,
		DataQueue:      make(chan model.Data),
		Connector:      connector,
	}

	usecases.Disconnect("1")

	if !repo.RemoveCalled {
		t.Fatalf("Expected connection 1 to be removed")
	}
}
