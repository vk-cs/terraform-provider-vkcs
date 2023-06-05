package sharedfilesystem_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"

	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
)

func TestAccSFSShare_basic(t *testing.T) {
	var share shares.Share

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckSFSShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccSFSShareConfigBasic, map[string]string{"TestAccSFSShareNetworkConfigBasic": acctest.AccTestRenderConfig(testAccSFSShareNetworkConfigBasic, map[string]string{"TestAccSFSShareNetworkConfig": testAccSFSShareNetworkConfig})}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSShareExists("vkcs_sharedfilesystem_share.share_1", &share),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "name", "nfs_share"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "description", "test share description"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "share_proto", "NFS"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccSFSShareConfigUpdate, map[string]string{"TestAccSFSShareNetworkConfigBasic": acctest.AccTestRenderConfig(testAccSFSShareNetworkConfigBasic, map[string]string{"TestAccSFSShareNetworkConfig": testAccSFSShareNetworkConfig})}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSShareExists("vkcs_sharedfilesystem_share.share_1", &share),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "name", "nfs_share_updated"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "description", ""),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "share_proto", "NFS"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccSFSShareConfigExtend, map[string]string{"TestAccSFSShareNetworkConfigBasic": acctest.AccTestRenderConfig(testAccSFSShareNetworkConfigBasic, map[string]string{"TestAccSFSShareNetworkConfig": testAccSFSShareNetworkConfig})}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSShareExists("vkcs_sharedfilesystem_share.share_1", &share),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "name", "nfs_share_extended"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "share_proto", "NFS"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "size", "2"),
				),
			},
		},
	})
}

func testAccCheckSFSShareDestroy(s *terraform.State) error {
	config := acctest.AccTestProvider.Meta().(clients.Config)
	sfsClient, err := config.SharedfilesystemV2Client(acctest.OsRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS sharedfilesystem client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_sharedfilesystem_securityservice" {
			continue
		}

		_, err := shares.Get(sfsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Manila share still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckSFSShareExists(n string, share *shares.Share) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.AccTestProvider.Meta().(clients.Config)
		sfsClient, err := config.SharedfilesystemV2Client(acctest.OsRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS sharedfilesystem client: %s", err)
		}

		found, err := shares.Get(sfsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Share not found")
		}

		*share = *found

		return nil
	}
}

const testAccSFSShareConfigBasic = `
{{.TestAccSFSShareNetworkConfigBasic}}

resource "vkcs_sharedfilesystem_share" "share_1" {
  name             = "nfs_share"
  description      = "test share description"
  share_proto      = "NFS"
  share_type       = "default_share_type"
  size             = 1
  share_network_id = vkcs_sharedfilesystem_sharenetwork.sharenetwork_1.id
}
`

const testAccSFSShareConfigUpdate = `
{{.TestAccSFSShareNetworkConfigBasic}}

resource "vkcs_sharedfilesystem_share" "share_1" {
  name             = "nfs_share_updated"
  share_proto      = "NFS"
  share_type       = "default_share_type"
  size             = 1
  share_network_id = vkcs_sharedfilesystem_sharenetwork.sharenetwork_1.id
}
`

const testAccSFSShareConfigExtend = `
{{.TestAccSFSShareNetworkConfigBasic}}

resource "vkcs_sharedfilesystem_share" "share_1" {
  name             = "nfs_share_extended"
  share_proto      = "NFS"
  share_type       = "default_share_type"
  size             = 2
  share_network_id = vkcs_sharedfilesystem_sharenetwork.sharenetwork_1.id
}
`
