package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// The helpers in this file back the attributes that only exist on Civo's
// /v3alpha1 database API: storage grow, extensions, server parameters, the
// read-only endpoint, the restore window and forking. The resource keeps
// using the /v2 calls for everything it did before, so regions that do not
// serve v3alpha1 yet keep working and simply leave those attributes empty.

// databaseReadyPollInterval is how often the status is polled while waiting
// for a database to come back to Ready after a v3alpha1 change.
const databaseReadyPollInterval = 5 * time.Second

// validateStorageChange enforces that storage can only grow. old and new are
// the planned values in GB; zero means "not set" and never fails.
func validateStorageChange(oldGB, newGB int) error {
	if oldGB > 0 && newGB > 0 && newGB < oldGB {
		return fmt.Errorf("storage_gb can only be increased: current %d GB, requested %d GB", oldGB, newGB)
	}
	return nil
}

// customizeDiffDatabase rejects plans that would shrink the storage before
// anything is sent to the API.
func customizeDiffDatabase(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	if diff.Id() == "" || !diff.HasChange("storage_gb") {
		return nil
	}
	oldValue, newValue := diff.GetChange("storage_gb")
	return validateStorageChange(oldValue.(int), newValue.(int))
}

// expandStringSet turns a TypeSet of strings into a slice.
func expandStringSet(value interface{}) []string {
	out := []string{}
	if value == nil {
		return out
	}
	for _, v := range value.(*schema.Set).List() {
		if s, ok := v.(string); ok && s != "" {
			out = append(out, s)
		}
	}
	return out
}

// expandStringMap turns a TypeMap of strings into a map.
func expandStringMap(value interface{}) map[string]string {
	out := map[string]string{}
	if value == nil {
		return out
	}
	for k, v := range value.(map[string]interface{}) {
		if s, ok := v.(string); ok {
			out[k] = s
		}
	}
	return out
}

// setDatabaseDetails reads the v3alpha1 view of the database and stores the
// attributes the /v2 API does not report. A region without v3alpha1 is not
// an error: the attributes are left as they are.
func setDatabaseDetails(d *schema.ResourceData, apiClient *civogo.Client, id string) {
	details, err := apiClient.GetDatabaseV3alpha1(id)
	if err != nil {
		log.Printf("[WARN] v3alpha1 details unavailable for database %s: %s", id, err)
		return
	}
	d.Set("storage_gb", details.StorageGB)
	d.Set("base_storage_gb", details.BaseStorageGB)
	d.Set("extensions", details.Extensions)
	d.Set("database_parameters", details.DatabaseParameters)
	if details.ReadEndpoint != nil {
		d.Set("read_endpoint", details.ReadEndpoint.Host)
		d.Set("read_endpoint_port", details.ReadEndpoint.Port)
		d.Set("read_endpoint_dns", details.ReadEndpoint.DNSEntry)
	} else {
		d.Set("read_endpoint", "")
		d.Set("read_endpoint_port", 0)
		d.Set("read_endpoint_dns", "")
	}
	if details.EarliestRestoreTime != nil {
		d.Set("earliest_restore_time", details.EarliestRestoreTime.UTC().Format(time.RFC3339))
	} else {
		d.Set("earliest_restore_time", "")
	}
}

// waitForDatabaseReady polls until the database reports Ready, the context
// expires, or the database reports a failure.
func waitForDatabaseReady(ctx context.Context, apiClient *civogo.Client, id string) error {
	// The status only flips to Pending on the operator's next reconcile, so
	// polling immediately would see the previous Ready and return too early.
	select {
	case <-ctx.Done():
		return fmt.Errorf("timed out waiting for database %s to be Ready", id)
	case <-time.After(databaseReadyPollInterval):
	}
	for {
		db, err := apiClient.GetDatabase(id)
		if err != nil {
			return err
		}
		switch db.Status {
		case "Ready":
			return nil
		case "Failed":
			return fmt.Errorf("database %s reported status %s", id, db.Status)
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for database %s to be Ready (last status %s)", id, db.Status)
		case <-time.After(databaseReadyPollInterval):
		}
	}
}

// applyDatabaseSettings pushes the v3alpha1-only settings that changed:
// storage, extensions and parameters. Each one waits for the database to be
// Ready again before the next is sent, since extensions and parameters can
// trigger a rolling restart.
//
// It must run before anything calls d.Set for these attributes (such as
// setDatabaseDetails): d.Get returns the value last set on the resource data,
// and once the API values are written in, the planned ones are gone.
func applyDatabaseSettings(ctx context.Context, d *schema.ResourceData, apiClient *civogo.Client, id string, force bool) error {
	if storageGB, ok := d.GetOk("storage_gb"); ok && (force || d.HasChange("storage_gb")) {
		gb := storageGB.(int)
		current := 0
		if force {
			// On create the plan's base storage is not in state yet; ask the API.
			if details, err := apiClient.GetDatabaseV3alpha1(id); err == nil {
				current = details.StorageGB
			}
		}
		if force && current > 0 && gb <= current {
			// Nothing to grow: the plan already includes this much.
			log.Printf("[INFO] database %s storage %d GB is within the current %d GB, skipping grow", id, gb, current)
		} else {
			log.Printf("[INFO] growing the storage of database %s to %d GB", id, gb)
			if _, err := apiClient.ResizeDatabaseStorage(id, gb); err != nil {
				return fmt.Errorf("failed to grow the storage: %w", err)
			}
			if err := waitForDatabaseReady(ctx, apiClient, id); err != nil {
				return err
			}
		}
	}

	if force || d.HasChange("extensions") {
		if _, ok := d.GetOk("extensions"); ok || d.HasChange("extensions") {
			exts := expandStringSet(d.Get("extensions"))
			log.Printf("[INFO] setting the extensions of database %s to %v", id, exts)
			if _, err := apiClient.UpdateDatabaseExtensions(id, exts); err != nil {
				return fmt.Errorf("failed to update the extensions: %w", err)
			}
			if err := waitForDatabaseReady(ctx, apiClient, id); err != nil {
				return err
			}
		}
	}

	if force || d.HasChange("database_parameters") {
		if _, ok := d.GetOk("database_parameters"); ok || d.HasChange("database_parameters") {
			params := expandStringMap(d.Get("database_parameters"))
			log.Printf("[INFO] setting the parameters of database %s to %v", id, params)
			if _, err := apiClient.UpdateDatabaseParameters(id, params); err != nil {
				return fmt.Errorf("failed to update the parameters: %w", err)
			}
			if err := waitForDatabaseReady(ctx, apiClient, id); err != nil {
				return err
			}
		}
	}

	return nil
}

// forkRequestFromResource builds the fork request when source_database_id is
// set. Size, nodes and network are the resource's own values, which the user
// normally copies from the source. Engine and version are sent too: the
// service does not derive them from the source and would fall back to its
// default engine otherwise.
func forkRequestFromResource(d *schema.ResourceData, region, networkID, firewallID string) (*civogo.ForkDatabaseRequest, error) {
	req := &civogo.ForkDatabaseRequest{
		Name:             d.Get("name").(string),
		Size:             d.Get("size").(string),
		Nodes:            d.Get("nodes").(int),
		Software:         d.Get("engine").(string),
		SoftwareVersion:  d.Get("version").(string),
		SourceDatabaseID: d.Get("source_database_id").(string),
		NetworkID:        networkID,
		FirewallID:       firewallID,
		Region:           region,
	}
	if target, ok := d.GetOk("source_target_time"); ok {
		t, err := time.Parse(time.RFC3339, target.(string))
		if err != nil {
			return nil, fmt.Errorf("source_target_time must be RFC 3339: %w", err)
		}
		t = t.UTC()
		req.SourceTargetTime = &t
	}
	if storageGB, ok := d.GetOk("storage_gb"); ok {
		gb := storageGB.(int)
		req.StorageGB = &gb
	}
	if _, ok := d.GetOk("extensions"); ok {
		req.Extensions = expandStringSet(d.Get("extensions"))
	}
	return req, nil
}
