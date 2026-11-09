package compute

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/restores"
)

const (
	restoreServerResourceKeyPrefix = "OS::Nova::Server#"
	restorePollInterval            = 5 * time.Second
)

type restoreGetter func(context.Context, string) (*restores.Restore, error)

func restoredComputeInstanceWaitState(state string, instanceFound bool) (string, error) {
	if state == "ACTIVE" {
		return state, nil
	}
	if state == "DELETED" && instanceFound {
		return "", fmt.Errorf("restored VM was deleted before becoming active")
	}
	return "BUILD", nil
}

func restoredComputeInstanceStateRefreshFunc(refresh retry.StateRefreshFunc, instanceID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		result, state, err := refresh()
		if err != nil {
			return result, state, err
		}
		instanceFound := true
		if server, ok := result.(*servers.Server); ok && server == nil {
			instanceFound = false
		}
		waitState, err := restoredComputeInstanceWaitState(state, instanceFound)
		if err != nil {
			return result, "", fmt.Errorf("restored VKCS instance %q: %w", instanceID, err)
		}
		return result, waitState, nil
	}
}

func getRestoredComputeInstanceID(ctx context.Context, backupClient *gophercloud.ServiceClient, restoreID string, timeout time.Duration) (string, error) {
	return waitForRestoredComputeInstanceID(ctx, restoreID, timeout, restorePollInterval, func(ctx context.Context, id string) (*restores.Restore, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return restores.Get(backupClient, id).Extract()
	})
}

func restoredComputeInstanceID(restore *restores.Restore) (string, bool, error) {
	if restore == nil {
		return "", false, fmt.Errorf("backup API returned an empty restore")
	}
	if len(restore.ResourcesStatus) == 0 {
		return "", false, nil
	}

	statusKeys := make([]string, 0, len(restore.ResourcesStatus))
	serverIDs := make([]string, 0, 1)
	for key := range restore.ResourcesStatus {
		statusKeys = append(statusKeys, key)
		if strings.HasPrefix(key, restoreServerResourceKeyPrefix) {
			serverID := strings.TrimPrefix(key, restoreServerResourceKeyPrefix)
			if serverID == "" {
				return "", false, fmt.Errorf("restore %q returned an empty Nova server ID", restore.ID)
			}
			serverIDs = append(serverIDs, serverID)
		}
	}
	sort.Strings(statusKeys)
	if len(serverIDs) > 1 {
		sort.Strings(serverIDs)
		return "", false, fmt.Errorf("restore %q contains multiple Nova servers (%s); vkcs_compute_instance.restore_id requires exactly one", restore.ID, strings.Join(serverIDs, ", "))
	}

	for _, key := range statusKeys {
		if restore.ResourcesStatus[key] == "error" {
			reason := restore.ResourcesReason[key]
			if reason == "" {
				reason = "no reason returned by the backup service"
			}
			return "", false, fmt.Errorf("restore %q failed for resource %q: %s", restore.ID, key, reason)
		}
	}
	for _, key := range statusKeys {
		if restore.ResourcesStatus[key] != "available" {
			return "", false, nil
		}
	}

	if len(serverIDs) == 0 {
		return "", false, fmt.Errorf("restore %q completed without a restored Nova server", restore.ID)
	}

	return serverIDs[0], true, nil
}

func waitForRestoredComputeInstanceID(ctx context.Context, restoreID string, timeout, pollInterval time.Duration, getRestore restoreGetter) (string, error) {
	if restoreID == "" {
		return "", fmt.Errorf("restore ID must not be empty")
	}
	if getRestore == nil {
		return "", fmt.Errorf("restore getter must not be nil")
	}
	if timeout <= 0 {
		return "", fmt.Errorf("restore wait timeout must be positive")
	}
	if pollInterval <= 0 {
		pollInterval = restorePollInterval
	}

	waitCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		restore, err := getRestore(waitCtx, restoreID)
		if err != nil {
			return "", fmt.Errorf("read backup restore %q: %w", restoreID, err)
		}

		serverID, ready, err := restoredComputeInstanceID(restore)
		if err != nil {
			return "", err
		}
		if ready {
			return serverID, nil
		}

		select {
		case <-waitCtx.Done():
			return "", fmt.Errorf("timed out waiting for backup restore %q: %w", restoreID, waitCtx.Err())
		case <-ticker.C:
		}
	}
}
