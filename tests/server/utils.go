package server_tests

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/nfwGytautas/oxylabs/internal/server"
)

func startServer(t *testing.T, ctx context.Context, port string) *server.TCPServer {
	tcpServer := server.NewTCPServer(ctx, "localhost:"+port, 1024, func(conn *server.Connection) {
		for {
			// Echo
			data := <-conn.In
			conn.Out <- data
		}
	})

	err := tcpServer.Run()
	if err != nil {
		t.Fatalf("Failed to start TCP server: %v", err)
	}

	// Dummy sleep to ensure the server is ready
	time.Sleep(100 * time.Millisecond)

	return tcpServer
}

func connectToServer(t *testing.T, port string) net.Conn {
	conn, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}

	return conn
}
