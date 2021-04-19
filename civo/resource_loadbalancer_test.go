package civo

import (
	"fmt"
	"testing"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var domainName = acctest.RandString(10)

// example.Widget represents a concrete Go type that represents an API resource
func TestAccCivoLoadBalancer_basic(t *testing.T) {
	var loadBalancer civogo.LoadBalancer

	// generate a random name for each test run
	resName := "civo_loadbalancer.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoLoadBalancerConfigBasic(domainName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckCivoLoadBalancerResourceExists(resName, &loadBalancer),
					// verify remote values
					testAccCheckCivoLoadBalancerValues(&loadBalancer, domainName),
					// verify local values
					resource.TestCheckResourceAttr(resName, "protocol", "http"),
					resource.TestCheckResourceAttr(resName, "port", "80"),
					resource.TestCheckResourceAttr(resName, "max_request_size", "30"),
					resource.TestCheckResourceAttr(resName, "policy", "round_robin"),
					resource.TestCheckResourceAttr(resName, "health_check_path", "/"),
					resource.TestCheckResourceAttr(resName, "max_conns", "10"),
				),
			},
		},
	})
}

func TestAccCivoLoadBalancer_update(t *testing.T) {
	var loadBalancer civogo.LoadBalancer

	// generate a random name for each test run
	resName := "civo_loadbalancer.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCivoLoadBalancerConfigBasic(domainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoLoadBalancerResourceExists(resName, &loadBalancer),
					resource.TestCheckResourceAttr(resName, "protocol", "http"),
					resource.TestCheckResourceAttr(resName, "port", "80"),
					resource.TestCheckResourceAttr(resName, "max_request_size", "30"),
					resource.TestCheckResourceAttr(resName, "policy", "round_robin"),
					resource.TestCheckResourceAttr(resName, "health_check_path", "/"),
					resource.TestCheckResourceAttr(resName, "max_conns", "10"),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoLoadBalancerConfigUpdates(domainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoLoadBalancerResourceExists(resName, &loadBalancer),
					testAccCheckCivoLoadBalancerUpdated(&loadBalancer, domainName),
					resource.TestCheckResourceAttr(resName, "protocol", "http"),
					resource.TestCheckResourceAttr(resName, "port", "80"),
					resource.TestCheckResourceAttr(resName, "max_request_size", "50"),
					resource.TestCheckResourceAttr(resName, "policy", "round_robin"),
					resource.TestCheckResourceAttr(resName, "health_check_path", "/"),
					resource.TestCheckResourceAttr(resName, "max_conns", "100"),
				),
			},
		},
	})
}

func testAccCheckCivoLoadBalancerValues(loadBalancer *civogo.LoadBalancer, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if loadBalancer.Hostname != name {
			return fmt.Errorf("bad protocol, expected \"%s\", got: %#v", name, loadBalancer.Hostname)
		}
		return nil
	}
}

// testAccCheckExampleResourceExists queries the API and retrieves the matching Widget.
func testAccCheckCivoLoadBalancerResourceExists(n string, loadBalancer *civogo.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := testAccProvider.Meta().(*civogo.Client)
		resp, err := client.FindLoadBalancer(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("LoadBalancer not found: (%s) %s", rs.Primary.ID, err)
		}

		// If no error, assign the response Widget attribute to the widget pointer
		*loadBalancer = *resp

		// return fmt.Errorf("Domain (%s) not found", rs.Primary.ID)
		return nil
	}
}

func testAccCheckCivoLoadBalancerUpdated(loadBalancer *civogo.LoadBalancer, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if loadBalancer.Hostname != fmt.Sprintf("rename-%s", name) {
			return fmt.Errorf("bad protocol, expected \"%s\", got: %#v", fmt.Sprintf("rename-%s", name), loadBalancer.Hostname)
		}
		return nil
	}
}

func testAccCheckCivoLoadBalancerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*civogo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "civo_loadbalancer" {
			continue
		}

		_, err := client.FindLoadBalancer(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("LoadBlanacer still exists")
		}
	}

	return nil
}

func testAccCheckCivoLoadBalancerConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "civo_instance" "vm" {
	hostname = "instance-%s"
}

resource "civo_loadbalancer" "foobar" {
	hostname = "%s"
	protocol = "http"
	port = 80
	max_request_size = 30
	policy = "round_robin"
	health_check_path = "/"
	max_conns = 10
	fail_timeout = 40
	depends_on = [civo_instance.vm]

	backend {
		instance_id = civo_instance.vm.id
		protocol =  "http"
		port = 80
	}
}
`, name, name)
}

func testAccCheckCivoLoadBalancerConfigUpdates(name string) string {
	return fmt.Sprintf(`
resource "civo_instance" "vm" {
	hostname = "instance-%s"
}

resource "civo_loadbalancer" "foobar" {
	hostname = "rename-%s"
	protocol = "http"
	port = 80
	max_request_size = 50
	policy = "round_robin"
	health_check_path = "/"
	max_conns = 100
	fail_timeout = 40
	depends_on = [civo_instance.vm]

	backend {
		instance_id = civo_instance.vm.id
		protocol =  "http"
		port = 80
	}
}
`, name, name)
}
