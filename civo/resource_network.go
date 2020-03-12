package civo

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

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
		//Importer: &schema.ResourceImporter{
		//	State: schema.ImportStatePassthrough,
		//},
	}
}

func resourceNetworkCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)
	network, err := apiClient.NewNetwork(d.Get("label").(string))
	if err != nil {
		fmt.Errorf("failed to create a new config: %s", err)
		return err
	}

	d.SetId(network.ID)

	return resourceNetworkRead(d, m)
}

func resourceNetworkRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	CurrentNetwork := civogo.Network{}

	resp, err := apiClient.ListNetworks()
	if err != nil {
		fmt.Errorf("failed to create a new config: %s", err)
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

func resourceNetworkUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	if d.HasChange("label") {
		_, err := apiClient.RenameNetwork(d.Get("label").(string), d.Id())
		if err != nil {
			log.Printf("[WARN] An error occurred while rename the network (%s)", d.Id())
		}

		return resourceNetworkRead(d, m)
	}

	return resourceNetworkRead(d, m)
}

func resourceNetworkDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	_, err := apiClient.DeleteNetwork(d.Id())
	if err != nil {
		log.Printf("[INFO] Civo network (%s) was delete", d.Id())
	}
	return nil
}
