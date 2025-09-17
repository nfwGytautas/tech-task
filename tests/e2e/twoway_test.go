package e2e_test

import (
	"context"
	"testing"
	"time"
)

func TestTwoWayConnection(t *testing.T) {
	const port = "15501"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go SpinServer(ctx, port)

	// Two way spammer
	{
		err := NewSpammer(ctx, "localhost:"+port, false)
		if err != nil {
			t.Fatalf("Failed to create spammer: %v", err)
		}
	}

	{
		err := NewSpammer(ctx, "localhost:"+port, false)
		if err != nil {
			t.Fatalf("Failed to create spammer: %v", err)
		}
	}

	time.Sleep(10 * time.Second)
}
