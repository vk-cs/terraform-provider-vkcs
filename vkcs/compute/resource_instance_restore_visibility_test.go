package compute

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWaitForRestoredServerVisibility(t *testing.T) {
	t.Run("retries a temporarily missing server", func(t *testing.T) {
		calls := 0
		getServer := func(context.Context, string) (bool, error) {
			calls++
			return calls > 1, nil
		}

		got, err := waitForRestoredServerVisibility(context.Background(), "server-id", time.Second, time.Millisecond, getServer)
		if err != nil {
			t.Fatalf("waitForRestoredServerVisibility() error = %v", err)
		}
		if !got || calls != 2 {
			t.Fatalf("waitForRestoredServerVisibility() = %t after %d calls, want visible after 2 calls", got, calls)
		}
	})

	t.Run("reports server absence after the grace period", func(t *testing.T) {
		getServer := func(context.Context, string) (bool, error) { return false, nil }

		got, err := waitForRestoredServerVisibility(context.Background(), "server-id", time.Millisecond, time.Hour, getServer)
		if err != nil {
			t.Fatalf("waitForRestoredServerVisibility() error = %v", err)
		}
		if got {
			t.Fatal("waitForRestoredServerVisibility() = true, want absent")
		}
	})

	t.Run("returns Nova API errors", func(t *testing.T) {
		wantErr := errors.New("Nova unavailable")
		getServer := func(context.Context, string) (bool, error) { return false, wantErr }

		_, err := waitForRestoredServerVisibility(context.Background(), "server-id", time.Second, time.Millisecond, getServer)
		if !errors.Is(err, wantErr) {
			t.Fatalf("waitForRestoredServerVisibility() error = %v, want wrapped %v", err, wantErr)
		}
	})
}
