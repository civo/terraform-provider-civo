package instances

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceInstance The instance resource represents an object of type instances
// and with it you can handle the instances created with Terraform
func ResourceInstance() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Civo instance resource. This can be used to create, modify, and delete instances.",
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The region for the instance, if not declare we use the region in declared in the provider",
			},
			"hostname": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "A fully qualified domain name that should be set as the instance's hostname",
				ValidateFunc: utils.ValidateNameSize,
			},
			"reverse_dns": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "A fully qualified domain name that should be used as the instance's IP's reverse DNS (optional, uses the hostname if unspecified)",
				ValidateFunc: utils.ValidateName,
			},
			"size": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "g3.xsmall",
				Description: "The name of the size, from the current list, e.g. g3.xsmall",
			},
			"public_ip_required": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "create",
				Description: "This should be either 'none' or 'create' (default: 'create')",
			},
			"network_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: utils.ValidateUUID,
				Description:  "This must be the ID of the network from the network listing (optional; default network used when not specified)",
			},
			"disk_image": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The ID for the disk image to use to build the instance",
				ForceNew:     true,
				ValidateFunc: utils.ValidateUUID,
			},
			"initial_user": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "civo",
				Description: "The name of the initial user created on the server (optional; this will default to the template's default_username and fallback to civo)",
			},
			"notes": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Add some notes to the instance",
			},
			"sshkey_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: utils.ValidateUUID,
				Description:  "The ID of an already uploaded SSH public key to use for login to the default user (optional; if one isn't provided a random password will be set and returned in the initial_password field)",
			},
			"firewall_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: utils.ValidateUUID,
				Description:  "The ID of the firewall to use, from the current list. If left blank or not sent, the default firewall will be used (open to all)",
			},
			"volume_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of volume to use, either 'ssd' or 'bssd' (optional; default 'ssd')",
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "An optional list of tags, represented as a key, value pair",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"script": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The contents of a script that will be uploaded to /usr/local/bin/civo-user-init-script on your instance, " +
					"read/write/executable only by root and then will be executed at the end of the cloud initialization",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			// Computed resource
			"cpu_cores": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Instance's CPU cores",
			},
			"ram_mb": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Instance's RAM (MB)",
			},
			"disk_gb": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Instance's disk (GB)",
			},
			"source_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance's source type",
			},
			"source_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance's source ID",
			},
			"initial_password": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Initial password for login",
			},
			"write_password": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          false,
				Description:      "If set to true, initial_password for instance will be saved to terraform state file",
				ValidateDiagFunc: utils.ValidateProviderVersion,
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance's private IP address",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance's public IP address",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance's status",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when the instance was created",
			},
			"private_ipv4": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The private IPv4 address for the instance (optional)",
			},
			"reserved_ipv4": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Can be either the UUID, name, or the IP address of the reserved IP",
			},
		},
		CreateContext: resourceInstanceCreate,
		ReadContext:   resourceInstanceRead,
		UpdateContext: resourceInstanceUpdate,
		DeleteContext: resourceInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		CustomizeDiff: customizeDiffInstance,
	}
}

// function to create an instance
func resourceInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is defined in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	log.Printf("[INFO] configuring the instance %s", d.Get("hostname").(string))
	config := &civogo.InstanceConfig{
		Count:            1,
		Hostname:         utils.RandomName(),
		Size:             "g3.medium",
		Region:           apiClient.Region,
		PublicIPRequired: "true",
		InitialUser:      "civo",
	}

	if hostname, ok := d.GetOk("hostname"); ok {
		config.Hostname = hostname.(string)
	} else {
		config.Hostname = utils.RandomName()
	}

	if attr, ok := d.GetOk("reverse_dns"); ok {
		config.ReverseDNS = attr.(string)
	}

	if attr, ok := d.GetOk("size"); ok {
		config.Size = attr.(string)
	}

	if attr, ok := d.GetOk("public_ip_required"); ok {
		config.PublicIPRequired = attr.(string)
	}

	if privateIPv4, ok := d.GetOk("private_ipv4"); ok {
		config.PrivateIPv4 = privateIPv4.(string)
	}

	if v, ok := d.GetOk("reserved_ipv4"); ok {
		config.ReservedIPv4 = v.(string)
	}

	if v, ok := d.GetOk("volume_type"); ok {
		config.VolumeType = v.(string)
	}

	if networtID, ok := d.GetOk("network_id"); ok {
		config.NetworkID = networtID.(string)
	} else {
		defaultNetwork, err := apiClient.GetDefaultNetwork()
		if err != nil {
			return diag.Errorf("[ERR] failed to get the default network: %s", err)
		}
		config.NetworkID = defaultNetwork.ID
	}

	if attr, ok := d.GetOk("disk_image"); ok {
		findDiskImage, err := apiClient.FindDiskImage(attr.(string))
		if err != nil {
			return diag.Errorf("[ERR] failed to get the disk image: %s", err)
		}
		config.TemplateID = findDiskImage.ID
	}

	if attr, ok := d.GetOk("initial_user"); ok {
		config.InitialUser = attr.(string)
	}

	if attr, ok := d.GetOk("sshkey_id"); ok {
		config.SSHKeyID = attr.(string)
	}

	if attr, ok := d.GetOk("script"); ok {
		config.Script = attr.(string)
	}

	tfTags := d.Get("tags").(*schema.Set).List()
	tags := make([]string, len(tfTags))
	for i, tfTag := range tfTags {
		tags[i] = tfTag.(string)
	}

	config.Tags = tags

	log.Printf("[INFO] creating the instance %s", d.Get("hostname").(string))

	// Initialize diagnostics
	diags := diag.Diagnostics{}

	isFirstInstance, err := checkNetworkFirstInstance(apiClient, config.NetworkID)
	if err != nil {
		return diag.Errorf("[ERR] failed to check network instances: %s", err)
	}

	if isFirstInstance {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "First Instance in Network",
			Detail:   fmt.Sprintf("The instance %s is the first instance in network %s and will be automatically assigned a public IP", config.Hostname, config.NetworkID),
		})
	}

	instance, err := apiClient.CreateInstance(config)
	if err != nil {
		customErr, parseErr := utils.ParseErrorResponse(err.Error())
		if parseErr == nil {
			err = customErr
		}
		// quota errors introduce new line after each missing quota, causing formatting issues:
		return diag.Errorf("[ERR] failed to create instance: %s", strings.ReplaceAll(err.Error(), "\n", " "))
	}

	d.SetId(instance.ID)

	createStateConf := &resource.StateChangeConf{
		Pending: []string{"BUILDING"},
		Target:  []string{"ACTIVE"},
		Refresh: func() (interface{}, string, error) {
			resp, err := apiClient.GetInstance(d.Id())
			if err != nil {
				return 0, "", err
			}
			return resp, resp.Status, nil
		},
		Timeout:        60 * time.Minute,
		Delay:          3 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 60,
	}
	_, err = createStateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for instance (%s) to be created: %s", d.Id(), err)
	}

	if attr, ok := d.GetOk("firewall_id"); ok {
		_, errInstance := apiClient.SetInstanceFirewall(d.Id(), attr.(string))
		if errInstance != nil {
			return diag.Errorf("[ERR] updating instance firewall: %s", err)
		}
	}

	if attr, ok := d.GetOk("notes"); ok {
		resp, err := apiClient.GetInstance(d.Id())
		if err != nil {
			return diag.Errorf("[ERR] getting instance: %s", err)
		}
		resp.Notes = attr.(string)
		_, errInstance := apiClient.UpdateInstance(resp)
		if errInstance != nil {
			return diag.Errorf("[ERR] updating instance notes: %s", err)
		}
	}

	// Append read resource diagnostics
	readDiags := resourceInstanceRead(ctx, d, m)
	diags = append(diags, readDiags...)

	return diags

}

// function to read the instance
func resourceInstanceRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is defined in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	log.Printf("[INFO] retriving the instance %s", d.Id())
	resp, err := apiClient.GetInstance(d.Id())
	if err != nil {
		if resp == nil {
			d.SetId("")
			return nil
		}

		return diag.Errorf("[ERR] failed to retriving the instance: %s", err)
	}

	diskImg, err := apiClient.GetDiskImageByName(resp.SourceID)
	if err != nil {
		return diag.Errorf("[ERR] failed to get the disk image: %s", err)
	}

	if d.Get("write_password").(bool) {
		d.Set("initial_password", resp.InitialPassword)
	} else {
		d.Set("initial_password", "")
	}

	if resp.Script == "" {
		d.Set("script", "")
	}

	decodedScript, err := base64.StdEncoding.DecodeString(resp.Script)
	if err != nil {
		return diag.Errorf("[ERR] failed to decode base64 script: %s", err)
	}

	d.Set("script", string(decodedScript))
	d.Set("hostname", resp.Hostname)
	d.Set("reverse_dns", resp.ReverseDNS)
	d.Set("size", resp.Size)
	d.Set("cpu_cores", resp.CPUCores)
	d.Set("ram_mb", resp.RAMMegabytes)
	d.Set("disk_gb", resp.DiskGigabytes)
	d.Set("initial_user", resp.InitialUser)
	d.Set("source_type", resp.SourceType)
	d.Set("source_id", resp.SourceID)
	d.Set("sshkey_id", resp.SSHKeyID)
	d.Set("tags", resp.Tags)
	d.Set("private_ip", resp.PrivateIP)
	d.Set("public_ip", resp.PublicIP)
	d.Set("network_id", resp.NetworkID)
	d.Set("firewall_id", resp.FirewallID)
	d.Set("status", resp.Status)
	d.Set("created_at", resp.CreatedAt.UTC().String())
	d.Set("notes", resp.Notes)
	d.Set("disk_image", diskImg.ID)
	d.Set("volume_type", resp.VolumeType)

	if resp.PublicIP != "" {
		d.Set("public_ip_required", "create")
	} else {
		d.Set("public_ip_required", "none")
	}

	if d.HasChange("reserved_ipv4") {
		_, new := d.GetChange("reserved_ipv4")
		newValue := new.(string)
		d.Set("reserved_ipv4", newValue)
	}

	return nil
}

// function to update an instance
func resourceInstanceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is defined in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	// check if the size change if change we send to resize the instance
	if d.HasChange("size") {
		newSize := d.Get("size").(string)

		log.Printf("[INFO] resizing the instance %s", d.Id())
		_, err := apiClient.UpgradeInstance(d.Id(), newSize)
		if err != nil {
			return diag.Errorf("[WARN] An error occurred while resizing the instance %s", d.Id())
		}

		createStateConf := &resource.StateChangeConf{
			Pending: []string{"BUILDING", "REBOOTING"},
			Target:  []string{"ACTIVE"},
			Refresh: func() (interface{}, string, error) {
				resp, err := apiClient.GetInstance(d.Id())
				if err != nil {
					return 0, "", err
				}
				return resp, resp.Status, nil
			},
			Timeout:        60 * time.Minute,
			Delay:          3 * time.Second,
			MinTimeout:     3 * time.Second,
			NotFoundChecks: 60,
		}
		_, err = createStateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for instance (%s) to be created: %s", d.Id(), err)
		}
	}

	// if notes or hostname have changed, add them to the instance
	if d.HasChange("notes") || d.HasChange("hostname") {
		notes := d.Get("notes").(string)
		hostname := d.Get("hostname").(string)

		instance, err := apiClient.GetInstance(d.Id())
		if err != nil {
			// check if the instance no longer exists.
			return diag.Errorf("[ERR] instance %s not found", d.Id())
		}

		if d.HasChange("notes") {
			instance.Notes = notes
		}
		if d.HasChange("hostname") {
			instance.Hostname = hostname
		}

		log.Printf("[INFO] updating instance %s", d.Id())
		_, err = apiClient.UpdateInstance(instance)
		if err != nil {
			return diag.Errorf("[ERR] an error occurred while updating notes or hostname of the instance %s", d.Id())
		}
	}

	// If reserved_ipv4 has changed, update the instance with the new reserved IP
	if d.HasChange("reserved_ipv4") {
		oldReservedIP, newReservedIP := d.GetChange("reserved_ipv4")
		instance, err := apiClient.GetInstance(d.Id())
		if err != nil {
			// Check if the instance no longer exists.
			return diag.Errorf("[ERR] instance %s not found", d.Id())
		}

		// Unassign the old reserved IP if it exists
		if oldReservedIP != "" {
			ip, err := apiClient.FindIP(oldReservedIP.(string))
			if err != nil {
				if errors.Is(err, civogo.ZeroMatchesError) {
					return diag.Errorf("sorry there is no %s IP in your account", oldReservedIP)
				} else if errors.Is(err, civogo.MultipleMatchesError) {
					return diag.Errorf("sorry we found more than one IP with that value in your account")
				} else {
					return diag.Errorf("error finding IP %s: %s", oldReservedIP, err)
				}
			}

			_, err = apiClient.UnassignIP(ip.ID, apiClient.Region)
			if err != nil {
				return diag.Errorf("[ERR] an error occurred while unassigning reserved IP %s from instance %s: %s", ip.ID, instance.ID, err)
			}
			log.Printf("[INFO] unassigned reserved IP %s from the instance %s", oldReservedIP, d.Id())
		}

		// Find the new reserved IP
		ip, err := apiClient.FindIP(newReservedIP.(string))
		if err != nil {
			if errors.Is(err, civogo.ZeroMatchesError) {
				return diag.Errorf("sorry there is no %s IP in your account", newReservedIP)
			} else if errors.Is(err, civogo.MultipleMatchesError) {
				return diag.Errorf("sorry we found more than one IP with that value in your account")
			} else {
				return diag.Errorf("error finding IP %s: %s", newReservedIP, err)
			}
		}

		// Assign the new reserved IP to the instance
		_, err = apiClient.AssignIP(ip.ID, instance.ID, "instance", apiClient.Region)
		if err != nil {
			return diag.Errorf("[ERR] an error occurred while assigning reserved IP %s to instance %s: %s", ip.ID, instance.ID, err)
		}

		log.Printf("[INFO] assigned reserved IP %s to the instance %s", newReservedIP, d.Id())
	}

	// if a firewall is declared we update the instance
	if d.HasChange("firewall_id") {
		firewallID := d.Get("firewall_id").(string)

		log.Printf("[INFO] adding firewall to the instance %s", d.Id())
		_, err := apiClient.SetInstanceFirewall(d.Id(), firewallID)
		if err != nil {
			// check if the instance no longer exists.
			return diag.Errorf("[ERR] an error occurred while set firewall to the instance %s", d.Id())
		}
	}

	if d.HasChange("initial_user") {
		return diag.Errorf("[ERR] updating initial_user is not supported")
	}

	if d.HasChange("sshkey_id") {
		return diag.Errorf("[ERR] updating sshkey_id is not supported")
	}

	// if tags is declare we update the instance with the tags
	if d.HasChange("tags") {
		tfTags := d.Get("tags").(*schema.Set).List()
		tags := make([]string, len(tfTags))
		for i, tfTag := range tfTags {
			tags[i] = tfTag.(string)
		}

		instance, err := apiClient.GetInstance(d.Id())
		if err != nil {
			// check if the instance no longer exists.
			return diag.Errorf("[ERR] instance %s not found", d.Id())
		}

		tagsToString := strings.Join(tags, " ")

		log.Printf("[INFO] adding tags to the instance %s", d.Id())
		_, err = apiClient.SetInstanceTags(instance, tagsToString)
		if err != nil {
			return diag.Errorf("[ERR] an error occurred while adding tags to the instance %s", d.Id())
		}

	}

	return resourceInstanceRead(ctx, d, m)
}

// function to delete instance.
func resourceInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is defined in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	log.Printf("[INFO] deleting the instance %s", d.Id())
	_, err := apiClient.DeleteInstance(d.Id())
	if err != nil {
		return diag.Errorf("[ERR] an error occurred while trying to delete instance %s", d.Id())
	}

	// Wait for the instance to be completely deleted
	deleteStateConf := &retry.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			resp, err := apiClient.GetInstance(d.Id())
			if err != nil {
				if errors.Is(err, civogo.DatabaseInstanceNotFoundError) {
					return 0, "DELETED", nil
				}
				return 0, "", err
			}
			return resp, resp.Status, nil
		},
		Timeout:        60 * time.Minute,
		Delay:          3 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 60,
	}
	_, err = deleteStateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for instance (%s) to be deleted: %s", d.Id(), err)
	}

	return nil
}

func customizeDiffInstance(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	if d.Id() != "" && d.HasChange("script") {
		return fmt.Errorf("the 'script' field is immutable")
	}
	return nil
}

// checkNetworkFirstInstance checks if this is the first instance in a given network
func checkNetworkFirstInstance(apiClient *civogo.Client, networkID string) (bool, error) {
	// List all instances
	instances, err := apiClient.ListAllInstances()
	if err != nil {
		return false, fmt.Errorf("failed to list instances: %v", err)
	}

	// Count instances in the specified network
	networkInstanceCount := 0
	for _, instance := range instances {
		if instance.NetworkID == networkID {
			networkInstanceCount++
		}
	}

	// Return true if this is the first instance in the network
	return networkInstanceCount == 0, nil
}
