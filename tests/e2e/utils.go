package e2e_tests

import (
	"context"
	"crypto/rand"
	"log"
	"math/big"
	"net"
	"sync"
	"time"
)

type Spammer struct {
	mu      sync.Mutex
	ctx     context.Context
	conn    net.Conn
	address string
}

func NewSpammer(ctx context.Context, address string) (*Spammer, error) {
	s := &Spammer{
		mu:      sync.Mutex{},
		ctx:     ctx,
		conn:    nil,
		address: address,
	}

	go s.background()

	return s, nil
}

func (s *Spammer) background() {
	s.reconnect()

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
		}

		if s.conn == nil {
			s.reconnect()
			continue // Continue, because we can potentially fail to connect
		}

		// Generate random payload
		payload := make([]byte, 50)
		_, err := rand.Read(payload)
		if err != nil {
			log.Println(err)
			return
		}

		// Random sleep
		interval, err := rand.Int(rand.Reader, big.NewInt(5))
		if err != nil {
			log.Println(err)
			return
		}

		time.Sleep(time.Duration(interval.Int64()) * 100 * time.Millisecond)

		s.mu.Lock()
		_, err = s.conn.Write(payload)
		if err != nil {
			log.Println(err)
			s.conn = nil
			s.mu.Unlock()
			continue
		}
		s.mu.Unlock()
	}
}

func (s *Spammer) reconnect() {
	s.mu.Lock()
	defer s.mu.Unlock()

	conn, err := net.Dial("tcp", s.address)
	if err != nil {
		log.Println(err)
		return
	}
	s.conn = conn
}

func (s *Spammer) Close() {
	s.conn.Close()
}
