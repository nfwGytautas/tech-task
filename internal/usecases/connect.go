package usecases

import "github.com/nfwGytautas/oxylabs/internal/model"

func (u *Usecases) Connect(id model.ConnectionID) {
	connection := &model.Connection{ID: id, IncomingBytes: 0, OutgoingBytes: 0}
	u.ConnectionRepo.AddConnection(connection)
}
