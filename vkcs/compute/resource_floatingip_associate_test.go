package compute_test

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/compute"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	iservers "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/compute/v2/servers"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	ifloatingips "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking/v2/floatingips"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func TestAccComputeFloatingIPAssociate_basic(t *testing.T) {
	var instance servers.Server
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeFloatingIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeFloatingIPAssociateBasic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckNetworkingFloatingIPExists("vkcs_networking_floatingip.fip_1", &fip),
					testAccCheckComputeFloatingIPAssociateAssociated(&fip, &instance, 1),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccComputeFloatingIPAssociateUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckComputeFloatingIPAssociateAssociated(&fip, &instance, 1),
				),
			},
		},
	})
}

func TestAccComputeV2FloatingIPAssociate_fixedIP(t *testing.T) {
	var instance servers.Server
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeFloatingIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeFloatingIPAssociateFixedIP),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckNetworkingFloatingIPExists("vkcs_networking_floatingip.fip_1", &fip),
					testAccCheckComputeFloatingIPAssociateAssociated(&fip, &instance, 1),
				),
			},
		},
	})
}

func TestAccComputeFloatingIPAssociate_attachNew(t *testing.T) {
	var instance servers.Server
	var floatingIP1 floatingips.FloatingIP
	var floatingIP2 floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeFloatingIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeFloatingIPAssociateAttachNew1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckNetworkingFloatingIPExists("vkcs_networking_floatingip.fip_1", &floatingIP1),
					testAccCheckNetworkingFloatingIPExists("vkcs_networking_floatingip.fip_2", &floatingIP2),
					testAccCheckComputeFloatingIPAssociateAssociated(&floatingIP1, &instance, 1),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccComputeFloatingIPAssociateAttachNew2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckNetworkingFloatingIPExists("vkcs_networking_floatingip.fip_1", &floatingIP1),
					testAccCheckNetworkingFloatingIPExists("vkcs_networking_floatingip.fip_2", &floatingIP2),
					testAccCheckComputeFloatingIPAssociateAssociated(&floatingIP2, &instance, 1),
				),
			},
		},
	})
}

func TestAccComputeFloatingIPAssociate_waitUntilAssociated(t *testing.T) {
	var instance servers.Server
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeFloatingIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeFloatingIPAssociateWaitUntilAssociated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckNetworkingFloatingIPExists("vkcs_networking_floatingip.fip_1", &fip),
					testAccCheckComputeFloatingIPAssociateAssociated(&fip, &instance, 1),
				),
			},
		},
	})
}

func testAccCheckNetworkingFloatingIPExists(n string, kp *floatingips.FloatingIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.AccTestProvider.Meta().(clients.Config)
		networkClient, err := config.NetworkingV2Client(acctest.OsRegionName, networking.DefaultSDN)
		if err != nil {
			return fmt.Errorf("Error creating VKCS networking client: %s", err)
		}

		found, err := ifloatingips.Get(networkClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Floating IP not found")
		}

		*kp = *found

		return nil
	}
}

func testAccCheckComputeFloatingIPAssociateDestroy(s *terraform.State) error {
	config := acctest.AccTestProvider.Meta().(clients.Config)
	computeClient, err := config.ComputeV2Client(acctest.OsRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_compute_floatingip_associate" {
			continue
		}

		floatingIP, instanceID, _, err := compute.ParseComputeFloatingIPAssociateID(rs.Primary.ID)
		if err != nil {
			return err
		}

		instance, err := iservers.Get(computeClient, instanceID).Extract()
		if err != nil {
			// If the error is a 404, then the instance does not exist,
			// and therefore the floating IP cannot be associated to it.
			if errutil.IsNotFound(err) {
				return nil
			}
			return err
		}

		// But if the instance still exists, then walk through its known addresses
		// and see if there's a floating IP.
		for _, networkAddresses := range instance.Addresses {
			for _, element := range networkAddresses.([]interface{}) {
				address := element.(map[string]interface{})
				if address["OS-EXT-IPS:type"] == "floating" {
					return fmt.Errorf("Floating IP %s is still attached to instance %s", floatingIP, instanceID)
				}
			}
		}
	}

	return nil
}

func testAccCheckComputeFloatingIPAssociateAssociated(
	fip *floatingips.FloatingIP, instance *servers.Server, n int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.AccTestProvider.Meta().(clients.Config)
		computeClient, err := config.ComputeV2Client(acctest.OsRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS compute client: %s", err)
		}

		newInstance, err := iservers.Get(computeClient, instance.ID).Extract()
		if err != nil {
			return err
		}

		// Walk through the instance's addresses and find the match
		i := 0
		for _, networkAddresses := range newInstance.Addresses {
			i++
			if i != n {
				continue
			}
			for _, element := range networkAddresses.([]interface{}) {
				address := element.(map[string]interface{})
				if address["OS-EXT-IPS:type"] == "floating" && address["addr"] == fip.FloatingIP {
					return nil
				}
			}
		}
		return fmt.Errorf("Floating IP %s was not attached to instance %s", fip.FloatingIP, instance.ID)
	}
}

const testAccComputeFloatingIPAssociateBasic = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = [vkcs_networking_router_interface.base]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}

resource "vkcs_networking_floatingip" "fip_1" {
  pool = "{{.ExtNetName}}"
}

resource "vkcs_compute_floatingip_associate" "fip_1" {
  floating_ip = vkcs_networking_floatingip.fip_1.address
  instance_id = vkcs_compute_instance.instance_1.id
}
`

const testAccComputeFloatingIPAssociateUpdate = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}

resource "vkcs_networking_floatingip" "fip_1" {
  pool = "{{.ExtNetName}}"
  description = "test"
}

resource "vkcs_compute_floatingip_associate" "fip_1" {
  floating_ip = vkcs_networking_floatingip.fip_1.address
  instance_id = vkcs_compute_instance.instance_1.id
}
`

const testAccComputeFloatingIPAssociateFixedIP = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}

resource "vkcs_networking_floatingip" "fip_1" {
	pool = "{{.ExtNetName}}"
}

resource "vkcs_compute_floatingip_associate" "fip_1" {
  floating_ip = vkcs_networking_floatingip.fip_1.address
  instance_id = vkcs_compute_instance.instance_1.id
  fixed_ip = vkcs_compute_instance.instance_1.access_ip_v4
}
`

const testAccComputeFloatingIPAssociateAttachNew1 = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}

resource "vkcs_networking_floatingip" "fip_1" {
  pool = "{{.ExtNetName}}"
}

resource "vkcs_networking_floatingip" "fip_2" {
  pool = "{{.ExtNetName}}"
}

resource "vkcs_compute_floatingip_associate" "fip_1" {
  floating_ip = vkcs_networking_floatingip.fip_1.address
  instance_id = vkcs_compute_instance.instance_1.id
}
`

const testAccComputeFloatingIPAssociateAttachNew2 = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}

resource "vkcs_networking_floatingip" "fip_1" {
  pool = "{{.ExtNetName}}"
}

resource "vkcs_networking_floatingip" "fip_2" {
  pool = "{{.ExtNetName}}"
}

resource "vkcs_compute_floatingip_associate" "fip_1" {
  floating_ip = vkcs_networking_floatingip.fip_2.address
  instance_id = vkcs_compute_instance.instance_1.id
}
`

const testAccComputeFloatingIPAssociateWaitUntilAssociated = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}

resource "vkcs_networking_floatingip" "fip_1" {
  pool = "{{.ExtNetName}}"
}

resource "vkcs_compute_floatingip_associate" "fip_1" {
  floating_ip = vkcs_networking_floatingip.fip_1.address
  instance_id = vkcs_compute_instance.instance_1.id

  wait_until_associated = true
}
`
