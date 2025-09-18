package main

import (
	"context"
	"log"
	"time"

	"github.com/nfwGytautas/oxylabs/internal/api"
	"github.com/nfwGytautas/oxylabs/internal/model"
	"github.com/nfwGytautas/oxylabs/internal/repo"
	"github.com/nfwGytautas/oxylabs/internal/usecases"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tcpServer := api.NewTCPServer(ctx, "localhost:9000", 100)

	usecases := usecases.Usecases{
		ConnectionRepo: &repo.ConnectionRepo{},
		DataLimit:      100,
		DataQueue:      make(chan model.Data),
		Connector:      tcpServer,
	}

	go usecases.SenderLoop()

	tcpServer.OnConnect = func(id model.ConnectionID) {
		usecases.Connect(id)
	}

	tcpServer.OnDisconnect = func(id model.ConnectionID) {
		usecases.Disconnect(id)
	}

	tcpServer.OnDataReceived = func(id model.ConnectionID, data []byte) {
		usecases.OnDataReceived(id, data)
	}

	go func() {
		for range time.Tick(10 * time.Millisecond) {
			usecases.Debug()
		}
	}()

	err := tcpServer.Run()
	if err != nil {
		log.Fatalf("Failed to start TCP server: %v", err)
	}
}
