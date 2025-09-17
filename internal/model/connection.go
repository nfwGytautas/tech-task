package model

type ConnectionID string
type DataChannel = chan []byte

type Connection struct {
	ID ConnectionID

	IncomingBytes int
	OutgoingBytes int
}
