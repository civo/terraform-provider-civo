package volume

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/civo/civogo"
)

func TestWaitForVolumeCreate(t *testing.T) {
	timeoutErr := fmt.Errorf("%w: problem connecting to the API", civogo.TimeoutError)

	tests := []struct {
		name    string
		prefer  time.Duration
		timeout time.Duration
		poll    time.Duration
		seq     []struct {
			status string
			err    error
		}
		wantErr bool
	}{
		{
			name:    "available before prefer window",
			prefer:  time.Minute,
			timeout: time.Minute,
			poll:    5 * time.Millisecond,
			seq: []struct {
				status string
				err    error
			}{
				{status: "creating"},
				{status: "available"},
			},
		},
		{
			name:    "creating accepted after prefer window",
			prefer:  15 * time.Millisecond,
			timeout: time.Minute,
			poll:    5 * time.Millisecond,
			seq: []struct {
				status string
				err    error
			}{
				{status: "creating"},
				{status: "creating"},
				{status: "creating"},
				{status: "creating"},
			},
		},
		{
			name:    "transient then available",
			prefer:  time.Minute,
			timeout: time.Minute,
			poll:    5 * time.Millisecond,
			seq: []struct {
				status string
				err    error
			}{
				{err: timeoutErr},
				{status: "available"},
			},
		},
		{
			name:    "empty status then available",
			prefer:  time.Minute,
			timeout: time.Minute,
			poll:    5 * time.Millisecond,
			seq: []struct {
				status string
				err    error
			}{
				{status: ""},
				{status: "available"},
			},
		},
		{
			name:    "unexpected status fails",
			prefer:  time.Minute,
			timeout: time.Minute,
			poll:    5 * time.Millisecond,
			seq: []struct {
				status string
				err    error
			}{
				{status: "failed"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldPrefer := volumeCreatePreferAvailable
			oldPoll := volumeCreatePollInterval
			volumeCreatePreferAvailable = tt.prefer
			volumeCreatePollInterval = tt.poll
			defer func() {
				volumeCreatePreferAvailable = oldPrefer
				volumeCreatePollInterval = oldPoll
			}()

			i := 0
			find := func(string) (*civogo.Volume, error) {
				r := tt.seq[i]
				if i < len(tt.seq)-1 {
					i++
				}
				if r.err != nil {
					return nil, r.err
				}
				return &civogo.Volume{Status: r.status}, nil
			}

			err := waitForVolumeCreate(context.Background(), find, "vol-1", tt.timeout)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
