package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type TCPServer struct {
	address    string
	bufferSize int
	onConnect  OnConnectedCallback

	listener net.Listener

	ctx context.Context
}

type Connection struct {
	conn net.Conn

	out chan []byte
	in  chan []byte

	Ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
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

	log.Printf("[Server] Listening on %s", s.address)

	go func() {
		for {
			select {
			case <-s.ctx.Done():
				err := s.listener.Close()
				if err != nil {
					log.Fatalf("[Server] failed to close listener: %v", err)
				}
				return
			default:
				conn, err := s.listener.Accept()
				if err != nil {
					if errors.Is(err, net.ErrClosed) {
						// Server is closing, clean exit
						return
					}

					log.Printf("[Server] failed to accept connection: %v", err)
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

	log.Printf("[Server] Accepted connection from %s", c.RemoteAddr())

	ctx, cancel := context.WithCancel(s.ctx)

	conn := &Connection{
		conn:   c,
		out:    make(chan []byte),
		in:     make(chan []byte),
		Ctx:    ctx,
		cancel: cancel,
		wg:     sync.WaitGroup{},
	}

	conn.wg.Add(2)

	// Start reader and writer as goroutines because they can happen simultaneously
	go func() {
		conn.readLoop(s.bufferSize)
		conn.wg.Done()
	}()

	go func() {
		conn.writeLoop()
		conn.wg.Done()
	}()

	go s.onConnect(conn)

	// Keep the connection alive until work is finished
	conn.wg.Wait()
}

func (c *Connection) Close() {
	c.cancel()
}

func (c *Connection) Send(data []byte) {
	c.out <- data
}

func (c *Connection) Receive() []byte {
	return <-c.in
}

func (c *Connection) readLoop(size int) {
	buffer := make([]byte, size)
	for {
		select {
		case <-c.Ctx.Done():
			// Connection is closing, clean exit
			return
		default:
			num, err := c.conn.Read(buffer)
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}

				if errors.Is(err, io.EOF) {
					return
				}

				log.Printf("[Server] [Error] Failed to read data from connection: %v", err)
				c.cancel()
				return
			}

			c.in <- buffer[:num]
		}
	}
}

func (c *Connection) writeLoop() {
	for {
		select {
		case <-c.Ctx.Done():
			// Connection is closing, flush remaining
			for data := range c.out {
				_, err := c.conn.Write(data)
				if err != nil {
					log.Printf("[Server] [Error] Failed to write data to connection: %v", err)
					c.cancel()
					return
				}
			}
			return
		case data := <-c.out:
			_, err := c.conn.Write(data)
			if err != nil {
				log.Printf("[Server] [Error] Failed to write data to connection: %v", err)
				c.cancel()
				return
			}
		}
	}
}
