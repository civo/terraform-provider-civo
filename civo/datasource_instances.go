package civo

import (
	"fmt"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source to get and filter all instances with filter
func dataSourceInstances() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema: instancesSchema(),
		ExtraQuerySchema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		ResultAttributeName: "instances",
		FlattenRecord:       flattenDataSourceInstances,
		GetRecords:          getDataSourceInstances,
	}

	return datalist.NewResource(dataListConfig)

}

func getDataSourceInstances(m interface{}, extra map[string]interface{}) ([]interface{}, error) {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	region, ok := extra["region"].(string)
	if !ok {
		return nil, fmt.Errorf("unable to find `region` key from query data")
	}

	if region != "" {
		apiClient.Region = region
	}

	var instance []interface{}
	partialInstances, err := apiClient.ListInstances(1, 200)
	if err != nil {
		return nil, fmt.Errorf("[ERR] error retrieving instances: %s", err)
	}

	for _, partialInstance := range partialInstances.Items {
		instance = append(instance, partialInstance)
	}

	return instance, nil
}

func flattenDataSourceInstances(instance, m interface{}, extra map[string]interface{}) (map[string]interface{}, error) {

	region, ok := extra["region"].(string)
	if !ok {
		return nil, fmt.Errorf("unable to find `region` key from query data")
	}

	i := instance.(civogo.Instance)

	flattenedInstance := map[string]interface{}{}
	flattenedInstance["id"] = i.ID
	flattenedInstance["hostname"] = i.Hostname
	flattenedInstance["region"] = region
	flattenedInstance["reverse_dns"] = i.ReverseDNS
	flattenedInstance["size"] = i.Size
	flattenedInstance["cpu_cores"] = i.CPUCores
	flattenedInstance["ram_mb"] = i.RAMMegabytes
	flattenedInstance["disk_gb"] = i.DiskGigabytes
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

func instancesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Description: "id of the instance",
		},
		"hostname": {
			Type:        schema.TypeString,
			Description: "hostname of the instance",
		},
		"region": {
			Type:        schema.TypeString,
			Description: "region of the instance",
		},
		"reverse_dns": {
			Type:        schema.TypeString,
			Description: "reverse DNS of the instance",
		},
		"size": {
			Type:        schema.TypeString,
			Description: "size of the instance",
		},
		"cpu_cores": {
			Type:        schema.TypeInt,
			Description: "CPU of the instance",
		},
		"ram_mb": {
			Type:        schema.TypeInt,
			Description: "RAM of the instance",
		},
		"disk_gb": {
			Type:        schema.TypeInt,
			Description: "SSD of the instance",
		},
		"network_id": {
			Type:        schema.TypeString,
			Description: "netwoerk id of the instance",
		},
		"template": {
			Type:        schema.TypeString,
			Description: "template of the instance",
		},
		"initial_user": {
			Type:        schema.TypeString,
			Description: "initial user of the instance",
		},
		"notes": {
			Type:        schema.TypeString,
			Description: "note of the instance",
		},
		"sshkey_id": {
			Type:        schema.TypeString,
			Description: "sshkey id of the instance",
		},
		"firewall_id": {
			Type:        schema.TypeString,
			Description: "firewall id of the instance",
		},
		"tags": {
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "tags of the instance",
		},
		"script": {
			Type:        schema.TypeString,
			Description: "script of the instance",
		},
		"initial_password": {
			Type:        schema.TypeString,
			Description: "initial password of the instance",
		},
		"private_ip": {
			Type:        schema.TypeString,
			Description: "private ip of the instance",
		},
		"public_ip": {
			Type:        schema.TypeString,
			Description: "public ip of the instance",
		},
		"pseudo_ip": {
			Type:        schema.TypeString,
			Description: "pseudo ip of the instance",
		},
		"status": {
			Type:        schema.TypeString,
			Description: "status of the instance",
		},
		"created_at": {
			Type:        schema.TypeString,
			Description: "creation date of the instance",
		},
	}
}
