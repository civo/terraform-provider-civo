package civo

import (
	"fmt"
	"log"
	"strings"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Data source to get from the api a specific instance
// using the id or the hostname
func dataSourceInstance() *schema.Resource {
	return &schema.Resource{
		Description: strings.Join([]string{
			"Get information on an instance for use in other resources. This data source provides all of the instance's properties as configured on your Civo account.",
			"Note: This data source returns a single instance. When specifying a hostname, an error will be raised if more than one instances found.",
		}, "\n\n"),
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
				Description:  "The hostname of the Instance",
			},
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "The region of an existing Instance",
			},
			// computed attributes
			"reverse_dns": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A fully qualified domain name",
			},
			"size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the size",
			},
			"cpu_cores": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total cpu of the inatance",
			},
			"ram_mb": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total ram of the instance",
			},
			"disk_gb": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of the disk",
			},
			"network_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "his will be the ID of the network",
			},
			"template": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID for the disk image/template to used to build the instance",
			},
			"initial_user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the initial user created on the server",
			},
			"notes": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The notes of the instance",
			},
			"sshkey_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID SSH key",
			},
			"firewall_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the firewall used",
			},
			"tags": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "An optional list of tags",
			},
			"script": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The contents of a script uploaded",
			},
			"initial_password": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance initial password",
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The private IP",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The public IP",
			},
			"pseudo_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Is the ip that is used to route the public ip from the internet to the instance using NAT",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the instance",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date of creation of the instance",
			},
		},
	}
}

func dataSourceInstanceRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	var foundImage *civogo.Instance

	if id, ok := d.GetOk("id"); ok {
		log.Printf("[INFO] Getting the instance by id")
		image, err := apiClient.FindInstance(id.(string))
		if err != nil {
			return fmt.Errorf("[ERR] failed to retrive instance: %s", err)
		}

		foundImage = image
	} else if hostname, ok := d.GetOk("hostname"); ok {
		log.Printf("[INFO] Getting the instance by hostname")
		image, err := apiClient.FindInstance(hostname.(string))
		if err != nil {
			return fmt.Errorf("[ERR] failed to retrive instance: %s", err)
		}

		foundImage = image
	}

	d.SetId(foundImage.ID)
	d.Set("hostname", foundImage.Hostname)
	d.Set("reverse_dns", foundImage.ReverseDNS)
	d.Set("size", foundImage.Size)
	d.Set("cpu_cores", foundImage.CPUCores)
	d.Set("ram_mb", foundImage.RAMMegabytes)
	d.Set("disk_gb", foundImage.DiskGigabytes)
	d.Set("initial_user", foundImage.InitialUser)
	d.Set("initial_password", foundImage.InitialPassword)
	d.Set("sshkey_id", foundImage.SSHKey)
	d.Set("tags", foundImage.Tags)
	d.Set("private_ip", foundImage.PrivateIP)
	d.Set("public_ip", foundImage.PublicIP)
	d.Set("pseudo_ip", foundImage.PseudoIP)
	d.Set("status", foundImage.Status)
	d.Set("region", apiClient.Region)
	d.Set("script", foundImage.Script)
	d.Set("created_at", foundImage.CreatedAt.UTC().String())
	d.Set("notes", foundImage.Notes)

	return nil
}
