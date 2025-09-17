package e2e_tests

import (
	"context"
	"testing"
	"time"

	"github.com/nfwGytautas/oxylabs/internal/room"
	"github.com/nfwGytautas/oxylabs/internal/server"
)

func TestE2EOneWay(t *testing.T) {
	const testCasePort = "10000"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	room := room.NewRoom(100)

	tcpServer := server.NewTCPServer(ctx, "localhost:"+testCasePort, 100, func(conn *server.Connection) {
		room.AddNewClient(conn)
	})

	err := tcpServer.Run()
	if err != nil {
		t.Fatalf("Failed to start TCP server: %v", err)
	}

	_, err = NewSpammer(ctx, "localhost:"+testCasePort)
	if err != nil {
		t.Fatalf("Failed to start spammer: %v", err)
	}

	time.Sleep(5 * time.Second)
}

func TestE2ETwoWay(t *testing.T) {
	const testCasePort = "10000"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	room := room.NewRoom(100)

	tcpServer := server.NewTCPServer(ctx, "localhost:"+testCasePort, 100, func(conn *server.Connection) {
		room.AddNewClient(conn)
	})

	err := tcpServer.Run()
	if err != nil {
		t.Fatalf("Failed to start TCP server: %v", err)
	}

	{
		_, err := NewSpammer(ctx, "localhost:"+testCasePort)
		if err != nil {
			t.Fatalf("Failed to start spammer: %v", err)
		}
	}

	{
		_, err := NewSpammer(ctx, "localhost:"+testCasePort)
		if err != nil {
			t.Fatalf("Failed to start spammer: %v", err)
		}
	}

	time.Sleep(5 * time.Second)
}

func TestE2EStress(t *testing.T) {
	const testCasePort = "10000"
	const numSpammers = 100

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	room := room.NewRoom(100)

	tcpServer := server.NewTCPServer(ctx, "localhost:"+testCasePort, 100, func(conn *server.Connection) {
		room.AddNewClient(conn)
	})

	err := tcpServer.Run()
	if err != nil {
		t.Fatalf("Failed to start TCP server: %v", err)
	}

	for i := 0; i < numSpammers; i++ {
		_, err := NewSpammer(ctx, "localhost:"+testCasePort)
		if err != nil {
			t.Fatalf("Failed to start spammer: %v", err)
		}
	}

	time.Sleep(5 * time.Second)
}
