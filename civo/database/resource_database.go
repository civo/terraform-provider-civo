package database

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceDatabase The Database resource represents an Database object
// and with it you can handle the Database created with Terraform.
func ResourceDatabase() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "Name of the database",
			},
			"size": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "Size of the database",
			},
			"engine": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The engine of the database",
				ValidateFunc: validation.NoZeroValues,
			},
			"version": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The version of the database",
				ValidateFunc: validation.NoZeroValues,
			},
			"network_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The id of the associated network",
				ForceNew:    true,
			},
			"nodes": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "Count of nodes",
			},
			"firewall_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the firewall to use, from the current list. If left blank or not sent, the default firewall will be used (open to all)",
			},
			"region": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				Description:      "The region where the database will be created.",
				DiffSuppressFunc: utils.IgnoreCaseDiff,
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The username of the database",
			},
			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The password of the database",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The endpoint of the database",
			},
			"dns_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The DNS endpoint of the database",
			},
			"port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The port of the database",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the database",
			},
			"private_ipv4": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The private IPv4 address for the database",
			},
			"storage_gb": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "Storage of every node in GB. It can be grown beyond the plan's base storage but never shrunk; the GB above `base_storage_gb` are billed as a storage add-on.",
			},
			"base_storage_gb": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Storage in GB included in the plan (size) of the database",
			},
			"extensions": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "PostgreSQL extensions to enable, replacing the current set. An empty list disables them all; removing the attribute keeps whatever is enabled. Extensions that need it trigger a brief rolling restart when applied.",
			},
			"database_parameters": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Custom PostgreSQL server parameters as key/value pairs, replacing the current set. An empty map resets them to the defaults; removing the attribute keeps the current ones.",
			},
			"source_database_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "ID of an existing database to fork. The new database is an independent copy built from the source's backups: by default as of its latest backup, or as of `source_target_time`. `engine` and `version` must match the source. Not returned by the API, so an imported fork plans a replacement when this is set.",
			},
			"source_target_time": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"source_database_id"},
				ValidateFunc: validation.IsRFC3339Time,
				Description:  "Point in time of `source_database_id` to fork from, as RFC 3339 (e.g. 2026-10-21T10:00:00Z). It must not be earlier than the source's `earliest_restore_time`. Leave unset to fork the source as of its latest backup.",
			},
			"read_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Host of the load-balanced read-only endpoint, available on databases with 3 or more nodes",
			},
			"read_endpoint_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Port of the read-only endpoint",
			},
			"read_endpoint_dns": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNS entry of the read-only endpoint",
			},
			"earliest_restore_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Earliest point in time (RFC 3339) the database can be restored or forked to. Empty until the first full backup has completed.",
			},
		},
		CustomizeDiff: customizeDiffDatabase,
		CreateContext: resourceDatabaseCreate,
		ReadContext:   resourceDatabaseRead,
		UpdateContext: resourceDatabaseUpdate,
		DeleteContext: resourceDatabaseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

	apiClient, err := utils.RegionalClient(apiClient, utils.ResolveRegion(d, utils.NetworkRef("network_id"), utils.FirewallRef("firewall_id")))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] configuring the database %s", d.Get("name").(string))

	config := &civogo.CreateDatabaseRequest{
		Name:            d.Get("name").(string),
		Software:        d.Get("engine").(string),
		SoftwareVersion: d.Get("version").(string),
		Region:          apiClient.Region,
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

	var database *civogo.Database
	if _, forking := d.GetOk("source_database_id"); forking {
		forkReq, err := forkRequestFromResource(d, apiClient.Region, config.NetworkID, config.FirewallID)
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("[INFO] forking the Database %s from %s", forkReq.Name, forkReq.SourceDatabaseID)
		database, err = apiClient.ForkDatabase(forkReq)
		if err != nil {
			return diag.Errorf("[ERR] failed to fork Database: %s", err)
		}
	} else {
		log.Printf("[INFO] creating the Database %s", d.Get("name").(string))
		database, err = apiClient.NewDatabase(config)
		if err != nil {
			return diag.Errorf("[ERR] failed to create Database: %s", err)
		}
	}

	d.SetId(database.ID)

	createStateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Ready"},
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

	// Storage, extensions and parameters are not part of the create call:
	// apply them once the database is Ready, reading the planned values
	// before Read overwrites them with the API's. A fork already carries
	// storage and extensions, so only the parameters are left for it.
	if _, forking := d.GetOk("source_database_id"); forking {
		if params, ok := d.GetOk("database_parameters"); ok && len(params.(map[string]interface{})) > 0 {
			if _, err := apiClient.UpdateDatabaseParameters(d.Id(), expandStringMap(params)); err != nil {
				return diag.Errorf("[ERR] failed to set the parameters of Database %s: %s", d.Id(), err)
			}
			if err := waitForDatabaseReady(ctx, apiClient, d.Id()); err != nil {
				return diag.FromErr(err)
			}
		}
	} else if err := applyDatabaseSettings(ctx, d, apiClient, d.Id(), true); err != nil {
		return diag.Errorf("[ERR] failed to configure Database %s: %s", d.Id(), err)
	}

	return resourceDatabaseRead(ctx, d, m)
}

// Function to Update the database
func resourceDatabaseUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	apiClient, err := utils.RegionalClient(apiClient, utils.ResolveRegion(d))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = apiClient.FindDatabase(d.Id())
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

	if d.HasChanges("nodes", "name", "firewall_id") {
		log.Printf("[INFO] updating the Database %s", d.Id())
		_, err = apiClient.UpdateDatabase(d.Id(), config)
		if err != nil {
			return diag.Errorf("[ERR] failed to update Database: %s", err)
		}
		if d.HasChange("nodes") {
			if err := waitForDatabaseReady(ctx, apiClient, d.Id()); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChanges("storage_gb", "extensions", "database_parameters") {
		if err := applyDatabaseSettings(ctx, d, apiClient, d.Id(), false); err != nil {
			return diag.Errorf("[ERR] failed to update Database %s: %s", d.Id(), err)
		}
	}

	return resourceDatabaseRead(ctx, d, m)
}

// Function to Read the database
func resourceDatabaseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	apiClient, err := utils.RegionalClient(apiClient, utils.ResolveRegion(d))
	if err != nil {
		return diag.FromErr(err)
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
	d.Set("engine", resp.Software)
	d.Set("version", resp.SoftwareVersion)
	d.Set("network_id", resp.NetworkID)
	d.Set("firewall_id", resp.FirewallID)
	d.Set("region", apiClient.Region)
	d.Set("username", resp.Username)
	d.Set("password", resp.Password)
	d.Set("endpoint", resp.PublicIPv4)
	d.Set("dns_endpoint", fmt.Sprintf("%s.db.civo.com", resp.ID))
	d.Set("port", resp.Port)
	d.Set("status", resp.Status)
	d.Set("private_ipv4", resp.PrivateIPv4)

	setDatabaseDetails(d, apiClient, d.Id())

	return nil
}

// Function to delete the database
func resourceDatabaseDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	apiClient, err := utils.RegionalClient(apiClient, utils.ResolveRegion(d))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] deleting the Database %s", d.Id())
	_, err = apiClient.DeleteDatabase(d.Id())
	if err != nil {
		return diag.Errorf("[ERR] an error occurred while trying to delete the Database %s", d.Id())
	}

	// Wait for the database to be completely gone. Deleting it is asynchronous,
	// so returning early races the teardown of the network and firewall it
	// references, which the API then refuses to delete as still in use.
	deleteStateConf := &resource.StateChangeConf{
		Pending: []string{"exists"},
		Target:  []string{"deleted"},
		Refresh: func() (interface{}, string, error) {
			resp, err := apiClient.GetDatabase(d.Id())
			if err != nil {
				if isDatabaseNotFound(err) {
					return 0, "deleted", nil
				}
				return 0, "", err
			}
			return resp, "exists", nil
		},
		Timeout:        60 * time.Minute,
		Delay:          3 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 60,
	}
	_, err = deleteStateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for Database (%s) to be deleted: %s", d.Id(), err)
	}

	return nil
}

// isDatabaseNotFound reports whether err from GetDatabase means the database no
// longer exists. civogo has no dedicated error for the API's 404 here, so the
// error code has to be matched directly.
func isDatabaseNotFound(err error) bool {
	return errors.Is(err, civogo.ZeroMatchesError) ||
		utils.APIErrorCode(err) == "database_database_find"
}
