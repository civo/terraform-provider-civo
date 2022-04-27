package civo

import (
	"context"
	"log"
	"time"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// The application resource represents a CivoApp object
// and with it you can handle the applications created with Terraform.
func resourceApplication() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a CivoApp resource. This can be used to create, modify, and delete applications.",
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region for the application, if not declared we use the region as declared in the provider (Defaults to LON1)",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: utils.ValidateNameSize,
				Description:  "The app name. Must be unique.",
			},
			"size": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "small",
				Description: "The name of the size, from the current list, e.g. small, medium",
			},
			"network_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "This must be the ID of the network from the network listing (optional; default network used when not specified)",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the application",
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
			"process_info": processInfoSchema(),
			"config":       configSchema(),
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the application",
			},
		},
		CreateContext: resourceAppCreate,
		ReadContext:   resourceAppRead,
		UpdateContext: resourceAppUpdate,
		DeleteContext: resourceAppDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},
	}
}

// schema for process info in the application
func processInfoSchema() *schema.Schema {
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
func configSchema() *schema.Schema {
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

// function to create an application
func resourceAppCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	log.Printf("[INFO] configuring the application %s", d.Get("name").(string))
	config, err := apiClient.NewApplicationConfig()
	if err != nil {
		return diag.Errorf("[ERR] failed to create a new config: %s", err)
	}

	if name, ok := d.GetOk("name"); ok {
		config.Name = name.(string)
	} else {
		config.Name = utils.RandomName()
	}

	if attr, ok := d.GetOk("size"); ok {
		config.Size = attr.(string)
	}

	if networkID, ok := d.GetOk("network_id"); ok {
		config.NetworkID = networkID.(string)
	} else {
		defaultNetwork, err := apiClient.GetDefaultNetwork()
		if err != nil {
			return diag.Errorf("[ERR] failed to get the default network: %s", err)
		}
		config.NetworkID = defaultNetwork.ID
	}

	log.Printf("[INFO] creating the application %s", d.Get("name").(string))
	app, err := apiClient.CreateApplication(config)
	if err != nil {
		return diag.Errorf("[ERR] failed to create application: %s", err)
	}

	d.SetId(app.ID)

	createStateConf := &resource.StateChangeConf{
		Pending: []string{"BUILDING"},
		Target:  []string{"ACTIVE"},
		Refresh: func() (interface{}, string, error) {
			resp, err := apiClient.GetApplication(d.Id())
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
		return diag.Errorf("error waiting for app (%s) to be created: %s", d.Id(), err)
	}

	return resourceAppRead(ctx, d, m)

}

// function to read the application
func resourceAppRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if it is defined in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	log.Printf("[INFO] retriving the application %s", d.Id())
	resp, err := apiClient.GetApplication(d.Id())
	if err != nil {
		if resp == nil {
			d.SetId("")
			return nil
		}

		return diag.Errorf("[ERR] failed to retriving the application: %s", err)
	}

	d.Set("name", resp.Name)
	d.Set("size", resp.Size)
	d.Set("network_id", resp.NetworkID)
	d.Set("ssh_key_ids", resp.SSHKeyIDs)
	d.Set("status", resp.Status)
	d.Set("description", resp.Description)
	d.Set("config", resp.Config)
	d.Set("domains", resp.Domains)
	d.Set("process_info", resp.ProcessInfo)

	if err := d.Set("process_info", flattenProcesses(resp.ProcessInfo)); err != nil {
		return diag.Errorf("[ERR] error retrieving the processes for the application error: %#v", err)
	}

	if err := d.Set("config", flattenEnvVar(resp.Config)); err != nil {
		return diag.Errorf("[ERR] error retrieving the config for the application error: %#v", err)
	}

	return nil
}

// function to update the application
func resourceAppUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if it is defined in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	config := &civogo.UpdateApplicationRequest{}

	if d.HasChange("name") {
		config.Name = d.Get("name").(string)
	}

	if d.HasChange("size") {
		config.Size = d.Get("size").(string)
	}

	if d.HasChange("advanced") {
		config.Advanced = d.Get("advanced").(bool)
	}

	if d.HasChange("image") {
		config.Image = d.Get("image").(string)
	}

	if d.HasChange("description") {
		config.Description = d.Get("description").(string)
	}

	log.Printf("[INFO] updating the application %s", d.Id())
	_, err := apiClient.UpdateApplication(d.Id(), config)
	if err != nil {
		return diag.Errorf("[ERR] failed to update application: %s", err)
	}

	return resourceAppRead(ctx, d, m)
}

// function to delete application
func resourceAppDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if it is defined in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	log.Printf("[INFO] deleting the application %s", d.Id())
	_, err := apiClient.DeleteApplication(d.Id())
	if err != nil {
		return diag.Errorf("[ERR] an error occurred while tring to delete application %s", d.Id())
	}
	return nil
}

// function to flatten all processes inside the application
func flattenProcesses(processes []civogo.ProcessInfo) []interface{} {
	if processes == nil {
		return nil
	}

	flattenedProcess := make([]interface{}, 0)
	for _, process := range processes {
		rawProcess := map[string]interface{}{
			"process_type":  process.ProcessType,
			"process_count": process.ProcessCount,
		}

		flattenedProcess = append(flattenedProcess, rawProcess)
	}

	return flattenedProcess
}

// function to flatten all config variables inside the application
func flattenEnvVar(variables []civogo.EnvVar) []interface{} {
	if variables == nil {
		return nil
	}

	flattenedVariable := make([]interface{}, 0)
	for _, variable := range variables {
		rawVariable := map[string]interface{}{
			"name":  variable.Name,
			"value": variable.Value,
		}

		flattenedVariable = append(flattenedVariable, rawVariable)
	}

	return flattenedVariable
}
