package network

import (
	"context"
	"log"
	"time"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceNetwork function returns a schema.Resource that represents a Network.
// This can be used to create, read, update, and delete operations for a Network in the infrastructure.
func ResourceNetwork() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Civo network resource. This can be used to create, modify, and delete networks.",
		Schema: map[string]*schema.Schema{
			"label": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name for the network",
				ValidateFunc: utils.ValidateName,
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The region of the network",
			},
			"cidr_v4": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The CIDR block for the network",
			},
			"nameservers_v4": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of nameservers for the network",
			},
			// Computed resource
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the network",
			},
			"default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If the network is default, this will be `true`",
			},
			// VLAN Network
			"vlan_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "VLAN ID for the network",
			},
			"vlan_cidr_v4": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CIDR for VLAN IPv4",
			},
			"vlan_gateway_ip_v4": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Gateway IP for VLAN IPv4",
			},
			"vlan_physical_interface": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Physical interface for VLAN",
			},
			"vlan_allocation_pool_v4_start": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Start of the IPv4 allocation pool for VLAN",
			},
			"vlan_allocation_pool_v4_end": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "End of the IPv4 allocation pool for VLAN",
			},
		},
		CreateContext: resourceNetworkCreate,
		ReadContext:   resourceNetworkRead,
		UpdateContext: resourceNetworkUpdate,
		DeleteContext: resourceNetworkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// function to create a new network
func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is defined in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	log.Printf("[INFO] creating the new network %s", d.Get("label").(string))
	vlanConfig := civogo.VLANConnectConfig{
		VlanID:                d.Get("vlan_id").(int),
		PhysicalInterface:     d.Get("vlan_physical_interface").(string),
		CIDRv4:                d.Get("vlan_cidr_v4").(string),
		GatewayIPv4:           d.Get("vlan_gateway_ip_v4").(string),
		AllocationPoolV4Start: d.Get("vlan_allocation_pool_v4_start").(string),
		AllocationPoolV4End:   d.Get("vlan_allocation_pool_v4_end").(string),
	}

	configs := civogo.NetworkConfig{
		Label:         d.Get("label").(string),
		CIDRv4:        d.Get("cidr_v4").(string),
		Region:        apiClient.Region,
		NameserversV4: expandStringList(d.Get("nameservers_v4")),
	}
	// Only add VLAN configuration if VLAN ID is set
	if vlanConfig.VlanID > 0 {
		configs.VLanConfig = &vlanConfig
	}

	// Retry the network creation using the utility function
	err := utils.RetryUntilSuccessOrTimeout(func() error {
		log.Printf("[INFO] Attempting to create the network %s", d.Get("label").(string))
		network, err := apiClient.CreateNetwork(configs)
		if err != nil {
			return err
		}
		d.SetId(network.ID)
		return nil
	}, 10*time.Second, 2*time.Minute)

	if err != nil {
		return diag.Errorf("[ERR] failed to create a new network after multiple attempts: %s", err)
	}

	return resourceNetworkRead(ctx, d, m)
}

// function to read a network
func resourceNetworkRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is defined in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	CurrentNetwork := civogo.Network{}

	log.Printf("[INFO] retriving the network %s", d.Id())
	resp, err := apiClient.ListNetworks()
	if err != nil {
		if resp == nil {
			d.SetId("")
			return nil
		}

		return diag.Errorf("[ERR] failed to list the network: %s", err)
	}

	for _, net := range resp {
		if net.ID == d.Id() {
			CurrentNetwork = net
		}
	}

	d.Set("name", CurrentNetwork.Name)
	d.Set("region", apiClient.Region)
	d.Set("label", CurrentNetwork.Label)
	d.Set("default", CurrentNetwork.Default)
	d.Set("cidr_v4", CurrentNetwork.CIDR)
	d.Set("nameservers_v4", CurrentNetwork.NameserversV4)

	return nil
}

// function to update the network
func resourceNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is defined in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	if d.HasChange("label") {
		log.Printf("[INFO] updating the network %s", d.Id())
		_, err := apiClient.RenameNetwork(d.Get("label").(string), d.Id())
		if err != nil {
			return diag.Errorf("[ERR] An error occurred while rename the network %s", d.Id())
		}
		return resourceNetworkRead(ctx, d, m)
	}
	return resourceNetworkRead(ctx, d, m)
}

// function to delete a network
func resourceNetworkDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is defined in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	netowrkID := d.Id()
	log.Printf("[INFO] Checking if firewall %s exists", netowrkID)
	_, err := apiClient.FindNetwork(netowrkID)
	if err != nil {
		log.Printf("[INFO] Unable to find network %s - probably it's been deleted", netowrkID)
		return nil
	}

	log.Printf("[INFO] deleting the network %s", netowrkID)

	deleteStateConf := &resource.StateChangeConf{
		Pending: []string{"failed"},
		Target:  []string{"success"},
		Refresh: func() (interface{}, string, error) {
			resp, err := apiClient.DeleteNetwork(netowrkID)
			if err != nil {
				return 0, "", err
			}
			return resp, string(resp.Result), nil
		},
		Timeout:        60 * time.Minute,
		Delay:          3 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 10,
	}
	_, err = deleteStateConf.WaitForStateContext(context.Background())
	if err != nil {
		return diag.Errorf("error waiting for network (%s) to be deleted: %s", netowrkID, err)
	}

	return nil
}

func expandStringList(input interface{}) []string {
	var result []string

	if inputList, ok := input.([]interface{}); ok {
		for _, item := range inputList {
			if str, ok := item.(string); ok {
				result = append(result, str)
			}
		}
	}
	return result
}
