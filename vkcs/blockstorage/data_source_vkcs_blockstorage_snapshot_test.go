package blockstorage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccBlockStorageSnapshotDataSource_basic(t *testing.T) {
	resourceName := "data.vkcs_blockstorage_snapshot.snapshot_1"
	snapshotName := "snapshot_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccBlockStorageSnapshotDataSourceBasic, map[string]string{"TestAccBlockStorageSnapshotBasic": acctest.AccTestRenderConfig(testAccBlockStorageSnapshotBasic, map[string]string{"TestAccBlockStorageVolumeBasic": acctest.AccTestRenderConfig(testAccBlockStorageVolumeBasic)})}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageSnapshotDataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", snapshotName),
				),
			},
		},
	})
}

func testAccCheckBlockStorageSnapshotDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find snapshot data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Snapshot data source ID not set")
		}

		return nil
	}
}

const testAccBlockStorageSnapshotDataSourceBasic = `
{{.TestAccBlockStorageSnapshotBasic}}

    data "vkcs_blockstorage_snapshot" "snapshot_1" {
      name = vkcs_blockstorage_snapshot.snapshot_1.name
    }
`
