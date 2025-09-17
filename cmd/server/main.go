package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nfwGytautas/oxylabs/internal/room"
	"github.com/nfwGytautas/oxylabs/internal/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	room := room.NewRoom(100)

	tcpServer := server.NewTCPServer(ctx, "localhost:9000", 100, func(conn *server.Connection) {
		room.AddNewClient(conn)
	})

	err := tcpServer.Run()
	if err != nil {
		log.Fatalf("Failed to start TCP server: %v", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
