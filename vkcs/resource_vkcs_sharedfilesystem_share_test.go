package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
)

func TestAccSFSShare_basic(t *testing.T) {
	var share shares.Share

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckSFS(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSFSShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSShareConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSShareExists("vkcs_sharedfilesystem_share.share_1", &share),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "name", "nfs_share"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "description", "test share description"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "share_proto", "NFS"),
				),
			},
			{
				Config: testAccSFSShareConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSShareExists("vkcs_sharedfilesystem_share.share_1", &share),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "name", "nfs_share_updated"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "description", ""),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "share_proto", "NFS"),
				),
			},
			{
				Config: testAccSFSShareConfigExtend(),
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
	config := testAccProvider.Meta().(*config)
	sfsClient, err := config.SharedfilesystemV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
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

		config := testAccProvider.Meta().(*config)
		sfsClient, err := config.SharedfilesystemV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
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

func testAccSFSShareConfigBasic() string {
	return fmt.Sprintf(`
%s

resource "vkcs_sharedfilesystem_share" "share_1" {
  name             = "nfs_share"
  description      = "test share description"
  share_proto      = "NFS"
  share_type       = "default_share_type"
  size             = 1
  share_network_id = "${vkcs_sharedfilesystem_sharenetwork.sharenetwork_1.id}"
}
`, testAccSFSShareNetworkConfigBasic())
}

func testAccSFSShareConfigUpdate() string {
	return fmt.Sprintf(`
%s

resource "vkcs_sharedfilesystem_share" "share_1" {
  name             = "nfs_share_updated"
  share_proto      = "NFS"
  share_type       = "default_share_type"
  size             = 1
  share_network_id = "${vkcs_sharedfilesystem_sharenetwork.sharenetwork_1.id}"
}
`, testAccSFSShareNetworkConfigBasic())
}

func testAccSFSShareConfigExtend() string {
	return fmt.Sprintf(`
%s

resource "vkcs_sharedfilesystem_share" "share_1" {
  name             = "nfs_share_extended"
  share_proto      = "NFS"
  share_type       = "default_share_type"
  size             = 2
  share_network_id = "${vkcs_sharedfilesystem_sharenetwork.sharenetwork_1.id}"
}
`, testAccSFSShareNetworkConfigBasic())
}
