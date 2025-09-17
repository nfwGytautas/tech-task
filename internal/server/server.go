package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
)

type TCPServer struct {
	address    string
	bufferSize int
	onConnect  OnConnectedCallback

	listener net.Listener

	ctx context.Context
}

type Connection struct {
	Out chan []byte
	In  chan []byte

	ctx    context.Context
	cancel context.CancelFunc
}

type OnConnectedCallback func(conn *Connection)

func NewTCPServer(ctx context.Context, address string, bufferSize int, onConnect OnConnectedCallback) *TCPServer {
	return &TCPServer{
		address:    address,
		bufferSize: bufferSize,
		onConnect:  onConnect,
		listener:   nil,
		ctx:        ctx,
	}
}

func (s *TCPServer) Run() error {
	var err error

	s.listener, err = net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.address, err)
	}

	log.Printf("Listening on %s", s.address)

	go func() {
		for {
			select {
			case <-s.ctx.Done():
				err := s.listener.Close()
				if err != nil {
					log.Fatalf("failed to close listener: %v", err)
				}
				return
			default:
				conn, err := s.listener.Accept()
				if err != nil {
					if errors.Is(err, net.ErrClosed) {
						// Server is closing, clean exit
						return
					}

					log.Printf("failed to accept connection: %v", err)
					continue
				}

				go s.handleConnection(conn)
			}
		}
	}()

	return nil
}

func (s *TCPServer) handleConnection(c net.Conn) {
	defer c.Close()

	log.Printf("Accepted connection from %s", c.RemoteAddr())

	ctx, cancel := context.WithCancel(s.ctx)

	conn := &Connection{
		Out:    make(chan []byte),
		In:     make(chan []byte),
		ctx:    ctx,
		cancel: cancel,
	}

	// Start reader and writer as goroutines because they can happen simultaneously
	go func() {
		buffer := make([]byte, s.bufferSize)
		for {
			select {
			case <-conn.ctx.Done():
				// Connection is closing, clean exit
				return
			default:
				num, err := c.Read(buffer)
				if err != nil {
					log.Printf("[Error] Failed to read data from connection: %v", err)
					cancel()
					return
				}

				conn.In <- buffer[:num]
			}

		}
	}()

	go func() {
		for {
			select {
			case <-conn.ctx.Done():
				// Connection is closing, clean exit
				return
			default:
				data := <-conn.Out
				_, err := c.Write(data)
				if err != nil {
					log.Printf("[Error] Failed to write data to connection: %v", err)
					cancel()
					return
				}
			}
		}
	}()

	go s.onConnect(conn)

	// Keep the connection alive until it is closed
	<-conn.ctx.Done()
}

func (c *Connection) Close() {
	c.cancel()
}
