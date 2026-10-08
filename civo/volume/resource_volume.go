package volume

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceVolume function returns a schema.Resource that represents a Volume.
// This can be used to create, read, update, and delete operations for a Volume in the infrastructure.
func ResourceVolume() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Civo volume which can be attached to an instance in order to provide expanded storage.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "A name that you wish to use to refer to this volume",
				ValidateFunc: utils.ValidateName,
			},
			"size_gb": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "A minimum of 1 and a maximum of your available disk space from your quota specifies the size of the volume in gigabytes. Increases are applied in place: on a detached volume, and on a volume attached to a running instance when its volume type supports online expansion. Otherwise the API refuses the resize with `volume_onlie_resize` (detach the volume or move it to an expandable volume type) or `volume_resize_requires_running_instance` (start the instance). Decreases are rejected at plan time; recreate the volume with `terraform apply -replace` for a smaller size.",
			},
			"region": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "The region for the volume, if not declare we use the region in declared in the provider.",
				DiffSuppressFunc: utils.IgnoreCaseDiff,
			},
			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The network that the volume belongs to",
			},
			// Computed resource
			"mount_point": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The mount point of the volume (from instance's perspective)",
			},
			"volume_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of the volume",
			},
		},
		CreateContext: resourceVolumeCreate,
		ReadContext:   resourceVolumeRead,
		UpdateContext: resourceVolumeUpdate,
		DeleteContext: resourceVolumeDelete,
		CustomizeDiff: customizeDiffVolume,
		Importer: &schema.ResourceImporter{
			State: resourceVolumeImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(60 * time.Minute),
		},
	}
}

// customizeDiffVolume rejects a size decrease at plan time. Volume size is additive only: the API
// refuses a shrink, so failing the plan is friendlier than failing the apply. A size that is not
// known yet at plan time (derived from another resource's computed attribute) reads as zero through
// GetChange, so it is left alone here and checked by the API at apply.
func customizeDiffVolume(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	if diff.Id() == "" || !diff.NewValueKnown("size_gb") || !diff.HasChange("size_gb") {
		return nil
	}
	oldSize, newSize := diff.GetChange("size_gb")
	if newSize.(int) < oldSize.(int) {
		return fmt.Errorf("size_gb cannot be decreased (from %d to %d): volume size is additive only. "+
			"To get a smaller volume, recreate it with `terraform apply -replace=<address of this civo_volume>`. "+
			"This also applies right after an import whose configuration declares less than the real size: set size_gb to the real size first",
			oldSize.(int), newSize.(int))
	}
	return nil
}

// function to create the new volume
func resourceVolumeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	apiClient, err := utils.RegionalClient(apiClient, utils.ResolveRegion(d, utils.NetworkRef("network_id")))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] configuring the volume %s", d.Get("name").(string))

	config := &civogo.VolumeConfig{
		Name:          d.Get("name").(string),
		SizeGigabytes: d.Get("size_gb").(int),
		NetworkID:     d.Get("network_id").(string),
		Region:        apiClient.Region,
	}

	if v, ok := d.GetOk("volume_type"); ok {
		config.VolumeType = v.(string)
	}

	_, err = apiClient.FindNetwork(config.NetworkID)
	if err != nil {
		return diag.Errorf("[ERR] Unable to find network ID %q in %q region", config.NetworkID, config.Region)
	}

	volume, err := apiClient.NewVolume(config)
	if err != nil {
		return diag.Errorf("[ERR] failed to create a new volume: %s", err)
	}

	d.SetId(volume.ID)

	createStateConf := &resource.StateChangeConf{
		Pending: []string{"creating"},
		Target:  []string{"available"},
		Refresh: func() (interface{}, string, error) {
			resp, err := apiClient.FindVolume(d.Id())
			if err != nil {
				return 0, "", err
			}
			return resp, resp.Status, nil
		},
		Timeout:        60 * time.Minute,
		Delay:          3 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 10,
	}
	_, err = createStateConf.WaitForStateContext(context.Background())
	if err != nil {
		return diag.Errorf("error waiting for volume (%s) to be created: %s", d.Id(), err)
	}

	return resourceVolumeRead(ctx, d, m)
}

// function to read the volume
func resourceVolumeRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	apiClient, err := utils.RegionalClient(apiClient, utils.ResolveRegion(d))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] retrieving the volume %s", d.Id())
	resp, err := apiClient.FindVolume(d.Id())
	if err != nil {
		if resp == nil {
			d.SetId("")
			return nil
		}
		return diag.Errorf("[ERR] failed retrieving the volume: %s", err)
	}

	d.Set("name", resp.Name)
	d.Set("network_id", resp.NetworkID)
	d.Set("size_gb", resp.SizeGigabytes)
	d.Set("mount_point", resp.MountPoint)
	d.Set("volume_type", resp.VolumeType)

	return nil
}

// function to update the volume
func resourceVolumeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiClient := m.(*civogo.Client)

	apiClient, err := utils.RegionalClient(apiClient, utils.ResolveRegion(d))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] retrieving the volume %s", d.Id())
	resp, err := apiClient.FindVolume(d.Id())
	if err != nil {
		return diag.Errorf("[ERR] failed retrieving the volume: %s", err)
	}

	if d.HasChange("size_gb") {
		// Resize in place, attached or not. The API grows a detached volume with a clone-and-swap,
		// and a volume attached to a running instance online when its volume type supports
		// expansion; when it cannot, it refuses with volume_onlie_resize or
		// volume_resize_requires_running_instance (see the size_gb description), and that verdict is
		// returned to the user rather than worked around with a detach / resize / re-attach cycle.
		newSize := d.Get("size_gb").(int)
		log.Printf("[INFO] resizing the volume %s to %dGB (attached to %q)", d.Id(), newSize, resp.InstanceID)
		if _, err := apiClient.ResizeVolume(d.Id(), newSize); err != nil {
			return diag.Errorf("[ERR] failed to resize the volume %s to %dGB: %s", d.Id(), newSize, err)
		}
		// The request is accepted from here on: record the new size now, so an error while waiting
		// below does not leave the old size in state and re-issue the resize on the next apply.
		if err := d.Set("size_gb", newSize); err != nil {
			return diag.FromErr(err)
		}

		w := resizeWait{
			find:      func() (*civogo.Volume, error) { return apiClient.FindVolume(d.Id()) },
			baseline:  resp,
			requested: newSize,
			grace:     resizeGraceWindow,
			gracePoll: resizeGracePoll,
			poll:      resizePoll,
		}
		confirmed, err := waitForVolumeResize(ctx, d.Id(), w)
		if err != nil {
			return diag.Errorf("[ERR] resizing the volume %s to %dGB: %s", d.Id(), newSize, err)
		}
		if !confirmed {
			diags := diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Resize of volume %s to %dGB accepted but not confirmed", d.Id(), newSize),
				Detail: fmt.Sprintf("The API accepted the resize but the volume did not report any progress within %s, and this API version does not report the delivered size. "+
					"Verify the volume's size before relying on it (civo volume ls, or the dashboard).", resizeGraceWindow),
			}}
			return append(diags, resourceVolumeRead(ctx, d, m)...)
		}
	}

	if d.HasChange("network_id") {
		return diag.Errorf("[ERR] Network change for volume is not supported at this moment")
	}

	if d.HasChange("name") {
		return diag.Errorf("[ERR] Name change for volume is not supported at this moment")
	}

	return resourceVolumeRead(ctx, d, m)
}

// Timings for waitForVolumeResize. The API accepts a resize before the platform picks it up, so
// the first polls after ResizeVolume returns usually still read the previous state.
const (
	resizeGraceWindow = 60 * time.Second
	resizeGracePoll   = 5 * time.Second
	resizePoll        = 10 * time.Second
)

// Volume status the platform reports while it is carrying a resize out.
const volumeStatusResizing = "resizing"

// resizeWait describes one resize to follow: the volume as it was read just before the request
// (baseline), the size that was requested, and the polling cadence. now and sleep are injectable
// for tests; nil means the real clock.
type resizeWait struct {
	find      func() (*civogo.Volume, error)
	baseline  *civogo.Volume
	requested int
	grace     time.Duration
	gracePoll time.Duration
	poll      time.Duration
	now       func() time.Time
	sleep     func(context.Context, time.Duration) error
}

// resizeOf returns the volume's resize outcome when it belongs to this request, nil when the
// volume reports none or only the outcome of an earlier resize (same timestamp as the baseline).
func (w resizeWait) resizeOf(v *civogo.Volume) *civogo.VolumeResize {
	r := v.Resize
	if r == nil {
		return nil
	}
	if b := w.baseline.Resize; b != nil && r.UpdatedAt.Equal(b.UpdatedAt) && r.State == b.State {
		return nil
	}
	return r
}

// delivered reports whether the volume says it now provides at least the requested size. Zero
// means the API does not report the delivered size at all.
func (w resizeWait) delivered(v *civogo.Volume) bool {
	return v.DeliveredSizeGigabytes > 0 && v.DeliveredSizeGigabytes >= w.requested
}

// evidence reports whether the volume shows that this request has been picked up: the resizing
// status, or a resize outcome newer than the baseline's. The delivered size alone is not evidence:
// the API falls back to the requested size when the platform has not reported one, so it can read
// as delivered before anything happened. A migrating status alone is not evidence either, it
// belongs to volume-type migrations.
func (w resizeWait) evidence(v *civogo.Volume) bool {
	return v.Status == volumeStatusResizing || w.resizeOf(v) != nil
}

// waitForVolumeResize follows a resize the API has accepted until it is delivered or has failed.
// It returns true when the volume confirms the new size (the delivered size reached the request,
// the resize outcome says completed, or, on an API that reports neither, the resizing status was
// seen and has cleared), an error when the resize settled without delivering, when the context
// ends, or when the volume cannot be read, and false with no error when the volume showed no sign
// of the request within the grace window and reports no delivered size: the caller then warns
// rather than claims success. A resize outcome that predates the request is never judged.
func waitForVolumeResize(ctx context.Context, id string, w resizeWait) (bool, error) {
	if w.now == nil {
		w.now = time.Now
	}
	if w.sleep == nil {
		w.sleep = sleepContext
	}
	if w.baseline == nil {
		w.baseline = &civogo.Volume{}
	}

	start := w.now()
	seen := false
	for {
		v, err := w.find()
		if err != nil {
			return false, fmt.Errorf("reading the volume while waiting for the resize: %w", err)
		}
		if w.evidence(v) {
			seen = true
		}

		switch {
		case seen && w.delivered(v):
			return true, nil
		case seen && w.resizeOf(v) != nil && w.resizeOf(v).Settled():
			r := w.resizeOf(v)
			if r.Succeeded() {
				return true, nil
			}
			msg := fmt.Sprintf("the resize %s (%s)", r.State, r.Reason)
			if r.Message != "" {
				msg += ": " + r.Message
			}
			if v.DeliveredSizeGigabytes > 0 {
				msg += fmt.Sprintf("; the volume still delivers %dGB of the %dGB requested", v.DeliveredSizeGigabytes, w.requested)
			}
			return false, errors.New(msg)
		case seen && v.Status != volumeStatusResizing && v.Resize == nil && v.DeliveredSizeGigabytes == 0:
			// An API that reports neither the delivered size nor the outcome: the resizing
			// status was seen and has cleared, which is all it can tell.
			return true, nil
		case !seen && w.now().Sub(start) >= w.grace:
			if w.baseline.DeliveredSizeGigabytes > 0 || v.DeliveredSizeGigabytes > 0 {
				// The API reports the delivered size, so the request will show up as evidence
				// when it is picked up: keep waiting for it within the context's deadline.
				break
			}
			log.Printf("[WARN] volume %s showed no sign of the resize within %s and the API reports no delivered size", id, w.grace)
			return false, nil
		}

		interval := w.poll
		if !seen {
			interval = w.gracePoll
		}
		if err := w.sleep(ctx, interval); err != nil {
			return false, fmt.Errorf("gave up waiting for the resize (last status %q, delivered %dGB of %dGB): %w", v.Status, v.DeliveredSizeGigabytes, w.requested, err)
		}
	}
}

// sleepContext sleeps for d or until ctx ends, returning ctx.Err() in the latter case.
func sleepContext(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

// function to delete the volume
func resourceVolumeDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	apiClient, err := utils.RegionalClient(apiClient, utils.ResolveRegion(d))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] deleting the volume %s", d.Id())
	_, err = apiClient.DeleteVolume(d.Id())
	if err != nil {
		return diag.Errorf("[ERR] an error occurred while trying to delete the volume %s", err)
	}
	return nil
}

// custom import to able to import a volume
func resourceVolumeImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	apiClient := m.(*civogo.Client)
	regions, err := apiClient.ListRegions()
	if err != nil {
		return nil, err
	}

	volumeFound := false
	for _, region := range regions {
		if volumeFound {
			break
		}

		currentRegion := region.Code
		apiClient, err = utils.RegionalClient(apiClient, utils.WithRegion(currentRegion))
		if err != nil {
			return nil, err
		}

		volumes, err := apiClient.ListVolumes()
		if err != nil {
			return nil, err
		}

		for _, volume := range volumes {
			if volume.ID == d.Id() {
				volumeFound = true
				d.SetId(volume.ID)
				d.Set("name", volume.Name)
				d.Set("network_id", volume.NetworkID)
				d.Set("region", currentRegion)
				d.Set("size_gb", volume.SizeGigabytes)
				d.Set("mount_point", volume.MountPoint)
			}
		}
	}

	if !volumeFound {
		return nil, fmt.Errorf("[ERR] Volume %s not found", d.Id())
	}

	return []*schema.ResourceData{d}, nil
}
