package civo

import (
	"fmt"
	"log"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Firewall resource with this we can create and manage all firewall
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

// function to create a firewall
func resourceFirewallCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] creating a new firewall %s", d.Get("name").(string))
	firewall, err := apiClient.NewFirewall(d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("[ERR] failed to create a new firewall: %s", err)
	}

	d.SetId(firewall.ID)

	return resourceFirewallRead(d, m)
}

// function to read a firewall
func resourceFirewallRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] retriving the firewall %s", d.Id())
	resp, err := apiClient.FindFirewall(d.Id())
	if err != nil {
		if resp == nil {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("[ERR] error retrieving firewall: %s", err)
	}

	d.Set("name", resp.Name)
	d.Set("region", resp.Region)

	return nil
}

// function to update the firewall
func resourceFirewallUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	if d.HasChange("name") {
		if d.Get("name").(string) != "" {
			log.Printf("[INFO] updating the firewall name, %s", d.Id())
			_, err := apiClient.RenameFirewall(d.Id(), d.Get("name").(string))
			if err != nil {
				return fmt.Errorf("[WARN] an error occurred while tring to rename the firewall %s, %s", d.Id(), err)
			}
		}
	}

	return resourceFirewallRead(d, m)
}

// function to delete a firewall
func resourceFirewallDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] deleting the firewall %s", d.Id())
	_, err := apiClient.DeleteFirewall(d.Id())
	if err != nil {
		return fmt.Errorf("[ERR] an error occurred while tring to delete the firewall %s, %s", d.Id(), err)
	}
	return nil
}
