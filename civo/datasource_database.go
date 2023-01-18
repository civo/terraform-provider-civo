package civo

import (
	"context"
	"log"
	"strings"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Data source to get from the api a specific Database
// using the id or the name
func dataSourceDatabase() *schema.Resource {
	return &schema.Resource{
		Description: strings.Join([]string{
			"Get information of an Database for use in other resources. This data source provides all of the Database's properties as configured on your Civo account.",
			"Note: This data source returns a single Database. When specifying a name, an error will be raised if more than one Databases with the same name found.",
		}, "\n\n"),
		ReadContext: dataSourceDatabaseRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the Database",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "The name of the Database",
			},
			"size": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "Size of the database",
			},
			"nodes": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "Count of nodes",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region of an existing Object Store",
			},
			"network_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The network id of the Database",
			},
			"firewall_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The secret access key of the Database",
			},
		},
	}
}

func dataSourceDatabaseRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	var foundDatabase *civogo.Database

	if name, ok := d.GetOk("name"); ok {
		log.Printf("[INFO] Getting the Database by name")
		database, err := apiClient.FindDatabase(name.(string))
		if err != nil {
			return diag.Errorf("[ERR] failed to retrive Database: %s", err)
		}

		foundDatabase = database
	}

	if id, ok := d.GetOk("id"); ok {
		log.Printf("[INFO] Getting the Database by name")
		database, err := apiClient.FindDatabase(id.(string))
		if err != nil {
			return diag.Errorf("[ERR] failed to retrive Database: %s", err)
		}

		foundDatabase = database
	}

	d.SetId(foundDatabase.ID)
	d.Set("name", foundDatabase.Name)
	d.Set("region", apiClient.Region)
	d.Set("size", foundDatabase.Size)
	d.Set("nodes", foundDatabase.Nodes)
	d.Set("network_id", foundDatabase.NetworkID)
	d.Set("firewall_id", foundDatabase.FirewallID)

	return nil
}
