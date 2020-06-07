package civo

import (
	"fmt"
	"log"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// The resource network represent a network inside the cloud
func resourceNetwork() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"label": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name for the network",
				ValidateFunc: validateName,
			},
			// Computed resource
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Create: resourceNetworkCreate,
		Read:   resourceNetworkRead,
		Update: resourceNetworkUpdate,
		Delete: resourceNetworkDelete,
		//Exists: resourceExistsItem,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// function to create a new network
func resourceNetworkCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] creating the new network %s", d.Get("label").(string))
	network, err := apiClient.NewNetwork(d.Get("label").(string))
	if err != nil {
		return fmt.Errorf("[ERR] failed to create a new network: %s", err)
	}

	d.SetId(network.ID)

	return resourceNetworkRead(d, m)
}

// function to read a network
func resourceNetworkRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	CurrentNetwork := civogo.Network{}

	log.Printf("[INFO] retriving the network %s", d.Id())
	resp, err := apiClient.ListNetworks()
	if err != nil {
		if resp == nil {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("[ERR] failed to list all network %s", err)
	}

	for _, net := range resp {
		if net.ID == d.Id() {
			CurrentNetwork = net
		}
	}

	d.Set("name", CurrentNetwork.Name)
	d.Set("region", CurrentNetwork.Region)
	d.Set("default", CurrentNetwork.Default)
	d.Set("cidr", CurrentNetwork.CIDR)
	d.Set("label", CurrentNetwork.Label)

	return nil
}

// function to update the network
func resourceNetworkUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	if d.HasChange("label") {
		log.Printf("[INFO] updating the network %s", d.Id())
		_, err := apiClient.RenameNetwork(d.Get("label").(string), d.Id())
		if err != nil {
			return fmt.Errorf("[ERR] An error occurred while rename the network %s", d.Id())
		}
		return resourceNetworkRead(d, m)
	}
	return resourceNetworkRead(d, m)
}

// function to delete a network
func resourceNetworkDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] deleting the network %s", d.Id())
	_, err := apiClient.DeleteNetwork(d.Id())
	if err != nil {
		return fmt.Errorf("[ERR] an error occurred while tring to delete the network %s", d.Id())
	}
	return nil
}
