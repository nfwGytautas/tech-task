package e2e_test

import (
	"context"
	"testing"
	"time"
)

func TestStressConnection(t *testing.T) {
	const port = "15502"
	const numSpammers = 10

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go SpinServer(ctx, port)

	for i := 0; i < numSpammers; i++ {
		err := NewSpammer(ctx, "localhost:"+port, false)
		if err != nil {
			t.Fatalf("Failed to create spammer: %v", err)
		}
	}

	time.Sleep(10 * time.Second)
}
