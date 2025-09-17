package repo_test

import (
	"testing"

	"github.com/nfwGytautas/oxylabs/internal/model"
	"github.com/nfwGytautas/oxylabs/internal/repo"
)

func TestConnectionRepo(t *testing.T) {
	repo := &repo.ConnectionRepo{}

	repo.AddConnection(&model.Connection{ID: "1"})
	repo.AddConnection(&model.Connection{ID: "2"})
	repo.AddConnection(&model.Connection{ID: "3"})

	if len(repo.GetAllConnections()) != 3 {
		t.Fatalf("Expected 3 connections, got %d", len(repo.GetAllConnections()))
	}

	repo.RemoveConnection("2")
	if len(repo.GetAllConnections()) != 2 {
		t.Fatalf("Expected 2 connections, got %d", len(repo.GetAllConnections()))
	}

	repo.GetConnection("1")
	if repo.GetConnection("1") == nil {
		t.Fatalf("Expected connection 1, got nil")
	}

	repo.GetConnection("2")
	if repo.GetConnection("2") != nil {
		t.Fatalf("Expected nil, got %v", repo.GetConnection("2"))
	}
}
