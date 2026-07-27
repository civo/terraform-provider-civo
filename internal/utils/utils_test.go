package utils

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

	scoped, err := RegionalClient(shared, WithRegion("fra1"))
	if err != nil {
		t.Fatalf("RegionalClient returned an unexpected error: %v", err)
	}

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
			scoped, err := RegionalClient(shared, WithRegion(region))
			if err != nil {
				t.Errorf("goroutine %d: unexpected error: %v", i, err)
				return
			}
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

// --- region resolution (issue #395) -----------------------------------------

// regionTestSchema is the slice of a real resource schema that region
// resolution touches.
func regionTestSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"region":     {Type: schema.TypeString, Optional: true, Computed: true},
		"network_id": {Type: schema.TypeString, Optional: true, Computed: true},
	}
}

// regionsClient returns a client whose ListRegions call is served locally.
// Region is deliberately empty: that is the state the provider is in when
// neither the provider block nor CIVO_REGION sets one, which is what issue #395
// depends on.
func regionsClient(t *testing.T) *civogo.Client {
	t.Helper()
	client, server, err := civogo.NewClientForTesting(map[string]string{
		"/v2/regions": `[{"code":"fra1"},{"code":"lon1"},{"code":"nyc1"}]`,
	})
	if err != nil {
		t.Fatalf("could not build the test client: %v", err)
	}
	t.Cleanup(server.Close)
	client.Region = ""
	return client
}

// networkIn builds a RegionRef whose network is only visible in wantRegion, and
// records every region probed so tests can assert on the lookup itself.
func networkIn(wantRegion string, probed *[]string) RegionRef {
	return RegionRef{Field: "network_id", Kind: "network", Probe: func(c *civogo.Client, _ string) (bool, error) {
		*probed = append(*probed, c.Region)
		return c.Region == wantRegion, nil
	}}
}

func TestRegionalClientExplicitRegionWins(t *testing.T) {
	var probed []string
	d := schema.TestResourceDataRaw(t, regionTestSchema(), map[string]interface{}{
		"region": "lon1", "network_id": "net-1",
	})

	scoped, err := RegionalClient(regionsClient(t), WithRegion("fra1"), ResolveRegion(d, networkIn("nyc1", &probed)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if scoped.Region != "fra1" {
		t.Errorf("Region = %q, want %q (the explicit argument must win)", scoped.Region, "fra1")
	}
	if len(probed) != 0 {
		t.Errorf("probed %v; an explicit region must not trigger a lookup", probed)
	}
}

func TestRegionalClientUsesResourceRegion(t *testing.T) {
	var probed []string
	d := schema.TestResourceDataRaw(t, regionTestSchema(), map[string]interface{}{
		"region": "lon1", "network_id": "net-1",
	})

	scoped, err := RegionalClient(regionsClient(t), ResolveRegion(d, networkIn("nyc1", &probed)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if scoped.Region != "lon1" {
		t.Errorf("Region = %q, want %q (the resource's own region must win over inference)", scoped.Region, "lon1")
	}
	if len(probed) != 0 {
		t.Errorf("probed %v; a configured region must not trigger a lookup", probed)
	}
}

// The regression from issue #395: region unset, network_id pointing at another
// region. The client must end up scoped to the network's region, not defaulted.
func TestRegionalClientInfersRegionFromRef(t *testing.T) {
	var probed []string
	d := schema.TestResourceDataRaw(t, regionTestSchema(), map[string]interface{}{
		"network_id": "net-1",
	})

	scoped, err := RegionalClient(regionsClient(t), ResolveRegion(d, networkIn("nyc1", &probed)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if scoped.Region != "nyc1" {
		t.Errorf("Region = %q, want %q (inferred from the referenced network)", scoped.Region, "nyc1")
	}
	if len(probed) == 0 {
		t.Error("no region was probed; inference did not run")
	}
	// The inferred region must be recorded so the Read after create, and every
	// later operation, scope themselves without repeating the lookup.
	if got := d.Get("region").(string); got != "nyc1" {
		t.Errorf("d.Get(\"region\") = %q, want %q (inferred region must be written back)", got, "nyc1")
	}
}

// No region and no populated ref is unambiguous, so the client must be handed
// back untouched and the account default left to apply. This is the path every
// single-region config relies on.
func TestRegionalClientLeavesClientAloneWithNothingToGoOn(t *testing.T) {
	var probed []string
	shared := regionsClient(t)
	d := schema.TestResourceDataRaw(t, regionTestSchema(), map[string]interface{}{})

	scoped, err := RegionalClient(shared, ResolveRegion(d, networkIn("nyc1", &probed)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if scoped != shared {
		t.Error("client was copied; with nothing to resolve it must be returned as-is")
	}
	if scoped.Region != "" {
		t.Errorf("Region = %q, want empty", scoped.Region)
	}
	if len(probed) != 0 {
		t.Errorf("probed %v; there was no network_id to look up", probed)
	}
}

func TestRegionalClientErrorsWhenRefNotFoundAnywhere(t *testing.T) {
	var probed []string
	d := schema.TestResourceDataRaw(t, regionTestSchema(), map[string]interface{}{
		"network_id": "net-missing",
	})

	_, err := RegionalClient(regionsClient(t), ResolveRegion(d, networkIn("no-such-region", &probed)))
	if err == nil {
		t.Fatal("expected an error when the referenced network is in no region")
	}
	// The message has to be actionable: name the thing, the regions searched,
	// and the way out.
	for _, want := range []string{"net-missing", "fra1", "lon1", "nyc1", "set `region`"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q does not mention %q", err.Error(), want)
		}
	}
	if len(probed) != 3 {
		t.Errorf("probed %v, want all 3 regions tried", probed)
	}
}

// A ref that is not set must be skipped so a later one still gets a chance.
func TestRegionalClientSkipsEmptyRefs(t *testing.T) {
	var probed []string
	d := schema.TestResourceDataRaw(t, regionTestSchema(), map[string]interface{}{
		"network_id": "net-1",
	})

	unset := RegionRef{Field: "firewall_id", Kind: "firewall", Probe: func(*civogo.Client, string) (bool, error) {
		t.Error("probed a ref whose field was not set")
		return false, nil
	}}

	scoped, err := RegionalClient(regionsClient(t), ResolveRegion(d, unset, networkIn("lon1", &probed)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if scoped.Region != "lon1" {
		t.Errorf("Region = %q, want %q", scoped.Region, "lon1")
	}
}

// A probe that fails is not the same as the resource being absent. Reporting
// "not found in any of your regions" when a region was simply unreachable sends
// people after the wrong problem, so the underlying failure has to surface.
func TestRegionalClientSurfacesProbeFailure(t *testing.T) {
	d := schema.TestResourceDataRaw(t, regionTestSchema(), map[string]interface{}{
		"network_id": "net-1",
	})
	broken := RegionRef{Field: "network_id", Kind: "network", Probe: func(*civogo.Client, string) (bool, error) {
		return false, fmt.Errorf("503 from the API")
	}}

	_, err := RegionalClient(regionsClient(t), ResolveRegion(d, broken))
	if err == nil {
		t.Fatal("expected an error when every probe fails")
	}
	if !strings.Contains(err.Error(), "503 from the API") {
		t.Errorf("error %q hides the underlying probe failure", err.Error())
	}
}
