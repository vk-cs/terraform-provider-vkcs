package firewall_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	igroups "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/firewall/v2/groups"
)

func TestAccFirewallSecGroup_basic(t *testing.T) {
	var securityGroup groups.SecGroup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccFirewallCheckSecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallSecGroupBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccFirewallCheckSecGroupExists("vkcs_networking_secgroup.secgroup_1", &securityGroup),
					testAccFirewallCheckSecGroupRuleCount(&securityGroup, 2),
				),
			},
			{
				Config: testAccFirewallSecGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr("vkcs_networking_secgroup.secgroup_1", "id", &securityGroup.ID),
					resource.TestCheckResourceAttr("vkcs_networking_secgroup.secgroup_1", "name", "security_group_2"),
				),
			},
		},
	})
}

func TestAccFirewallSecGroup_noDefaultRules(t *testing.T) {
	var securityGroup groups.SecGroup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccFirewallCheckSecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallSecGroupNoDefaultRules,
				Check: resource.ComposeTestCheckFunc(
					testAccFirewallCheckSecGroupExists(
						"vkcs_networking_secgroup.secgroup_1", &securityGroup),
					testAccFirewallCheckSecGroupRuleCount(&securityGroup, 0),
				),
			},
		},
	})
}

func TestAccFirewallSecGroup_timeout(t *testing.T) {
	var securityGroup groups.SecGroup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccFirewallCheckSecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallSecGroupTimeout,
				Check: resource.ComposeTestCheckFunc(
					testAccFirewallCheckSecGroupExists(
						"vkcs_networking_secgroup.secgroup_1", &securityGroup),
				),
			},
		},
	})
}

func testAccFirewallCheckSecGroupDestroy(s *terraform.State) error {
	config := acctest.AccTestProvider.Meta().(clients.Config)
	networkingClient, err := config.NetworkingV2Client(acctest.OsRegionName, networking.DefaultSDN)
	if err != nil {
		return fmt.Errorf("Error creating VKCS networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_networking_secgroup" {
			continue
		}

		_, err := igroups.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Security group still exists")
		}
	}

	return nil
}

func testAccFirewallCheckSecGroupExists(n string, sg *groups.SecGroup) resource.TestCheckFunc {
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

		found, err := igroups.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Security group not found")
		}

		*sg = *found

		return nil
	}
}

func testAccFirewallCheckSecGroupRuleCount(sg *groups.SecGroup, count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(sg.Rules) == count {
			return nil
		}

		return fmt.Errorf("Unexpected number of rules in group %s. Expected %d, got %d",
			sg.ID, count, len(sg.Rules))
	}
}

const testAccFirewallSecGroupBasic = `
resource "vkcs_networking_secgroup" "secgroup_1" {
  name = "security_group"
  description = "terraform security group acceptance test"
}
`

const testAccFirewallSecGroupUpdate = `
resource "vkcs_networking_secgroup" "secgroup_1" {
  name = "security_group_2"
  description = "terraform security group acceptance test"
}
`

const testAccFirewallSecGroupNoDefaultRules = `
resource "vkcs_networking_secgroup" "secgroup_1" {
	name = "security_group_1"
	description = "terraform security group acceptance test"
	delete_default_rules = true
}
`

const testAccFirewallSecGroupTimeout = `
resource "vkcs_networking_secgroup" "secgroup_1" {
  name = "security_group"
  description = "terraform security group acceptance test"

  timeouts {
    delete = "5m"
  }
}
`
