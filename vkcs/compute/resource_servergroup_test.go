package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/servergroups"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	iservergroups "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/compute/v2/servergroups"
)

func TestAccComputeServerGroup_basic(t *testing.T) {
	var sg servergroups.ServerGroup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
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
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeServerGroupAffinity),
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
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeServerGroupSoftAffinity),
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
	config := acctest.AccTestProvider.Meta().(clients.Config)
	computeClient, err := config.ComputeV2Client(acctest.OsRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_compute_servergroup" {
			continue
		}

		_, err := iservergroups.Get(computeClient, rs.Primary.ID).Extract()
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

		config := acctest.AccTestProvider.Meta().(clients.Config)
		computeClient, err := config.ComputeV2Client(acctest.OsRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS compute client: %s", err)
		}

		found, err := iservergroups.Get(computeClient, rs.Primary.ID).Extract()
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

const testAccComputeServerGroupAffinity = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_servergroup" "sg_1" {
  name = "sg_1"
  policies = ["affinity"]
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  scheduler_hints {
    group = vkcs_compute_servergroup.sg_1.id
  }
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeServerGroupSoftAffinity = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_servergroup" "sg_1" {
  name = "sg_1"
  policies = ["soft-affinity"]
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  scheduler_hints {
    group = vkcs_compute_servergroup.sg_1.id
  }
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`
