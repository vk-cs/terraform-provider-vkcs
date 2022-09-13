package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSFSShareDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSFSShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSShareDataSourceBasic(),
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

func testAccSFSShareDataSourceBasic() string {
	return fmt.Sprintf(`
%s

resource "vkcs_sharedfilesystem_share" "share_1" {
  name        = "nfs_share"
  description = "test share description"
  share_proto = "NFS"
  share_type  = "default_share_type"
  size        = 1
  share_network_id = vkcs_sharedfilesystem_sharenetwork.sharenetwork_1.id
}

data "vkcs_sharedfilesystem_share" "share_1" {
  name        = vkcs_sharedfilesystem_share.share_1.name
  description = vkcs_sharedfilesystem_share.share_1.description
  share_network_id = vkcs_sharedfilesystem_sharenetwork.sharenetwork_1.id
}
`, testAccSFSShareNetworkConfigBasic())
}
