package server_tests

import (
	"context"
	"testing"

	test_utils "github.com/nfwGytautas/oxylabs/tests"
)

func TestSingleConnection(t *testing.T) {
	const testCasePort = "9000"
	const testCaseMessage = "Hello, world!"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_ = test_utils.StartServer(t, ctx, testCasePort)
	defer cancel()

	conn := test_utils.ConnectToServer(t, testCasePort)
	defer conn.Close()

	frame := make([]byte, 1024)

	copy(frame, testCaseMessage)

	// Dummy message
	_, err := conn.Write(frame)
	if err != nil {
		t.Fatalf("Failed to write to server: %v", err)
	}

	_, err = conn.Read(frame)
	if err != nil {
		t.Fatalf("Failed to read from server: %v", err)
	}

	if string(frame[:len(testCaseMessage)]) != testCaseMessage {
		t.Fatalf("Expected response to be '%s', got '%s'", testCaseMessage, string(frame))
	}
}
