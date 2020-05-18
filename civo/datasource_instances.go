package civo

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source to get and filter all instances with filter
func dataSourceInstances() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reverse_dns": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"initial_user": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sshkey_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"firewall_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"script": {
				Type:     schema.TypeString,
				Computed: true,
			},
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
		FilterKeys: []string{
			"id",
			"hostname",
			"public_ip",
			"private_ip",
			"pseudo_ip",
			"size",
			"template",
			"created_at",
		},
		SortKeys: []string{
			"id",
			"hostname",
			"public_ip",
			"private_ip",
			"pseudo_ip",
			"size",
			"template",
			"created_at",
		},
		ResultAttributeName: "instances",
		FlattenRecord:       flattenDataSourceInstances,
		GetRecords:          getDataSourceInstances,
	}

	return datalist.NewResource(dataListConfig)

}

func getDataSourceInstances(m interface{}) ([]interface{}, error) {
	apiClient := m.(*civogo.Client)

	instance := []interface{}{}
	partialInstances, err := apiClient.ListInstances(1, 200)
	if err != nil {
		return nil, fmt.Errorf("[ERR] error retrieving sizes: %s", err)
	}

	for _, partialInstance := range partialInstances.Items {
		instance = append(instance, partialInstance)
	}

	return instance, nil
}

func flattenDataSourceInstances(instance, m interface{}) (map[string]interface{}, error) {

	i := instance.(civogo.Instance)

	flattenedInstance := map[string]interface{}{}
	flattenedInstance["id"] = i.ID
	flattenedInstance["hostname"] = i.Hostname
	flattenedInstance["reverse_dns"] = i.ReverseDNS
	flattenedInstance["size"] = i.Size
	flattenedInstance["network_id"] = i.NetworkID
	flattenedInstance["template"] = i.TemplateID
	flattenedInstance["initial_user"] = i.InitialUser
	flattenedInstance["notes"] = i.Notes
	flattenedInstance["sshkey_id"] = i.SSHKey
	flattenedInstance["firewall_id"] = i.FirewallID
	flattenedInstance["tags"] = i.Tags
	flattenedInstance["script"] = i.Script
	flattenedInstance["initial_password"] = i.InitialPassword
	flattenedInstance["private_ip"] = i.PublicIP
	flattenedInstance["public_ip"] = i.PrivateIP
	flattenedInstance["pseudo_ip"] = i.PseudoIP
	flattenedInstance["status"] = i.Status
	flattenedInstance["created_at"] = i.CreatedAt.UTC().String()

	return flattenedInstance, nil
}
