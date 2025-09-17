package usecases

import "github.com/nfwGytautas/oxylabs/internal/model"

func (u *Usecases) Disconnect(id model.ConnectionID) {
	connection := u.ConnectionRepo.GetConnection(id)
	if connection == nil {
		return
	}

	u.ConnectionRepo.RemoveConnection(connection.ID)
}
