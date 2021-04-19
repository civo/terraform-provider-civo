package civo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCivoLoadBalancer_basic(t *testing.T) {
	datasourceName := "data.civo_loadbalancer.foobar"
	name := acctest.RandomWithPrefix("lb-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoLoadBalancerConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "hostname", name),
					resource.TestCheckResourceAttr(datasourceName, "protocol", "http"),
					resource.TestCheckResourceAttr(datasourceName, "port", "80"),
					resource.TestCheckResourceAttr(datasourceName, "max_request_size", "30"),
					resource.TestCheckResourceAttr(datasourceName, "policy", "round_robin"),
					resource.TestCheckResourceAttr(datasourceName, "health_check_path", "/"),
					resource.TestCheckResourceAttr(datasourceName, "max_conns", "10"),
					resource.TestCheckResourceAttr(datasourceName, "backend.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceCivoLoadBalancerConfig(name string) string {
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

data "civo_loadbalancer" "foobar" {
	hostname = civo_loadbalancer.foobar.hostname
}
`, name, name)
}
