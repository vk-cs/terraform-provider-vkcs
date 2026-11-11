package compute

import (
	"context"
	"fmt"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/checkpoints"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/providers"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/restores"
)

const (
	restoreServerResourceKeyPrefix = "OS::Nova::Server#"
	novaBackupProviderType         = "OS::Nova"
	novaBackupResourceType         = "OS::Nova::Server"
	restorePollInterval            = 5 * time.Second
	restoreServerVisibilityGrace   = 30 * time.Second
)

type restoreGetter func(context.Context, string) (*restores.Restore, error)
type restoredServerExistsGetter func(context.Context, string) (bool, error)

func findNovaBackupProviderID(client *gophercloud.ServiceClient) (string, error) {
	allProviders, err := providers.List(client).Extract()
	if err != nil {
		return "", fmt.Errorf("list backup providers: %w", err)
	}
	for _, provider := range allProviders {
		if provider.Name == novaBackupProviderType {
			return provider.ID, nil
		}
	}
	return "", fmt.Errorf("backup provider %q was not found", novaBackupProviderType)
}

func selectCheckpointRestoreSourceID(resources []checkpoints.CheckpointDetailResource, requestedID string) (string, error) {
	var novaResourceIDs []string
	for _, resource := range resources {
		if resource.Type == novaBackupResourceType {
			novaResourceIDs = append(novaResourceIDs, resource.OriginalID)
		}
	}

	if requestedID != "" {
		for _, resourceID := range novaResourceIDs {
			if resourceID == requestedID {
				return resourceID, nil
			}
		}
		return "", fmt.Errorf("checkpoint does not contain Nova server %q", requestedID)
	}

	switch len(novaResourceIDs) {
	case 0:
		return "", fmt.Errorf("checkpoint does not contain a Nova server")
	case 1:
		return novaResourceIDs[0], nil
	default:
		return "", fmt.Errorf("checkpoint contains multiple Nova servers; set source_resource_id to choose one")
	}
}

func computeRestoreNetworkID(d *schema.ResourceData, meta interface{}) (string, error) {
	if networkMode := d.Get("network_mode").(string); networkMode != "" {
		return "", fmt.Errorf("network_mode is not supported for checkpoint restore; set one network block")
	}

	networks, err := getAllInstanceNetworks(d, meta)
	if err != nil {
		return "", err
	}
	return singleCheckpointRestoreNetworkID(networks)
}

func singleCheckpointRestoreNetworkID(networks []InstanceNetwork) (string, error) {
	if len(networks) != 1 {
		return "", fmt.Errorf("checkpoint restore requires exactly one network, got %d", len(networks))
	}
	if networks[0].UUID == "" {
		return "", fmt.Errorf("checkpoint restore network must resolve to a network UUID")
	}
	return networks[0].UUID, nil
}

func computeRestoreParameters(d *schema.ResourceData, networkID, flavorID string) (restores.RestoreParameters, error) {
	name := d.Get("name").(string)
	if name == "" {
		return restores.RestoreParameters{}, fmt.Errorf("name must be set when restoring a checkpoint")
	}
	if networkID == "" {
		return restores.RestoreParameters{}, fmt.Errorf("network must be set when restoring a checkpoint")
	}
	keyPair := d.Get("key_pair").(string)
	if keyPair == "" {
		return restores.RestoreParameters{}, fmt.Errorf("key_pair must be set when restoring a checkpoint")
	}

	parameters := restores.RestoreParameters{
		RestoreName:    name,
		Destination:    "new",
		RestoreKeypair: keyPair,
		RestoreNetID:   networkID,
		Flavor:         flavorID,
		SecurityGroups: resourceInstanceRestoreSecurityGroupIDs(d),
	}
	if availabilityZone := d.Get("availability_zone").(string); availabilityZone != "" {
		parameters.AvailabilityZone = availabilityZone
	}
	return parameters, nil
}

func resourceInstanceRestoreSecurityGroupIDs(d *schema.ResourceData) []string {
	securityGroups := d.Get("security_group_ids").(*schema.Set).List()
	ids := make([]string, 0, len(securityGroups))
	for _, securityGroup := range securityGroups {
		ids = append(ids, securityGroup.(string))
	}
	return ids
}

func findMatchingRestore(existing []restores.Restore, desired restores.RestoreCreateOpts) []*restores.Restore {
	matches := make([]*restores.Restore, 0)
	for i := range existing {
		candidate := &existing[i]
		if candidate.ProviderID != desired.ProviderID || candidate.CheckpointID != desired.CheckpointID {
			continue
		}
		if sameRestoreParameters(candidate.Parameters, desired.Parameters) {
			matches = append(matches, candidate)
		}
	}
	return matches
}

func getOrCreateRestore(
	existing []restores.Restore,
	desired restores.RestoreCreateOpts,
	canReuse func(*restores.Restore) (bool, error),
	create func() (*restores.Restore, error),
) (*restores.Restore, error) {
	for _, candidate := range findMatchingRestore(existing, desired) {
		reuse, err := canReuse(candidate)
		if err != nil {
			return nil, err
		}
		if reuse {
			return candidate, nil
		}
	}
	return create()
}

func canReuseExistingRestore(restore *restores.Restore, serverExists func(string) (bool, error)) (bool, error) {
	if restore == nil || !restoreStatusesComplete(restore.ResourcesStatus) {
		return true, nil
	}
	for _, status := range restore.ResourcesStatus {
		if status == "error" {
			return true, nil
		}
	}

	serverID, ready, err := restoredComputeInstanceID(restore)
	if err != nil {
		return false, err
	}
	if !ready {
		return true, nil
	}
	if serverExists == nil {
		return false, fmt.Errorf("server existence checker must not be nil")
	}
	return serverExists(serverID)
}

func waitForRestoredServerVisibility(ctx context.Context, serverID string, timeout, pollInterval time.Duration, getServer restoredServerExistsGetter) (bool, error) {
	if serverID == "" {
		return false, fmt.Errorf("restored server ID must not be empty")
	}
	if getServer == nil {
		return false, fmt.Errorf("server existence checker must not be nil")
	}
	if timeout <= 0 {
		return false, nil
	}
	if pollInterval <= 0 {
		pollInterval = restorePollInterval
	}

	waitCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		if err := waitCtx.Err(); err != nil {
			if ctx.Err() != nil {
				return false, ctx.Err()
			}
			return false, nil
		}
		exists, err := getServer(waitCtx, serverID)
		if err != nil {
			return false, err
		}
		if exists {
			return true, nil
		}

		select {
		case <-waitCtx.Done():
			if ctx.Err() != nil {
				return false, ctx.Err()
			}
			return false, nil
		case <-ticker.C:
		}
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
		sameRestoreStrings(left.SecurityGroups, right.SecurityGroups) &&
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

func sameRestoreStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	leftCopy := slices.Clone(left)
	rightCopy := slices.Clone(right)
	sort.Strings(leftCopy)
	sort.Strings(rightCopy)
	return slices.Equal(leftCopy, rightCopy)
}

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
		return "", false, fmt.Errorf("restore %q contains multiple Nova servers (%s); vkcs_compute_instance.checkpoint_id requires exactly one", restore.ID, strings.Join(serverIDs, ", "))
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
