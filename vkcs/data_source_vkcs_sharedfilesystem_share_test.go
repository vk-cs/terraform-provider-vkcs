package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSFSShareDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckSFS(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSFSShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSShareDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSShareDataSourceID("data.vkcs_sharedfilesystem_share.share_1"),
					resource.TestCheckResourceAttr("data.vkcs_sharedfilesystem_share.share_1", "name", "nfs_share"),
					resource.TestCheckResourceAttr("data.vkcs_sharedfilesystem_share.share_1", "description", "test share description"),
					resource.TestCheckResourceAttr("data.vkcs_sharedfilesystem_share.share_1", "share_proto", "NFS"),
					resource.TestCheckResourceAttr("data.vkcs_sharedfilesystem_share.share_1", "size", "1"),
				),
			},
		},
	})
}

func testAccCheckSFSShareDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find share data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Share data source ID not set")
		}

		return nil
	}
}

const testAccSFSShareDataSourceBasic = `
resource "vkcs_sharedfilesystem_share" "share_1" {
  name        = "nfs_share"
  description = "test share description"
  share_proto = "NFS"
  share_type  = "dhss_false"
  size        = 1
}

data "vkcs_sharedfilesystem_share" "share_1" {
  name        = "${vkcs_sharedfilesystem_share.share_1.name}"
  description = "${vkcs_sharedfilesystem_share.share_1.description}"
}
`
