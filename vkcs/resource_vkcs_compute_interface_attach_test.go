package vkcs

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/attachinterfaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccComputeInterfaceAttach_basic(t *testing.T) {
	var ai attachinterfaces.Interface

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckCompute(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInterfaceAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInterfaceAttachBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInterfaceAttachExists("vkcs_compute_interface_attach.ai_1", &ai),
				),
			},
		},
	})
}

func TestAccComputeInterfaceAttach_IP(t *testing.T) {
	var ai attachinterfaces.Interface

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckCompute(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInterfaceAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInterfaceAttachIP(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInterfaceAttachExists("vkcs_compute_interface_attach.ai_1", &ai),
					testAccCheckComputeInterfaceAttachIP(&ai, "192.168.1.100"),
				),
			},
		},
	})
}

func testAccCheckComputeInterfaceAttachDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(configer)
	computeClient, err := config.ComputeV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_compute_interface_attach" {
			continue
		}

		instanceID, portID, err := computeInterfaceAttachParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = attachinterfaces.Get(computeClient, instanceID, portID).Extract()
		if err == nil {
			return fmt.Errorf("Volume attachment still exists")
		}
	}

	return nil
}

func testAccCheckComputeInterfaceAttachExists(n string, ai *attachinterfaces.Interface) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(configer)
		computeClient, err := config.ComputeV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack compute client: %s", err)
		}

		instanceID, portID, err := computeInterfaceAttachParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, err := attachinterfaces.Get(computeClient, instanceID, portID).Extract()
		if err != nil {
			return err
		}

		//if found.instanceID != instanceID || found.PortID != portID {
		if found.PortID != portID {
			return fmt.Errorf("InterfaceAttach not found")
		}

		*ai = *found

		return nil
	}
}

func testAccCheckComputeInterfaceAttachIP(
	ai *attachinterfaces.Interface, ip string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, i := range ai.FixedIPs {
			if i.IPAddress == ip {
				return nil
			}
		}
		return fmt.Errorf("Requested ip (%s) does not exist on port", ip)
	}
}

func testAccComputeInterfaceAttachBasic() string {
	return fmt.Sprintf(`
resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  network_id = "%s"
  admin_state_up = "true"
}

resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}

resource "vkcs_compute_interface_attach" "ai_1" {
  instance_id = "${vkcs_compute_instance.instance_1.id}"
  port_id = "${vkcs_networking_port.port_1.id}"
}
`, osNetworkID, osNetworkID)
}

func testAccComputeInterfaceAttachIP() string {
	return fmt.Sprintf(`
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  network_id = "${vkcs_networking_network.network_1.id}"
  cidr = "192.168.1.0/24"
  ip_version = 4
  enable_dhcp = true
  no_gateway = true
}

resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}

resource "vkcs_compute_interface_attach" "ai_1" {
  instance_id = "${vkcs_compute_instance.instance_1.id}"
  network_id = "${vkcs_networking_network.network_1.id}"
  fixed_ip = "192.168.1.100"
}
`, osNetworkID)
}
