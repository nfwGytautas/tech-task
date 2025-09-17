package e2e_test

import (
	"context"
	"testing"
	"time"
)

func TestSingleConnection(t *testing.T) {
	const port = "15500"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go SpinServer(ctx, port)

	// Single spammer
	{
		err := NewSpammer(ctx, "localhost:"+port, false)
		if err != nil {
			t.Fatalf("Failed to create spammer: %v", err)
		}
	}

	time.Sleep(10 * time.Second)
}
