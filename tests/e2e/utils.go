package e2e_tests

import (
	"context"
	"crypto/rand"
	"log"
	"math/big"
	"net"
	"strconv"
	"time"
)

type Spammer struct {
	ctx     context.Context
	address string
	id      int
}

func NewSpammer(ctx context.Context, address string) (*Spammer, error) {
	id, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return nil, err
	}

	s := &Spammer{
		ctx:     ctx,
		address: address,
		id:      int(id.Int64()),
	}

	go s.watcher()

	return s, nil
}

func (s *Spammer) watcher() {
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
		}

		log.Printf("[%d] New instance", s.id)
		err := s.background()
		if err != nil {
			log.Printf("[%d] [Error] Failed to background: %v", s.id, err)
			return
		}
	}
}

func (s *Spammer) background() error {
	conn, err := net.Dial("tcp", s.address)
	if err != nil {
		return err
	}
	defer conn.Close()

	go func() {
		buffer := make([]byte, 1024)
		for {
			select {
			case <-s.ctx.Done():
				return
			default:
			}

			n, err := conn.Read(buffer)
			if err != nil {
				log.Printf("[%d] [Error] Failed to read from connection: %v", s.id, err)
				return
			}

			log.Printf("[%d] Message received: %v", s.id, string(buffer[:n]))
		}
	}()

	for {
		select {
		case <-s.ctx.Done():
			return nil
		default:
		}

		// Random sleep
		interval, err := rand.Int(rand.Reader, big.NewInt(5))
		if err != nil {
			log.Printf("[%d] [Error] Failed to generate random interval: %v", s.id, err)
			return err
		}

		time.Sleep(time.Duration(interval.Int64()) * 100 * time.Millisecond)

		_, err = conn.Write([]byte(strconv.Itoa(s.id)))
		if err != nil {
			log.Printf("[%d] [Error] Failed to write to connection: %v", s.id, err)
			return err
		}
	}
}
