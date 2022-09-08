package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/servergroups"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

func TestAccComputeServerGroup_basic(t *testing.T) {
	var sg servergroups.ServerGroup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeServerGroupBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeServerGroupExists("vkcs_compute_servergroup.sg_1", &sg),
					resource.TestCheckResourceAttr(
						"vkcs_compute_servergroup.sg_1", "policies.#", "1"),
					resource.TestCheckResourceAttr(
						"vkcs_compute_servergroup.sg_1", "policies.0", "affinity"),
				),
			},
		},
	})
}

func TestAccComputeServerGroup_affinity(t *testing.T) {
	var instance servers.Server
	var sg servergroups.ServerGroup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeServerGroupAffinity(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeServerGroupExists("vkcs_compute_servergroup.sg_1", &sg),
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckComputeInstanceInServerGroup(&instance, &sg),
					resource.TestCheckResourceAttr(
						"vkcs_compute_servergroup.sg_1", "policies.#", "1"),
					resource.TestCheckResourceAttr(
						"vkcs_compute_servergroup.sg_1", "policies.0", "affinity"),
				),
			},
		},
	})
}

func TestAccComputeServerGroup_soft_affinity(t *testing.T) {
	var instance servers.Server
	var sg servergroups.ServerGroup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeServerGroupSoftAffinity(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeServerGroupExists("vkcs_compute_servergroup.sg_1", &sg),
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckComputeInstanceInServerGroup(&instance, &sg),
					resource.TestCheckResourceAttr(
						"vkcs_compute_servergroup.sg_1", "policies.#", "1"),
					resource.TestCheckResourceAttr(
						"vkcs_compute_servergroup.sg_1", "policies.0", "soft-affinity"),
				),
			},
		},
	})
}

func testAccCheckComputeServerGroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(configer)
	computeClient, err := config.ComputeV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_compute_servergroup" {
			continue
		}

		_, err := servergroups.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("ServerGroup still exists")
		}
	}

	return nil
}

func testAccCheckComputeServerGroupExists(n string, kp *servergroups.ServerGroup) resource.TestCheckFunc {
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
			return fmt.Errorf("Error creating VKCS compute client: %s", err)
		}

		found, err := servergroups.Get(computeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("ServerGroup not found")
		}

		*kp = *found

		return nil
	}
}

func testAccCheckComputeInstanceInServerGroup(instance *servers.Server, sg *servergroups.ServerGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(sg.Members) > 0 {
			for _, m := range sg.Members {
				if m == instance.ID {
					return nil
				}
			}
		}

		return fmt.Errorf("Instance %s is not part of Server Group %s", instance.ID, sg.ID)
	}
}

const testAccComputeServerGroupBasic = `
resource "vkcs_compute_servergroup" "sg_1" {
  name = "sg_1"
  policies = ["affinity"]
}
`

func testAccComputeServerGroupAffinity() string {
	return fmt.Sprintf(`
%s

%s

%s

resource "vkcs_compute_servergroup" "sg_1" {
  name = "sg_1"
  policies = ["affinity"]
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.base"]
  name = "instance_1"
  security_groups = ["default"]
  scheduler_hints {
    group = "${vkcs_compute_servergroup.sg_1.id}"
  }
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`, testAccBaseFlavor, testAccBaseImage, testAccBaseNetwork)
}

func testAccComputeServerGroupSoftAffinity() string {
	return fmt.Sprintf(`
%s

%s

%s

resource "vkcs_compute_servergroup" "sg_1" {
  name = "sg_1"
  policies = ["soft-affinity"]
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.base"]
  name = "instance_1"
  security_groups = ["default"]
  scheduler_hints {
    group = "${vkcs_compute_servergroup.sg_1.id}"
  }
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`, testAccBaseFlavor, testAccBaseImage, testAccBaseNetwork)
}
