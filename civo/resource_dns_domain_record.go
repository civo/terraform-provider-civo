package civo

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
)

const (
	// DNSRecordTypeA represents an A record
	DNSRecordTypeA = "a"

	// DNSRecordTypeCName represents an CNAME record
	DNSRecordTypeCName = "cname"

	// DNSRecordTypeMX represents an MX record
	DNSRecordTypeMX = "mx"

	// DNSRecordTypeTXT represents an TXT record
	DNSRecordTypeTXT = "txt"
)

func resourceDnsDomainRecord() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id from domain name",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The choice of RR type from a, cname, mx or txt",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The portion before the domain name (e.g. www) or an @ for the apex/root domain (you cannot use an A record with an amex/root domain)",
			},
			"value": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The IP address (A or MX), hostname (CNAME or MX) or text value (TXT) to serve for this record",
				ValidateFunc: validation.NoZeroValues,
			},
			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "Useful for MX records only, the priority mail should be attempted it (defaults to 10)",
			},
			"ttl": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "How long caching DNS servers should cache this record for, in seconds (the minimum is 600 and the default if unspecified is 600)",
			},
			// Computed resource
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
		Create: resourceDnsDomainRecordCreate,
		Read:   resourceDnsDomainRecordRead,
		Update: resourceDnsDomainRecordUpdate,
		Delete: resourceDnsDomainRecordDelete,
		//Exists: resourceExistsItem,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceDnsDomainRecordCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	config := &civogo.DNSRecordConfig{
		Name:  d.Get("name").(string),
		Value: d.Get("value").(string),
		TTL:   d.Get("ttl").(int),
	}

	if attr, ok := d.GetOk("priority"); ok {
		if d.Get("type").(string) != "mx" {
			return fmt.Errorf("[WARN] warning priority value is only allow in the MX records")
		}
		config.Priority = attr.(int)
	}

	if d.Get("type").(string) == "a" {
		config.Type = DNSRecordTypeA
	}

	if d.Get("type").(string) == "cname" {
		config.Type = DNSRecordTypeCName
	}

	if d.Get("type").(string) == "mx" {
		config.Type = DNSRecordTypeMX
	}

	if d.Get("type").(string) == "txt" {
		config.Type = DNSRecordTypeTXT
	}

	dnsDomainRecord, err := apiClient.CreateDNSRecord(d.Get("domain_id").(string), config)
	if err != nil {
		fmt.Errorf("failed to create a new record: %s", err)
		return err
	}

	d.SetId(dnsDomainRecord.ID)

	return resourceDnsDomainRecordRead(d, m)
}

func resourceDnsDomainRecordRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	resp, err := apiClient.GetDNSRecord(d.Get("domain_id").(string), d.Id())
	if err != nil {
		log.Printf("[WARN] civo domain record (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("name", resp.Name)
	d.Set("account_id", resp.AccountID)
	d.Set("domain_id", resp.DNSDomainID)
	d.Set("name", resp.Name)
	d.Set("value", resp.Value)
	d.Set("type", resp.Type)
	d.Set("priority", resp.Priority)
	d.Set("ttl", resp.TTL)
	d.Set("created_at", resp.CreatedAt.String())
	d.Set("updated_at", resp.UpdatedAt.String())

	return nil
}

func resourceDnsDomainRecordUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	resp, err := apiClient.GetDNSRecord(d.Get("domain_id").(string), d.Id())
	if err != nil {
		log.Printf("[WARN] civo domain record (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	config := &civogo.DNSRecordConfig{}

	if d.HasChange("name") || d.HasChange("value") || d.HasChange("priority") || d.HasChange("ttl") || d.HasChange("type") {
		config.Name = d.Get("name").(string)
		config.Value = d.Get("value").(string)
		config.Priority = d.Get("priority").(int)
		config.TTL = d.Get("ttl").(int)

		if d.Get("type").(string) == "a" {
			config.Type = DNSRecordTypeA
		}

		if d.Get("type").(string) == "cname" {
			config.Type = DNSRecordTypeCName
		}

		if d.Get("type").(string) == "mx" {
			config.Type = DNSRecordTypeMX
		}

		if d.Get("type").(string) == "txt" {
			config.Type = DNSRecordTypeTXT
		}
	}

	_, err = apiClient.UpdateDNSRecord(resp, config)
	if err != nil {
		log.Printf("[WARN] an error occurred while renamed the domain record (%s)", d.Id())
	}

	return resourceDnsDomainRecordRead(d, m)
}

func resourceDnsDomainRecordDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	resp, err := apiClient.GetDNSRecord(d.Get("domain_id").(string), d.Id())
	if err != nil {
		log.Printf("[WARN] civo domain record (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	_, err = apiClient.DeleteDNSRecord(resp)
	if err != nil {
		log.Printf("[WARN] civo domain record (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	return nil
}
