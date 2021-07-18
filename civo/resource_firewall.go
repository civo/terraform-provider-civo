package civo

import (
	"fmt"
	"log"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Firewall resource with this we can create and manage all firewall
func resourceFirewall() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: utils.ValidateName,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
		Create: resourceFirewallCreate,
		Read:   resourceFirewallRead,
		Update: resourceFirewallUpdate,
		Delete: resourceFirewallDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// function to create a firewall
func resourceFirewallCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)
	networkID := ""

	if attr, ok := d.GetOk("region"); ok {
		apiClient.Region = attr.(string)
	}

	if attr, ok := d.GetOk("network_id"); ok {
		networkID = attr.(string)
	} else {
		network, err := apiClient.GetDefaultNetwork()
		if err != nil {
			return fmt.Errorf("[ERR] failed to get the default network: %s", err)
		}
		networkID = network.ID
	}

	log.Printf("[INFO] creating a new firewall %s", d.Get("name").(string))
	firewall, err := apiClient.NewFirewall(d.Get("name").(string), networkID)
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
	d.Set("network_id", resp.NetworkID)

	return nil
}

// function to update the firewall
func resourceFirewallUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	if d.HasChange("name") {
		if d.Get("name").(string) != "" {
			firewall := civogo.FirewallConfig{
				Name: d.Get("name").(string),
			}
			log.Printf("[INFO] updating the firewall name, %s", d.Id())
			_, err := apiClient.RenameFirewall(d.Id(), &firewall)
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
