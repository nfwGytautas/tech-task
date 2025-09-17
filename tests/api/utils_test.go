package api_test

import (
	"bytes"
	"context"

	"github.com/nfwGytautas/oxylabs/internal/api"
	"github.com/nfwGytautas/oxylabs/internal/model"
)

func runServer(ctx context.Context, port string, expectedPayload []byte) {
	server := api.NewTCPServer(ctx, "localhost:"+port, 1024)

	server.OnDataReceived = func(id model.ConnectionID, data []byte) {
		if !bytes.Equal(data, expectedPayload) {
			panic("Incorrect data received: " + string(data))
		}
	}

	err := server.Run()
	if err != nil {
		panic(err)
	}
}
