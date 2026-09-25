package baremetal_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

// TestManualAccBareMetalServerResource_basic is a manual acceptance test.
// It cannot be executed in an automated environment because provisioning
// a vkcs_baremetal_server requires allocation of a real physical server,
// which is costly, and not suitable for routine test runs. This is too inefficient.
func TestManualAccBareMetalServerResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccServerResourceBasic, map[string]string{
					"TestAccBaremetalOSDataSourceBasic":     acctest.AccTestRenderConfig(testAccOSDataSourceBasic),
					"TestAccBaremetalFlavorDataSourceBasic": acctest.AccTestRenderConfig(testAccFlavorDataSourceBasic),
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vkcs_baremetal_server.basic", "id"),
					resource.TestCheckResourceAttr("vkcs_baremetal_server.basic", "name", "TerraformAccTestBasic"),
					resource.TestCheckResourceAttr("vkcs_baremetal_server.basic", "availability_zone", "ME1"),
				),
			},
		},
	})
}

const testAccServerResourceBasic = `
{{.TestAccBaremetalOSDataSourceBasic}}
{{.TestAccBaremetalFlavorDataSourceBasic}}

resource "vkcs_compute_keypair" "baremetal" {
  name = "tf-baremetal-test-key-pair"
}

resource "vkcs_networking_network" "network" {
  name        = "tf-baremetal-test-network"
  description = "my network description"
}

resource "vkcs_networking_subnet" "subnet" {
  name       = "tf-baremetal-test-subnet"
  network_id = vkcs_networking_network.network.id
}

resource "vkcs_baremetal_server" "basic" {
  name              = "TerraformAccTestBasic"
  os_id             = data.vkcs_baremetal_os.basic.id
  flavor_id         = data.vkcs_baremetal_flavor.basic.id
  key_pair          = vkcs_compute_keypair.baremetal.name
  availability_zone = "ME1"

  nic {
    name = "nic0"
    vlan {
      native     = true
      network_id = vkcs_networking_network.network.id
      subnet_id  = vkcs_networking_subnet.subnet.id
    }
  }
}
`

// TestManualAccBareMetalServerResource_storageLayout is a manual acceptance
// test covering the custom storage layout blocks.
func TestManualAccBareMetalServerResource_storageLayout(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccServerResourceStorageLayout, map[string]string{
					"TestAccBaremetalOSDataSourceBasic":     acctest.AccTestRenderConfig(testAccOSDataSourceBasic),
					"TestAccBaremetalFlavorDataSourceBasic": acctest.AccTestRenderConfig(testAccFlavorDataSourceBasic),
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vkcs_baremetal_server.layout", "id"),
					resource.TestCheckResourceAttr("vkcs_baremetal_server.layout", "storage_layout.disk[0].id", "system-0"),
					resource.TestCheckResourceAttr("vkcs_baremetal_server.layout", "storage_layout.disk[0].type", "ssd"),
					resource.TestCheckResourceAttr("vkcs_baremetal_server.layout", "storage_layout.disk[0].size", "447"),
					resource.TestCheckResourceAttr("vkcs_baremetal_server.layout", "storage_layout.disk[0].partition[0].mount", "/boot/efi"),
					resource.TestCheckResourceAttr("vkcs_baremetal_server.layout", "storage_layout.disk[0].partition[0].fs", "vfat"),
					resource.TestCheckResourceAttr("vkcs_baremetal_server.layout", "storage_layout.disk[0].partition[1].mount", "/"),
					resource.TestCheckResourceAttr("vkcs_baremetal_server.layout", "storage_layout.disk[0].partition[1].fs", "ext4"),
					resource.TestCheckResourceAttr("vkcs_baremetal_server.layout", "storage_layout.disk[0].partition[2].mount", ""),
					resource.TestCheckResourceAttr("vkcs_baremetal_server.layout", "storage_layout.disk[0].partition[3].fs", "swap"),
				),
			},
		},
	})
}

const testAccServerResourceStorageLayout = `
{{.TestAccBaremetalOSDataSourceBasic}}
{{.TestAccBaremetalFlavorDataSourceBasic}}

resource "vkcs_compute_keypair" "baremetal" {
  name = "tf-baremetal-test-key-pair"
}

resource "vkcs_networking_network" "network" {
  name        = "tf-baremetal-test-network"
  description = "my network description"
}

resource "vkcs_networking_subnet" "subnet" {
  name       = "tf-baremetal-test-subnet"
  network_id = vkcs_networking_network.network.id
}

resource "vkcs_baremetal_server" "layout" {
  name              = "TerraformAccTestStorageLayout"
  os_id             = data.vkcs_baremetal_os.basic.id
  flavor_id         = data.vkcs_baremetal_flavor.basic.id
  key_pair          = vkcs_compute_keypair.baremetal.name
  availability_zone = "ME1"

  storage_layout {
    disk {
      id   = "system-0"
      type = "SSD"
      size = 447

      partition {
        mount = "/boot/efi"
        fs    = "vfat"
        size  = "512MiB"
      }

      partition {
        mount = "/"
        fs    = "ext4"
        size  = "50GiB"
      }

      partition {
        mount = ""
        fs    = "ext4"
        size  = "60GiB"
      }

      partition {
        fs   = "swap"
        size = "512MiB"
      }
    }
  }

  nic {
    name = "nic0"
    vlan {
      native     = true
      network_id = vkcs_networking_network.network.id
      subnet_id  = vkcs_networking_subnet.subnet.id
    }
  }
}
`
