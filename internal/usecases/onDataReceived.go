package usecases

import (
	"github.com/nfwGytautas/oxylabs/internal/model"
)

func (u *Usecases) OnDataReceived(id model.ConnectionID, payload []byte) {
	connection := u.ConnectionRepo.GetConnection(id)
	if connection == nil {
		return
	}

	u.rw.Lock()
	defer u.rw.Unlock()

	connection.IncomingBytes += len(payload)

	if connection.IncomingBytes > u.DataLimit {
		// log.Printf("[ Usecases ] Connection %s exceeded upload data limit, dropping...", id)
		u.Connector.Send(connection.ID, []byte("Exceeded data limit"))
		u.Connector.Close(connection.ID)
		u.ConnectionRepo.RemoveConnection(connection.ID)
		return
	}

	// log.Printf("[ Usecases ] OnDataReceived: %s, %d", connection.ID, connection.IncomingBytes)

	// Enqueue data for processing
	go func() {
		u.DataQueue <- model.Data{
			Sender: connection.ID,
			Data:   payload,
		}
	}()
}
