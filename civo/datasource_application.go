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

// Data source to get from the api a specific Application
// using the id or the name
func dataSourceApplication() *schema.Resource {
	return &schema.Resource{
		Description: strings.Join([]string{
			"Get information on an application for use in other resources. This data source provides all of the application's properties as configured on your Civo account.",
			"Note: This data source returns a single application. When specifying a name, an error will be raised if more than one applications with the same name found.",
		}, "\n\n"),
		ReadContext: dataSourceAppRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the Application",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				//ExactlyOneOf: []string{"id", "name"},
				Description: "The name of the Application",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region of an existing Application",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of an existing Application",
			},
			"domains": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Space separated list of application domains",
			},
			"ssh_key_ids": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Space separated list of SSH key IDs",
			},
			"process_info": dataSourceProcessInfoSchema(),
			"config":       dataSourceConfigSchema(),
			"size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the size",
			},
			"network_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This will be the ID of the network",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the application",
			},
		},
	}
}

// schema for process info in the application
func dataSourceProcessInfoSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"process_type": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The process type",
				},
				"process_count": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "The process count",
				},
			},
		},
	}
}

// schema for the application config
func dataSourceConfigSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The environment variable name",
				},
				"value": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The environment variable value",
				},
			},
		},
	}
}

func dataSourceAppRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	var foundApp *civogo.Application

	if name, ok := d.GetOk("name"); ok {
		log.Printf("[INFO] Getting the application by iname")
		app, err := apiClient.FindApplication(name.(string))
		if err != nil {
			return diag.Errorf("[ERR] failed to retrive application: %s", err)
		}

		foundApp = app
	}

	d.SetId(foundApp.ID)
	d.Set("name", foundApp.Name)
	d.Set("size", foundApp.Size)
	d.Set("network_id", foundApp.NetworkID)
	d.Set("ssh_key_ids", foundApp.SSHKeyIDs)
	d.Set("status", foundApp.Status)
	d.Set("description", foundApp.Description)
	d.Set("config", foundApp.Config)
	d.Set("domains", foundApp.Domains)
	d.Set("process_info", foundApp.ProcessInfo)

	if err := d.Set("process_info", flattenProcesses(foundApp.ProcessInfo)); err != nil {
		return diag.Errorf("[ERR] error retrieving the processes for the application error: %#v", err)
	}

	if err := d.Set("config", flattenEnvVar(foundApp.Config)); err != nil {
		return diag.Errorf("[ERR] error retrieving the config for the application error: %#v", err)
	}

	return nil
}
