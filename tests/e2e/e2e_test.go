package e2e_tests

import (
	"context"
	"testing"
	"time"

	"github.com/nfwGytautas/oxylabs/internal/room"
	test_utils "github.com/nfwGytautas/oxylabs/tests"
)

func TestE2ESimple(t *testing.T) {
	const testCasePort = "10000"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	room := room.NewRoom(2000)

	_ = test_utils.StartServer(t, ctx, testCasePort)
	defer cancel()

	{
		conn := test_utils.ConnectToServer(t, testCasePort)
		defer conn.Close()

		driver := &test_utils.NetConnDriver{
			Conn: conn,
		}

		room.AddNewClient(driver)
	}

	{
		conn := test_utils.ConnectToServer(t, testCasePort)
		defer conn.Close()

		driver := &test_utils.NetConnDriver{
			Conn: conn,
		}

		room.AddNewClient(driver)

		driver.Send([]byte("Hello, world!"))
	}

	time.Sleep(100 * time.Millisecond)
}
