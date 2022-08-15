package vkcs

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccComputeFloatingIPAssociate_basic(t *testing.T) {
	var instance servers.Server
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeFloatingIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeFloatingIPAssociateBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckNetworkingFloatingIPExists("vkcs_networking_floatingip.fip_1", &fip),
					testAccCheckComputeFloatingIPAssociateAssociated(&fip, &instance, 1),
				),
			},
			{
				Config: testAccComputeFloatingIPAssociateUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckNetworkingFloatingIPExists("vkcs_networking_floatingip.fip_1", &fip),
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
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeFloatingIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeFloatingIPAssociateFixedIP(),
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
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeFloatingIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeFloatingIPAssociateAttachNew1(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckNetworkingFloatingIPExists("vkcs_networking_floatingip.fip_1", &floatingIP1),
					testAccCheckNetworkingFloatingIPExists("vkcs_networking_floatingip.fip_2", &floatingIP2),
					testAccCheckComputeFloatingIPAssociateAssociated(&floatingIP1, &instance, 1),
				),
			},
			{
				Config: testAccComputeFloatingIPAssociateAttachNew2(),
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
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeFloatingIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeFloatingIPAssociateWaitUntilAssociated(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckNetworkingFloatingIPExists("vkcs_networking_floatingip.fip_1", &fip),
					testAccCheckComputeFloatingIPAssociateAssociated(&fip, &instance, 1),
				),
			},
		},
	})
}

func testAccCheckComputeFloatingIPAssociateDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(configer)
	computeClient, err := config.ComputeV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_compute_floatingip_associate" {
			continue
		}

		floatingIP, instanceID, _, err := parseComputeFloatingIPAssociateID(rs.Primary.ID)
		if err != nil {
			return err
		}

		instance, err := servers.Get(computeClient, instanceID).Extract()
		if err != nil {
			// If the error is a 404, then the instance does not exist,
			// and therefore the floating IP cannot be associated to it.
			if _, ok := err.(gophercloud.ErrDefault404); ok {
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
		config := testAccProvider.Meta().(configer)
		computeClient, err := config.ComputeV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS compute client: %s", err)
		}

		newInstance, err := servers.Get(computeClient, instance.ID).Extract()
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

func testAccComputeFloatingIPAssociateBasic() string {
	return fmt.Sprintf(`
%s

%s

%s

resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}

resource "vkcs_networking_floatingip" "fip_1" {
  pool = data.vkcs_networking_network.extnet.name
}

resource "vkcs_compute_floatingip_associate" "fip_1" {
  floating_ip = "${vkcs_networking_floatingip.fip_1.address}"
  instance_id = "${vkcs_compute_instance.instance_1.id}"
}
`, testAccBaseFlavor, testAccBaseImage, testAccBaseNetwork)
}

func testAccComputeFloatingIPAssociateUpdate() string {
	return fmt.Sprintf(`
%s

%s

%s

resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}

resource "vkcs_networking_floatingip" "fip_1" {
  pool = data.vkcs_networking_network.extnet.name
  description = "test"
}

resource "vkcs_compute_floatingip_associate" "fip_1" {
  floating_ip = "${vkcs_networking_floatingip.fip_1.address}"
  instance_id = "${vkcs_compute_instance.instance_1.id}"
}
`, testAccBaseFlavor, testAccBaseImage, testAccBaseNetwork)
}

func testAccComputeFloatingIPAssociateFixedIP() string {
	return fmt.Sprintf(`
%s

%s

%s

resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}

resource "vkcs_networking_floatingip" "fip_1" {
	pool = data.vkcs_networking_network.extnet.name
}

resource "vkcs_compute_floatingip_associate" "fip_1" {
  floating_ip = "${vkcs_networking_floatingip.fip_1.address}"
  instance_id = "${vkcs_compute_instance.instance_1.id}"
  fixed_ip = "${vkcs_compute_instance.instance_1.access_ip_v4}"
}
`, testAccBaseFlavor, testAccBaseImage, testAccBaseNetwork)
}

func testAccComputeFloatingIPAssociateAttachNew1() string {
	return fmt.Sprintf(`
%s

%s

%s

resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}

resource "vkcs_networking_floatingip" "fip_1" {
  pool = data.vkcs_networking_network.extnet.name
}

resource "vkcs_networking_floatingip" "fip_2" {
  pool = data.vkcs_networking_network.extnet.name
}

resource "vkcs_compute_floatingip_associate" "fip_1" {
  floating_ip = "${vkcs_networking_floatingip.fip_1.address}"
  instance_id = "${vkcs_compute_instance.instance_1.id}"
}
`, testAccBaseFlavor, testAccBaseImage, testAccBaseNetwork)
}

func testAccComputeFloatingIPAssociateAttachNew2() string {
	return fmt.Sprintf(`
%s

%s

%s

resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}

resource "vkcs_networking_floatingip" "fip_1" {
  pool = data.vkcs_networking_network.extnet.name
}

resource "vkcs_networking_floatingip" "fip_2" {
  pool = data.vkcs_networking_network.extnet.name
}

resource "vkcs_compute_floatingip_associate" "fip_1" {
  floating_ip = "${vkcs_networking_floatingip.fip_2.address}"
  instance_id = "${vkcs_compute_instance.instance_1.id}"
}
`, testAccBaseFlavor, testAccBaseImage, testAccBaseNetwork)
}

func testAccComputeFloatingIPAssociateWaitUntilAssociated() string {
	return fmt.Sprintf(`
%s

%s

%s

resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}

resource "vkcs_networking_floatingip" "fip_1" {
  pool = data.vkcs_networking_network.extnet.name
}

resource "vkcs_compute_floatingip_associate" "fip_1" {
  floating_ip = "${vkcs_networking_floatingip.fip_1.address}"
  instance_id = "${vkcs_compute_instance.instance_1.id}"

  wait_until_associated = true
}
`, testAccBaseFlavor, testAccBaseImage, testAccBaseNetwork)
}
