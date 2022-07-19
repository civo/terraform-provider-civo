package civo

import (
	"fmt"
	"testing"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCivoReservedIP_basic(t *testing.T) {
	var ip civogo.IP

	// generate a random name for each test run
	resName := "civo_reserved_ip.foobar"
	var nameIP = acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoReservedIPDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoReservedIPConfigBasic(nameIP),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckCivoReservedIPResourceExists(resName, &ip),
					// verify remote values
					testAccCheckCivoReservedIPValues(&ip, nameIP),
					// verify local values
					resource.TestCheckResourceAttr(resName, "name", nameIP),
				),
			},
		},
	})
}

func TestAccCivoReservedIP_update(t *testing.T) {
	var ip civogo.IP

	// generate a random name for each test run
	resName := "civo_reserved_ip.foobar"
	var nameIP = acctest.RandomWithPrefix("rename-tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoReservedIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCivoReservedIPConfigUpdates(nameIP),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoReservedIPResourceExists(resName, &ip),
					testAccCheckCivoReservedIPValues(&ip, nameIP),
					resource.TestCheckResourceAttr(resName, "name", nameIP),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoReservedIPConfigUpdates(nameIP),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoReservedIPResourceExists(resName, &ip),
					testAccCheckCivoReservedIPUpdated(&ip, nameIP),
					resource.TestCheckResourceAttr(resName, "name", nameIP),
				),
			},
		},
	})
}

func testAccCheckCivoReservedIPValues(ip *civogo.IP, name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if ip.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, ip.Name)
		}
		return nil
	}
}

// testAccCheckCivoReservedIPResourceExists queries the API and retrieves the matching Widget.
func testAccCheckCivoReservedIPResourceExists(n string, ip *civogo.IP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := testAccProvider.Meta().(*civogo.Client)
		resp, err := client.FindIP(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("IP not found: (%s) %s", rs.Primary.ID, err)
		}

		*ip = *resp

		return nil
	}
}

func testAccCheckCivoReservedIPUpdated(ip *civogo.IP, name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if ip.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, ip.Name)
		}
		return nil
	}
}

func testAccCheckCivoReservedIPDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*civogo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "civo_reserved_ip" {
			continue
		}

		_, err := client.FindNetwork(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("IP still exists")
		}
	}

	return nil
}

func testAccCheckCivoReservedIPConfigBasic(label string) string {
	return fmt.Sprintf(`
resource "civo_reserved_ip" "foobar" {
	name = "%s"
}`, label)
}

func testAccCheckCivoReservedIPConfigUpdates(label string) string {
	return fmt.Sprintf(`
resource "civo_reserved_ip" "foobar" {
	name = "%s"
}`, label)
}
