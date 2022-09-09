package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBlockStorageSnapshotDataSource_basic(t *testing.T) {
	resourceName := "data.vkcs_blockstorage_snapshot.snapshot_1"
	snapshotName := "snapshot_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageSnapshotDataSourceBasic(),
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

func testAccBlockStorageSnapshotDataSourceBasic() string {
	return fmt.Sprintf(`
%s

    data "vkcs_blockstorage_snapshot" "snapshot_1" {
      name = vkcs_blockstorage_snapshot.snapshot_1.name
    }
  `, testAccBlockStorageSnapshotBasic())
}
