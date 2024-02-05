package ip

import (
	"context"
	"log"
	"strings"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// DataSourceReservedIP function returns a schema.Resource that represents a reserved IP.
// This can be used to query and retrieve details about a specific reserved IP in the infrastructure.
func DataSourceReservedIP() *schema.Resource {
	return &schema.Resource{
		Description: strings.Join([]string{
			"Get information on a reserved IP. This data source provides the region and Instance id as configured on your Civo account.",
			"This is useful if the reserved IP in question is not managed by Terraform or you need to find the instance the IP is attached to.",
			"An error will be raised if the provided domain name is not in your Civo account.",
		}, "\n\n"),
		Schema: map[string]*schema.Schema{
			// Computed resource
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
				Description:  "ID for the ip address",
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "Name for the ip address",
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			// Computed resource
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP Address requested",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The region the ip address is in",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the instance the IP is attached to",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the instance the IP is attached to",
			},
		},
		ReadContext: dataSourceReservedIPRead,
	}
}

// function to read a the IP resource
func dataSourceReservedIPRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] retriving the ip address %s", d.Id())

	var foundIP *civogo.IP

	if id, ok := d.GetOk("id"); ok {
		log.Printf("[INFO] Getting the ip by id")
		resp, err := apiClient.FindIP(id.(string))
		if err != nil {
			return diag.Errorf("[ERR] failed to retrive ip: %s", err)
		}

		foundIP = resp
	} else if name, ok := d.GetOk("name"); ok {
		log.Printf("[INFO] Getting the ip by name")
		resp, err := apiClient.FindIP(name.(string))
		if err != nil {
			return diag.Errorf("[ERR] failed to retrive ip: %s", err)
		}

		foundIP = resp
	}

	d.SetId(foundIP.ID)
	d.Set("name", foundIP.Name)
	d.Set("region", apiClient.Region)
	d.Set("ip", foundIP.IP)

	if foundIP.AssignedTo.ID != "" {
		d.Set("instance_id", foundIP.AssignedTo.ID)
		d.Set("instance_name", foundIP.AssignedTo.Name)
	}

	return nil
}
