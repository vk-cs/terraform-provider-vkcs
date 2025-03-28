package compute_test

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/attachinterfaces"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/compute"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	iattachinterfaces "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/compute/v2/attachinterfaces"
)

func TestAccComputeInterfaceAttach_basic(t *testing.T) {
	var ai attachinterfaces.Interface

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInterfaceAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInterfaceAttachBasic),
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
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInterfaceAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInterfaceAttachIP),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInterfaceAttachExists("vkcs_compute_interface_attach.ai_1", &ai),
					testAccCheckComputeInterfaceAttachIP(&ai, "192.168.1.100"),
				),
			},
		},
	})
}

func testAccCheckComputeInterfaceAttachDestroy(s *terraform.State) error {
	config := acctest.AccTestProvider.Meta().(clients.Config)
	computeClient, err := config.ComputeV2Client(acctest.OsRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_compute_interface_attach" {
			continue
		}

		instanceID, portID, err := compute.ComputeInterfaceAttachParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = iattachinterfaces.Get(computeClient, instanceID, portID).Extract()
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

		config := acctest.AccTestProvider.Meta().(clients.Config)
		computeClient, err := config.ComputeV2Client(acctest.OsRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS compute client: %s", err)
		}

		instanceID, portID, err := compute.ComputeInterfaceAttachParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, err := iattachinterfaces.Get(computeClient, instanceID, portID).Extract()
		if err != nil {
			return err
		}

		// if found.instanceID != instanceID || found.PortID != portID {
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

const testAccComputeInterfaceAttachBasic = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_networking_port" "port_1" {
  depends_on = ["vkcs_networking_subnet.base"]

  name = "port_1"
  network_id = vkcs_networking_network.base.id
  admin_state_up = "true"
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}

resource "vkcs_compute_interface_attach" "ai_1" {
  instance_id = vkcs_compute_instance.instance_1.id
  port_id = vkcs_networking_port.port_1.id
}
`

const testAccComputeInterfaceAttachIP = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_networking_network" "network_1" {
  name = "network_1"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  network_id = vkcs_networking_network.network_1.id
  cidr = "192.168.1.0/24"
  enable_dhcp = true
  no_gateway = true
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}

resource "vkcs_compute_interface_attach" "ai_1" {
  instance_id = vkcs_compute_instance.instance_1.id
  network_id = vkcs_networking_network.network_1.id
  fixed_ip = "192.168.1.100"
}
`
