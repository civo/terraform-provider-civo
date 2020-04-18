package civo

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

// Dns Domain resource, with this we can create and manage DNS Domain
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
			State: resourceDnsDomainImport,
		},
	}
}

// function to create a new domain in your account
func resourceDnsDomainNameCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] Creating the domain %s", d.Get("name").(string))
	dnsDomain, err := apiClient.CreateDNSDomain(d.Get("name").(string))
	if err != nil {
		fmt.Errorf("failed to create a new domains: %s", err)
		return err
	}

	d.SetId(dnsDomain.ID)

	return resourceDnsDomainNameRead(d, m)
}

// function to read a domain from your account
func resourceDnsDomainNameRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] retriving the domain %s", d.Get("name").(string))
	resp, err := apiClient.GetDNSDomain(d.Get("name").(string))
	if err != nil {
		if resp != nil {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("[ERR] error retrieving domain: %s", err)
	}

	d.Set("name", resp.Name)
	d.Set("account_id", resp.AccountID)

	return nil
}

// function to update a specific domain
func resourceDnsDomainNameUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] Searching the domain %s", d.Get("name").(string))
	resp, err := apiClient.FindDNSDomain(d.Id())
	if err != nil {
		log.Printf("[WARN] domain (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	if d.HasChange("name") {
		name := d.Get("name").(string)
		log.Printf("[INFO] Renaming the domain to %s", d.Get("name").(string))
		_, err := apiClient.UpdateDNSDomain(resp, name)
		if err != nil {
			return fmt.Errorf("[WARN] an error occurred while renamed the domain (%s)", d.Id())
		}

	}

	return resourceDnsDomainNameRead(d, m)
}

// function to delete a specific domain
func resourceDnsDomainNameDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] Searching the domain to %s", d.Get("name").(string))
	resp, err := apiClient.FindDNSDomain(d.Id())
	if err != nil {
		log.Printf("[WARN] domain (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	log.Printf("[INFO] Deleting the domain %s", d.Get("name").(string))
	_, err = apiClient.DeleteDNSDomain(resp)
	if err != nil {
		return fmt.Errorf("[ERR] an error occurred while tring to delete the domain %s", d.Id())
	}
	return nil
}

// custom import to able add a main domain to the terraform
func resourceDnsDomainImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] Searching the domain %s", d.Get("name").(string))
	resp, err := apiClient.GetDNSDomain(d.Id())
	if err != nil {
		if resp != nil {
			return nil, err
		}
	}

	d.SetId(resp.ID)
	d.Set("name", resp.Name)
	d.Set("account_id", resp.AccountID)

	return []*schema.ResourceData{d}, nil
}
