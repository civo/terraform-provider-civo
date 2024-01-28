package firewall

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceFirewall Firewall resource with this we can create and manage all firewall
func ResourceFirewall() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Civo firewall resource. This can be used to create, modify, and delete firewalls.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: utils.ValidateName,
				Description:  "The firewall name",
			},
			"network_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The firewall network, if is not defined we use the default network",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The firewall region, if is not defined we use the global defined in the provider",
			},
			"create_default_rules": {
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				ForceNew:    true,
				Description: "The create rules flag is used to create the default firewall rules, if is not defined will be set to true, and if you set to false you need to define at least one ingress or egress rule",
			},
			"ingress_rule": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        firewallRuleSchema(),
				Description: "The ingress rules, this is a list of rules that will be applied to the firewall",
			},
			"egress_rule": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        firewallRuleSchema(),
				Description: "The egress rules, this is a list of rules that will be applied to the firewall",
			},
		},
		CreateContext: resourceFirewallCreate,
		ReadContext:   resourceFirewallRead,
		UpdateContext: resourceFirewallUpdate,
		DeleteContext: resourceFirewallDelete,
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {

			ingressRules := diff.Get("ingress_rule")
			egressRules := diff.Get("egress_rule")

			for _, v := range ingressRules.(*schema.Set).List() {
				ingress := v.(map[string]interface{})
				protocol := ingress["protocol"]

				port := ingress["port_range"]
				if protocol != "icmp" && port == "" {
					return fmt.Errorf("`ports` of ingress rules is required if protocol is `tcp` or `udp`")
				}
			}

			for _, v := range egressRules.(*schema.Set).List() {
				egress := v.(map[string]interface{})
				protocol := egress["protocol"]

				port := egress["port_range"]
				if protocol != "icmp" && port == "" {
					return fmt.Errorf("`ports` of egress rules is required if protocol is `tcp` or `udp`")
				}
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// function to create a firewall
func resourceFirewallCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if it's defined
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	createDefaultRules := d.Get("create_default_rules").(bool)

	// Added verification
	if !createDefaultRules {
		var ingressRules bool
		var egressRules bool

		if _, ok := d.GetOk("ingress_rule"); ok {
			ingressRules = ok
		}
		if _, ok := d.GetOk("egress_rule"); ok {
			egressRules = ok
		}
		if !ingressRules && !egressRules {
			return diag.Errorf("if you set create_default_rules to false you need to define at least one ingress or egress rule")
		}
	}

	log.Printf("[INFO] creating a new firewall %s", d.Get("name").(string))

	firewallConfig, err := firewallRequestBuild(d, apiClient)
	if err != nil {
		return diag.Errorf("[ERR] an error occurred while trying to build the firewall request, %s", err)
	}

	// Create StateChangeConf to wait for the firewall to be created
	createStateConf := &resource.StateChangeConf{
		Pending: []string{"failed"},
		Target:  []string{"success"},
		Refresh: func() (interface{}, string, error) {
			resp, err := apiClient.NewFirewall(firewallConfig)
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
	_, err = createStateConf.WaitForStateContext(context.Background())
	if err != nil {
		return diag.Errorf("[ERR] failed to create a new firewall: %s, err: %s", firewallConfig.Name, err)
	}

	// Get the firewall
	firewall, err := apiClient.FindFirewall(firewallConfig.Name)
	if err != nil {
		return diag.Errorf("[ERR] error retrieving firewall: %s, err: %s", firewallConfig.Name, err)
	}

	d.SetId(firewall.ID)

	return resourceFirewallRead(ctx, d, m)
}

// function to read a firewall
func resourceFirewallRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if it's defined
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	log.Printf("[INFO] retriving the firewall %s", d.Id())
	resp, err := apiClient.FindFirewall(d.Id())
	if err != nil {
		if resp == nil {
			d.SetId("")
			return nil
		}
		return diag.Errorf("[ERR] error retrieving firewall: %s", err)
	}

	d.Set("name", resp.Name)
	d.Set("network_id", resp.NetworkID)
	d.Set("region", apiClient.Region)
	d.Set("create_default_rules", d.Get("create_default_rules").(bool))

	for _, rule := range resp.Rules {
		if rule.Direction == "ingress" {
			if err := d.Set("ingress_rule", flattenFirewallRules(resp.Rules, rule.Direction)); err != nil {
				return diag.Errorf("[ERR] error setting ingress rules: %s", err)
			}
		} else {
			if err := d.Set("egress_rule", flattenFirewallRules(resp.Rules, rule.Direction)); err != nil {
				return diag.Errorf("[ERR] error setting egress rules: %s", err)
			}
		}
	}

	return nil
}

// function to update the firewall
func resourceFirewallUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if it's defined
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	if d.HasChange("name") {
		if d.Get("name").(string) != "" {
			firewall := civogo.FirewallConfig{
				Name: d.Get("name").(string),
			}
			log.Printf("[INFO] updating the firewall name, %s", d.Id())
			_, err := apiClient.RenameFirewall(d.Id(), &firewall)
			if err != nil {
				return diag.Errorf("[WARN] an error occurred while trying to rename the firewall %s, %s", d.Id(), err)
			}
		}
	}

	if d.HasChange("ingress_rule") || d.HasChange("egress_rule") {
		var ingressRules []interface{}
		var egressRules []interface{}

		if value, ok := d.GetOk("ingress_rule"); ok {
			ingressRules = value.(*schema.Set).List()
		}

		if value, ok := d.GetOk("egress_rule"); ok {
			egressRules = value.(*schema.Set).List()
		}

		// call the api to get the current rules
		allRules, err := apiClient.ListFirewallRules(d.Id())
		if err != nil {
			return diag.Errorf("[ERR] an error occurred while trying to list the firewall rules, %s", err)
		}

		// remove the rules that are not in terraform
		for _, rule := range allRules {
			if rule.Direction == "ingress" {
				if !ingressRulesContains(ingressRules, rule) {
					log.Printf("[INFO] removing the ingress rule %s", rule.ID)
					_, err := apiClient.DeleteFirewallRule(d.Id(), rule.ID)
					if err != nil {
						return diag.Errorf("[WARN] an error occurred while trying to delete the ingress rule %s, %s", rule.ID, err)
					}
				}
			} else {
				if !egressRulesContains(egressRules, rule) {
					log.Printf("[INFO] removing the egress rule %s", rule.ID)
					_, err := apiClient.DeleteFirewallRule(d.Id(), rule.ID)
					if err != nil {
						return diag.Errorf("[WARN] an error occurred while trying to delete the egress rule %s, %s", rule.ID, err)
					}
				}
			}
		}

		// add the new ingressRules that are not in the current rules
		if len(ingressRules) > 0 {
			for _, ingressRule := range ingressRules {
				if ingressRule.(map[string]interface{})["id"] == "" {
					fwRule := firewallUpdateBuild(ingressRule, apiClient.Region, "ingress", d)
					resp, err := apiClient.NewFirewallRule(fwRule)
					if err != nil {
						return diag.Errorf("[WARN] an error occurred while trying to create the ingress rule %s, %s", fwRule, err)
					}
					log.Printf("[INFO] creating a new ingress rule %s", resp.ID)
				}
			}
		}

		// add the new egressRules that are not in the current rules
		if len(egressRules) > 0 {
			for _, egressRule := range egressRules {
				if egressRule.(map[string]interface{})["id"] == "" {
					fwRule := firewallUpdateBuild(egressRule, apiClient.Region, "egress", d)
					resp, err := apiClient.NewFirewallRule(fwRule)
					if err != nil {
						return diag.Errorf("[WARN] an error occurred while trying to create the egress rule %s, %s", fwRule, err)
					}
					log.Printf("[INFO] creating a new egress rule %s", resp.ID)
				}
			}
		}
	}

	return resourceFirewallRead(ctx, d, m)
}

// function to delete a firewall
func resourceFirewallDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if it's defined
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	firewallID := d.Id()
	log.Printf("[INFO] Checking if firewall %s exists", firewallID)
	_, err := apiClient.FindFirewall(firewallID)
	if err != nil {
		log.Printf("[INFO] Unable to find firewall %s - probably it's been deleted", firewallID)
		return nil
	}

	log.Printf("[INFO] deleting the firewall %s", firewallID)

	deleteStateConf := &resource.StateChangeConf{
		Pending: []string{"failed"},
		Target:  []string{"success"},
		Refresh: func() (interface{}, string, error) {
			resp, err := apiClient.DeleteFirewall(firewallID)
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
		return diag.Errorf("error waiting for firewall (%s) to be deleted: %s", firewallID, err)
	}

	return nil
}

// ingressRulesContains check if the ingress rules contains the rule
func ingressRulesContains(ingressRules []interface{}, rule civogo.FirewallRule) bool {
	for _, ingressRule := range ingressRules {
		if ingressRule.(map[string]interface{})["id"] == rule.ID {
			return true
		}
	}
	return false
}

// egressRulesContains check if the egress rules contains the rule
func egressRulesContains(egressRules []interface{}, rule civogo.FirewallRule) bool {
	for _, egressRule := range egressRules {
		if egressRule.(map[string]interface{})["id"] == rule.ID {
			return true
		}
	}
	return false
}

func firewallRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    false,
				Computed:    true,
				Description: "The ID of the firewall rule. This is only set when the rule is created by terraform.",
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
				// Computed:     true,
				Description:  "A string that will be the displayed name/reference for this rule",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"protocol": {
				Type:        schema.TypeString,
				Default:     "tcp",
				Optional:    true,
				Description: "The protocol choice from `tcp`, `udp` or `icmp` (the default if unspecified is `tcp`)",
				ValidateFunc: validation.StringInSlice([]string{
					"tcp",
					"udp",
					"icmp",
				}, false),
			},
			"port_range": {
				Type:     schema.TypeString,
				Optional: true,
				// Computed:     true,
				Description:  "The port or port range to open, can be a single port or a range separated by a dash (`-`), e.g. `80` or `80-443`",
				ValidateFunc: validation.NoZeroValues,
			},
			"cidr": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The CIDR notation of the other end to affect, or a valid network CIDR (e.g. 0.0.0.0/0 to open for everyone or 1.2.3.4/32 to open just for a specific IP address)",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The action of the rule can be allow or deny. When we set the `action = 'allow'`, this is going to add a rule to allow traffic. Similarly, setting `action = 'deny'` will deny the traffic.",
				ValidateFunc: validation.StringInSlice([]string{
					"allow", "deny",
				}, false),
			},
		},
	}
}

// firewallRequestBuild builds the request body for a firewall
func firewallRequestBuild(d *schema.ResourceData, client *civogo.Client) (*civogo.FirewallConfig, error) {
	var networkID string

	if attr, ok := d.GetOk("network_id"); ok {
		networkID = attr.(string)
	} else {
		network, err := client.GetDefaultNetwork()
		if err != nil {
			return nil, fmt.Errorf("[ERR] failed to get the default network: %s", err)
		}
		networkID = network.ID
	}

	createFirewallRules := d.Get("create_default_rules").(bool)

	firewallCofig := &civogo.FirewallConfig{
		Name:        d.Get("name").(string),
		Region:      client.Region,
		NetworkID:   networkID,
		CreateRules: &createFirewallRules,
	}

	// Get ingress_rule
	if v, ok := d.GetOk("ingress_rule"); ok {
		firewallCofig.Rules = append(firewallCofig.Rules, expandFirewallRules(v.(*schema.Set).List(), "ingress")...)
	}

	// Get egress_rule
	if v, ok := d.GetOk("egress_rule"); ok {
		firewallCofig.Rules = append(firewallCofig.Rules, expandFirewallRules(v.(*schema.Set).List(), "egress")...)
	}

	return firewallCofig, nil
}

// expandFirewallIngressRules expands the ingress rules
func expandFirewallRules(rules []interface{}, direction string) []civogo.FirewallRule {
	var firewallRules []civogo.FirewallRule
	for _, rule := range rules {
		rule := rule.(map[string]interface{})
		fwRule := civogo.FirewallRule{
			Label:     rule["label"].(string),
			Protocol:  rule["protocol"].(string),
			Direction: direction,
			Cidr:      expandFirewallRuleCIDR(rule["cidr"].(*schema.Set).List()),
			Action:    rule["action"].(string),
		}

		if rule["port_range"].(string) != "" {
			fwRule.Ports = rule["port_range"].(string)
		}

		firewallRules = append(firewallRules, fwRule)
	}
	return firewallRules
}

// expandFirewallRuleCIDR expands the cidr rules
func expandFirewallRuleCIDR(strings []interface{}) []string {
	expandedStrings := make([]string, len(strings))
	for i, v := range strings {
		expandedStrings[i] = v.(string)
	}

	return expandedStrings
}

// flattenFirewallRules flattens the firewall rules
func flattenFirewallRules(rules []civogo.FirewallRule, direction string) []interface{} {
	if rules == nil {
		return nil
	}

	// We need to do this because our rules come all together
	// and we need to split them up into ingress and egress rules
	rulesCount := 0
	rulesObject := []civogo.FirewallRule{}
	for _, rule := range rules {
		if rule.Direction == direction {
			rulesCount++
			rulesObject = append(rulesObject, rule)
		}
	}

	log.Printf("[INFO] retriving the firewall rules %+v", rulesObject)

	flattenedRules := make([]interface{}, rulesCount)
	for i, rule := range rulesObject {
		flattenedRules[i] = map[string]interface{}{
			"id":         rule.ID,
			"label":      rule.Label,
			"protocol":   rule.Protocol,
			"port_range": rule.Ports,
			"action":     rule.Action,
			"cidr":       flattenFirewallRuleCIDR(rule.Cidr),
		}
	}

	log.Printf("[INFO] retriving the flattenedRules %+v", flattenedRules)

	return flattenedRules
}

func flattenFirewallRuleCIDR(strings []string) *schema.Set {
	flattenedStrings := schema.NewSet(schema.HashString, []interface{}{})
	for _, v := range strings {
		flattenedStrings.Add(v)
	}

	return flattenedStrings
}

// firewallUpdateBuild builds the request body for a firewall update
func firewallUpdateBuild(object interface{}, region, direction string, d *schema.ResourceData) *civogo.FirewallRuleConfig {
	firewalObject := &civogo.FirewallRuleConfig{
		FirewallID: d.Id(),
		Region:     region,
		Protocol:   object.(map[string]interface{})["protocol"].(string),
		Cidr:       expandFirewallRuleCIDR(object.(map[string]interface{})["cidr"].(*schema.Set).List()),
		Direction:  direction,
		Action:     object.(map[string]interface{})["action"].(string),
		Label:      object.(map[string]interface{})["label"].(string),
		Ports:      object.(map[string]interface{})["port_range"].(string),
	}
	return firewalObject
}
