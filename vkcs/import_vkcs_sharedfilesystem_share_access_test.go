package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSFSShareAccess_importBasic(t *testing.T) {
	shareName := "vkcs_sharedfilesystem_share.share_1"
	shareAccessName := "vkcs_sharedfilesystem_share_access.share_access_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckSFS(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSFSShareAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSShareAccessConfigBasic(),
			},

			{
				ResourceName:      shareAccessName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSFSShareAccessImportID(shareName, shareAccessName),
			},
		},
	})
}

func testAccSFSShareAccessImportID(shareResource, shareAccessResource string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		share, ok := s.RootModule().Resources[shareResource]
		if !ok {
			return "", fmt.Errorf("Share not found: %s", shareResource)
		}

		shareAccess, ok := s.RootModule().Resources[shareAccessResource]
		if !ok {
			return "", fmt.Errorf("Share access not found: %s", shareAccessResource)
		}

		return fmt.Sprintf("%s/%s", share.Primary.ID, shareAccess.Primary.ID), nil
	}
}
