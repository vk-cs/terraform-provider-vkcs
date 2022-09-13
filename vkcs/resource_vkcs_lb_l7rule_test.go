package vkcs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/l7policies"
)

func TestAccLBL7Rule_basic(t *testing.T) {
	var l7rule l7policies.Rule

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBL7RuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccCheckLbL7RuleConfigBasic, map[string]string{"TestAccCheckLbL7RuleConfig": testAccCheckLbL7RuleConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBL7RuleExists("vkcs_lb_l7rule.l7rule_1", &l7rule),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "type", "PATH"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "compare_type", "EQUAL_TO"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "value", "/api"),
					resource.TestMatchResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "listener_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
					resource.TestMatchResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "l7policy_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
				),
			},
			{
				Config: testAccRenderConfig(testAccCheckLbL7RuleConfigUpdate1, map[string]string{"TestAccCheckLbL7RuleConfig": testAccCheckLbL7RuleConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBL7RuleExists("vkcs_lb_l7rule.l7rule_1", &l7rule),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "type", "HOST_NAME"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "compare_type", "EQUAL_TO"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "value", "www.example.com"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "invert", "true"),
					resource.TestMatchResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "listener_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
					resource.TestMatchResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "l7policy_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
				),
			},
			{
				Config: testAccRenderConfig(testAccCheckLbL7RuleConfigUpdate2, map[string]string{"TestAccCheckLbL7RuleConfig": testAccCheckLbL7RuleConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBL7RuleExists("vkcs_lb_l7rule.l7rule_1", &l7rule),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "type", "HOST_NAME"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "compare_type", "EQUAL_TO"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "value", "www.example.com"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "invert", "true"),
				),
			},
			{
				Config: testAccRenderConfig(testAccCheckLbL7RuleConfigUpdate3, map[string]string{"TestAccCheckLbL7RuleConfig": testAccCheckLbL7RuleConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBL7RuleExists("vkcs_lb_l7rule.l7rule_1", &l7rule),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "type", "HEADER"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "compare_type", "EQUAL_TO"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "key", "Host"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "value", "www.example.com"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "invert", "false"),
				),
			},
			{
				Config: testAccRenderConfig(testAccCheckLbL7RuleConfigUpdate4, map[string]string{"TestAccCheckLbL7RuleConfig": testAccCheckLbL7RuleConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBL7RuleExists("vkcs_lb_l7rule.l7rule_1", &l7rule),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "type", "HOST_NAME"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "compare_type", "EQUAL_TO"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "key", ""),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "value", "www.example.com"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "invert", "false"),
				),
			},
			{
				Config: testAccRenderConfig(testAccCheckLbL7RuleConfigUpdate5, map[string]string{"TestAccCheckLbL7RuleConfig": testAccCheckLbL7RuleConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBL7RuleExists("vkcs_lb_l7rule.l7rule_1", &l7rule),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "type", "COOKIE"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "compare_type", "EQUAL_TO"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "key", "X-Ref"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "value", "foo"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "invert", "false"),
				),
			},
			{
				Config: testAccRenderConfig(testAccCheckLbL7RuleConfigUpdate6, map[string]string{"TestAccCheckLbL7RuleConfig": testAccCheckLbL7RuleConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBL7RuleExists("vkcs_lb_l7rule.l7rule_1", &l7rule),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "type", "PATH"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "compare_type", "STARTS_WITH"),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "key", ""),
					resource.TestCheckResourceAttr(
						"vkcs_lb_l7rule.l7rule_1", "value", "/images"),
				),
			},
		},
	})
}

func testAccCheckLBL7RuleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config)
	lbClient, err := config.LoadBalancerV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS load balancing client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_lb_l7rule" {
			continue
		}

		l7policyID := ""
		for k, v := range rs.Primary.Attributes {
			if k == "l7policy_id" {
				l7policyID = v
				break
			}
		}

		if l7policyID == "" {
			return fmt.Errorf("Unable to find l7policy_id")
		}

		_, err := l7policies.GetRule(lbClient, l7policyID, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("L7 Rule still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBL7RuleExists(n string, l7rule *l7policies.Rule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config)
		lbClient, err := config.LoadBalancerV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS load balancing client: %s", err)
		}

		l7policyID := ""
		for k, v := range rs.Primary.Attributes {
			if k == "l7policy_id" {
				l7policyID = v
				break
			}
		}

		if l7policyID == "" {
			return fmt.Errorf("Unable to find l7policy_id")
		}

		found, err := l7policies.GetRule(lbClient, l7policyID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Policy not found")
		}

		*l7rule = *found

		return nil
	}
}

const testAccCheckLbL7RuleConfig = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
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

resource "vkcs_lb_l7policy" "l7policy_1" {
  name         = "test"
  action       = "REDIRECT_TO_URL"
  description  = "test description"
  position     = 1
  listener_id  = vkcs_lb_listener.listener_1.id
  redirect_url = "http://www.example.com"
}
`

const testAccCheckLbL7RuleConfigBasic = `
{{.TestAccCheckLbL7RuleConfig}}

resource "vkcs_lb_l7rule" "l7rule_1" {
  l7policy_id  = vkcs_lb_l7policy.l7policy_1.id
  type         = "PATH"
  compare_type = "EQUAL_TO"
  value        = "/api"
}
`

const testAccCheckLbL7RuleConfigUpdate1 = `
{{.TestAccCheckLbL7RuleConfig}}

resource "vkcs_lb_l7rule" "l7rule_1" {
  l7policy_id  = vkcs_lb_l7policy.l7policy_1.id
  type         = "HOST_NAME"
  compare_type = "EQUAL_TO"
  value        = "www.example.com"
  invert       = true
}
`

const testAccCheckLbL7RuleConfigUpdate2 = `
{{.TestAccCheckLbL7RuleConfig}}

resource "vkcs_lb_l7rule" "l7rule_1" {
  l7policy_id  = vkcs_lb_l7policy.l7policy_1.id
  type         = "HOST_NAME"
  compare_type = "EQUAL_TO"
  value        = "www.example.com"
  invert       = true
}
`

const testAccCheckLbL7RuleConfigUpdate3 = `
{{.TestAccCheckLbL7RuleConfig}}

resource "vkcs_lb_l7rule" "l7rule_1" {
  l7policy_id  = vkcs_lb_l7policy.l7policy_1.id
  type         = "HEADER"
  compare_type = "EQUAL_TO"
  key          = "Host"
  value        = "www.example.com"
}
`

const testAccCheckLbL7RuleConfigUpdate4 = `
{{.TestAccCheckLbL7RuleConfig}}

resource "vkcs_lb_l7rule" "l7rule_1" {
  l7policy_id  = vkcs_lb_l7policy.l7policy_1.id
  type         = "HOST_NAME"
  compare_type = "EQUAL_TO"
  value        = "www.example.com"
}
`

const testAccCheckLbL7RuleConfigUpdate5 = `
{{.TestAccCheckLbL7RuleConfig}}

resource "vkcs_lb_l7rule" "l7rule_1" {
  l7policy_id  = vkcs_lb_l7policy.l7policy_1.id
  type         = "COOKIE"
  compare_type = "EQUAL_TO"
  key          = "X-Ref"
  value        = "foo"
}
`

const testAccCheckLbL7RuleConfigUpdate6 = `
{{.TestAccCheckLbL7RuleConfig}}

resource "vkcs_lb_l7rule" "l7rule_1" {
  l7policy_id  = vkcs_lb_l7policy.l7policy_1.id
  type         = "PATH"
  compare_type = "STARTS_WITH"
  value        = "/images"
}
`
