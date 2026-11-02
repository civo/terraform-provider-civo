package volume

import (
	"fmt"
	"testing"

	"github.com/civo/civogo"
)

func TestVolumeCreateRefresh(t *testing.T) {
	timeout := fmt.Errorf("%w: problem connecting to the API", civogo.TimeoutError)
	notFound := fmt.Errorf("%w: gone", civogo.DatabaseVolumeNotFoundError)

	type read struct {
		status string
		err    error
	}
	tests := []struct {
		name    string
		reads   []read
		want    []string
		wantErr bool
	}{
		{
			name:  "creating then available",
			reads: []read{{status: "creating"}, {status: "available"}},
			want:  []string{"creating", "available"},
		},
		{
			// #418: a client timeout after the volume ID was assigned.
			name:  "timeout keeps waiting",
			reads: []read{{status: "creating"}, {err: timeout}, {status: "available"}},
			want:  []string{"creating", "creating", "available"},
		},
		{
			name:    "persistent timeout gives up",
			reads:   []read{{err: timeout}, {err: timeout}, {err: timeout}, {err: timeout}, {err: timeout}, {err: timeout}},
			want:    []string{"creating", "creating", "creating", "creating", "creating", ""},
			wantErr: true,
		},
		{
			name:  "counter resets after a good read",
			reads: []read{{err: timeout}, {err: timeout}, {err: timeout}, {err: timeout}, {err: timeout}, {status: "creating"}, {err: timeout}, {status: "available"}},
			want:  []string{"creating", "creating", "creating", "creating", "creating", "creating", "creating", "available"},
		},
		{
			name:    "not found fails",
			reads:   []read{{err: notFound}},
			want:    []string{""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := 0
			find := func(string) (*civogo.Volume, error) {
				r := tt.reads[i]
				i++
				if r.err != nil {
					return nil, r.err
				}
				return &civogo.Volume{Status: r.status}, nil
			}
			refresh := volumeCreateRefresh(find, "vol-1")
			var err error
			for n := range tt.reads {
				var state string
				_, state, err = refresh()
				if state != tt.want[n] {
					t.Fatalf("poll %d: state = %q, want %q", n+1, state, tt.want[n])
				}
				if err != nil && n < len(tt.reads)-1 {
					t.Fatalf("poll %d: unexpected error %v", n+1, err)
				}
			}
			if (err != nil) != tt.wantErr {
				t.Fatalf("last poll error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
