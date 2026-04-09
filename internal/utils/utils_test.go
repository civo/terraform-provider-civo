package utils

import "testing"

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
