package usecases_test

import (
	"testing"

	"github.com/nfwGytautas/oxylabs/internal/model"
	"github.com/nfwGytautas/oxylabs/internal/usecases"
)

func TestConnect(t *testing.T) {
	repo := &TestRepo{
		OnAddConnection: func(connection *model.Connection) {
			if connection.ID != "1" {
				t.Fatalf("Expected connection id 1, got %s", connection.ID)
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

	usecases.Connect("1")

	if !repo.AddCalled {
		t.Fatalf("Expected connection to be added")
	}
}
