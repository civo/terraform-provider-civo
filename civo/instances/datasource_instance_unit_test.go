package instances_test

import (
	"testing"

	"github.com/civo/terraform-provider-civo/civo/instances"
)

// TestDataSourceInstance_initialPasswordSensitive keeps the instance's initial password out of
// plan output when it is read through the civo_instance data source, matching the resource.
func TestDataSourceInstance_initialPasswordSensitive(t *testing.T) {
	if !instances.DataSourceInstance().Schema["initial_password"].Sensitive {
		t.Error("initial_password on the civo_instance data source should be Sensitive")
	}
	if !instances.ResourceInstance().Schema["initial_password"].Sensitive {
		t.Error("initial_password on the civo_instance resource should be Sensitive")
	}
}
