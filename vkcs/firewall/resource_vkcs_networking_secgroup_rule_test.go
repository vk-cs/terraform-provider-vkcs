package firewall_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/rules"
)

func TestAccFirewallSecGroupRule_basic(t *testing.T) {
	var secgroup1 groups.SecGroup
	var secgroup2 groups.SecGroup
	var secgroupRule1 rules.SecGroupRule
	var secgroupRule2 rules.SecGroupRule

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccFirewallCheckSecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallSecGroupRuleBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccFirewallCheckSecGroupExists(
						"vkcs_networking_secgroup.secgroup_1", &secgroup1),
					testAccFirewallCheckSecGroupExists(
						"vkcs_networking_secgroup.secgroup_2", &secgroup2),
					testAccFirewallCheckSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_1", &secgroupRule1),
					testAccFirewallCheckSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_2", &secgroupRule2),
					resource.TestCheckResourceAttr(
						"vkcs_networking_secgroup_rule.secgroup_rule_1", "description", "secgroup_rule_1"),
					resource.TestCheckResourceAttr(
						"vkcs_networking_secgroup_rule.secgroup_rule_2", "description", ""),
				),
			},
		},
	})
}

func TestAccFirewallSecGroupRule_timeout(t *testing.T) {
	var secgroup1 groups.SecGroup
	var secgroup2 groups.SecGroup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccFirewallCheckSecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallSecGroupRuleTimeout,
				Check: resource.ComposeTestCheckFunc(
					testAccFirewallCheckSecGroupExists(
						"vkcs_networking_secgroup.secgroup_1", &secgroup1),
					testAccFirewallCheckSecGroupExists(
						"vkcs_networking_secgroup.secgroup_2", &secgroup2),
				),
			},
		},
	})
}

func TestAccFirewallSecGroupRule_protocols(t *testing.T) {
	var secgroup1 groups.SecGroup
	var secgroupRuleAh rules.SecGroupRule
	var secgroupRuleDccp rules.SecGroupRule
	var secgroupRuleEgp rules.SecGroupRule
	var secgroupRuleEsp rules.SecGroupRule
	var secgroupRuleGre rules.SecGroupRule
	var secgroupRuleIgmp rules.SecGroupRule
	var secgroupRuleOspf rules.SecGroupRule
	var secgroupRulePgm rules.SecGroupRule
	var secgroupRuleRsvp rules.SecGroupRule
	var secgroupRuleSctp rules.SecGroupRule
	var secgroupRuleUdplite rules.SecGroupRule
	var secgroupRuleVrrp rules.SecGroupRule

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccFirewallCheckSecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallSecGroupRuleProtocols,
				Check: resource.ComposeTestCheckFunc(
					testAccFirewallCheckSecGroupExists(
						"vkcs_networking_secgroup.secgroup_1", &secgroup1),
					testAccFirewallCheckSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_ah", &secgroupRuleAh),
					testAccFirewallCheckSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_dccp", &secgroupRuleDccp),
					testAccFirewallCheckSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_egp", &secgroupRuleEgp),
					testAccFirewallCheckSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_esp", &secgroupRuleEsp),
					testAccFirewallCheckSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_gre", &secgroupRuleGre),
					testAccFirewallCheckSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_igmp", &secgroupRuleIgmp),
					testAccFirewallCheckSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_ospf", &secgroupRuleOspf),
					testAccFirewallCheckSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_pgm", &secgroupRulePgm),
					testAccFirewallCheckSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_rsvp", &secgroupRuleRsvp),
					testAccFirewallCheckSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_sctp", &secgroupRuleSctp),
					testAccFirewallCheckSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_udplite", &secgroupRuleUdplite),
					testAccFirewallCheckSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_vrrp", &secgroupRuleVrrp),
					resource.TestCheckResourceAttr(
						"vkcs_networking_secgroup_rule.secgroup_rule_ah", "protocol", "ah"),
					resource.TestCheckResourceAttr(
						"vkcs_networking_secgroup_rule.secgroup_rule_dccp", "protocol", "dccp"),
					resource.TestCheckResourceAttr(
						"vkcs_networking_secgroup_rule.secgroup_rule_egp", "protocol", "egp"),
					resource.TestCheckResourceAttr(
						"vkcs_networking_secgroup_rule.secgroup_rule_esp", "protocol", "esp"),
					resource.TestCheckResourceAttr(
						"vkcs_networking_secgroup_rule.secgroup_rule_gre", "protocol", "gre"),
					resource.TestCheckResourceAttr(
						"vkcs_networking_secgroup_rule.secgroup_rule_igmp", "protocol", "igmp"),
					resource.TestCheckResourceAttr(
						"vkcs_networking_secgroup_rule.secgroup_rule_ospf", "protocol", "ospf"),
					resource.TestCheckResourceAttr(
						"vkcs_networking_secgroup_rule.secgroup_rule_pgm", "protocol", "pgm"),
					resource.TestCheckResourceAttr(
						"vkcs_networking_secgroup_rule.secgroup_rule_rsvp", "protocol", "rsvp"),
					resource.TestCheckResourceAttr(
						"vkcs_networking_secgroup_rule.secgroup_rule_sctp", "protocol", "sctp"),
					resource.TestCheckResourceAttr(
						"vkcs_networking_secgroup_rule.secgroup_rule_udplite", "protocol", "udplite"),
					resource.TestCheckResourceAttr(
						"vkcs_networking_secgroup_rule.secgroup_rule_vrrp", "protocol", "vrrp"),
				),
			},
		},
	})
}

func TestAccFirewallSecGroupRule_numericProtocol(t *testing.T) {
	var secgroup1 groups.SecGroup
	var secgroupRule1 rules.SecGroupRule

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccFirewallCheckSecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallSecGroupRuleNumericProtocol,
				Check: resource.ComposeTestCheckFunc(
					testAccFirewallCheckSecGroupExists(
						"vkcs_networking_secgroup.secgroup_1", &secgroup1),
					testAccFirewallCheckSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_1", &secgroupRule1),
					resource.TestCheckResourceAttr(
						"vkcs_networking_secgroup_rule.secgroup_rule_1", "protocol", "6"),
				),
			},
		},
	})
}

func testAccFirewallCheckSecGroupRuleDestroy(s *terraform.State) error {
	config := acctest.AccTestProvider.Meta().(clients.Config)
	networkingClient, err := config.NetworkingV2Client(acctest.OsRegionName, networking.DefaultSDN)
	if err != nil {
		return fmt.Errorf("Error creating VKCS networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_networking_secgroup_rule" {
			continue
		}

		_, err := rules.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Security group rule still exists")
		}
	}

	return nil
}

func testAccFirewallCheckSecGroupRuleExists(n string, securityGroupRule *rules.SecGroupRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.AccTestProvider.Meta().(clients.Config)
		networkingClient, err := config.NetworkingV2Client(acctest.OsRegionName, networking.DefaultSDN)
		if err != nil {
			return fmt.Errorf("Error creating VKCS networking client: %s", err)
		}

		found, err := rules.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Security group rule not found")
		}

		*securityGroupRule = *found

		return nil
	}
}

const testAccFirewallSecGroupRuleBasic = `
resource "vkcs_networking_secgroup" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group rule acceptance test"
}

resource "vkcs_networking_secgroup" "secgroup_2" {
  name = "secgroup_2"
  description = "terraform security group rule acceptance test"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_1" {
  direction = "ingress"
  ethertype = "IPv4"
  port_range_max = 22
  port_range_min = 22
  protocol = "tcp"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = vkcs_networking_secgroup.secgroup_1.id
	description = "secgroup_rule_1"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_2" {
  direction = "ingress"
  ethertype = "IPv4"
  port_range_max = 80
  port_range_min = 80
  protocol = "tcp"
  remote_group_id = vkcs_networking_secgroup.secgroup_1.id
  security_group_id = vkcs_networking_secgroup.secgroup_2.id
}
`

const testAccFirewallSecGroupRuleTimeout = `
resource "vkcs_networking_secgroup" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group rule acceptance test"
}

resource "vkcs_networking_secgroup" "secgroup_2" {
  name = "secgroup_2"
  description = "terraform security group rule acceptance test"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_1" {
  direction = "ingress"
  ethertype = "IPv4"
  port_range_max = 22
  port_range_min = 22
  protocol = "tcp"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = vkcs_networking_secgroup.secgroup_1.id

  timeouts {
    delete = "5m"
  }
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_2" {
  direction = "ingress"
  ethertype = "IPv4"
  port_range_max = 80
  port_range_min = 80
  protocol = "tcp"
  remote_group_id = vkcs_networking_secgroup.secgroup_1.id
  security_group_id = vkcs_networking_secgroup.secgroup_2.id

  timeouts {
    delete = "5m"
  }
}
`

const testAccFirewallSecGroupRuleProtocols = `
resource "vkcs_networking_secgroup" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group rule acceptance test"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_ah" {
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "ah"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = vkcs_networking_secgroup.secgroup_1.id
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_dccp" {
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "dccp"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = vkcs_networking_secgroup.secgroup_1.id
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_egp" {
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "egp"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = vkcs_networking_secgroup.secgroup_1.id
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_esp" {
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "esp"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = vkcs_networking_secgroup.secgroup_1.id
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_gre" {
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "gre"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = vkcs_networking_secgroup.secgroup_1.id
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_igmp" {
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "igmp"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = vkcs_networking_secgroup.secgroup_1.id
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_ospf" {
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "ospf"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = vkcs_networking_secgroup.secgroup_1.id
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_pgm" {
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "pgm"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = vkcs_networking_secgroup.secgroup_1.id
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_rsvp" {
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "rsvp"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = vkcs_networking_secgroup.secgroup_1.id
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_sctp" {
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "sctp"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = vkcs_networking_secgroup.secgroup_1.id
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_udplite" {
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "udplite"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = vkcs_networking_secgroup.secgroup_1.id
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_vrrp" {
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "vrrp"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = vkcs_networking_secgroup.secgroup_1.id
}
`

const testAccFirewallSecGroupRuleNumericProtocol = `
resource "vkcs_networking_secgroup" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group rule acceptance test"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_1" {
  direction = "ingress"
  ethertype = "IPv4"
  port_range_max = 22
  port_range_min = 22
  protocol = "6"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = vkcs_networking_secgroup.secgroup_1.id
}
`
