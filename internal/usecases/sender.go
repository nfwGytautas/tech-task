package usecases

import (
	"github.com/nfwGytautas/oxylabs/internal/model"
)

func (u *Usecases) SenderLoop() {
	for data := range u.DataQueue {
		u.handleData(data)
	}
}

func (u *Usecases) handleData(data model.Data) {
	u.rw.RLock()
	defer u.rw.RUnlock()

	for _, conn := range u.ConnectionRepo.GetAllConnections() {
		if conn.ID == data.Sender {
			continue
		}

		conn.OutgoingBytes += len(data.Data)

		if conn.OutgoingBytes > u.DataLimit {
			// log.Printf("[ Usecases ] Connection %s exceeded download data limit, dropping...", conn.ID)
			u.Connector.Send(conn.ID, []byte("Exceeded data limit"))
			u.Connector.Close(conn.ID)
			u.ConnectionRepo.RemoveConnection(conn.ID)
			continue
		}

		u.Connector.Send(conn.ID, data.Data)
	}
}
