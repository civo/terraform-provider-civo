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
		Read: dataSourceSnapshotRead,
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
			},
			// Computed resource
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"safe": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cron_timing": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size_gb": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"next_execution": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"requested_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"completed_at": {
				Type:     schema.TypeString,
				Computed: true,
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
