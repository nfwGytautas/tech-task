package usecases

import "fmt"

func (u *Usecases) Debug() {
	u.rw.RLock()
	defer u.rw.RUnlock()

	fmt.Print("\033[H\033[2J")

	fmt.Printf("Connections: %d\n", len(u.ConnectionRepo.GetAllConnections()))
	for _, conn := range u.ConnectionRepo.GetAllConnections() {
		fmt.Printf("\t %s, in: %d, out: %d\n", conn.ID, conn.IncomingBytes, conn.OutgoingBytes)
	}
}
