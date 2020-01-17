package provider

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
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
			"hostname": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "A fully qualified domain name that should be set as the instance's hostname (required)",
				ForceNew:     true,
				ValidateFunc: validateName,
			},
			"reverse_dns": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A fully qualified domain name that should be used as the instance's IP's reverse DNS (optional, uses the hostname if unspecified)",
			},
			"size": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the size, from the current list, e.g. g2.small (required)",
			},
			"public_ip_requiered": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This should be either none, create",
			},
			"network_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This must be the ID of the network from the network listing (optional; default network used when not specified)",
			},
			"template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID for the template to use to build the instance",
			},
			"initial_user": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the initial user created on the server (optional; this will default to the template's default_username and fallback to civo)",
			},
			"sshkey_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of an already uploaded SSH public key to use for login to the default user (optional; if one isn't provided a random password will be set and returned in the initial_password field)",
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional tags",
			},
			//"tags": {
			//	Type:        schema.TypeSet,
			//	Optional:    true,
			//	Description: "An optional list of tags, represented as a key, value pair",
			//	Elem:        &schema.Schema{Type: schema.TypeString},
			//},
			// Computed resource
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pseudo_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
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

	config.Hostname = d.Get("hostname").(string)

	if attr, ok := d.GetOk("reverse_dns"); ok {
		config.ReverseDNS = attr.(string)
	}

	if attr, ok := d.GetOk("size"); ok {
		config.Size = attr.(string)
	}

	if attr, ok := d.GetOk("size"); ok {
		config.Size = attr.(string)
	}

	if attr, ok := d.GetOk("network_id"); ok {
		config.NetworkID = attr.(string)
	}

	if attr, ok := d.GetOk("template"); ok {
		config.TemplateID = attr.(string)
	}

	if attr, ok := d.GetOk("initial_user"); ok {
		config.InitialUser = attr.(string)
	}

	if attr, ok := d.GetOk("sshkey_id"); ok {
		config.SSHKeyID = attr.(string)
	}

	if attr, ok := d.GetOk("tags"); ok {
		//config.Tags = attr.(*schema.Set).List()
		config.Tags = attr.(string)
	}

	instance, err := apiClient.CreateInstance(config)
	if err != nil {
		fmt.Errorf("failed to create instance: %s", err)
		return err
	}

	d.SetId(instance.ID)

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		resp, err := apiClient.GetInstance(instance.ID)

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error geting instance: %s", err))
		}

		if resp.Status != "ACTIVE" {
			return resource.RetryableError(fmt.Errorf("expected instance to be created but was in state %s", resp.Status))
		}

		return resource.NonRetryableError(resourceInstanceRead(d, m))
	})
}

func resourceInstanceRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	resp, err := apiClient.GetInstance(d.Id())
	if err != nil {
		// check if the droplet no longer exists.
		log.Printf("[WARN] Civo instance (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("hostname", resp.Hostname)
	d.Set("reverse_dns", resp.ReverseDNS)
	d.Set("size", resp.Size)
	d.Set("network_id", resp.NetworkID)
	d.Set("template", resp.TemplateID)
	d.Set("initial_user", resp.InitialUser)
	d.Set("sshkey_id", resp.SSHKey)
	d.Set("tags", resp.Tags)
	d.Set("private_ip", resp.PrivateIP)
	d.Set("public_ip", resp.PublicIP)
	d.Set("pseudo_ip", resp.PseudoIP)
	d.Set("status", resp.Status)
	d.Set("created_at", resp.CreatedAt)

	return nil
}

func resourceInstanceUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceInstanceRead(d, m)
}

func resourceInstanceDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	_, err := apiClient.DeleteInstance(d.Id())
	if err != nil {
		log.Printf("[INFO] Civo instance (%s) was delete", d.Id())
	}
	return nil
}
