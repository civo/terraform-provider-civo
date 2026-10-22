package volume

import (
	"errors"
	"strings"
	"testing"
)

func TestResizeRefusalHint(t *testing.T) {
	refused := errors.New(`unknown error response - status: 422 Unprocessable Entity, code: 0, reason: {"code":"volume_onlie_resize","reason":"Volume is attached to an instance and its volume type does not support online expansion"}`)
	if hint := resizeRefusalHint(refused); !strings.Contains(hint, "no longer detaches the volume") || !strings.Contains(hint, "civo_volume_attachment") {
		t.Fatalf("hint on the detach-to-resize refusal = %q", hint)
	}
	for _, err := range []error{nil, errors.New("volume_resize_requires_running_instance"), errors.New("boom")} {
		if hint := resizeRefusalHint(err); hint != "" {
			t.Errorf("hint for %v = %q, want none", err, hint)
		}
	}
}
