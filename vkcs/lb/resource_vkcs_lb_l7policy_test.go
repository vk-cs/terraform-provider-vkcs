package lb_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"

	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/l7policies"
)

func TestAccLBL7Policy_basic(t *testing.T) {
	var l7Policy l7policies.L7Policy

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckLBL7PolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccCheckLbL7PolicyConfigBasic, map[string]string{"TestAccCheckLbL7PolicyConfig": testAccCheckLbL7PolicyConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBL7PolicyExists("vkcs_lb_l7policy.l7policy_1", &l7Policy),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7policy.l7policy_1", "name", "test"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7policy.l7policy_1", "description", "test description"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7policy.l7policy_1", "action", "REJECT"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7policy.l7policy_1", "position", "1"),
					resource.TestMatchResourceAttr(
						"vkcs_lb_l7policy.l7policy_1", "listener_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccCheckLbL7PolicyConfigUpdate1, map[string]string{"TestAccCheckLbL7PolicyConfig": testAccCheckLbL7PolicyConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBL7PolicyExists("vkcs_lb_l7policy.l7policy_1", &l7Policy),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7policy.l7policy_1", "name", "test"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7policy.l7policy_1", "description", "test description"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7policy.l7policy_1", "action", "REDIRECT_TO_URL"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7policy.l7policy_1", "position", "1"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7policy.l7policy_1", "redirect_url", "http://www.example.com"),
					resource.TestMatchResourceAttr(
						"vkcs_lb_l7policy.l7policy_1", "listener_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccCheckLbL7PolicyConfigUpdate2, map[string]string{"TestAccCheckLbL7PolicyConfig": testAccCheckLbL7PolicyConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBL7PolicyExists("vkcs_lb_l7policy.l7policy_1", &l7Policy),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7policy.l7policy_1", "name", "test_updated"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7policy.l7policy_1", "description", ""),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7policy.l7policy_1", "action", "REDIRECT_TO_POOL"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7policy.l7policy_1", "position", "1"),
					resource.TestMatchResourceAttr(
						"vkcs_lb_l7policy.l7policy_1", "listener_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
					resource.TestMatchResourceAttr(
						"vkcs_lb_l7policy.l7policy_1", "redirect_pool_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccCheckLbL7PolicyConfigUpdate3, map[string]string{"TestAccCheckLbL7PolicyConfig": testAccCheckLbL7PolicyConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBL7PolicyExists("vkcs_lb_l7policy.l7policy_1", &l7Policy),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7policy.l7policy_1", "name", "test_updated"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7policy.l7policy_1", "description", ""),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7policy.l7policy_1", "action", "REJECT"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7policy.l7policy_1", "position", "1"),
					resource.TestMatchResourceAttr(
						"vkcs_lb_l7policy.l7policy_1", "listener_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
				),
			},
		},
	})
}

func testAccCheckLBL7PolicyDestroy(s *terraform.State) error {
	config := acctest.AccTestProvider.Meta().(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(acctest.OsRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS load balancing client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_lb_l7policy" {
			continue
		}

		_, err := l7policies.Get(lbClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("L7 Policy still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBL7PolicyExists(n string, l7Policy *l7policies.L7Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.AccTestProvider.Meta().(clients.Config)
		lbClient, err := config.LoadBalancerV2Client(acctest.OsRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS load balancing client: %s", err)
		}

		found, err := l7policies.Get(lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Policy not found")
		}

		*l7Policy = *found

		return nil
	}
}

const testAccCheckLbL7PolicyConfig = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  network_id = vkcs_networking_network.network_1.id
}

resource "vkcs_lb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = vkcs_networking_subnet.subnet_1.id

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

resource "vkcs_lb_listener" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = vkcs_lb_loadbalancer.loadbalancer_1.id
}
`

const testAccCheckLbL7PolicyConfigBasic = `
{{.TestAccCheckLbL7PolicyConfig}}

resource "vkcs_lb_l7policy" "l7policy_1" {
  name         = "test"
  action       = "REJECT"
  description  = "test description"
  position     = 1
  listener_id  = vkcs_lb_listener.listener_1.id
}
`

const testAccCheckLbL7PolicyConfigUpdate1 = `
{{.TestAccCheckLbL7PolicyConfig}}

resource "vkcs_lb_l7policy" "l7policy_1" {
  name         = "test"
  action       = "REDIRECT_TO_URL"
  description  = "test description"
  position     = 1
  listener_id  = vkcs_lb_listener.listener_1.id
  redirect_url = "http://www.example.com"
}
`

const testAccCheckLbL7PolicyConfigUpdate2 = `
{{.TestAccCheckLbL7PolicyConfig}}

resource "vkcs_lb_pool" "pool_1" {
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = vkcs_lb_loadbalancer.loadbalancer_1.id
}

resource "vkcs_lb_l7policy" "l7policy_1" {
  name             = "test_updated"
  action           = "REDIRECT_TO_POOL"
  position         = 1
  listener_id      = vkcs_lb_listener.listener_1.id
  redirect_pool_id = vkcs_lb_pool.pool_1.id
}
`

const testAccCheckLbL7PolicyConfigUpdate3 = `
{{.TestAccCheckLbL7PolicyConfig}}

resource "vkcs_lb_pool" "pool_1" {
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = vkcs_lb_loadbalancer.loadbalancer_1.id
}

resource "vkcs_lb_l7policy" "l7policy_1" {
  name             = "test_updated"
  action           = "REJECT"
  position         = 1
  listener_id      = vkcs_lb_listener.listener_1.id
}
`
