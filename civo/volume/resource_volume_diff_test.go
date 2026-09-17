package volume_test

import (
	"context"
	"strings"
	"testing"

	"github.com/civo/terraform-provider-civo/civo/volume"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestCivoVolume_sizeDecreaseRejectedAtPlan: volume size is additive only, so a plan that shrinks
// size_gb must fail before anything reaches the API, while a grow and an unchanged size plan cleanly.
func TestCivoVolume_sizeDecreaseRejectedAtPlan(t *testing.T) {
	state := &terraform.InstanceState{
		ID: "vol-1",
		Attributes: map[string]string{
			"id":         "vol-1",
			"name":       "data",
			"size_gb":    "10",
			"network_id": "net-1",
			"region":     "LON1",
		},
	}
	cases := []struct {
		name    string
		sizeGB  int
		wantErr string
	}{
		{name: "shrink is rejected", sizeGB: 5, wantErr: "size_gb cannot be decreased (from 10 to 5)"},
		{name: "grow plans", sizeGB: 20},
		{name: "same size plans", sizeGB: 10},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			config := terraform.NewResourceConfigRaw(map[string]interface{}{
				"name":       "data",
				"size_gb":    tc.sizeGB,
				"network_id": "net-1",
				"region":     "LON1",
			})
			_, err := volume.ResourceVolume().Diff(context.Background(), state, config, nil)
			switch {
			case tc.wantErr == "" && err != nil:
				t.Fatalf("unexpected plan error: %v", err)
			case tc.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tc.wantErr)):
				t.Fatalf("plan error = %v, want it to contain %q", err, tc.wantErr)
			}
		})
	}
}
