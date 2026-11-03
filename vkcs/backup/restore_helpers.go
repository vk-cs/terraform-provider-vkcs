package backup

import (
	"context"
	"fmt"
	"sort"
	"slices"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/checkpoints"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/restores"
)

const (
	restoreCreateTimeout = time.Hour
	restorePollInterval  = 5 * time.Second
)

func selectRestoreResourceID(resources []checkpoints.CheckpointDetailResource, requestedID string) (string, error) {
	var novaResourceIDs []string
	for _, resource := range resources {
		if resource.Type == NovaInstance {
			novaResourceIDs = append(novaResourceIDs, resource.OriginalID)
		}
	}

	if requestedID != "" {
		for _, resourceID := range novaResourceIDs {
			if resourceID == requestedID {
				return resourceID, nil
			}
		}
		return "", fmt.Errorf("checkpoint does not contain Nova server resource %q", requestedID)
	}

	switch len(novaResourceIDs) {
	case 0:
		return "", fmt.Errorf("checkpoint does not contain a Nova server resource")
	case 1:
		return novaResourceIDs[0], nil
	default:
		return "", fmt.Errorf("checkpoint contains multiple Nova servers; set source_resource_id to choose one")
	}
}

func restoreStatusesComplete(statuses map[string]string) bool {
	if len(statuses) == 0 {
		return false
	}
	for _, status := range statuses {
		if status != "available" && status != "error" {
			return false
		}
	}
	return true
}

func restoreFailureDetails(statuses, reasons map[string]string) []string {
	resourceIDs := make([]string, 0, len(statuses))
	for resourceID, status := range statuses {
		if status == "error" {
			resourceIDs = append(resourceIDs, resourceID)
		}
	}
	sort.Strings(resourceIDs)

	failures := make([]string, 0, len(resourceIDs))
	for _, resourceID := range resourceIDs {
		reason := reasons[resourceID]
		if reason == "" {
			reason = "restore failed"
		}
		failures = append(failures, fmt.Sprintf("%s: %s", resourceID, reason))
	}
	return failures
}

func findMatchingRestore(existing []restores.Restore, desired restores.RestoreCreateOpts) *restores.Restore {
	for i := range existing {
		candidate := &existing[i]
		if candidate.ProviderID != desired.ProviderID || candidate.CheckpointID != desired.CheckpointID {
			continue
		}
		if sameRestoreParameters(candidate.Parameters, desired.Parameters) {
			return candidate
		}
	}
	return nil
}

func getOrCreateRestore(
	existing []restores.Restore,
	desired restores.RestoreCreateOpts,
	create func() (*restores.Restore, error),
) (*restores.Restore, error) {
	if match := findMatchingRestore(existing, desired); match != nil {
		return match, nil
	}
	return create()
}

func sameRestoreParameters(existing, desired map[string]restores.RestoreParameters) bool {
	if len(existing) != len(desired) {
		return false
	}
	for key, desiredParameters := range desired {
		existingParameters, ok := existing[key]
		if !ok || !sameRestoreParameterValues(existingParameters, desiredParameters) {
			return false
		}
	}
	return true
}

func sameRestoreParameterValues(left, right restores.RestoreParameters) bool {
	return left.RestoreName == right.RestoreName &&
		left.Destination == right.Destination &&
		left.RestoreKeypair == right.RestoreKeypair &&
		left.AvailabilityZone == right.AvailabilityZone &&
		left.RestoreNetID == right.RestoreNetID &&
		left.VolumeType == right.VolumeType &&
		int64PointerValue(left.VolumeSize) == int64PointerValue(right.VolumeSize) &&
		boolPointerValue(left.EnableFIP) == boolPointerValue(right.EnableFIP) &&
		left.Flavor == right.Flavor &&
		left.VolumeTargetVMID == right.VolumeTargetVMID &&
		int64PointerValue(left.VolumeAutoresizeEnabled) == int64PointerValue(right.VolumeAutoresizeEnabled) &&
		int64PointerValue(left.VolumeAutoresizeMaxSize) == int64PointerValue(right.VolumeAutoresizeMaxSize) &&
		sameStrings(left.SecurityGroups, right.SecurityGroups) &&
		boolPointerValue(left.CloudMonitoringEnabled) == boolPointerValue(right.CloudMonitoringEnabled)
}

func int64PointerValue(value *int64) int64 {
	if value == nil {
		return 0
	}
	return *value
}

func boolPointerValue(value *bool) bool {
	return value != nil && *value
}

func sameStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	leftCopy := slices.Clone(left)
	rightCopy := slices.Clone(right)
	sort.Strings(leftCopy)
	sort.Strings(rightCopy)
	return slices.Equal(leftCopy, rightCopy)
}

func waitForRestore(ctx context.Context, client *gophercloud.ServiceClient, restoreID string) (*restores.Restore, error) {
	ticker := time.NewTicker(restorePollInterval)
	defer ticker.Stop()

	var latest *restores.Restore
	for {
		if err := ctx.Err(); err != nil {
			return latest, fmt.Errorf("waiting for restore %s completion: %w", restoreID, err)
		}

		response, err := restores.Get(client, restoreID).Extract()
		if err != nil {
			return latest, fmt.Errorf("read restore status for %s: %w", restoreID, err)
		}
		latest = response
		if restoreStatusesComplete(response.ResourcesStatus) {
			return response, nil
		}

		select {
		case <-ctx.Done():
			return latest, fmt.Errorf("waiting for restore %s completion: %w", restoreID, ctx.Err())
		case <-ticker.C:
		}
	}
}
