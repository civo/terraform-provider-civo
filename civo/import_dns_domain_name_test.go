package civo

import (
	"testing"

	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccCivoDNSDomainName_importBasic(t *testing.T) {
	resourceName := "civo_dns_domain_name.foobar"
	domainName := fmt.Sprintf("foobar-test-terraform-%s.com", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoDNSDomainNameDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCivoDNSDomainNameConfigBasic(domainName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprint(domainName),
			},
		},
	})
}
