package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLBLoadBalancerDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccLBLoadBalancerDataSourceBasic),
			},
			{
				Config: testAccRenderConfig(testAccLBLoadBalancerDataSourceSource, map[string]string{"TestAccLBLoadBalancerDataSourceBasic": testAccRenderConfig(testAccLBLoadBalancerDataSourceBasic)}),
				Check: resource.ComposeTestCheckFunc(
					testAccLBCheckLoadBalancerDataSourceID("data.vkcs_lb_loadbalancer.source_1"),
					resource.TestCheckResourceAttr("data.vkcs_lb_loadbalancer.source_1", "name", "loadbalancer_1"),
					resource.TestCheckResourceAttrPair("data.vkcs_lb_loadbalancer.source_1", "vip_subnet_id", "vkcs_lb_loadbalancer.loadbalancer_1", "vip_subnet_id"),
					resource.TestCheckResourceAttr("data.vkcs_lb_loadbalancer.source_1", "tags.#", "1"),
				),
			},
		},
	})
}

func testAccLBCheckLoadBalancerDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find loadbalancer data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Loadbalancer data source ID not set")
		}

		return nil
	}
}

const testAccLBLoadBalancerDataSourceBasic = `
{{.BaseNetwork}}

resource "vkcs_lb_loadbalancer" "loadbalancer_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "loadbalancer_1"
  vip_subnet_id = vkcs_networking_subnet.base.id
  tags = ["tag1"]

  timeouts {
	create = "15m"
	update = "15m"
	delete = "15m"
  }
}
`

const testAccLBLoadBalancerDataSourceSource = `
{{.TestAccLBLoadBalancerDataSourceBasic}}

data "vkcs_lb_loadbalancer" "source_1" {
  id = vkcs_lb_loadbalancer.loadbalancer_1.id
}
`
