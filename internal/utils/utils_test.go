package utils

import (
	"fmt"
	"sync"
	"testing"

	"github.com/civo/civogo"
)

func TestIgnoreCaseDiff(t *testing.T) {
	tests := []struct {
		name     string
		old      string
		new      string
		suppress bool
	}{
		{
			name:     "same case lowercase",
			old:      "fra1",
			new:      "fra1",
			suppress: true,
		},
		{
			name:     "same case uppercase",
			old:      "FRA1",
			new:      "FRA1",
			suppress: true,
		},
		{
			name:     "upper to lower",
			old:      "FRA1",
			new:      "fra1",
			suppress: true,
		},
		{
			name:     "lower to upper",
			old:      "fra1",
			new:      "FRA1",
			suppress: true,
		},
		{
			name:     "mixed case",
			old:      "Fra1",
			new:      "fRA1",
			suppress: true,
		},
		{
			name:     "different regions",
			old:      "fra1",
			new:      "lon1",
			suppress: false,
		},
		{
			name:     "different regions different case",
			old:      "FRA1",
			new:      "lon1",
			suppress: false,
		},
		{
			name:     "empty strings",
			old:      "",
			new:      "",
			suppress: true,
		},
		{
			name:     "empty old",
			old:      "",
			new:      "fra1",
			suppress: false,
		},
		{
			name:     "empty new",
			old:      "fra1",
			new:      "",
			suppress: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IgnoreCaseDiff("region", tt.old, tt.new, nil)
			if got != tt.suppress {
				t.Errorf("IgnoreCaseDiff(%q, %q) = %v, want %v", tt.old, tt.new, got, tt.suppress)
			}
		})
	}
}

func TestRegionalClient(t *testing.T) {
	shared := &civogo.Client{Region: "lon1", APIKey: "secret"}

	scoped := RegionalClient(shared, "fra1")

	if scoped.Region != "fra1" {
		t.Errorf("scoped.Region = %q, want %q", scoped.Region, "fra1")
	}
	if scoped == shared {
		t.Error("RegionalClient returned the same pointer; it must return a copy")
	}
	if scoped.APIKey != shared.APIKey {
		t.Errorf("scoped.APIKey = %q, want %q (copy must carry over other fields)", scoped.APIKey, shared.APIKey)
	}
	if shared.Region != "lon1" {
		t.Errorf("shared.Region was mutated to %q; RegionalClient must not touch the shared client", shared.Region)
	}
}

// TestRegionalClientConcurrent reproduces the scenario from issue #395: many
// resources scoping the single shared client to different regions at the same
// time (as Terraform does under -parallelism). Each goroutine must see its own
// region. Run with -race to catch any regression back to mutating the shared
// client's Region field.
func TestRegionalClientConcurrent(t *testing.T) {
	shared := &civogo.Client{Region: "lon1", APIKey: "secret"}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			region := fmt.Sprintf("region-%d", i)
			scoped := RegionalClient(shared, region)
			if scoped.Region != region {
				t.Errorf("goroutine %d: scoped.Region = %q, want %q", i, scoped.Region, region)
			}
		}(i)
	}
	wg.Wait()

	if shared.Region != "lon1" {
		t.Errorf("shared.Region = %q after concurrent scoping, want %q", shared.Region, "lon1")
	}
}
