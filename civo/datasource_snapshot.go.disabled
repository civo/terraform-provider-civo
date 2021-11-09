package civo

import (
	"fmt"
	"log"
	"time"

	"github.com/civo/civogo"
	"github.com/gorhill/cronexpr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Data source to get from the api a specific snapshot
// using the id or the name
func dataSourceSnapshot() *schema.Resource {
	return &schema.Resource{
		Description: "Snapshots are saved instances of a block storage volume. Use this data source to retrieve the ID of a Civo snapshot for use in other resources.",
		Read:        dataSourceSnapshotRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
				Description:  "The name of the snapshot",
			},
			// Computed resource
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the instance from which the snapshot was be taken",
			},
			"safe": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If `true`, the instance will be shut down during the snapshot",
			},
			"cron_timing": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A string with the cron format",
			},
			"hostname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The hostname of the instance",
			},
			"template_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The disk image/template ID",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The region where the snapshot was taken",
			},
			"size_gb": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of the snapshot in GB",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the snapshot",
			},
			"next_execution": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If cron was define this date will be the next execution date",
			},
			"requested_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date where the snapshot was requested",
			},
			"completed_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date where the snapshot was completed",
			},
		},
	}
}

func dataSourceSnapshotRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	var searchBy string

	if id, ok := d.GetOk("id"); ok {
		log.Printf("[INFO] Getting the snapshot key by id")
		searchBy = id.(string)
	} else if name, ok := d.GetOk("name"); ok {
		log.Printf("[INFO] Getting the snapshot key by label")
		searchBy = name.(string)
	}

	snapShot, err := apiClient.FindSnapshot(searchBy)
	if err != nil {
		return fmt.Errorf("[ERR] failed to retrive snapshot: %s", err)
	}

	d.SetId(snapShot.ID)
	d.Set("name", snapShot.Name)
	d.Set("instance_id", snapShot.InstanceID)
	d.Set("safe", snapShot.Safe)
	d.Set("cron_timing", snapShot.Cron)
	d.Set("hostname", snapShot.Hostname)
	d.Set("template_id", snapShot.Template)
	d.Set("region", snapShot.Region)
	d.Set("size_gb", snapShot.SizeGigabytes)
	d.Set("state", snapShot.State)

	nextExecution := time.Time{}
	if snapShot.Cron != "" {
		nextExecution = cronexpr.MustParse(snapShot.Cron).Next(time.Now().UTC())
	}

	d.Set("next_execution", nextExecution.String())
	d.Set("requested_at", snapShot.RequestedAt.UTC().String())
	d.Set("completed_at", snapShot.CompletedAt.UTC().String())

	return nil
}
