package blockstorage_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccBlockStorageVolumeDataSource_basic(t *testing.T) {
	resourceName := "data.vkcs_blockstorage_volume.volume_1"
	volumeName := "volume_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccBlockStorageVolumeDataSourceBasic, map[string]string{"TestAccBlockStorageVolumeBasic": acctest.AccTestRenderConfig(testAccBlockStorageVolumeBasic)}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", volumeName),
					resource.TestCheckResourceAttr(resourceName, "size", "1"),
				),
			},
		},
	})
}

const testAccBlockStorageVolumeDataSourceBasic = `
{{.TestAccBlockStorageVolumeBasic}}

    data "vkcs_blockstorage_volume" "volume_1" {
      name = vkcs_blockstorage_volume.volume_1.name
    }
`
