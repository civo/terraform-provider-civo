package dns_test

import (
	"testing"

	"fmt"

	"github.com/civo/terraform-provider-civo/civo/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func CivoDNSDomainName_importBasic(t *testing.T) {
	resourceName := "civo_dns_domain_name.foobar"
	domainName := fmt.Sprintf("foobar-test-terraform-%s.com", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CivoDNSDomainNameDestroy,
		Steps: []resource.TestStep{
			{
				Config: CivoDNSDomainNameConfigBasic(domainName),
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
