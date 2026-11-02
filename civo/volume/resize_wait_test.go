package volume

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/civo/civogo"
)

// fakeClock drives waitForVolumeResize without real sleeps: every sleep advances the clock, and a
// deadline turns the sleep into a context error, the way the Update timeout would.
type fakeClock struct {
	t        time.Time
	deadline time.Time
	slept    time.Duration
}

func (c *fakeClock) now() time.Time { return c.t }

func (c *fakeClock) sleep(ctx context.Context, d time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	c.t = c.t.Add(d)
	c.slept += d
	if !c.deadline.IsZero() && c.t.After(c.deadline) {
		return context.DeadlineExceeded
	}
	return nil
}

// reads returns a find function that serves the given volumes in order and repeats the last one.
func reads(volumes ...*civogo.Volume) func() (*civogo.Volume, error) {
	i := 0
	return func() (*civogo.Volume, error) {
		v := volumes[i]
		if i < len(volumes)-1 {
			i++
		}
		return v, nil
	}
}

// repeat returns n copies of v, for a read sequence that stays unchanged for a while.
func repeat(n int, v *civogo.Volume) []*civogo.Volume {
	out := make([]*civogo.Volume, n)
	for i := range out {
		out[i] = v
	}
	return out
}

func TestWaitForVolumeResize(t *testing.T) {
	t0 := time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)
	earlier := t0.Add(-time.Hour)
	stale := &civogo.VolumeResize{State: civogo.VolumeResizeFailed, Reason: "RequiresDetach", UpdatedAt: earlier}
	inProgress := &civogo.VolumeResize{State: civogo.VolumeResizeInProgress, Reason: "Expanding", UpdatedAt: t0.Add(10 * time.Second)}
	completed := &civogo.VolumeResize{State: civogo.VolumeResizeCompleted, Reason: "Completed", UpdatedAt: t0.Add(40 * time.Second)}
	failed := &civogo.VolumeResize{State: civogo.VolumeResizeFailed, Reason: "InsufficientStorage", Message: "pool is full", UpdatedAt: t0.Add(40 * time.Second)}

	cases := []struct {
		name          string
		baseline      *civogo.Volume
		reads         []*civogo.Volume
		deadline      time.Duration
		wantConfirmed bool
		wantErr       string
		wantMaxSlept  time.Duration
	}{
		{
			name:     "evidence then delivered",
			baseline: &civogo.Volume{Status: "attached", SizeGigabytes: 20, DeliveredSizeGigabytes: 20},
			reads: []*civogo.Volume{
				{Status: "attached", SizeGigabytes: 60, DeliveredSizeGigabytes: 20},
				{Status: "resizing", SizeGigabytes: 60, DeliveredSizeGigabytes: 20, Resize: inProgress},
				{Status: "attached", SizeGigabytes: 60, DeliveredSizeGigabytes: 60, Resize: completed},
			},
			wantConfirmed: true,
			wantMaxSlept:  30 * time.Second,
		},
		{
			name:     "evidence then a failed outcome is an error",
			baseline: &civogo.Volume{Status: "attached", SizeGigabytes: 20, DeliveredSizeGigabytes: 20},
			reads: []*civogo.Volume{
				{Status: "resizing", SizeGigabytes: 60, DeliveredSizeGigabytes: 20, Resize: inProgress},
				{Status: "attached", SizeGigabytes: 60, DeliveredSizeGigabytes: 20, Resize: failed},
			},
			wantErr: "the resize failed (InsufficientStorage): pool is full; the volume still delivers 20GB of the 60GB requested",
		},
		{
			name:     "a stale failed outcome from before the request is ignored while the new one completes",
			baseline: &civogo.Volume{Status: "available", SizeGigabytes: 20, DeliveredSizeGigabytes: 20, Resize: stale},
			reads: []*civogo.Volume{
				{Status: "available", SizeGigabytes: 60, DeliveredSizeGigabytes: 20, Resize: stale},
				{Status: "available", SizeGigabytes: 60, DeliveredSizeGigabytes: 20, Resize: stale},
				{Status: "available", SizeGigabytes: 60, DeliveredSizeGigabytes: 60, Resize: completed},
			},
			wantConfirmed: true,
			wantMaxSlept:  30 * time.Second,
		},
		{
			name:     "the delivered size alone is not evidence: the API falls back to the request",
			baseline: &civogo.Volume{Status: "attached", SizeGigabytes: 20, DeliveredSizeGigabytes: 20, Resize: stale},
			reads: []*civogo.Volume{
				{Status: "attached", SizeGigabytes: 60, DeliveredSizeGigabytes: 60, Resize: stale},
			},
			deadline: 5 * time.Minute,
			wantErr:  "gave up waiting for the resize (last status \"attached\", delivered 60GB of 60GB): context deadline exceeded",
		},
		{
			name:     "old API never picks the request up: unconfirmed, not an error",
			baseline: &civogo.Volume{Status: "available", SizeGigabytes: 20},
			reads: []*civogo.Volume{
				{Status: "available", SizeGigabytes: 60},
			},
			wantConfirmed: false,
		},
		{
			name:     "old API: resizing seen and cleared confirms",
			baseline: &civogo.Volume{Status: "available", SizeGigabytes: 20},
			reads: []*civogo.Volume{
				{Status: "available", SizeGigabytes: 60},
				{Status: "resizing", SizeGigabytes: 60},
				{Status: "resizing", SizeGigabytes: 60},
				{Status: "available", SizeGigabytes: 60},
			},
			wantConfirmed: true,
			wantMaxSlept:  40 * time.Second,
		},
		{
			name:     "completed before the first poll, with its own outcome",
			baseline: &civogo.Volume{Status: "available", SizeGigabytes: 20, DeliveredSizeGigabytes: 20},
			reads: []*civogo.Volume{
				{Status: "available", SizeGigabytes: 60, DeliveredSizeGigabytes: 60, Resize: completed},
			},
			wantConfirmed: true,
			wantMaxSlept:  0,
		},
		{
			name:     "completed between two polls after the grace window, resizing never seen: the outcome is the evidence",
			baseline: &civogo.Volume{Status: "attached", SizeGigabytes: 20, DeliveredSizeGigabytes: 20},
			// 14 unchanged reads at the 5s grace poll carry the wait past the 60s window.
			reads: append(repeat(14, &civogo.Volume{Status: "attached", SizeGigabytes: 60, DeliveredSizeGigabytes: 20}),
				&civogo.Volume{Status: "attached", SizeGigabytes: 60, DeliveredSizeGigabytes: 60, Resize: completed}),
			wantConfirmed: true,
			wantMaxSlept:  resizeGraceWindow + 2*resizeGracePoll,
		},
		{
			name:     "fallback delivered size while resizing keeps waiting and follows the real outcome",
			baseline: &civogo.Volume{Status: "attached", SizeGigabytes: 20, DeliveredSizeGigabytes: 20},
			reads: []*civogo.Volume{
				{Status: "resizing", SizeGigabytes: 60, DeliveredSizeGigabytes: 60},
				{Status: "resizing", SizeGigabytes: 60, DeliveredSizeGigabytes: 60},
				{Status: "attached", SizeGigabytes: 60, DeliveredSizeGigabytes: 20, Resize: failed},
			},
			wantErr: "the resize failed (InsufficientStorage): pool is full; the volume still delivers 20GB of the 60GB requested",
		},
		{
			name:     "fallback delivered size with a newer failed outcome is an error, not delivered",
			baseline: &civogo.Volume{Status: "attached", SizeGigabytes: 20, DeliveredSizeGigabytes: 20},
			reads: []*civogo.Volume{
				{Status: "resizing", SizeGigabytes: 60, DeliveredSizeGigabytes: 60, Resize: failed},
			},
			wantErr: "the resize failed (InsufficientStorage): pool is full; the volume still delivers 60GB of the 60GB requested",
		},
		{
			name:     "stuck in progress hits the deadline",
			baseline: &civogo.Volume{Status: "attached", SizeGigabytes: 20, DeliveredSizeGigabytes: 20},
			reads: []*civogo.Volume{
				{Status: "resizing", SizeGigabytes: 60, DeliveredSizeGigabytes: 20, Resize: inProgress},
			},
			deadline: 5 * time.Minute,
			wantErr:  "context deadline exceeded",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			clock := &fakeClock{t: t0}
			if tc.deadline > 0 {
				clock.deadline = t0.Add(tc.deadline)
			}
			w := resizeWait{
				find:      reads(tc.reads...),
				baseline:  tc.baseline,
				requested: 60,
				grace:     resizeGraceWindow,
				gracePoll: resizeGracePoll,
				poll:      resizePoll,
				now:       clock.now,
				sleep:     clock.sleep,
			}
			confirmed, err := waitForVolumeResize(context.Background(), "vol-1", w)
			switch {
			case tc.wantErr == "" && err != nil:
				t.Fatalf("unexpected error: %v", err)
			case tc.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tc.wantErr)):
				t.Fatalf("error = %v, want it to contain %q", err, tc.wantErr)
			}
			if confirmed != tc.wantConfirmed {
				t.Errorf("confirmed = %v, want %v", confirmed, tc.wantConfirmed)
			}
			if err == nil && tc.wantConfirmed && clock.slept > tc.wantMaxSlept {
				t.Errorf("slept %s before confirming, want at most %s", clock.slept, tc.wantMaxSlept)
			}
			if err == nil && !tc.wantConfirmed && clock.slept < resizeGraceWindow {
				t.Errorf("gave up after %s, before the %s grace window", clock.slept, resizeGraceWindow)
			}
		})
	}
}

func TestWaitForVolumeResize_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	clock := &fakeClock{t: time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)}
	w := resizeWait{
		find:      reads(&civogo.Volume{Status: "attached", SizeGigabytes: 60, DeliveredSizeGigabytes: 20}),
		baseline:  &civogo.Volume{Status: "attached", SizeGigabytes: 20, DeliveredSizeGigabytes: 20},
		requested: 60,
		grace:     resizeGraceWindow, gracePoll: resizeGracePoll, poll: resizePoll,
		now: clock.now, sleep: clock.sleep,
	}
	if _, err := waitForVolumeResize(ctx, "vol-1", w); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected the context error, got %v", err)
	}
}

func TestWaitForVolumeResize_LookupError(t *testing.T) {
	w := resizeWait{
		find:      func() (*civogo.Volume, error) { return nil, errors.New("boom") },
		baseline:  &civogo.Volume{},
		requested: 60,
	}
	if _, err := waitForVolumeResize(context.Background(), "vol-1", w); err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("expected the lookup error, got %v", err)
	}
}

func TestWaitForVolumeToSettle(t *testing.T) {
	t0 := time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)
	inProgress := &civogo.VolumeResize{State: civogo.VolumeResizeInProgress, Reason: "Expanding", UpdatedAt: t0}
	failed := &civogo.VolumeResize{State: civogo.VolumeResizeFailed, Reason: "Timeout", UpdatedAt: t0.Add(20 * time.Second)}

	t.Run("an earlier resize still running is waited out and the settled read is returned", func(t *testing.T) {
		clock := &fakeClock{t: t0}
		settled := &civogo.Volume{Status: "attached", SizeGigabytes: 60, DeliveredSizeGigabytes: 20, Resize: failed}
		w := resizeWait{
			find: reads(
				&civogo.Volume{Status: "resizing", SizeGigabytes: 60, DeliveredSizeGigabytes: 20, Resize: inProgress},
				&civogo.Volume{Status: "attached", SizeGigabytes: 60, DeliveredSizeGigabytes: 20, Resize: inProgress},
				settled,
			),
			poll: resizePoll, sleep: clock.sleep,
		}
		v, err := waitForVolumeToSettle(context.Background(), "vol-1", w)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v != settled {
			t.Errorf("returned %+v, want the settled read", v)
		}
		if clock.slept != 2*resizePoll {
			t.Errorf("slept %s, want %s", clock.slept, 2*resizePoll)
		}
	})

	t.Run("a settled volume returns at once", func(t *testing.T) {
		clock := &fakeClock{t: t0}
		w := resizeWait{find: reads(&civogo.Volume{Status: "available", SizeGigabytes: 20}), poll: resizePoll, sleep: clock.sleep}
		if _, err := waitForVolumeToSettle(context.Background(), "vol-1", w); err != nil || clock.slept != 0 {
			t.Fatalf("err %v, slept %s", err, clock.slept)
		}
	})

	t.Run("the deadline ends the wait", func(t *testing.T) {
		clock := &fakeClock{t: t0, deadline: t0.Add(time.Minute)}
		w := resizeWait{find: reads(&civogo.Volume{Status: "resizing", SizeGigabytes: 60}), poll: resizePoll, sleep: clock.sleep}
		if _, err := waitForVolumeToSettle(context.Background(), "vol-1", w); !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("expected the deadline error, got %v", err)
		}
	})
}
