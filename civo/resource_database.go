package civo

import (
	"context"
	"log"
	"time"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// The Database resource represents an Database object
// and with it you can handle the Database created with Terraform.
func resourceDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabaseCreate,
		ReadContext:   resourceDatabaseRead,
		UpdateContext: resourceDatabaseUpdate,
		DeleteContext: resourceDatabaseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "Name of the database",
			},

			"nodes": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "Count of nodes",
			},

			"size": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "Size of the database",
			},

			"network_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The id of the associated network",
			},

			"firewall_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the firewall to use, from the current list. If left blank or not sent, the default firewall will be used (open to all)",
			},

			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region where the database will be created.",
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
	}
}

// function to create a database
func resourceDatabaseCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if it is defined in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	log.Printf("[INFO] configuring the database %s", d.Get("name").(string))

	config := &civogo.CreateDatabaseRequest{
		Name:   d.Get("name").(string),
		Region: apiClient.Region,
	}

	if attr, ok := d.GetOk("nodes"); ok {
		config.Nodes = attr.(int)
	}

	if attr, ok := d.GetOk("size"); ok {
		config.Size = attr.(string)
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

	if attr, ok := d.GetOk("firewall_id"); ok {
		firewallID := attr.(string)
		firewall, err := apiClient.FindFirewall(firewallID)
		if err != nil {
			return diag.Errorf("[ERR] unable to find firewall - %s", err)
		}

		if firewall.NetworkID != config.NetworkID {
			return diag.Errorf("[ERR] firewall %s is not part of network %s", firewall.ID, config.NetworkID)
		}

		config.FirewallID = firewallID
	}

	log.Printf("[INFO] creating the Database %s", d.Get("name").(string))
	database, err := apiClient.NewDatabase(config)
	if err != nil {
		return diag.Errorf("[ERR] failed to create Database: %s", err)
	}

	d.SetId(database.ID)

	createStateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"ready"},
		Refresh: func() (interface{}, string, error) {
			resp, err := apiClient.GetDatabase(d.Id())
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
		return diag.Errorf("error waiting for Database (%s) to be created: %s", d.Id(), err)
	}

	return resourceDatabaseRead(ctx, d, m)
}

// Function to Update the database
func resourceDatabaseUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if it is defined in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	_, err := apiClient.FindDatabase(d.Id())
	if err != nil {
		return diag.Errorf("[ERR] failed to find Database: %s", err)
	}

	config := &civogo.UpdateDatabaseRequest{
		Region: apiClient.Region,
	}

	if d.HasChange("nodes") {
		nodes := d.Get("nodes").(int)
		config.Nodes = &nodes
	}

	if d.HasChange("name") {
		name := d.Get("name").(string)
		config.Name = name
	}

	if d.HasChange("firewall_id") {
		firewallID := d.Get("firewall_id").(string)
		config.FirewallID = firewallID
	}

	log.Printf("[INFO] updating the Database %s", d.Id())
	_, err = apiClient.UpdateDatabase(d.Id(), config)
	if err != nil {
		return diag.Errorf("[ERR] failed to update Database: %s", err)
	}

	return resourceDatabaseRead(ctx, d, m)
}

// Function to Read the database
func resourceDatabaseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if it is defined in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	log.Printf("[INFO] retriving the Database %s", d.Id())
	resp, err := apiClient.GetDatabase(d.Id())
	if err != nil {
		if resp == nil {
			d.SetId("")
			return nil
		}
		return diag.Errorf("[ERR] failed to retrive the Database: %s", err)
	}

	d.Set("name", resp.Name)
	d.Set("size", resp.Size)
	d.Set("nodes", resp.Nodes)
	d.Set("network_id", resp.NetworkID)
	d.Set("firewall_id", resp.FirewallID)

	return nil
}

// Function to delete the database
func resourceDatabaseDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if it is defined in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	log.Printf("[INFO] deleting the Database %s", d.Id())
	_, err := apiClient.DeleteDatabase(d.Id())
	if err != nil {
		return diag.Errorf("[ERR] an error occurred while trying to delete the Database %s", d.Id())
	}
	return nil
}
