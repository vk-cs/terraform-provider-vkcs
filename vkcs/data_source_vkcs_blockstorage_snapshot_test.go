package vkcs

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/snapshots"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
)

func TestAccBlockStorageSnapshotDataSource_basic(t *testing.T) {
	resourceName := "data.vkcs_blockstorage_snapshot.snapshot_1"
	volumeName := acctest.RandomWithPrefix("tf-acc-volume")
	snapshotName := acctest.RandomWithPrefix("tf-acc-snapshot")

	var volumeID, snapshotID string
	if os.Getenv("TF_ACC") != "" {
		var err error
		volumeID, snapshotID, err = testAccBlockStorageCreateVolumeAndSnapshot(volumeName, snapshotName)
		if err != nil {
			t.Fatal(err)
		}
		defer testAccBlockStorageDeleteVolumeAndSnapshot(t, volumeID, snapshotID)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageSnapshotDataSourceBasic(snapshotName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageSnapshotDataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", snapshotName),
					resource.TestCheckResourceAttr(resourceName, "volume_id", volumeID),
				),
			},
		},
	})
}

func testAccBlockStorageCreateVolumeAndSnapshot(volumeName, snapshotName string) (string, string, error) {
	config, err := testAccAuthFromEnv()
	if err != nil {
		return "", "", err
	}

	bsClient, err := config.BlockStorageV3Client(osRegionName)
	if err != nil {
		return "", "", err
	}

	volCreateOpts := volumes.CreateOpts{
		Size: 1,
		Name: volumeName,
	}

	volume, err := volumes.Create(bsClient, volCreateOpts).Extract()
	if err != nil {
		return "", "", err
	}

	err = volumes.WaitForStatus(bsClient, volume.ID, "available", 60)
	if err != nil {
		return "", "", err
	}

	snapCreateOpts := snapshots.CreateOpts{
		VolumeID: volume.ID,
		Name:     snapshotName,
	}

	snapshot, err := snapshots.Create(bsClient, snapCreateOpts).Extract()
	if err != nil {
		return volume.ID, "", err
	}

	err = snapshots.WaitForStatus(bsClient, snapshot.ID, "available", 60)
	if err != nil {
		return volume.ID, "", err
	}

	return volume.ID, snapshot.ID, nil
}

func testAccBlockStorageDeleteVolumeAndSnapshot(t *testing.T, volumeID, snapshotID string) {
	config, err := testAccAuthFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	bsClient, err := config.BlockStorageV3Client(osRegionName)
	if err != nil {
		t.Fatal(err)
	}

	err = snapshots.Delete(bsClient, snapshotID).ExtractErr()
	if err != nil {
		t.Fatal(err)
	}

	err = snapshots.WaitForStatus(bsClient, snapshotID, "DELETED", 60)
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); !ok {
			t.Fatal(err)
		}
	}

	err = volumes.Delete(bsClient, volumeID, nil).ExtractErr()
	if err != nil {
		t.Fatal(err)
	}

	err = volumes.WaitForStatus(bsClient, volumeID, "DELETED", 60)
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); !ok {
			t.Fatal(err)
		}
	}
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

func testAccBlockStorageSnapshotDataSourceBasic(snapshotName string) string {
	return fmt.Sprintf(`
    data "vkcs_blockstorage_snapshot" "snapshot_1" {
      name = "%s"
    }
  `, snapshotName)
}
