package vkcs

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
)

func TestAccNetworkingPortSecGroupAssociate_update(t *testing.T) {
	var port ports.Port

	if os.Getenv("TF_ACC") != "" {
		hiddenPort, err := testAccCheckNetworkingPortSecGroupCreatePort(t, "hidden_port", true)
		if err != nil {
			t.Fatal(err)
		}
		defer testAccCheckNetworkingPortSecGroupDeletePort(t, hiddenPort) //nolint:errcheck
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			// enforce = false
			{ // step 0
				Config: testAccRenderConfig(testAccNetworkingPortSecGroupAssociateManifestUpdate0, map[string]string{"TestAccNetworkingPortSecGroupAssociate": testAccNetworkingPortSecGroupAssociate}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortSecGroupAssociateExists("vkcs_networking_port_secgroup_associate.port_1", &port),
					testAccCheckNetworkingPortSecGroupAssociateCountSecurityGroups(&port, 3),
				),
			},
			{ // step 1
				Config: testAccRenderConfig(testAccNetworkingPortSecGroupAssociateManifestUpdate1, map[string]string{"TestAccNetworkingPortSecGroupAssociate": testAccNetworkingPortSecGroupAssociate}), // unset user defined security groups only
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortSecGroupAssociateExists("data.vkcs_networking_port.hidden_port_1", &port),
					testAccCheckNetworkingPortSecGroupAssociateCountSecurityGroups(&port, 2),
				),
			},
			// enforce = true
			{ // step 2
				Config: testAccRenderConfig(testAccNetworkingPortSecGroupAssociateManifestUpdate2, map[string]string{"TestAccNetworkingPortSecGroupAssociate": testAccNetworkingPortSecGroupAssociate}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortSecGroupAssociateExists("vkcs_networking_port_secgroup_associate.port_1", &port),
					testAccCheckNetworkingPortSecGroupAssociateCountSecurityGroups(&port, 1),
				),
			},
			{ // step 3
				Config: testAccRenderConfig(testAccNetworkingPortSecGroupAssociateManifestUpdate3, map[string]string{"TestAccNetworkingPortSecGroupAssociate": testAccNetworkingPortSecGroupAssociate}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortSecGroupAssociateExists("vkcs_networking_port_secgroup_associate.port_1", &port),
					testAccCheckNetworkingPortSecGroupAssociateCountSecurityGroups(&port, 3),
				),
			},
			{ // step 4
				Config: testAccRenderConfig(testAccNetworkingPortSecGroupAssociateManifestUpdate4, map[string]string{"TestAccNetworkingPortSecGroupAssociate": testAccNetworkingPortSecGroupAssociate}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortSecGroupAssociateExists("vkcs_networking_port_secgroup_associate.port_1", &port),
					testAccCheckNetworkingPortSecGroupAssociateCountSecurityGroups(&port, 1),
				),
			},
			{ // step 5
				Config: testAccRenderConfig(testAccNetworkingPortSecGroupAssociateManifestUpdate5, map[string]string{"TestAccNetworkingPortSecGroupAssociate": testAccNetworkingPortSecGroupAssociate}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortSecGroupAssociateExists("vkcs_networking_port_secgroup_associate.port_1", &port),
					testAccCheckNetworkingPortSecGroupAssociateCountSecurityGroups(&port, 0),
				),
			},
			{ // step 6
				Config: testAccRenderConfig(testAccNetworkingPortSecGroupAssociateManifestUpdate6, map[string]string{"TestAccNetworkingPortSecGroupAssociate": testAccNetworkingPortSecGroupAssociate}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortSecGroupAssociateExists("vkcs_networking_port_secgroup_associate.port_1", &port),
					testAccCheckNetworkingPortSecGroupAssociateCountSecurityGroups(&port, 1),
				),
			},
			// enforce = false
			{ // step 7
				Config: testAccRenderConfig(testAccNetworkingPortSecGroupAssociateManifestUpdate7, map[string]string{"TestAccNetworkingPortSecGroupAssociate": testAccNetworkingPortSecGroupAssociate}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortSecGroupAssociateExists("vkcs_networking_port_secgroup_associate.port_1", &port),
					testAccCheckNetworkingPortSecGroupAssociateCountSecurityGroups(&port, 1),
				),
			},
			{ // step 8
				Config: testAccRenderConfig(testAccNetworkingPortSecGroupAssociateManifestUpdate8, map[string]string{"TestAccNetworkingPortSecGroupAssociate": testAccNetworkingPortSecGroupAssociate}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortSecGroupAssociateExists("vkcs_networking_port_secgroup_associate.port_1", &port),
					testAccCheckNetworkingPortSecGroupAssociateCountSecurityGroups(&port, 3),
				),
			},
			{ // step 9
				Config: testAccRenderConfig(testAccNetworkingPortSecGroupAssociateManifestUpdate9, map[string]string{"TestAccNetworkingPortSecGroupAssociate": testAccNetworkingPortSecGroupAssociate}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortSecGroupAssociateExists("vkcs_networking_port_secgroup_associate.port_1", &port),
					testAccCheckNetworkingPortSecGroupAssociateCountSecurityGroups(&port, 1),
				),
			},
			{ // step 10
				Config: testAccRenderConfig(testAccNetworkingPortSecGroupAssociateManifestUpdate10, map[string]string{"TestAccNetworkingPortSecGroupAssociate": testAccNetworkingPortSecGroupAssociate}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortSecGroupAssociateExists("vkcs_networking_port_secgroup_associate.port_1", &port),
					testAccCheckNetworkingPortSecGroupAssociateCountSecurityGroups(&port, 0),
				),
			},
			{ // step 11
				Config: testAccRenderConfig(testAccNetworkingPortSecGroupAssociateManifestUpdate11, map[string]string{"TestAccNetworkingPortSecGroupAssociate": testAccNetworkingPortSecGroupAssociate}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortSecGroupAssociateExists("vkcs_networking_port_secgroup_associate.port_1", &port),
					testAccCheckNetworkingPortSecGroupAssociateCountSecurityGroups(&port, 1),
				),
			},
			{ // step 12
				Config: testAccRenderConfig(testAccNetworkingPortSecGroupAssociateManifestUpdate12, map[string]string{"TestAccNetworkingPortSecGroupAssociate": testAccNetworkingPortSecGroupAssociate}), // cleanup all the ports
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortSecGroupAssociateExists("data.vkcs_networking_port.hidden_port_1", &port),
					testAccCheckNetworkingPortSecGroupAssociateCountSecurityGroups(&port, 0),
				),
			},
		},
	})
}

func testAccCheckNetworkingPortSecGroupCreatePort(t *testing.T, portName string, defaultSecGroups bool) (*ports.Port, error) {
	config, err := testAccAuthFromEnv()
	if err != nil {
		return nil, err
	}

	client, err := config.NetworkingV2Client(osRegionName, defaultSDN)
	if err != nil {
		return nil, err
	}

	createNetOpts := networks.CreateOpts{
		Name:         "test",
		AdminStateUp: gophercloud.Enabled,
	}

	network, err := networks.Create(client, createNetOpts).Extract()
	if err != nil {
		return nil, err
	}

	t.Logf("Network %s created", network.ID)

	var securityGroups []string
	if defaultSecGroups {
		// create default security groups
		createSecGroupOpts := groups.CreateOpts{
			Name: "default_1",
		}

		secGroup1, err := groups.Create(client, createSecGroupOpts).Extract()
		if err != nil {
			return nil, err
		}

		t.Logf("Default security group 1 %s created", secGroup1.ID)

		createSecGroupOpts.Name = "default_2"

		secGroup2, err := groups.Create(client, createSecGroupOpts).Extract()
		if err != nil {
			return nil, err
		}

		t.Logf("Default security group 2 %s created", secGroup2.ID)

		// reversed order, just in case
		securityGroups = append(securityGroups, secGroup2.ID)
		securityGroups = append(securityGroups, secGroup1.ID)
	}

	// create port with default security groups assigned
	createOpts := ports.CreateOpts{
		NetworkID:      network.ID,
		Name:           portName,
		SecurityGroups: &securityGroups,
		AdminStateUp:   gophercloud.Enabled,
	}

	port, err := ports.Create(client, createOpts).Extract()
	if err != nil {
		nErr := networks.Delete(client, network.ID).ExtractErr()
		if nErr != nil {
			return nil, fmt.Errorf("Unable to create port (%s) and delete network (%s: %s)", err, network.ID, nErr)
		}
		return nil, err
	}

	t.Logf("Port %s created", port.ID)

	return port, nil
}

func testAccCheckNetworkingPortSecGroupDeletePort(t *testing.T, port *ports.Port) error {
	config, err := testAccAuthFromEnv()
	if err != nil {
		return err
	}

	client, err := config.NetworkingV2Client(osRegionName, defaultSDN)
	if err != nil {
		return err
	}

	err = ports.Delete(client, port.ID).ExtractErr()
	if err != nil {
		return err
	}

	t.Logf("Port %s deleted", port.ID)

	// delete default security groups
	for _, secGroupID := range port.SecurityGroups {
		err = groups.Delete(client, secGroupID).ExtractErr()
		if err != nil {
			return err
		}
		t.Logf("Default security group %s deleted", secGroupID)
	}

	err = networks.Delete(client, port.NetworkID).ExtractErr()
	if err != nil {
		return err
	}

	t.Logf("Network %s deleted", port.NetworkID)

	return nil
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

		config := testAccProvider.Meta().(configer)
		networkingClient, err := config.NetworkingV2Client(osRegionName, defaultSDN)
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

const testAccNetworkingPortSecGroupAssociate = `
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
