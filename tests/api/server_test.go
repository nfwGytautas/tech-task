package api_test

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	const port = "15000"
	payload := []byte("Testing")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go runServer(ctx, port, payload)

	// Create test connection
	conn, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}

	_, err = conn.Write(payload)
	if err != nil {
		t.Fatalf("Failed to write data to server: %v", err)
	}

	time.Sleep(10 * time.Millisecond)

	conn.Close()

	cancel()
}
