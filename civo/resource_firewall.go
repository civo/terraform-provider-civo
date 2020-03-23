package civo

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceFirewall() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateName,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Create: resourceFirewallCreate,
		Read:   resourceFirewallRead,
		Update: resourceFirewallUpdate,
		Delete: resourceFirewallDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceFirewallCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	firewall, err := apiClient.NewFirewall(d.Get("name").(string))
	if err != nil {
		fmt.Errorf("failed to create a new firewall: %s", err)
		return err
	}

	d.SetId(firewall.ID)

	return resourceFirewallRead(d, m)
}

func resourceFirewallRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	resp, err := apiClient.FindFirewall(d.Id())
	if err != nil {
		if resp != nil {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("[ERR] error retrieving firewall: %s", err)
	}

	d.Set("name", resp.Name)
	d.Set("region", resp.Region)

	return nil
}

func resourceFirewallUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	if d.HasChange("name") {
		if d.Get("name").(string) != "" {
			_, err := apiClient.RenameFirewall(d.Id(), d.Get("name").(string))
			if err != nil {
				log.Printf("[WARN] an error occurred while trying to rename the firewall (%s)", d.Id())
			}
		}
	}

	return resourceFirewallRead(d, m)
}

func resourceFirewallDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	_, err := apiClient.DeleteFirewall(d.Id())
	if err != nil {
		log.Printf("[INFO] civo firewall (%s) was delete", d.Id())
	}
	return nil
}
