package civo

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceDnsDomainName() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "A fully qualified domain name",
				ValidateFunc: validateName,
			},
			// Computed resource
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Create: resourceDnsDomainNameCreate,
		Read:   resourceDnsDomainNameRead,
		Update: resourceDnsDomainNameUpdate,
		Delete: resourceDnsDomainNameDelete,
		//Exists: resourceExistsItem,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceDnsDomainNameCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	dnsDomain, err := apiClient.CreateDNSDomain(d.Get("name").(string))
	if err != nil {
		fmt.Errorf("failed to create a new domains: %s", err)
		return err
	}

	d.SetId(dnsDomain.ID)

	return resourceDnsDomainNameRead(d, m)
}

func resourceDnsDomainNameRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	resp, err := apiClient.GetDNSDomain(d.Get("name").(string))
	if err != nil {
		if resp != nil {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("error retrieving domain: %s", err)
	}

	d.Set("name", resp.Name)
	d.Set("account_id", resp.AccountID)

	return nil
}

func resourceDnsDomainNameUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	resp, err := apiClient.FindDNSDomain(d.Id())
	if err != nil {
		log.Printf("[WARN] Civo domain (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	if d.HasChange("name") {
		name := d.Get("name").(string)
		_, err := apiClient.UpdateDNSDomain(resp, name)
		if err != nil {
			log.Printf("[WARN] An error occurred while renamed the domain (%s)", d.Id())
		}

	}

	return resourceDnsDomainNameRead(d, m)
}

func resourceDnsDomainNameDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	resp, err := apiClient.FindDNSDomain(d.Id())
	if err != nil {
		log.Printf("[WARN] Civo domain (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	_, err = apiClient.DeleteDNSDomain(resp)
	if err != nil {
		log.Printf("[INFO] Civo domain (%s) was delete", d.Id())
	}
	return nil
}
