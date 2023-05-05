package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
)

func TestAccNetworkingPortSecGroupAssociate_update(t *testing.T) {
	var port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			// enforce = false
			{ // preconfig
				Config: testAccRenderConfig(testAccNetworkingPortSecGroupAssociatePreTest),
			},
			{ // step 0
				Config: testAccRenderConfig(testAccNetworkingPortSecGroupAssociateManifestUpdate0, map[string]string{"TestAccNetworkingPortSecGroupAssociate": testAccRenderConfig(testAccNetworkingPortSecGroupAssociate, map[string]string{"TestAccNetworkingPortSecGroupAssociatePreTest": testAccRenderConfig(testAccNetworkingPortSecGroupAssociatePreTest)})}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortSecGroupAssociateExists("vkcs_networking_port_secgroup_associate.port_1", &port),
					testAccCheckNetworkingPortSecGroupAssociateCountSecurityGroups(&port, 3),
				),
			},
			{ // step 1
				Config: testAccRenderConfig(testAccNetworkingPortSecGroupAssociateManifestUpdate1, map[string]string{"TestAccNetworkingPortSecGroupAssociate": testAccRenderConfig(testAccNetworkingPortSecGroupAssociate, map[string]string{"TestAccNetworkingPortSecGroupAssociatePreTest": testAccRenderConfig(testAccNetworkingPortSecGroupAssociatePreTest)})}), // unset user defined security groups only
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortSecGroupAssociateExists("data.vkcs_networking_port.hidden_port_1", &port),
					testAccCheckNetworkingPortSecGroupAssociateCountSecurityGroups(&port, 2),
				),
			},
			// enforce = true
			{ // step 2
				Config: testAccRenderConfig(testAccNetworkingPortSecGroupAssociateManifestUpdate2, map[string]string{"TestAccNetworkingPortSecGroupAssociate": testAccRenderConfig(testAccNetworkingPortSecGroupAssociate, map[string]string{"TestAccNetworkingPortSecGroupAssociatePreTest": testAccRenderConfig(testAccNetworkingPortSecGroupAssociatePreTest)})}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortSecGroupAssociateExists("vkcs_networking_port_secgroup_associate.port_1", &port),
					testAccCheckNetworkingPortSecGroupAssociateCountSecurityGroups(&port, 1),
				),
			},
			{ // step 3
				Config: testAccRenderConfig(testAccNetworkingPortSecGroupAssociateManifestUpdate3, map[string]string{"TestAccNetworkingPortSecGroupAssociate": testAccRenderConfig(testAccNetworkingPortSecGroupAssociate, map[string]string{"TestAccNetworkingPortSecGroupAssociatePreTest": testAccRenderConfig(testAccNetworkingPortSecGroupAssociatePreTest)})}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortSecGroupAssociateExists("vkcs_networking_port_secgroup_associate.port_1", &port),
					testAccCheckNetworkingPortSecGroupAssociateCountSecurityGroups(&port, 3),
				),
			},
			{ // step 4
				Config: testAccRenderConfig(testAccNetworkingPortSecGroupAssociateManifestUpdate4, map[string]string{"TestAccNetworkingPortSecGroupAssociate": testAccRenderConfig(testAccNetworkingPortSecGroupAssociate, map[string]string{"TestAccNetworkingPortSecGroupAssociatePreTest": testAccRenderConfig(testAccNetworkingPortSecGroupAssociatePreTest)})}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortSecGroupAssociateExists("vkcs_networking_port_secgroup_associate.port_1", &port),
					testAccCheckNetworkingPortSecGroupAssociateCountSecurityGroups(&port, 1),
				),
			},
			{ // step 5
				Config: testAccRenderConfig(testAccNetworkingPortSecGroupAssociateManifestUpdate5, map[string]string{"TestAccNetworkingPortSecGroupAssociate": testAccRenderConfig(testAccNetworkingPortSecGroupAssociate, map[string]string{"TestAccNetworkingPortSecGroupAssociatePreTest": testAccRenderConfig(testAccNetworkingPortSecGroupAssociatePreTest)})}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortSecGroupAssociateExists("vkcs_networking_port_secgroup_associate.port_1", &port),
					testAccCheckNetworkingPortSecGroupAssociateCountSecurityGroups(&port, 0),
				),
			},
			{ // step 6
				Config: testAccRenderConfig(testAccNetworkingPortSecGroupAssociateManifestUpdate6, map[string]string{"TestAccNetworkingPortSecGroupAssociate": testAccRenderConfig(testAccNetworkingPortSecGroupAssociate, map[string]string{"TestAccNetworkingPortSecGroupAssociatePreTest": testAccRenderConfig(testAccNetworkingPortSecGroupAssociatePreTest)})}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortSecGroupAssociateExists("vkcs_networking_port_secgroup_associate.port_1", &port),
					testAccCheckNetworkingPortSecGroupAssociateCountSecurityGroups(&port, 1),
				),
			},
			// enforce = false
			{ // step 7
				Config: testAccRenderConfig(testAccNetworkingPortSecGroupAssociateManifestUpdate7, map[string]string{"TestAccNetworkingPortSecGroupAssociate": testAccRenderConfig(testAccNetworkingPortSecGroupAssociate, map[string]string{"TestAccNetworkingPortSecGroupAssociatePreTest": testAccRenderConfig(testAccNetworkingPortSecGroupAssociatePreTest)})}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortSecGroupAssociateExists("vkcs_networking_port_secgroup_associate.port_1", &port),
					testAccCheckNetworkingPortSecGroupAssociateCountSecurityGroups(&port, 1),
				),
			},
			{ // step 8
				Config: testAccRenderConfig(testAccNetworkingPortSecGroupAssociateManifestUpdate8, map[string]string{"TestAccNetworkingPortSecGroupAssociate": testAccRenderConfig(testAccNetworkingPortSecGroupAssociate, map[string]string{"TestAccNetworkingPortSecGroupAssociatePreTest": testAccRenderConfig(testAccNetworkingPortSecGroupAssociatePreTest)})}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortSecGroupAssociateExists("vkcs_networking_port_secgroup_associate.port_1", &port),
					testAccCheckNetworkingPortSecGroupAssociateCountSecurityGroups(&port, 3),
				),
			},
			{ // step 9
				Config: testAccRenderConfig(testAccNetworkingPortSecGroupAssociateManifestUpdate9, map[string]string{"TestAccNetworkingPortSecGroupAssociate": testAccRenderConfig(testAccNetworkingPortSecGroupAssociate, map[string]string{"TestAccNetworkingPortSecGroupAssociatePreTest": testAccRenderConfig(testAccNetworkingPortSecGroupAssociatePreTest)})}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortSecGroupAssociateExists("vkcs_networking_port_secgroup_associate.port_1", &port),
					testAccCheckNetworkingPortSecGroupAssociateCountSecurityGroups(&port, 1),
				),
			},
			{ // step 10
				Config: testAccRenderConfig(testAccNetworkingPortSecGroupAssociateManifestUpdate10, map[string]string{"TestAccNetworkingPortSecGroupAssociate": testAccRenderConfig(testAccNetworkingPortSecGroupAssociate, map[string]string{"TestAccNetworkingPortSecGroupAssociatePreTest": testAccRenderConfig(testAccNetworkingPortSecGroupAssociatePreTest)})}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortSecGroupAssociateExists("vkcs_networking_port_secgroup_associate.port_1", &port),
					testAccCheckNetworkingPortSecGroupAssociateCountSecurityGroups(&port, 0),
				),
			},
			{ // step 11
				Config: testAccRenderConfig(testAccNetworkingPortSecGroupAssociateManifestUpdate11, map[string]string{"TestAccNetworkingPortSecGroupAssociate": testAccRenderConfig(testAccNetworkingPortSecGroupAssociate, map[string]string{"TestAccNetworkingPortSecGroupAssociatePreTest": testAccRenderConfig(testAccNetworkingPortSecGroupAssociatePreTest)})}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortSecGroupAssociateExists("vkcs_networking_port_secgroup_associate.port_1", &port),
					testAccCheckNetworkingPortSecGroupAssociateCountSecurityGroups(&port, 1),
				),
			},
			{ // step 12
				Config: testAccRenderConfig(testAccNetworkingPortSecGroupAssociateManifestUpdate12, map[string]string{"TestAccNetworkingPortSecGroupAssociate": testAccRenderConfig(testAccNetworkingPortSecGroupAssociate, map[string]string{"TestAccNetworkingPortSecGroupAssociatePreTest": testAccRenderConfig(testAccNetworkingPortSecGroupAssociatePreTest)})}), // cleanup all the ports
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortSecGroupAssociateExists("data.vkcs_networking_port.hidden_port_1", &port),
					testAccCheckNetworkingPortSecGroupAssociateCountSecurityGroups(&port, 0),
				),
			},
		},
	})
}

func testAccCheckNetworkingPortSecGroupAssociateExists(n string, port *ports.Port) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(clients.Config)
		networkingClient, err := config.NetworkingV2Client(osRegionName, networking.DefaultSDN)
		if err != nil {
			return fmt.Errorf("Error creating VKCS networking client: %s", err)
		}

		found, err := ports.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Port not found")
		}

		*port = *found

		return nil
	}
}

func testAccCheckNetworkingPortSecGroupAssociateCountSecurityGroups(port *ports.Port, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(port.SecurityGroups) != expected {
			return fmt.Errorf("Expected %d Security Groups, got %d", expected, len(port.SecurityGroups))
		}

		return nil
	}
}

const testAccNetworkingPortSecGroupAssociatePreTest = `
{{.BaseNetwork}}

resource "vkcs_networking_secgroup" "default_1" {
  name = "default_1"
}

resource "vkcs_networking_secgroup" "default_2" {
  name = "default_2"
}

resource "vkcs_networking_port" "hidden_port_1" {
  name = "hidden_port"
  admin_state_up = true
  network_id = vkcs_networking_network.base.id
  security_group_ids = [vkcs_networking_secgroup.default_1.id, vkcs_networking_secgroup.default_2.id]
}
`

const testAccNetworkingPortSecGroupAssociate = `
{{.TestAccNetworkingPortSecGroupAssociatePreTest}}

resource "vkcs_networking_secgroup" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "vkcs_networking_secgroup" "secgroup_2" {
  name = "secgroup_2"
  description = "terraform security group acceptance test"
}

data "vkcs_networking_secgroup" "default_1" {
  name = "default_1"
}

data "vkcs_networking_secgroup" "default_2" {
  name = "default_2"
}

data "vkcs_networking_port" "hidden_port_1" {
  name = "hidden_port"
}
`

const testAccNetworkingPortSecGroupAssociateManifestUpdate0 = `
{{.TestAccNetworkingPortSecGroupAssociate}}

resource "vkcs_networking_port_secgroup_associate" "port_1" {
  port_id = data.vkcs_networking_port.hidden_port_1.id
  enforce = "false"
  security_group_ids = [
    vkcs_networking_secgroup.secgroup_1.id,
  ]
}
`

const testAccNetworkingPortSecGroupAssociateManifestUpdate1 = `
{{.TestAccNetworkingPortSecGroupAssociate}}
`

const testAccNetworkingPortSecGroupAssociateManifestUpdate2 = `
{{.TestAccNetworkingPortSecGroupAssociate}}

resource "vkcs_networking_port_secgroup_associate" "port_1" {
  port_id = data.vkcs_networking_port.hidden_port_1.id
  enforce = "true"
  security_group_ids = [
    vkcs_networking_secgroup.secgroup_1.id,
  ]
}
`

const testAccNetworkingPortSecGroupAssociateManifestUpdate3 = `
{{.TestAccNetworkingPortSecGroupAssociate}}

resource "vkcs_networking_port_secgroup_associate" "port_1" {
  port_id = data.vkcs_networking_port.hidden_port_1.id
  enforce = "true"
  security_group_ids = [
    vkcs_networking_secgroup.secgroup_1.id,
    vkcs_networking_secgroup.secgroup_2.id,
    data.vkcs_networking_secgroup.default_2.id,
  ]
}
`

const testAccNetworkingPortSecGroupAssociateManifestUpdate4 = `
{{.TestAccNetworkingPortSecGroupAssociate}}

resource "vkcs_networking_port_secgroup_associate" "port_1" {
  port_id = data.vkcs_networking_port.hidden_port_1.id
  enforce = "true"
  security_group_ids = [
    vkcs_networking_secgroup.secgroup_2.id,
  ]
}
`

const testAccNetworkingPortSecGroupAssociateManifestUpdate5 = `
{{.TestAccNetworkingPortSecGroupAssociate}}

data "vkcs_networking_port" "port_1" {
  port_id = vkcs_networking_port_secgroup_associate.port_1.id
}

resource "vkcs_networking_port_secgroup_associate" "port_1" {
  port_id = data.vkcs_networking_port.hidden_port_1.id
  enforce = "true"
  security_group_ids = []
}
`

const testAccNetworkingPortSecGroupAssociateManifestUpdate6 = `
{{.TestAccNetworkingPortSecGroupAssociate}}

resource "vkcs_networking_port_secgroup_associate" "port_1" {
  port_id = data.vkcs_networking_port.hidden_port_1.id
  enforce = "true"
  security_group_ids = [
    data.vkcs_networking_secgroup.default_2.id,
  ]
}
`

const testAccNetworkingPortSecGroupAssociateManifestUpdate7 = `
{{.TestAccNetworkingPortSecGroupAssociate}}

resource "vkcs_networking_port_secgroup_associate" "port_1" {
  port_id = data.vkcs_networking_port.hidden_port_1.id
  enforce = "false"
  security_group_ids = [
    vkcs_networking_secgroup.secgroup_1.id,
  ]
}
`

const testAccNetworkingPortSecGroupAssociateManifestUpdate8 = `
{{.TestAccNetworkingPortSecGroupAssociate}}

resource "vkcs_networking_port_secgroup_associate" "port_1" {
  port_id = data.vkcs_networking_port.hidden_port_1.id
  enforce = "false"
  security_group_ids = [
    vkcs_networking_secgroup.secgroup_1.id,
    vkcs_networking_secgroup.secgroup_2.id,
    data.vkcs_networking_secgroup.default_2.id,
  ]
}
`

const testAccNetworkingPortSecGroupAssociateManifestUpdate9 = `
{{.TestAccNetworkingPortSecGroupAssociate}}

resource "vkcs_networking_port_secgroup_associate" "port_1" {
  port_id = data.vkcs_networking_port.hidden_port_1.id
  enforce = "false"
  security_group_ids = [
    vkcs_networking_secgroup.secgroup_2.id,
  ]
}
`

const testAccNetworkingPortSecGroupAssociateManifestUpdate10 = `
{{.TestAccNetworkingPortSecGroupAssociate}}

resource "vkcs_networking_port_secgroup_associate" "port_1" {
  port_id = data.vkcs_networking_port.hidden_port_1.id
  enforce = "false"
  security_group_ids = []
}
`

const testAccNetworkingPortSecGroupAssociateManifestUpdate11 = `
{{.TestAccNetworkingPortSecGroupAssociate}}

resource "vkcs_networking_port_secgroup_associate" "port_1" {
  port_id = data.vkcs_networking_port.hidden_port_1.id
  enforce = "false"
  security_group_ids = [
    data.vkcs_networking_secgroup.default_2.id,
  ]
}
`

const testAccNetworkingPortSecGroupAssociateManifestUpdate12 = `
{{.TestAccNetworkingPortSecGroupAssociate}}
`
