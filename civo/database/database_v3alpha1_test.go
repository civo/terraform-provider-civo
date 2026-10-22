package database

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestValidateStorageChange(t *testing.T) {
	cases := []struct {
		name     string
		old, new int
		wantErr  bool
	}{
		{"unset to set", 0, 120, false},
		{"grow", 40, 120, false},
		{"same", 120, 120, false},
		{"shrink", 120, 40, true},
		{"set to unset", 120, 0, false},
	}
	for _, c := range cases {
		err := validateStorageChange(c.old, c.new)
		if (err != nil) != c.wantErr {
			t.Errorf("%s: expected error %v, got %v", c.name, c.wantErr, err)
		}
	}
}

func TestExpandStringSet(t *testing.T) {
	set := schema.NewSet(schema.HashString, []interface{}{"postgis", "pg_cron", ""})
	got := expandStringSet(set)
	if len(got) != 2 {
		t.Errorf("expected 2 extensions, got %v", got)
	}
	if got := expandStringSet(nil); len(got) != 0 {
		t.Errorf("expected no extensions for nil, got %v", got)
	}
}

func TestExpandStringMap(t *testing.T) {
	got := expandStringMap(map[string]interface{}{"max_connections": "200", "ignored": 1})
	if len(got) != 1 || got["max_connections"] != "200" {
		t.Errorf("unexpected map %v", got)
	}
	if got := expandStringMap(nil); len(got) != 0 {
		t.Errorf("expected an empty map for nil, got %v", got)
	}
}

func TestForkRequestFromResource(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceDatabase().Schema, map[string]interface{}{
		"name":               "copy",
		"size":               "g3.db.small",
		"engine":             "PostgreSQL",
		"version":            "17",
		"nodes":              3,
		"source_database_id": "11111111-1111-1111-1111-111111111111",
		"source_target_time": "2026-10-21T12:00:00+02:00",
		"storage_gb":         120,
		"extensions":         []interface{}{"pg_cron"},
	})

	req, err := forkRequestFromResource(d, "LON1", "net-1", "fw-1")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if req.Name != "copy" || req.Size != "g3.db.small" || req.Nodes != 3 || req.SourceDatabaseID != "11111111-1111-1111-1111-111111111111" {
		t.Errorf("unexpected request %+v", req)
	}
	// The service does not derive the engine from the source: without these
	// it would create its default engine instead of a copy.
	if req.Software != "PostgreSQL" || req.SoftwareVersion != "17" {
		t.Errorf("expected the engine and version to be sent, got %q %q", req.Software, req.SoftwareVersion)
	}
	if req.NetworkID != "net-1" || req.FirewallID != "fw-1" || req.Region != "LON1" {
		t.Errorf("unexpected placement %+v", req)
	}
	if req.SourceTargetTime == nil || req.SourceTargetTime.Format("2006-01-02T15:04:05Z07:00") != "2026-10-21T10:00:00Z" {
		t.Errorf("expected the target time in UTC, got %v", req.SourceTargetTime)
	}
	if req.StorageGB == nil || *req.StorageGB != 120 || len(req.Extensions) != 1 || req.Extensions[0] != "pg_cron" {
		t.Errorf("unexpected storage/extensions %+v", req)
	}

	d = schema.TestResourceDataRaw(t, ResourceDatabase().Schema, map[string]interface{}{
		"name": "copy", "size": "g3.db.small", "engine": "PostgreSQL", "version": "17", "nodes": 1,
		"source_database_id": "11111111-1111-1111-1111-111111111111",
		"source_target_time": "not a time",
	})
	if _, err := forkRequestFromResource(d, "LON1", "", ""); err == nil {
		t.Errorf("expected an error for an invalid source_target_time")
	}
}
