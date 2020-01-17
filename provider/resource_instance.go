package provider

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"regexp"
)

func validateName(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("Expected name to be string"))
		return warns, errs
	}
	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(value)) {
		errs = append(errs, fmt.Errorf("name cannot contain whitespace. Got %s", value))
		return warns, errs
	}
	return warns, errs
}

func resourceInstance() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the resource, also acts as it's unique ID",
				ForceNew:     true,
				ValidateFunc: validateName,
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A description of an item",
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "An optional list of tags, represented as a key, value pair",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Create: resourceInstanceCreate,
		Read:   resourceInstanceRead,
		Update: resourceInstanceUpdate,
		Delete: resourceInstanceDelete,
		//Exists: resourceExistsItem,
		//Importer: &schema.ResourceImporter{
		//	State: schema.ImportStatePassthrough,
		//},
	}
}

func resourceInstanceCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	config, err := apiClient.NewInstanceConfig()
	if err != nil {
		fmt.Errorf("failed to create a new config: %s", err)
		return err
	}

	config.Hostname = d.Get("name").(string)

	_, err = apiClient.CreateInstance(config)
	if err != nil {
		fmt.Errorf("failed to create instance: %s", err)
		return err
	}

	d.SetId(config.Hostname)
	return nil
}

func resourceInstanceRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceInstanceUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceInstanceRead(d, m)
}

func resourceInstanceDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
