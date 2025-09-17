package test_utils

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/nfwGytautas/oxylabs/internal/server"
)

type NetConnDriver struct {
	Conn net.Conn
}

func (d *NetConnDriver) Receive() []byte {
	frame := make([]byte, 1024)
	_, err := d.Conn.Read(make([]byte, 1024))
	if err != nil {
		panic(err)
	}
	return frame
}

func (d *NetConnDriver) Send(data []byte) {
	_, err := d.Conn.Write(data)
	if err != nil {
		panic(err)
	}
}

func (d *NetConnDriver) Close() {
	d.Conn.Close()
}

func StartServer(t *testing.T, ctx context.Context, port string) *server.TCPServer {
	tcpServer := server.NewTCPServer(ctx, "localhost:"+port, 1024, func(conn *server.Connection) {
		for {
			// Echo
			data := conn.Receive()
			conn.Send(data)
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

func ConnectToServer(t *testing.T, port string) net.Conn {
	conn, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}

	return conn
}
