package civo

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
)

func resourceFirewallRule() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"firewall_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateName,
			},
			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "",
				ValidateFunc: validation.StringInSlice([]string{
					"tcp",
					"udp",
					"icmp",
				}, false),
			},
			"start_port": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "",
				ValidateFunc: validation.NoZeroValues,
			},
			"end_port": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "",
				ValidateFunc: validation.NoZeroValues,
			},
			"cird": {
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				Description: "",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"direction": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "",
				ValidateFunc: validation.StringInSlice([]string{
					"inbound",
					"outbound",
				}, false),
			},
			"label": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
		Create: resourceFirewallRuleCreate,
		Read:   resourceFirewallRuleRead,
		Delete: resourceFirewallRuleDelete,
		//Importer: &schema.ResourceImporter{
		//	State: schema.ImportStatePassthrough,
		//},
	}
}

func resourceFirewallRuleCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	tfCidrs := d.Get("cird").(*schema.Set).List()
	cird := make([]string, len(tfCidrs))
	for i, tfCird := range tfCidrs {
		cird[i] = tfCird.(string)
	}

	config := &civogo.FirewallRuleConfig{
		FirewallID: d.Get("firewall_id").(string),
		Protocol:   d.Get("protocol").(string),
		StartPort:  d.Get("start_port").(string),
		Direction:  d.Get("direction").(string),
		Cidr:       cird,
	}

	if attr, ok := d.GetOk("end_port"); ok {
		config.EndPort = attr.(string)
	}

	if attr, ok := d.GetOk("label"); ok {
		config.Label = attr.(string)
	}

	firewallRule, err := apiClient.NewFirewallRule(config)
	if err != nil {
		fmt.Errorf("[ERR] failed to create a new firewall: %s", err)
		return err
	}

	d.SetId(firewallRule.ID)

	return resourceFirewallRuleRead(d, m)
}

func resourceFirewallRuleRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	resp, err := apiClient.FindFirewallRule(d.Get("firewall_id").(string), d.Id())
	if err != nil {
		if resp != nil {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("[ERR] error retrieving firewall Rule: %s", err)
	}

	d.Set("firewall_id", resp.FirewallID)
	d.Set("protocol", resp.Protocol)
	d.Set("start_port", resp.StartPort)
	d.Set("end_port", resp.EndPort)
	d.Set("cird", resp.Cidr)
	d.Set("direction", resp.Direction)
	d.Set("label", resp.Label)

	return nil
}

func resourceFirewallRuleDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	_, err := apiClient.DeleteFirewallRule(d.Get("firewall_id").(string), d.Id())
	if err != nil {
		log.Printf("[INFO] civo firewall rule (%s) was delete", d.Id())
	}
	return nil
}
