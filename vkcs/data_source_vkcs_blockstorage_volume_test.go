package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBlockStorageVolumeDataSource_basic(t *testing.T) {
	resourceName := "data.vkcs_blockstorage_volume.volume_1"
	volumeName := "volume_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageVolumeDataSourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", volumeName),
					resource.TestCheckResourceAttr(resourceName, "size", "1"),
				),
			},
		},
	})
}

func testAccBlockStorageVolumeDataSourceBasic() string {
	return fmt.Sprintf(`
%s

    data "vkcs_blockstorage_volume" "volume_1" {
      name = "${vkcs_blockstorage_volume.volume_1.name}"
    }
  `, testAccRenderConfig(testAccBlockStorageVolumeBasic, testAccValues))
}
