package civo

import (
	"fmt"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source to get and filter all instances with filter
func dataSourceInstances() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		Description: strings.Join([]string{
			"Get information on instances for use in other resources, with the ability to filter and sort the results. If no filters are specified, all instances will be returned.",
			"Note: You can use the `civo_instance` data source to obtain metadata about a single instance if you already know the id, unique hostname, or unique tag to retrieve.",
		}, "\n\n"),
		RecordSchema: instancesSchema(),
		ExtraQuerySchema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If used, all instances will be from the provided region",
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
			Description: "ID of the instance",
		},
		"hostname": {
			Type:        schema.TypeString,
			Description: "Hostname of the instance",
		},
		"region": {
			Type:        schema.TypeString,
			Description: "Region of the instance",
		},
		"reverse_dns": {
			Type:        schema.TypeString,
			Description: "Reverse DNS of the instance",
		},
		"size": {
			Type:        schema.TypeString,
			Description: "Size of the instance",
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
			Description: "SSD size of the instance",
		},
		"network_id": {
			Type:        schema.TypeString,
			Description: "Network id of the instance",
		},
		"template": {
			Type:        schema.TypeString,
			Description: "Disk image/template of the instance",
		},
		"initial_user": {
			Type:        schema.TypeString,
			Description: "Initial user of the instance",
		},
		"notes": {
			Type:        schema.TypeString,
			Description: "Note of the instance",
		},
		"sshkey_id": {
			Type:        schema.TypeString,
			Description: "SSH key id of the instance",
		},
		"firewall_id": {
			Type:        schema.TypeString,
			Description: "Firewall ID of the instance",
		},
		"tags": {
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Tags of the instance",
		},
		"script": {
			Type:        schema.TypeString,
			Description: "Script of the instance",
		},
		"initial_password": {
			Type:        schema.TypeString,
			Description: "Initial password of the instance",
		},
		"private_ip": {
			Type:        schema.TypeString,
			Description: "Private IP of the instance",
		},
		"public_ip": {
			Type:        schema.TypeString,
			Description: "Public IP of the instance",
		},
		"pseudo_ip": {
			Type:        schema.TypeString,
			Description: "Pseudo IP of the instance",
		},
		"status": {
			Type:        schema.TypeString,
			Description: "Status of the instance",
		},
		"created_at": {
			Type:        schema.TypeString,
			Description: "Creation date of the instance",
		},
	}
}
