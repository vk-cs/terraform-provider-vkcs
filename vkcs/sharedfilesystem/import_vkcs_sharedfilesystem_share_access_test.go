package sharedfilesystem_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccSFSShareAccess_importBasic(t *testing.T) {
	shareName := "vkcs_sharedfilesystem_share.share_1"
	shareAccessName := "vkcs_sharedfilesystem_share_access.share_access_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckSFSShareAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccSFSShareAccessConfigBasic, map[string]string{"TestAccSFSShareNetworkConfigBasic": acctest.AccTestRenderConfig(testAccSFSShareNetworkConfigBasic, map[string]string{"TestAccSFSShareNetworkConfig": testAccSFSShareNetworkConfig}), "TestAccSFSShareAccessConfig": testAccSFSShareAccessConfig}),
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
