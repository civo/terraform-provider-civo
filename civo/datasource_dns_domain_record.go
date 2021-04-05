package civo

import (
	"fmt"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Data source to get from the api a specific domain record
// using the id or the name of the domain
func dataSourceDNSDomainRecord() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDNSDomainRecordRead,
		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			// Computed resource
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDNSDomainRecordRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)
	domain := d.Get("domain_id").(string)
	name := d.Get("name").(string)

	allRecords, err := apiClient.ListDNSRecords(domain)
	if err != nil {
		return fmt.Errorf("error retrieving all domain records: %s", err)
	}

	record, err := getRecordByName(allRecords, name)
	if err != nil {
		return err
	}

	d.SetId(record.ID)
	d.Set("name", record.Name)
	d.Set("type", record.Type)
	d.Set("value", record.Value)
	d.Set("priority", record.Priority)
	d.Set("ttl", record.TTL)
	d.Set("account_id", record.AccountID)
	d.Set("created_at", record.CreatedAt.UTC().String())
	d.Set("updated_at", record.UpdatedAt.UTC().String())

	return nil
}

func getRecordByName(allRecord []civogo.DNSRecord, name string) (*civogo.DNSRecord, error) {
	results := make([]civogo.DNSRecord, 0)
	for _, v := range allRecord {
		if v.Name == name {
			results = append(results, v)
		}
	}
	if len(results) == 1 {
		return &results[0], nil
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no records found with name %s", name)
	}
	return nil, fmt.Errorf("too many records found (found %d, expected 1)", len(results))
}
