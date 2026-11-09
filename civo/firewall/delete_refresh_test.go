package firewall

import (
	"fmt"
	"testing"

	"github.com/civo/civogo"
)

// deletes returns a delete function that serves the given results in order
// and repeats the last one.
func deletes(results ...error) func(string) (*civogo.SimpleResponse, error) {
	i := 0
	return func(string) (*civogo.SimpleResponse, error) {
		err := results[i]
		if i < len(results)-1 {
			i++
		}
		if err != nil {
			return nil, err
		}
		return &civogo.SimpleResponse{Result: "success"}, nil
	}
}

func TestFirewallDeleteRefresh(t *testing.T) {
	inUse := fmt.Errorf("unknown error response - status: 409 Conflict, code: 409, reason: {\"code\":\"database_firewall_inuse_by_instance\",\"result\":\"failed\",\"reason\":\"in use\"}")
	notFound := fmt.Errorf("%w: Failed to find that firewall in the internal database", civogo.DatabaseFirewallNotFoundError)
	timeout := fmt.Errorf("%w: problem connecting to the API", civogo.TimeoutError)
	other := fmt.Errorf("%w: boom", civogo.InternalServerError)

	tests := []struct {
		name    string
		results []error
		polls   int
		want    []string
		wantErr bool
	}{
		{name: "deleted first time", results: []error{nil}, polls: 1, want: []string{"success"}},
		{name: "in use then deleted", results: []error{inUse, nil}, polls: 2, want: []string{"in-use", "success"}},
		{
			// #418: the DELETE that went through is followed by one that
			// cannot find the firewall any more.
			name:    "in use then not found",
			results: []error{inUse, notFound},
			polls:   2,
			want:    []string{"in-use", "success"},
		},
		{name: "timeout then deleted", results: []error{timeout, nil}, polls: 2, want: []string{"retrying", "success"}},
		{
			// Continuous transient errors stay in "retrying"; the wait's
			// Timeout, not a consecutive-error cap, decides when to stop.
			name:    "many timeouts keep waiting",
			results: []error{timeout, timeout, timeout, timeout, timeout, timeout, nil},
			polls:   7,
			want:    []string{"retrying", "retrying", "retrying", "retrying", "retrying", "retrying", "success"},
		},
		{name: "other error fails", results: []error{other}, polls: 1, want: []string{""}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			refresh := firewallDeleteRefresh(deletes(tt.results...), "fw-1")
			var err error
			for i := 0; i < tt.polls; i++ {
				var state string
				_, state, err = refresh()
				if state != tt.want[i] {
					t.Fatalf("poll %d: state = %q, want %q", i+1, state, tt.want[i])
				}
				if err != nil && i < tt.polls-1 {
					t.Fatalf("poll %d: unexpected error %v", i+1, err)
				}
			}
			if (err != nil) != tt.wantErr {
				t.Fatalf("last poll error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
