package civo

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"strings"
)

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
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "A fully qualified domain name that should be used as the instance's IP's reverse DNS (optional, uses the hostname if unspecified)",
				ValidateFunc: validateName,
			},
			"size": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the size, from the current list, e.g. g2.small (required)",
			},
			"public_ip_requiered": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This should be either false, true or `move_ip_from:intances_id`",
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
			"notes": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Add some notes to the instance",
			},
			"sshkey_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of an already uploaded SSH public key to use for login to the default user (optional; if one isn't provided a random password will be set and returned in the initial_password field)",
			},
			"firewall_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the firewall to use, from the current list. If left blank or not sent, the default firewall will be used (open to all)",
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "An optional list of tags, represented as a key, value pair",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			// Computed resource
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
		Create: resourceInstanceCreate,
		Read:   resourceInstanceRead,
		Update: resourceInstanceUpdate,
		Delete: resourceInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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

	tfTags := d.Get("tags").(*schema.Set).List()
	tags := make([]string, len(tfTags))
	for i, tfTag := range tfTags {
		tags[i] = tfTag.(string)
	}

	config.Tags = tags

	instance, err := apiClient.CreateInstance(config)
	if err != nil {
		fmt.Errorf("[WARN] failed to create instance: %s", err)
		return err
	}

	d.SetId(instance.ID)

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		resp, err := apiClient.GetInstance(instance.ID)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error geting instance: %s", err))
		}

		if resp.Status != "ACTIVE" {
			return resource.RetryableError(fmt.Errorf("[WARN] expected instance to be created but was in state %s", resp.Status))
		} else {
			/*
				Once the instance is created, we check if the object firewall_id,
				if it is, then we set the firewall id to the instances
			*/
			if attr, ok := d.GetOk("firewall_id"); ok {
				_, errInstance := apiClient.SetInstanceFirewall(instance.ID, attr.(string))
				if errInstance != nil {
					return resource.NonRetryableError(fmt.Errorf("[WARN] failed to set firewall to the instance: %s", errInstance))
				}
			}

			/*
				Once the instance is created, we check if the object notes,
				if it is, then we add the note to the instances
			*/
			if attr, ok := d.GetOk("notes"); ok {
				resp.Notes = attr.(string)
				_, errInstance := apiClient.UpdateInstance(resp)
				if errInstance != nil {
					return resource.NonRetryableError(fmt.Errorf("[WARN] failed to set note to the instance: %s", errInstance))
				}
			}
		}

		return resource.NonRetryableError(resourceInstanceRead(d, m))
	})
}

func resourceInstanceRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	resp, err := apiClient.GetInstance(d.Id())
	if err != nil {
		// check if the instance no longer exists.
		log.Printf("[WARN] civo instance (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("hostname", resp.Hostname)
	d.Set("reverse_dns", resp.ReverseDNS)
	d.Set("size", resp.Size)
	d.Set("initial_user", resp.InitialUser)
	d.Set("initial_password", resp.InitialPassword)
	d.Set("sshkey_id", resp.SSHKey)
	d.Set("tags", resp.Tags)
	d.Set("private_ip", resp.PrivateIP)
	d.Set("public_ip", resp.PublicIP)
	d.Set("pseudo_ip", resp.PseudoIP)
	d.Set("status", resp.Status)
	d.Set("created_at", resp.CreatedAt.String())
	d.Set("notes", resp.Notes)

	if _, ok := d.GetOk("network_id"); ok {
		d.Set("network_id", resp.NetworkID)
	}

	if _, ok := d.GetOk("template"); ok {
		d.Set("template", resp.TemplateID)
	}

	if attr, ok := d.GetOk("firewall_id"); ok {
		d.Set("firewall_id", attr.(string))
	}

	return nil
}

func resourceInstanceUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	if d.HasChange("size") {
		newSize := d.Get("size").(string)
		_, err := apiClient.UpgradeInstance(d.Id(), newSize)
		if err != nil {
			log.Printf("[WARN] An error occurred while resizing the instance (%s)", d.Id())
		}

		return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			resp, err := apiClient.GetInstance(d.Id())

			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("[WARN] error geting instance: %s", err))
			}

			if resp.Status != "ACTIVE" {
				return resource.RetryableError(fmt.Errorf("[WARN] expected instance to be resizing but was in state %s", resp.Status))
			}

			return resource.NonRetryableError(resourceInstanceRead(d, m))
		})
	}

	if d.HasChange("notes") {
		notes := d.Get("notes").(string)
		instance, err := apiClient.GetInstance(d.Id())
		if err != nil {
			// check if the instance no longer exists.
			log.Printf("[WARN] civo instance (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		instance.Notes = notes

		_, err = apiClient.UpdateInstance(instance)
		if err != nil {
			log.Printf("[WARN] an error occurred while adding a note to the instance (%s)", d.Id())
		}
	}

	if d.HasChange("firewall_id") {
		firewallID := d.Get("firewall_id").(string)
		_, err := apiClient.SetInstanceFirewall(d.Id(), firewallID)
		if err != nil {
			// check if the instance no longer exists.
			log.Printf("[WARN] an error occurred while set firewall to the instance (%s)", d.Id())
			d.SetId("")
			return nil
		}
	}

	if d.HasChange("tags") {
		tfTags := d.Get("tags").(*schema.Set).List()
		tags := make([]string, len(tfTags))
		for i, tfTag := range tfTags {
			tags[i] = tfTag.(string)
		}

		instance, err := apiClient.GetInstance(d.Id())
		if err != nil {
			// check if the instance no longer exists.
			log.Printf("[WARN] civo instance (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		tagsToString := strings.Join(tags, " ")

		_, err = apiClient.SetInstanceTags(instance, tagsToString)
		if err != nil {
			log.Printf("[WARN] an error occurred while adding tags to the instance (%s)", d.Id())
		}

	}

	return resourceInstanceRead(d, m)
}

func resourceInstanceDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	_, err := apiClient.DeleteInstance(d.Id())
	if err != nil {
		log.Printf("[INFO] civo instance (%s) was delete", d.Id())
	}
	return nil
}
