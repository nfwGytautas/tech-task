package e2e_test

import (
	"context"
	"crypto/rand"
	"log"
	"math/big"
	"net"
	"strconv"
	"time"

	"github.com/nfwGytautas/oxylabs/internal/api"
	"github.com/nfwGytautas/oxylabs/internal/model"
	"github.com/nfwGytautas/oxylabs/internal/repo"
	"github.com/nfwGytautas/oxylabs/internal/usecases"
)

type Spammer struct {
	ctx     context.Context
	address string
	id      int
	logging bool

	backgroundCancel context.CancelFunc
}

func SpinServer(ctx context.Context, port string) {
	tcpServer := api.NewTCPServer(ctx, "localhost:"+port, 100)

	usecases := usecases.Usecases{
		ConnectionRepo: &repo.ConnectionRepo{},
		DataLimit:      100,
		DataQueue:      make(chan model.Data),
		Connector:      tcpServer,
	}

	go usecases.SenderLoop()

	tcpServer.OnConnect = func(id model.ConnectionID) {
		usecases.Connect(id)
	}

	tcpServer.OnDisconnect = func(id model.ConnectionID) {
		usecases.Disconnect(id)
	}

	tcpServer.OnDataReceived = func(id model.ConnectionID, data []byte) {
		usecases.OnDataReceived(id, data)
	}

	go func() {
		for range time.Tick(10 * time.Millisecond) {
			usecases.Debug()
		}
	}()

	err := tcpServer.Run()
	if err != nil {
		panic(err)
	}
}

func NewSpammer(ctx context.Context, address string, logging bool) error {
	id, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return err
	}

	s := &Spammer{
		ctx:     ctx,
		address: address,
		id:      int(id.Int64()),
		logging: logging,
	}

	go s.watcher()

	return nil
}

func (s *Spammer) watcher() {
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
		}

		ctx, cancel := context.WithCancel(s.ctx)
		s.backgroundCancel = cancel

		s.log("[%d] New instance", s.id)
		err := s.background(ctx)
		s.backgroundCancel()
		if err != nil {
			s.log("[%d] [Error] Failed to background: %v", s.id, err)
			return
		}
	}
}

func (s *Spammer) background(ctx context.Context) error {
	conn, err := net.Dial("tcp", s.address)
	if err != nil {
		return err
	}
	defer conn.Close()

	go func() {
		buffer := make([]byte, 1024)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			n, err := conn.Read(buffer)
			if err != nil {
				s.log("[%d] [Error] Failed to read from connection: %v", s.id, err)
				s.backgroundCancel()
				return
			}

			s.log("[%d] Message received: %v", s.id, string(buffer[:n]))
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		// Random sleep
		interval, err := rand.Int(rand.Reader, big.NewInt(5))
		if err != nil {
			s.log("[%d] [Error] Failed to generate random interval: %v", s.id, err)
			return err
		}

		time.Sleep(time.Duration(interval.Int64()) * 100 * time.Millisecond)

		_, err = conn.Write([]byte(strconv.Itoa(s.id)))
		if err != nil {
			s.log("[%d] [Error] Failed to write to connection: %v", s.id, err)
			return err
		}
	}
}

func (s *Spammer) log(format string, v ...interface{}) {
	if s.logging {
		log.Printf(format, v...)
	}
}
