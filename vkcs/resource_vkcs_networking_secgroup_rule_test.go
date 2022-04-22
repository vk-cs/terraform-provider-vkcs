package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/rules"
)

func TestAccNetworkingSecGroupRule_basic(t *testing.T) {
	var secgroup1 groups.SecGroup
	var secgroup2 groups.SecGroup
	var secgroupRule1 rules.SecGroupRule
	var secgroupRule2 rules.SecGroupRule

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckNetworking(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingSecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRuleBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupExists(
						"vkcs_networking_secgroup.secgroup_1", &secgroup1),
					testAccCheckNetworkingSecGroupExists(
						"vkcs_networking_secgroup.secgroup_2", &secgroup2),
					testAccCheckNetworkingSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_1", &secgroupRule1),
					testAccCheckNetworkingSecGroupRuleExists(
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

func TestAccNetworkingSecGroupRule_lowerCaseCIDR(t *testing.T) {
	var secgroup1 groups.SecGroup
	var secgroupRule1 rules.SecGroupRule

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckNetworking(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingSecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRuleLowerCaseCidr,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupExists(
						"vkcs_networking_secgroup.secgroup_1", &secgroup1),
					testAccCheckNetworkingSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_1", &secgroupRule1),
					resource.TestCheckResourceAttr(
						"vkcs_networking_secgroup_rule.secgroup_rule_1", "remote_ip_prefix", "2001:558:fc00::/39"),
				),
			},
		},
	})
}

func TestAccNetworkingSecGroupRule_timeout(t *testing.T) {
	var secgroup1 groups.SecGroup
	var secgroup2 groups.SecGroup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckNetworking(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingSecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRuleTimeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupExists(
						"vkcs_networking_secgroup.secgroup_1", &secgroup1),
					testAccCheckNetworkingSecGroupExists(
						"vkcs_networking_secgroup.secgroup_2", &secgroup2),
				),
			},
		},
	})
}

func TestAccNetworkingSecGroupRule_protocols(t *testing.T) {
	var secgroup1 groups.SecGroup
	var secgroupRuleAh rules.SecGroupRule
	var secgroupRuleDccp rules.SecGroupRule
	var secgroupRuleEgp rules.SecGroupRule
	var secgroupRuleEsp rules.SecGroupRule
	var secgroupRuleGre rules.SecGroupRule
	var secgroupRuleIgmp rules.SecGroupRule
	var secgroupRuleIPv6Encap rules.SecGroupRule
	var secgroupRuleIPv6Frag rules.SecGroupRule
	var secgroupRuleIPv6Icmp rules.SecGroupRule
	var secgroupRuleIPv6Nonxt rules.SecGroupRule
	var secgroupRuleIPv6Opts rules.SecGroupRule
	var secgroupRuleIPv6Route rules.SecGroupRule
	var secgroupRuleOspf rules.SecGroupRule
	var secgroupRulePgm rules.SecGroupRule
	var secgroupRuleRsvp rules.SecGroupRule
	var secgroupRuleSctp rules.SecGroupRule
	var secgroupRuleUdplite rules.SecGroupRule
	var secgroupRuleVrrp rules.SecGroupRule

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckNetworking(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingSecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRuleProtocols,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupExists(
						"vkcs_networking_secgroup.secgroup_1", &secgroup1),
					testAccCheckNetworkingSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_ah", &secgroupRuleAh),
					testAccCheckNetworkingSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_dccp", &secgroupRuleDccp),
					testAccCheckNetworkingSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_egp", &secgroupRuleEgp),
					testAccCheckNetworkingSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_esp", &secgroupRuleEsp),
					testAccCheckNetworkingSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_gre", &secgroupRuleGre),
					testAccCheckNetworkingSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_igmp", &secgroupRuleIgmp),
					testAccCheckNetworkingSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_ipv6_encap", &secgroupRuleIPv6Encap),
					testAccCheckNetworkingSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_ipv6_frag", &secgroupRuleIPv6Frag),
					testAccCheckNetworkingSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_ipv6_icmp", &secgroupRuleIPv6Icmp),
					testAccCheckNetworkingSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_ipv6_nonxt", &secgroupRuleIPv6Nonxt),
					testAccCheckNetworkingSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_ipv6_opts", &secgroupRuleIPv6Opts),
					testAccCheckNetworkingSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_ipv6_route", &secgroupRuleIPv6Route),
					testAccCheckNetworkingSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_ospf", &secgroupRuleOspf),
					testAccCheckNetworkingSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_pgm", &secgroupRulePgm),
					testAccCheckNetworkingSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_rsvp", &secgroupRuleRsvp),
					testAccCheckNetworkingSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_sctp", &secgroupRuleSctp),
					testAccCheckNetworkingSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_udplite", &secgroupRuleUdplite),
					testAccCheckNetworkingSecGroupRuleExists(
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
						"vkcs_networking_secgroup_rule.secgroup_rule_ipv6_encap", "protocol", "ipv6-encap"),
					resource.TestCheckResourceAttr(
						"vkcs_networking_secgroup_rule.secgroup_rule_ipv6_frag", "protocol", "ipv6-frag"),
					resource.TestCheckResourceAttr(
						"vkcs_networking_secgroup_rule.secgroup_rule_ipv6_icmp", "protocol", "ipv6-icmp"),
					resource.TestCheckResourceAttr(
						"vkcs_networking_secgroup_rule.secgroup_rule_ipv6_nonxt", "protocol", "ipv6-nonxt"),
					resource.TestCheckResourceAttr(
						"vkcs_networking_secgroup_rule.secgroup_rule_ipv6_opts", "protocol", "ipv6-opts"),
					resource.TestCheckResourceAttr(
						"vkcs_networking_secgroup_rule.secgroup_rule_ipv6_route", "protocol", "ipv6-route"),
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

func TestAccNetworkingSecGroupRule_numericProtocol(t *testing.T) {
	var secgroup1 groups.SecGroup
	var secgroupRule1 rules.SecGroupRule

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckNetworking(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingSecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRuleNumericProtocol,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupExists(
						"vkcs_networking_secgroup.secgroup_1", &secgroup1),
					testAccCheckNetworkingSecGroupRuleExists(
						"vkcs_networking_secgroup_rule.secgroup_rule_1", &secgroupRule1),
					resource.TestCheckResourceAttr(
						"vkcs_networking_secgroup_rule.secgroup_rule_1", "protocol", "6"),
				),
			},
		},
	})
}

func testAccCheckNetworkingSecGroupRuleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(configer)
	networkingClient, err := config.NetworkingV2Client(osRegionName, defaultSDN)
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

func testAccCheckNetworkingSecGroupRuleExists(n string, securityGroupRule *rules.SecGroupRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(configer)
		networkingClient, err := config.NetworkingV2Client(osRegionName, defaultSDN)
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

const testAccNetworkingSecGroupRuleBasic = `
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
  security_group_id = "${vkcs_networking_secgroup.secgroup_1.id}"
	description = "secgroup_rule_1"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_2" {
  direction = "ingress"
  ethertype = "IPv4"
  port_range_max = 80
  port_range_min = 80
  protocol = "tcp"
  remote_group_id = "${vkcs_networking_secgroup.secgroup_1.id}"
  security_group_id = "${vkcs_networking_secgroup.secgroup_2.id}"
}
`

const testAccNetworkingSecGroupRuleLowerCaseCidr = `
resource "vkcs_networking_secgroup" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group rule acceptance test"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_1" {
  direction = "ingress"
  ethertype = "IPv6"
  port_range_max = 22
  port_range_min = 22
  protocol = "tcp"
  remote_ip_prefix = "2001:558:FC00::/39"
  security_group_id = "${vkcs_networking_secgroup.secgroup_1.id}"
}
`

const testAccNetworkingSecGroupRuleTimeout = `
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
  security_group_id = "${vkcs_networking_secgroup.secgroup_1.id}"

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
  remote_group_id = "${vkcs_networking_secgroup.secgroup_1.id}"
  security_group_id = "${vkcs_networking_secgroup.secgroup_2.id}"

  timeouts {
    delete = "5m"
  }
}
`

const testAccNetworkingSecGroupRuleProtocols = `
resource "vkcs_networking_secgroup" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group rule acceptance test"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_ah" {
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "ah"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = "${vkcs_networking_secgroup.secgroup_1.id}"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_dccp" {
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "dccp"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = "${vkcs_networking_secgroup.secgroup_1.id}"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_egp" {
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "egp"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = "${vkcs_networking_secgroup.secgroup_1.id}"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_esp" {
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "esp"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = "${vkcs_networking_secgroup.secgroup_1.id}"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_gre" {
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "gre"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = "${vkcs_networking_secgroup.secgroup_1.id}"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_igmp" {
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "igmp"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = "${vkcs_networking_secgroup.secgroup_1.id}"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_ipv6_encap" {
  direction = "ingress"
  ethertype = "IPv6"
  protocol = "ipv6-encap"
  remote_ip_prefix = "::/0"
  security_group_id = "${vkcs_networking_secgroup.secgroup_1.id}"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_ipv6_frag" {
  direction = "ingress"
  ethertype = "IPv6"
  protocol = "ipv6-frag"
  remote_ip_prefix = "::/0"
  security_group_id = "${vkcs_networking_secgroup.secgroup_1.id}"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_ipv6_icmp" {
  direction = "ingress"
  ethertype = "IPv6"
  protocol = "ipv6-icmp"
  remote_ip_prefix = "::/0"
  security_group_id = "${vkcs_networking_secgroup.secgroup_1.id}"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_ipv6_nonxt" {
  direction = "ingress"
  ethertype = "IPv6"
  protocol = "ipv6-nonxt"
  remote_ip_prefix = "::/0"
  security_group_id = "${vkcs_networking_secgroup.secgroup_1.id}"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_ipv6_opts" {
  direction = "ingress"
  ethertype = "IPv6"
  protocol = "ipv6-opts"
  remote_ip_prefix = "::/0"
  security_group_id = "${vkcs_networking_secgroup.secgroup_1.id}"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_ipv6_route" {
  direction = "ingress"
  ethertype = "IPv6"
  protocol = "ipv6-route"
  remote_ip_prefix = "::/0"
  security_group_id = "${vkcs_networking_secgroup.secgroup_1.id}"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_ospf" {
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "ospf"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = "${vkcs_networking_secgroup.secgroup_1.id}"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_pgm" {
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "pgm"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = "${vkcs_networking_secgroup.secgroup_1.id}"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_rsvp" {
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "rsvp"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = "${vkcs_networking_secgroup.secgroup_1.id}"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_sctp" {
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "sctp"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = "${vkcs_networking_secgroup.secgroup_1.id}"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_udplite" {
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "udplite"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = "${vkcs_networking_secgroup.secgroup_1.id}"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_vrrp" {
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "vrrp"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = "${vkcs_networking_secgroup.secgroup_1.id}"
}
`

const testAccNetworkingSecGroupRuleNumericProtocol = `
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
  security_group_id = "${vkcs_networking_secgroup.secgroup_1.id}"
}
`
