package civo

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
)

// Data source to get from the api a specific template
// using the code of the image
func dataSourceInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceInstanceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "hostname"},
			},
			"hostname": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "hostname"},
			},
			// computed attributes
			"reverse_dns": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"initial_user": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sshkey_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"firewall_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"script": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"initial_password": {
				Type:     schema.TypeString,
				Computed: true,
			},
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
	}
}

func dataSourceInstanceRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	var foundImage *civogo.Instance

	if id, ok := d.GetOk("id"); ok {
		log.Printf("[INFO] Getting the instance by id")
		image, err := apiClient.FindInstance(id.(string))
		if err != nil {
			fmt.Errorf("[ERR] failed to retrive instance: %s", err)
			return err
		}

		foundImage = image
	} else if hostname, ok := d.GetOk("hostname"); ok {
		log.Printf("[INFO] Getting the instance by hostname")
		image, err := apiClient.FindInstance(hostname.(string))
		if err != nil {
			fmt.Errorf("[ERR] failed to retrive instance: %s", err)
			return err
		}

		foundImage = image
	}

	d.SetId(foundImage.ID)
	d.Set("hostname", foundImage.Hostname)
	d.Set("reverse_dns", foundImage.ReverseDNS)
	d.Set("size", foundImage.Size)
	d.Set("initial_user", foundImage.InitialUser)
	d.Set("initial_password", foundImage.InitialPassword)
	d.Set("sshkey_id", foundImage.SSHKey)
	d.Set("tags", foundImage.Tags)
	d.Set("private_ip", foundImage.PrivateIP)
	d.Set("public_ip", foundImage.PublicIP)
	d.Set("pseudo_ip", foundImage.PseudoIP)
	d.Set("status", foundImage.Status)
	d.Set("script", foundImage.Script)
	d.Set("created_at", foundImage.CreatedAt.UTC().String())
	d.Set("notes", foundImage.Notes)

	return nil
}
