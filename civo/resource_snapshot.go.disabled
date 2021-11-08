package civo

import (
	"fmt"
	"log"
	"time"

	"github.com/civo/civogo"
	"github.com/gorhill/cronexpr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Snapshot resource, with this we can create and manage all Snapshot
// this resource dont have update option, we used ForceNew so any change
// in any value will be recreate the resource
func resourceSnapshot() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a resource which can be used to create a snapshot from an existing Civo instance.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "This is a unqiue, alphanumerical, short, human readable code for the snapshot",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"instance_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The ID of the instance to snapshot",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"safe": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
				Description: "If true the instance will be shut down during the snapshot to ensure all files" +
					"are in a consistent state (e.g. database tables aren't in the middle of being optimised" +
					"and hence risking corruption). The default is false so you experience no interruption" +
					"of service, but a small risk of corruption.",
			},
			"cron_timing": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "If a valid cron string is passed, the snapshot will be saved as an automated snapshot," +
					"continuing to automatically update based on the schedule of the cron sequence provided." +
					"The default is nil meaning the snapshot will be saved as a one-off snapshot.",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			// Computed resource
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
		Create: resourceSnapshotCreate,
		Read:   resourceSnapshotRead,
		Delete: resourceSnapshotDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// function to create a new snapshot
func resourceSnapshotCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] configuring the new snapshot %s", d.Get("name").(string))
	config := &civogo.SnapshotConfig{
		InstanceID: d.Get("instance_id").(string),
	}

	if attr, ok := d.GetOk("safe"); ok {
		config.Safe = attr.(bool)
	}

	if attr, ok := d.GetOk("cron_timing"); ok {
		config.Cron = attr.(string)
	}

	log.Printf("[INFO] creating the new snapshot %s", d.Get("name").(string))
	resp, err := apiClient.CreateSnapshot(d.Get("name").(string), config)
	if err != nil {
		return fmt.Errorf("[ERR] failed to create snapshot: %s", err)
	}

	d.SetId(resp.ID)

	_, hasCronTiming := d.GetOk("cron_timing")

	if hasCronTiming {
		/*
			if hasCronTiming is declare them we no need to wait the state from the backend
		*/
		return resourceSnapshotRead(d, m)
	}
	/*
		if hasCronTiming is not declare them we need to wait the state from the backend
		and made a resource retry
	*/
	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		resp, err := apiClient.FindSnapshot(d.Id())
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error geting snapshot: %s", err))
		}

		if resp.State != "complete" {
			return resource.RetryableError(fmt.Errorf("[WARN] expected snapshot to be created but was in state %s", resp.State))
		}

		return resource.NonRetryableError(resourceSnapshotRead(d, m))
	})

}

// function to read the snapshot
func resourceSnapshotRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] retrieving the snapshot %s", d.Get("name").(string))
	resp, err := apiClient.FindSnapshot(d.Id())
	if err != nil {
		if resp == nil {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("[ERR] failed retrieving the snapshot: %s", err)
	}

	safeValue := false
	nextExecution := time.Time{}

	if resp.Safe == 1 {
		safeValue = true
	}

	if resp.Cron != "" {
		nextExecution = cronexpr.MustParse(resp.Cron).Next(time.Now().UTC())
	}

	d.Set("instance_id", resp.InstanceID)
	d.Set("hostname", resp.Hostname)
	d.Set("template_id", resp.Template)
	d.Set("region", resp.Region)
	d.Set("name", resp.Name)
	d.Set("safe", safeValue)
	d.Set("size_gb", resp.SizeGigabytes)
	d.Set("state", resp.State)
	d.Set("cron_timing", resp.Cron)
	d.Set("next_execution", nextExecution.String())
	d.Set("requested_at", resp.RequestedAt.UTC().String())
	d.Set("completed_at", resp.CompletedAt.UTC().String())

	return nil
}

// function to delete snapshot
func resourceSnapshotDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] deleting the snapshot %s", d.Id())
	_, err := apiClient.DeleteSnapshot(d.Id())
	if err != nil {
		return fmt.Errorf("[ERR] an error occurred while tring to delete the snapshot %s", d.Id())
	}

	return nil
}
