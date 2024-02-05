package acceptance

import (
	"fmt"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// CivoInstanceDestroy is used to destroy the instance created during the test
func CivoInstanceDestroy(s *terraform.State) error {
	client := TestAccProvider.Meta().(*civogo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "civo_instance" {
			continue
		}

		_, err := client.GetInstance(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("instance still exists")
		}
	}

	return nil
}

// CivoInstanceResourceExists queries the API and retrieves the matching Widget.
func CivoInstanceResourceExists(n string, instance *civogo.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := TestAccProvider.Meta().(*civogo.Client)
		resp, err := client.GetInstance(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("instance not found: (%s) %s", rs.Primary.ID, err)
		}

		// If no error, assign the response Widget attribute to the widget pointer
		*instance = *resp

		// return fmt.Errorf("Domain (%s) not found", rs.Primary.ID)
		return nil
	}
}
