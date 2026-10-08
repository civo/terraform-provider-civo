package volume

import (
	"context"
	"errors"
	"testing"
	"time"
)

// TestWaitForVolumeResize covers the two phases of the post-resize wait: the grace window in which
// the transitional state must be observed once, and the longer wait for it to clear.
func TestWaitForVolumeResize(t *testing.T) {
	fast := resizeWaitTimings{grace: 50 * time.Millisecond, gracePoll: 5 * time.Millisecond, timeout: time.Second, poll: 5 * time.Millisecond}

	sequence := func(states ...string) (func() (string, error), *int) {
		calls := 0
		return func() (string, error) {
			i := calls
			if i >= len(states) {
				i = len(states) - 1
			}
			calls++
			return states[i], nil
		}, &calls
	}

	t.Run("waits through resizing until it clears", func(t *testing.T) {
		status, calls := sequence("attached", "attached", "resizing", "resizing", "attached")
		if err := waitForVolumeResize(context.Background(), "vol-1", status, fast); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if *calls < 5 {
			t.Fatalf("expected the wait to poll past the resizing state, got %d polls", *calls)
		}
	})

	t.Run("settles when resizing is never observed in the grace window", func(t *testing.T) {
		status, calls := sequence("attached")
		if err := waitForVolumeResize(context.Background(), "vol-1", status, fast); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if *calls < 2 {
			t.Fatalf("expected the grace window to poll more than once, got %d polls", *calls)
		}
	})

	t.Run("propagates a lookup error", func(t *testing.T) {
		status := func() (string, error) { return "", errors.New("boom") }
		if err := waitForVolumeResize(context.Background(), "vol-1", status, fast); err == nil || err.Error() != "boom" {
			t.Fatalf("expected the lookup error, got %v", err)
		}
	})

	t.Run("times out while still resizing", func(t *testing.T) {
		status, _ := sequence("resizing")
		if err := waitForVolumeResize(context.Background(), "vol-1", status, resizeWaitTimings{grace: 50 * time.Millisecond, gracePoll: 5 * time.Millisecond, timeout: 30 * time.Millisecond, poll: 5 * time.Millisecond}); err == nil {
			t.Fatal("expected a timeout error")
		}
	})
}
