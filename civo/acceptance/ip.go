package acceptance

import (
	"fmt"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// CivoReservedIPResourceExists queries the API and retrieves the matching Widget.
func CivoReservedIPResourceExists(n string, ip *civogo.IP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := TestAccProvider.Meta().(*civogo.Client)
		resp, err := client.FindIP(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("ip not found: (%s) %s", rs.Primary.ID, err)
		}

		*ip = *resp

		return nil
	}
}
