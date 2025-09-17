package room_tests

import (
	"testing"
	"time"

	"github.com/nfwGytautas/oxylabs/internal/room"
)

func TestLimitTwoClients(t *testing.T) {
	room := room.NewRoom(100)

	dd1 := NewDummyDriver()
	dd2 := NewDummyDriver()

	room.AddNewClient(dd1)
	room.AddNewClient(dd2)

	frame := make([]byte, 101)
	for i := 0; i < 101; i++ {
		frame[i] = 'a'
	}

	dd1.received <- frame

	time.Sleep(100 * time.Millisecond)

	if !dd1.closed {
		t.Fatalf("Driver 1 should be closed")
	}

	if !dd2.closed {
		t.Fatalf("Driver 2 should be closed")
	}

	// 2 Messages, 1 is the frame, 1 is the "Exceeded data limit" message
	if len(dd2.sent) != 2 {
		t.Fatalf("Driver 2 should have sent 1 message, got %d", len(dd2.sent))
	}

	if string(dd2.sent[1]) != "Exceeded data limit" {
		t.Fatalf("Driver 2 should have sent 'Exceeded data limit' message, got %s", string(dd2.sent[1]))
	}
}

func TestLimitThreeClients(t *testing.T) {
	room := room.NewRoom(100)

	dd1 := NewDummyDriver()
	dd2 := NewDummyDriver()
	dd3 := NewDummyDriver()

	room.AddNewClient(dd1)
	room.AddNewClient(dd2)

	frame := make([]byte, 51)
	for i := 0; i < 51; i++ {
		frame[i] = 'a'
	}

	dd1.received <- frame

	time.Sleep(100 * time.Millisecond)

	room.AddNewClient(dd3)

	time.Sleep(100 * time.Millisecond)

	dd1.received <- frame

	time.Sleep(100 * time.Millisecond)

	if !dd1.closed {
		t.Fatalf("Driver 1 should be closed")
	}

	if !dd2.closed {
		t.Fatalf("Driver 2 should be closed")
	}

	if dd3.closed {
		t.Fatalf("Driver 3 should not be closed")
	}

	// Last frame
	if len(dd3.sent) != 1 {
		t.Fatalf("Driver 3 should have sent 1 message, got %d", len(dd3.sent))
	}

	// First, second and exceeded
	if len(dd2.sent) != 3 {
		t.Fatalf("Driver 2 should have sent 3 message, got %d", len(dd2.sent))
	}

	if string(dd2.sent[2]) != "Exceeded data limit" {
		t.Fatalf("Driver 2 should have sent 'Exceeded data limit' message, got %s", string(dd2.sent[1]))
	}
}
