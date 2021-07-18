package civo

import (
	"fmt"
	"log"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Dns Domain resource, with this we can create and manage DNS Domain
func resourceDNSDomainName() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "A fully qualified domain name",
				ValidateFunc: utils.ValidateName,
			},
			// Computed resource
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Create: resourceDNSDomainNameCreate,
		Read:   resourceDNSDomainNameRead,
		Update: resourceDNSDomainNameUpdate,
		Delete: resourceDNSDomainNameDelete,
		//Exists: resourceExistsItem,
		Importer: &schema.ResourceImporter{
			State: resourceDNSDomainImport,
		},
	}
}

// function to create a new domain in your account
func resourceDNSDomainNameCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] Creating the domain %s", d.Get("name").(string))
	dnsDomain, err := apiClient.CreateDNSDomain(d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("failed to create a new domains: %s", err)
	}

	d.SetId(dnsDomain.ID)

	return resourceDNSDomainNameRead(d, m)
}

// function to read a domain from your account
func resourceDNSDomainNameRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] retriving the domain %s", d.Get("name").(string))
	resp, err := apiClient.GetDNSDomain(d.Get("name").(string))
	if err != nil {
		if resp == nil {
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
func resourceDNSDomainNameUpdate(d *schema.ResourceData, m interface{}) error {
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

	return resourceDNSDomainNameRead(d, m)
}

// function to delete a specific domain
func resourceDNSDomainNameDelete(d *schema.ResourceData, m interface{}) error {
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
func resourceDNSDomainImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] Searching the domain %s", d.Id())
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
