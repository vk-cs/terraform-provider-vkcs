package compute_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccComputeServerGroupDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeServerGroupDataSourceBasic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.vkcs_compute_servergroup.sg_1", "id",
						"vkcs_compute_servergroup.sg_1", "id"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_servergroup.sg_1", "name", "sg_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_servergroup.sg_1", "policies.#", "1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_servergroup.sg_1", "policies.0", "soft-affinity"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_servergroup.sg_1", "members.#", "1"),
				),
			},
		},
	})
}

func TestAccComputeServerGroupDataSource_notFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccComputeServerGroupDataSourceNotFound,
				ExpectError: regexp.MustCompile(`Could not find any server group with this name`),
			},
		},
	})
}

func TestAccComputeServerGroupDataSource_duplicateName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccComputeServerGroupDataSourceDuplicateName,
				ExpectError: regexp.MustCompile(`More than one server group found with this name`),
			},
		},
	})
}

const testAccComputeServerGroupDataSourceBasic = `
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

data "vkcs_compute_servergroup" "sg_1" {
  name       = vkcs_compute_servergroup.sg_1.name
  depends_on = [vkcs_compute_instance.instance_1]
}
`

const testAccComputeServerGroupDataSourceNotFound = `
data "vkcs_compute_servergroup" "sg_1" {
  name = "does-not-exist-servergroup"
}
`

const testAccComputeServerGroupDataSourceDuplicateName = `
resource "vkcs_compute_servergroup" "sg_1" {
  name     = "duplicate-sg"
  policies = ["soft-affinity"]
}

resource "vkcs_compute_servergroup" "sg_2" {
  name     = "duplicate-sg"
  policies = ["soft-affinity"]
}

data "vkcs_compute_servergroup" "sg" {
  name = "duplicate-sg"
}
`
