package blockstorage

import (
	"fmt"
	"sort"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/snapshots"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func blockStorageSnapshotStateRefreshFunc(client *gophercloud.ServiceClient, volumeSnapshotID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := snapshots.Get(client, volumeSnapshotID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return v, bsSnapshotStatusDeleted, nil
			}
			return nil, "", err
		}
		if v.Status == "error" {
			return v, v.Status, fmt.Errorf("there was an error creating the block storage volume snapshot")
		}

		return v, v.Status, nil
	}
}

// blockStorageSnapshotSort represents a sortable slice of block storage
// v3 snapshots.
type blockStorageSnapshotSort []snapshots.Snapshot

func (snaphot blockStorageSnapshotSort) Len() int {
	return len(snaphot)
}

func (snaphot blockStorageSnapshotSort) Swap(i, j int) {
	snaphot[i], snaphot[j] = snaphot[j], snaphot[i]
}

func (snaphot blockStorageSnapshotSort) Less(i, j int) bool {
	itime := snaphot[i].CreatedAt
	jtime := snaphot[j].CreatedAt
	return itime.Unix() < jtime.Unix()
}

func dataSourceBlockStorageMostRecentSnapshot(snapshots []snapshots.Snapshot) snapshots.Snapshot {
	sortedSnapshots := snapshots
	sort.Sort(blockStorageSnapshotSort(sortedSnapshots))
	return sortedSnapshots[len(sortedSnapshots)-1]
}
