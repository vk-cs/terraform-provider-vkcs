package compute_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/servergroups"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	iservergroups "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/compute/v2/servergroups"
)

func TestAccComputeInstance_serverGroupMembership(t *testing.T) {
	var instance servers.Server
	var sg1, sg2 servergroups.ServerGroup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceServerGroupCreate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckComputeServerGroupExists("vkcs_compute_servergroup.sg_1", &sg1),
					testAccCheckComputeInstanceInServerGroup(&instance, &sg1),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceServerGroupMove),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckComputeServerGroupExists("vkcs_compute_servergroup.sg_2", &sg2),
					testAccCheckComputeInstanceInServerGroup(&instance, &sg2),
					resource.TestCheckResourceAttrPair(
						"vkcs_compute_instance.instance_1", "server_group_id",
						"vkcs_compute_servergroup.sg_2", "id"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceServerGroupRemoved),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckComputeServerGroupExists("vkcs_compute_servergroup.sg_2", &sg2),
					testAccCheckComputeInstanceNotInServerGroup(&instance, &sg2),
				),
			},
		},
	})
}

func TestAccComputeInstance_serverGroupMembershipActive(t *testing.T) {
	var instance servers.Server
	var sg1 servergroups.ServerGroup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceServerGroupCreate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckComputeServerGroupExists("vkcs_compute_servergroup.sg_1", &sg1),
					testAccCheckComputeInstanceInServerGroup(&instance, &sg1),
				),
			},
			{
				Config:      acctest.AccTestRenderConfig(testAccComputeInstanceServerGroupMoveActive),
				ExpectError: regexp.MustCompile(`stop_before_server_group_change`),
			},
		},
	})
}

func TestAccComputeInstance_serverGroupMembershipActiveWithOption(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceServerGroupCreate),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceServerGroupMoveActiveWithOption),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckComputeInstanceState(&instance, "active"),
					resource.TestCheckResourceAttrPair(
						"vkcs_compute_instance.instance_1", "server_group_id",
						"vkcs_compute_servergroup.sg_2", "id"),
				),
			},
		},
	})
}

func TestAccComputeInstance_serverGroupMembershipIdempotent(t *testing.T) {
	var instance servers.Server
	var sg1 servergroups.ServerGroup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceServerGroupCreateShutoff),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckComputeServerGroupExists("vkcs_compute_servergroup.sg_1", &sg1),
					testAccCheckComputeInstanceManuallyAddedToServerGroup(
						"vkcs_compute_instance.instance_1", "vkcs_compute_servergroup.sg_1"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceServerGroupAdopt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
					testAccCheckComputeServerGroupExists("vkcs_compute_servergroup.sg_1", &sg1),
					testAccCheckComputeInstanceInServerGroup(&instance, &sg1),
				),
			},
		},
	})
}

func TestAccComputeInstance_serverGroupMembershipImport(t *testing.T) {
	resourceName := "vkcs_compute_instance.instance_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceServerGroupCreate),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"stop_before_destroy",
					"force_delete",
					"server_group_id",
				},
			},
		},
	})
}

func testAccCheckComputeInstanceNotInServerGroup(instance *servers.Server, sg *servergroups.ServerGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, m := range sg.Members {
			if m == instance.ID {
				return fmt.Errorf("Instance %s is still part of Server Group %s", instance.ID, sg.ID)
			}
		}

		return nil
	}
}

// testAccCheckComputeInstanceManuallyAddedToServerGroup adds the instance to
// the server group outside of Terraform to verify that the provider reconciles
// membership idempotently.
func testAccCheckComputeInstanceManuallyAddedToServerGroup(instanceResource, groupResource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		instanceRs, ok := s.RootModule().Resources[instanceResource]
		if !ok {
			return fmt.Errorf("Not found: %s", instanceResource)
		}

		groupRs, ok := s.RootModule().Resources[groupResource]
		if !ok {
			return fmt.Errorf("Not found: %s", groupResource)
		}

		config := acctest.AccTestProvider.Meta().(clients.Config)
		computeClient, err := config.ComputeV2Client(acctest.OsRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS compute client: %s", err)
		}

		if err := iservergroups.AddMember(computeClient, groupRs.Primary.ID, instanceRs.Primary.ID).ExtractErr(); err != nil {
			return fmt.Errorf("Error adding instance %s to server group %s: %s", instanceRs.Primary.ID, groupRs.Primary.ID, err)
		}

		return nil
	}
}

const testAccComputeInstanceServerGroupCreate = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_servergroup" "sg_1" {
  name     = "sg_1"
  policies = ["soft-affinity"]
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on        = ["vkcs_networking_router_interface.base"]
  name              = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"

  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  server_group_id    = vkcs_compute_servergroup.sg_1.id

  network {
    uuid = vkcs_networking_network.base.id
  }

  image_id  = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceServerGroupCreateShutoff = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_servergroup" "sg_1" {
  name     = "sg_1"
  policies = ["soft-affinity"]
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on        = ["vkcs_networking_router_interface.base"]
  name              = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  power_state       = "shutoff"

  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]

  network {
    uuid = vkcs_networking_network.base.id
  }

  image_id  = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceServerGroupAdopt = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_servergroup" "sg_1" {
  name     = "sg_1"
  policies = ["soft-affinity"]
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on        = ["vkcs_networking_router_interface.base"]
  name              = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  power_state       = "shutoff"

  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  server_group_id    = vkcs_compute_servergroup.sg_1.id

  network {
    uuid = vkcs_networking_network.base.id
  }

  image_id  = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceServerGroupMove = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_servergroup" "sg_1" {
  name     = "sg_1"
  policies = ["soft-affinity"]
}

resource "vkcs_compute_servergroup" "sg_2" {
  name     = "sg_2"
  policies = ["soft-affinity"]
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on        = ["vkcs_networking_router_interface.base"]
  name              = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  power_state       = "shutoff"

  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  server_group_id    = vkcs_compute_servergroup.sg_2.id

  network {
    uuid = vkcs_networking_network.base.id
  }

  image_id  = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceServerGroupMoveActive = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_servergroup" "sg_1" {
  name     = "sg_1"
  policies = ["soft-affinity"]
}

resource "vkcs_compute_servergroup" "sg_2" {
  name     = "sg_2"
  policies = ["soft-affinity"]
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on        = ["vkcs_networking_router_interface.base"]
  name              = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"

  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  server_group_id    = vkcs_compute_servergroup.sg_2.id

  network {
    uuid = vkcs_networking_network.base.id
  }

  image_id  = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceServerGroupMoveActiveWithOption = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_servergroup" "sg_1" {
  name     = "sg_1"
  policies = ["soft-affinity"]
}

resource "vkcs_compute_servergroup" "sg_2" {
  name     = "sg_2"
  policies = ["soft-affinity"]
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on        = ["vkcs_networking_router_interface.base"]
  name              = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"

  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  server_group_id    = vkcs_compute_servergroup.sg_2.id

  vendor_options {
    stop_before_server_group_change = true
  }

  network {
    uuid = vkcs_networking_network.base.id
  }

  image_id  = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceServerGroupRemoved = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_compute_servergroup" "sg_1" {
  name     = "sg_1"
  policies = ["soft-affinity"]
}

resource "vkcs_compute_servergroup" "sg_2" {
  name     = "sg_2"
  policies = ["soft-affinity"]
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on        = ["vkcs_networking_router_interface.base"]
  name              = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  power_state       = "shutoff"

  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]

  network {
    uuid = vkcs_networking_network.base.id
  }

  image_id  = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`
