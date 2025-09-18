package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/nfwGytautas/oxylabs/internal/model"
)

type TCPServer struct {
	rw sync.RWMutex

	address    string
	bufferSize int

	listener    net.Listener
	ctx         context.Context
	connections map[model.ConnectionID]net.Conn

	OnConnect      func(id model.ConnectionID)
	OnDisconnect   func(id model.ConnectionID)
	OnDataReceived func(id model.ConnectionID, data []byte)
}

func NewTCPServer(ctx context.Context, address string, bufferSize int) *TCPServer {
	return &TCPServer{
		rw:          sync.RWMutex{},
		address:     address,
		bufferSize:  bufferSize,
		listener:    nil,
		ctx:         ctx,
		connections: make(map[model.ConnectionID]net.Conn),
	}
}

func (s *TCPServer) Run() error {
	var err error

	s.listener, err = net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.address, err)
	}

	log.Printf("[Server] Listening on %s", s.address)

	for {
		select {
		case <-s.ctx.Done():
			err := s.listener.Close()
			if err != nil {
				return fmt.Errorf("failed to close listener: %w", err)
			}
			return nil
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					// Server is closing, clean exit
					return nil
				}

				continue
			}

			go s.handleConnection(conn)
		}
	}
}

func (s *TCPServer) Send(id model.ConnectionID, data []byte) {
	s.rw.RLock()
	defer s.rw.RUnlock()

	conn, ok := s.connections[id]
	if !ok {
		log.Printf("[Server] [Error] Connection not found: %v", conn)
		return
	}

	_, err := conn.Write(data)
	if err != nil {
		log.Printf("[Server] [Error] Failed to write data to connection: %v", err)
		s.Close(id)
	}
}

func (s *TCPServer) Close(id model.ConnectionID) {
	s.rw.Lock()
	defer s.rw.Unlock()

	conn, ok := s.connections[id]
	if !ok {
		log.Printf("[Server] [Error] Connection not found: %v", id)
		return
	}
	conn.Close()

	delete(s.connections, id)
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	s.rw.Lock()
	defer s.rw.Unlock()

	s.connections[model.ConnectionID(conn.RemoteAddr().String())] = conn

	if s.OnConnect != nil {
		s.OnConnect(model.ConnectionID(conn.RemoteAddr().String()))
	}

	// Reader loop
	go func() {
		buffer := make([]byte, s.bufferSize)
		for {
			n, err := conn.Read(buffer)

			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}

				if errors.Is(err, io.EOF) {
					if s.OnDisconnect != nil {
						s.OnDisconnect(model.ConnectionID(conn.RemoteAddr().String()))
					}
					return
				}

				continue
			}

			if s.OnDataReceived != nil {
				s.OnDataReceived(model.ConnectionID(conn.RemoteAddr().String()), buffer[:n])
			}
		}
	}()
}
