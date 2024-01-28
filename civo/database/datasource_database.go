package database

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// DataSourceDatabase Data source to get from the api a specific Database
// using the id or the name
func DataSourceDatabase() *schema.Resource {
	return &schema.Resource{
		Description: strings.Join([]string{
			"Get information of an Database for use in other resources. This data source provides all of the Database's properties as configured on your Civo account.",
			"Note: This data source returns a single Database. When specifying a name, an error will be raised if more than one Databases with the same name found.",
		}, "\n\n"),
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
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Size of the database",
			},
			"engine": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The engine of the database",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the database",
			},
			"nodes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Count of nodes",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region of an existing Database",
			},
			"network_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The network id of the Database",
			},
			"firewall_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The firewall id of the Database",
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The username of the database",
			},
			"password": {
				Type:        schema.TypeString,
				Computed:    true,
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
		},
		ReadContext: dataSourceDatabaseRead,
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
		log.Printf("[INFO] Getting the Database by id")
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
	d.Set("engine", foundDatabase.Software)
	d.Set("version", foundDatabase.SoftwareVersion)
	d.Set("nodes", foundDatabase.Nodes)
	d.Set("network_id", foundDatabase.NetworkID)
	d.Set("firewall_id", foundDatabase.FirewallID)
	d.Set("username", foundDatabase.Username)
	d.Set("password", foundDatabase.Password)
	d.Set("endpoint", foundDatabase.PublicIPv4)
	d.Set("dns_endpoint", fmt.Sprintf("%s.db.civo.com", foundDatabase.ID))
	d.Set("port", foundDatabase.Port)
	d.Set("status", foundDatabase.Status)

	return nil
}
