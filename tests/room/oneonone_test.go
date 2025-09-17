package room_tests

import (
	"testing"
	"time"

	"github.com/nfwGytautas/oxylabs/internal/room"
)

func TestOneOnOne(t *testing.T) {
	room := room.NewRoom(100)

	dd1 := NewDummyDriver()
	dd2 := NewDummyDriver()

	room.AddNewClient(dd1)
	room.AddNewClient(dd2)

	dd1.received <- []byte("Hello, world!")

	time.Sleep(100 * time.Millisecond)

	if dd1.closed {
		t.Fatalf("Driver 1 should not be closed")
	}

	if dd2.closed {
		t.Fatalf("Driver 2 should not be closed")
	}

	if len(dd2.sent) != 1 {
		t.Fatalf("Driver 2 should have sent 1 message, got %d", len(dd2.sent))
	}
}
