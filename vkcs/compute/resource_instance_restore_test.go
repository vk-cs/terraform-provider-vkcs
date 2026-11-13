package compute

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/checkpoints"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/restores"
)

func TestRestoredComputeInstanceID(t *testing.T) {
	tests := []struct {
		name          string
		restore       *restores.Restore
		wantID        string
		wantReady     bool
		wantErrSubstr string
	}{
		{
			name: "restore operation is still pending",
			restore: &restores.Restore{
				ID:              "restore-1",
				ResourcesStatus: map[string]string{"OS::Nova::Server#server-1": "restoring"},
			},
		},
		{
			name: "returns the created Nova server ID",
			restore: &restores.Restore{
				ID:              "restore-2",
				ResourcesStatus: map[string]string{"OS::Nova::Server#server-2": "available"},
			},
			wantID:    "server-2",
			wantReady: true,
		},
		{
			name: "reports a failed restore reason",
			restore: &restores.Restore{
				ID:              "restore-3",
				ResourcesStatus: map[string]string{"OS::Nova::Server#server-3": "error"},
				ResourcesReason: map[string]string{"OS::Nova::Server#server-3": "quota exceeded"},
			},
			wantErrSubstr: "quota exceeded",
		},
		{
			name: "reports an error even when another resource is still pending",
			restore: &restores.Restore{
				ID: "restore-8",
				ResourcesStatus: map[string]string{
					"OS::Cinder::Volume#volume-8": "restoring",
					"OS::Nova::Server#server-8":   "error",
				},
				ResourcesReason: map[string]string{"OS::Nova::Server#server-8": "server restore failed"},
			},
			wantErrSubstr: "server restore failed",
		},
		{
			name: "rejects a restore containing multiple VMs",
			restore: &restores.Restore{
				ID: "restore-4",
				ResourcesStatus: map[string]string{
					"OS::Nova::Server#server-4": "available",
					"OS::Nova::Server#server-5": "available",
				},
			},
			wantErrSubstr: "multiple Nova servers",
		},
		{
			name: "rejects a completed restore without a Nova server",
			restore: &restores.Restore{
				ID:              "restore-6",
				ResourcesStatus: map[string]string{"OS::Cinder::Volume#volume-6": "available"},
			},
			wantErrSubstr: "without a restored Nova server",
		},
		{
			name:    "waits when restore statuses are not populated yet",
			restore: &restores.Restore{ID: "restore-7"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotID, gotReady, err := restoredComputeInstanceID(test.restore)
			if test.wantErrSubstr != "" {
				if err == nil || !strings.Contains(err.Error(), test.wantErrSubstr) {
					t.Fatalf("restoredComputeInstanceID() error = %v, want substring %q", err, test.wantErrSubstr)
				}
				return
			}
			if err != nil {
				t.Fatalf("restoredComputeInstanceID() error = %v", err)
			}
			if gotID != test.wantID {
				t.Errorf("restoredComputeInstanceID() ID = %q, want %q", gotID, test.wantID)
			}
			if gotReady != test.wantReady {
				t.Errorf("restoredComputeInstanceID() ready = %t, want %t", gotReady, test.wantReady)
			}
		})
	}
}

func TestComputeInstanceCheckpointRestoreSchema(t *testing.T) {
	resource := ResourceComputeInstance()
	resourceSchema := resource.Schema
	checkpointID := resourceSchema["checkpoint_id"]
	if checkpointID == nil {
		t.Fatal("checkpoint_id schema attribute is missing")
	}
	if !checkpointID.Optional || !checkpointID.ForceNew {
		t.Fatalf("checkpoint_id must be optional and force replacement, got Optional=%t ForceNew=%t", checkpointID.Optional, checkpointID.ForceNew)
	}

	sourceResourceID := resourceSchema["source_resource_id"]
	if sourceResourceID == nil || !sourceResourceID.Optional || !sourceResourceID.Computed || !sourceResourceID.ForceNew {
		t.Fatalf("source_resource_id must be optional, computed, and force replacement: %#v", sourceResourceID)
	}
	if len(sourceResourceID.RequiredWith) != 1 || sourceResourceID.RequiredWith[0] != "checkpoint_id" {
		t.Fatalf("source_resource_id must require checkpoint_id, got RequiredWith=%v", sourceResourceID.RequiredWith)
	}

	conflicts := make(map[string]bool, len(checkpointID.ConflictsWith))
	for _, attribute := range checkpointID.ConflictsWith {
		conflicts[attribute] = true
	}
	for _, attribute := range []string{
		"image_id", "image_name", "block_device", "user_data", "network_mode", "config_drive", "personality",
		"security_group_ids", "flavor_id", "flavor_name", "metadata",
	} {
		if !conflicts[attribute] {
			t.Errorf("checkpoint_id must conflict with %q", attribute)
		}
	}

	restoreID := resourceSchema["restore_id"]
	if restoreID == nil || !restoreID.Computed || restoreID.Optional {
		t.Fatalf("restore_id must be computed and not user-configurable: %#v", restoreID)
	}
}

func TestValidateCheckpointRestoreConfigRejectsUnsupportedNetworkAttributes(t *testing.T) {
	for _, attribute := range []string{"port", "fixed_ip_v4"} {
		for _, value := range []struct {
			name  string
			value cty.Value
		}{
			{name: "known value", value: cty.StringVal("unsupported-value")},
			{name: "unknown expression", value: cty.UnknownVal(cty.String)},
		} {
			t.Run(attribute+"/"+value.name, func(t *testing.T) {
				networkAttributes := map[string]cty.Value{
					"uuid":           cty.StringVal("network-id"),
					"name":           cty.NullVal(cty.String),
					"port":           cty.NullVal(cty.String),
					"fixed_ip_v4":    cty.NullVal(cty.String),
					"mac":            cty.NullVal(cty.String),
					"access_network": cty.BoolVal(false),
				}
				networkAttributes[attribute] = value.value
				network := cty.ObjectVal(networkAttributes)
				config := cty.ObjectVal(map[string]cty.Value{
					"checkpoint_id": cty.StringVal("checkpoint-id"),
					"network":       cty.ListVal([]cty.Value{network}),
				})

				err := validateCheckpointRestoreConfig(config)
				if err == nil || !strings.Contains(err.Error(), "network."+attribute) {
					t.Fatalf("validateCheckpointRestoreConfig() error = %v, want unsupported network.%s error", err, attribute)
				}
			})
		}
	}
}

func TestSelectCheckpointRestoreSourceID(t *testing.T) {
	tests := []struct {
		name          string
		resources     []checkpoints.CheckpointDetailResource
		requestedID   string
		wantID        string
		wantErrSubstr string
	}{
		{
			name: "selects the only Nova VM",
			resources: []checkpoints.CheckpointDetailResource{
				{OriginalID: "server-one", Type: "OS::Nova::Server"},
				{OriginalID: "volume-one", Type: "OS::Cinder::Volume"},
			},
			wantID: "server-one",
		},
		{
			name: "requires a choice when multiple Nova VMs exist",
			resources: []checkpoints.CheckpointDetailResource{
				{OriginalID: "server-one", Type: "OS::Nova::Server"},
				{OriginalID: "server-two", Type: "OS::Nova::Server"},
			},
			wantErrSubstr: "multiple Nova servers",
		},
		{
			name: "selects the requested Nova VM",
			resources: []checkpoints.CheckpointDetailResource{
				{OriginalID: "server-one", Type: "OS::Nova::Server"},
				{OriginalID: "server-two", Type: "OS::Nova::Server"},
			},
			requestedID: "server-two",
			wantID:      "server-two",
		},
		{
			name: "rejects a resource that is not in the checkpoint",
			resources: []checkpoints.CheckpointDetailResource{
				{OriginalID: "server-one", Type: "OS::Nova::Server"},
			},
			requestedID:   "server-missing",
			wantErrSubstr: "does not contain Nova server",
		},
		{
			name:          "rejects checkpoints without a Nova VM",
			resources:     []checkpoints.CheckpointDetailResource{{OriginalID: "volume-one", Type: "OS::Cinder::Volume"}},
			wantErrSubstr: "does not contain a Nova server",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotID, err := selectCheckpointRestoreSourceID(test.resources, test.requestedID)
			if test.wantErrSubstr != "" {
				if err == nil || !strings.Contains(err.Error(), test.wantErrSubstr) {
					t.Fatalf("selectCheckpointRestoreSourceID() error = %v, want substring %q", err, test.wantErrSubstr)
				}
				return
			}
			if err != nil {
				t.Fatalf("selectCheckpointRestoreSourceID() error = %v", err)
			}
			if gotID != test.wantID {
				t.Fatalf("selectCheckpointRestoreSourceID() = %q, want %q", gotID, test.wantID)
			}
		})
	}
}

func TestComputeRestoreParametersUsesComputeInstanceConfiguration(t *testing.T) {
	data := schema.TestResourceDataRaw(t, ResourceComputeInstance().Schema, map[string]interface{}{
		"name":              "restored-server",
		"key_pair":          "ssh-key",
		"availability_zone": "MS1",
	})

	got, err := computeRestoreParameters(data, "network-id")
	if err != nil {
		t.Fatalf("computeRestoreParameters() error = %v", err)
	}
	if got.RestoreName != "restored-server" || got.Destination != "new" || got.RestoreKeypair != "ssh-key" || got.AvailabilityZone != "MS1" || got.RestoreNetID != "network-id" {
		t.Fatalf("unexpected compute restore parameters: %#v", got)
	}
	if got.Flavor != "" || len(got.SecurityGroups) != 0 {
		t.Fatalf("restore must not send unsupported flavor/security group overrides: %#v", got)
	}
}

func TestComputeRestoreParametersRequiresTargetNetworkAndKeyPair(t *testing.T) {
	data := schema.TestResourceDataRaw(t, ResourceComputeInstance().Schema, map[string]interface{}{"name": "restored-server"})
	if _, err := computeRestoreParameters(data, ""); err == nil || !strings.Contains(err.Error(), "network") {
		t.Fatalf("computeRestoreParameters() error = %v, want missing network error", err)
	}
	if _, err := computeRestoreParameters(data, "network-id"); err == nil || !strings.Contains(err.Error(), "key_pair") {
		t.Fatalf("computeRestoreParameters() error = %v, want missing key_pair error", err)
	}
}

func TestSingleCheckpointRestoreNetworkID(t *testing.T) {
	tests := []struct {
		name          string
		networks      []InstanceNetwork
		wantID        string
		wantErrSubstr string
	}{
		{
			name:     "uses the configured network UUID",
			networks: []InstanceNetwork{{UUID: "network-id"}},
			wantID:   "network-id",
		},
		{
			name:          "requires a configured network",
			wantErrSubstr: "exactly one network",
		},
		{
			name:          "rejects multiple networks",
			networks:      []InstanceNetwork{{UUID: "network-one"}, {UUID: "network-two"}},
			wantErrSubstr: "exactly one network",
		},
		{
			name:          "requires a resolved network UUID",
			networks:      []InstanceNetwork{{Name: "network-name"}},
			wantErrSubstr: "network UUID",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := singleCheckpointRestoreNetworkID(test.networks)
			if test.wantErrSubstr != "" {
				if err == nil || !strings.Contains(err.Error(), test.wantErrSubstr) {
					t.Fatalf("singleCheckpointRestoreNetworkID() error = %v, want substring %q", err, test.wantErrSubstr)
				}
				return
			}
			if err != nil {
				t.Fatalf("singleCheckpointRestoreNetworkID() error = %v", err)
			}
			if got != test.wantID {
				t.Fatalf("singleCheckpointRestoreNetworkID() = %q, want %q", got, test.wantID)
			}
		})
	}
}

func TestGetOrCreateRestoreAdoptsMatchingRequest(t *testing.T) {
	desired := restores.RestoreCreateOpts{
		ProviderID:   "nova-provider",
		CheckpointID: "checkpoint-id",
		Parameters: map[string]restores.RestoreParameters{
			"OS::Nova::Server#server-id": {RestoreName: "restored-server", Destination: "new"},
		},
	}
	existing := []restores.Restore{{
		ID:           "restore-id",
		ProviderID:   desired.ProviderID,
		CheckpointID: desired.CheckpointID,
		Parameters: map[string]restores.RestoreParameters{
			"OS::Nova::Server#server-id": {RestoreName: "restored-server", Destination: "new", Flavor: "ignored-flavor", SecurityGroups: []string{"ignored-group"}},
		},
	}}
	createCalls := 0

	got, err := getOrCreateRestore(existing, desired, func(candidate *restores.Restore) (bool, error) {
		return canReuseExistingRestore(candidate, nil)
	}, func() (*restores.Restore, error) {
		createCalls++
		return &restores.Restore{ID: "new-restore-id"}, nil
	})
	if err != nil {
		t.Fatalf("getOrCreateRestore() error = %v", err)
	}
	if got.ID != "restore-id" || createCalls != 0 {
		t.Fatalf("getOrCreateRestore() = %#v with %d create calls, want existing restore and no create", got, createCalls)
	}
}

func TestGetOrCreateRestoreCreatesAgainWhenPreviousVMWasDeleted(t *testing.T) {
	desired := restores.RestoreCreateOpts{
		ProviderID:   "nova-provider",
		CheckpointID: "checkpoint-id",
		Parameters: map[string]restores.RestoreParameters{
			"OS::Nova::Server#source-id": {RestoreName: "restored-server", Destination: "new"},
		},
	}
	existing := []restores.Restore{{
		ID:              "old-restore-id",
		ProviderID:      desired.ProviderID,
		CheckpointID:    desired.CheckpointID,
		Parameters:      desired.Parameters,
		ResourcesStatus: map[string]string{"OS::Nova::Server#old-restored-id": "available"},
	}}
	createCalls := 0

	got, err := getOrCreateRestore(existing, desired, func(candidate *restores.Restore) (bool, error) {
		return canReuseExistingRestore(candidate, func(string) (bool, error) { return false, nil })
	}, func() (*restores.Restore, error) {
		createCalls++
		return &restores.Restore{ID: "new-restore-id"}, nil
	})
	if err != nil {
		t.Fatalf("getOrCreateRestore() error = %v", err)
	}
	if got.ID != "new-restore-id" || createCalls != 1 {
		t.Fatalf("getOrCreateRestore() = %#v with %d create calls, want a new restore after the prior VM was deleted", got, createCalls)
	}
}

func TestCanReuseExistingRestore(t *testing.T) {
	tests := []struct {
		name           string
		restore        *restores.Restore
		serverExists   bool
		wantReuse      bool
		wantServerCall bool
	}{
		{
			name:      "reuses an in-progress restore",
			restore:   &restores.Restore{ResourcesStatus: map[string]string{"OS::Nova::Server#vm": "restoring"}},
			wantReuse: true,
		},
		{
			name:           "reuses a completed restore while the VM exists",
			restore:        &restores.Restore{ResourcesStatus: map[string]string{"OS::Nova::Server#vm": "available"}},
			serverExists:   true,
			wantReuse:      true,
			wantServerCall: true,
		},
		{
			name:           "does not reuse a completed restore after its VM was deleted",
			restore:        &restores.Restore{ResourcesStatus: map[string]string{"OS::Nova::Server#vm": "available"}},
			serverExists:   false,
			wantReuse:      false,
			wantServerCall: true,
		},
		{
			name:      "adopts a failed restore so an identical request cannot create a duplicate",
			restore:   &restores.Restore{ResourcesStatus: map[string]string{"OS::Nova::Server#vm": "error"}},
			wantReuse: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			serverCalls := 0
			got, err := canReuseExistingRestore(test.restore, func(serverID string) (bool, error) {
				serverCalls++
				if serverID != "vm" {
					t.Fatalf("server ID = %q, want vm", serverID)
				}
				return test.serverExists, nil
			})
			if err != nil {
				t.Fatalf("canReuseExistingRestore() error = %v", err)
			}
			if got != test.wantReuse {
				t.Errorf("canReuseExistingRestore() = %t, want %t", got, test.wantReuse)
			}
			if (serverCalls > 0) != test.wantServerCall {
				t.Errorf("server existence checks = %d, want call=%t", serverCalls, test.wantServerCall)
			}
		})
	}
}

func TestRestoredComputeInstanceWaitState(t *testing.T) {
	tests := []struct {
		name          string
		state         string
		instanceFound bool
		wantState     string
		wantErrSubstr string
	}{
		{name: "active VM is ready", state: "ACTIVE", instanceFound: true, wantState: "ACTIVE"},
		{name: "build waits", state: "BUILD", instanceFound: true, wantState: "BUILD"},
		{name: "provisioning waits", state: "PROVISIONING", instanceFound: true, wantState: "BUILD"},
		{name: "shutoff waits for active state", state: "SHUTOFF", instanceFound: true, wantState: "BUILD"},
		{name: "not-yet-visible VM waits", state: "DELETED", wantState: "BUILD"},
		{name: "deleted VM fails", state: "DELETED", instanceFound: true, wantErrSubstr: "deleted before becoming active"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotState, err := restoredComputeInstanceWaitState(test.state, test.instanceFound)
			if test.wantErrSubstr != "" {
				if err == nil || !strings.Contains(err.Error(), test.wantErrSubstr) {
					t.Fatalf("restoredComputeInstanceWaitState() error = %v, want substring %q", err, test.wantErrSubstr)
				}
				return
			}
			if err != nil {
				t.Fatalf("restoredComputeInstanceWaitState() error = %v", err)
			}
			if gotState != test.wantState {
				t.Errorf("restoredComputeInstanceWaitState(%q, %t) = %q, want %q", test.state, test.instanceFound, gotState, test.wantState)
			}
		})
	}
}

func TestRestoredComputeInstanceStateRefreshFunc(t *testing.T) {
	t.Run("waits while the restored VM is not visible yet", func(t *testing.T) {
		var missingServer *servers.Server
		refresh := restoredComputeInstanceStateRefreshFunc(func() (interface{}, string, error) {
			return missingServer, "DELETED", nil
		}, "server-1")

		_, state, err := refresh()
		if err != nil {
			t.Fatalf("restored refresh returned error: %v", err)
		}
		if state != "BUILD" {
			t.Fatalf("restored refresh state = %q, want BUILD while Nova is not yet visible", state)
		}
	})

	t.Run("fails when a visible restored VM has been deleted", func(t *testing.T) {
		refresh := restoredComputeInstanceStateRefreshFunc(func() (interface{}, string, error) {
			return &servers.Server{ID: "server-2"}, "DELETED", nil
		}, "server-2")

		_, _, err := refresh()
		if err == nil || !strings.Contains(err.Error(), "deleted before becoming active") {
			t.Fatalf("restored refresh error = %v, want deleted VM error", err)
		}
	})
}

func TestWaitForRestoredComputeInstanceID(t *testing.T) {
	t.Run("polls until a VM is available", func(t *testing.T) {
		responses := []*restores.Restore{
			{
				ID:              "restore-1",
				ResourcesStatus: map[string]string{"OS::Nova::Server#server-1": "restoring"},
			},
			{
				ID:              "restore-1",
				ResourcesStatus: map[string]string{"OS::Nova::Server#server-1": "available"},
			},
		}
		calls := 0
		getter := func(_ context.Context, restoreID string) (*restores.Restore, error) {
			if restoreID != "restore-1" {
				t.Fatalf("restore ID = %q, want %q", restoreID, "restore-1")
			}
			responseIndex := calls
			calls++
			if responseIndex >= len(responses) {
				responseIndex = len(responses) - 1
			}
			return responses[responseIndex], nil
		}

		gotID, err := waitForRestoredComputeInstanceID(context.Background(), "restore-1", time.Second, time.Millisecond, getter)
		if err != nil {
			t.Fatalf("waitForRestoredComputeInstanceID() error = %v", err)
		}
		if gotID != "server-1" {
			t.Fatalf("waitForRestoredComputeInstanceID() ID = %q, want %q", gotID, "server-1")
		}
		if calls != 2 {
			t.Fatalf("restore getter called %d times, want 2", calls)
		}
	})

	t.Run("returns API errors", func(t *testing.T) {
		wantErr := errors.New("backup API unavailable")
		getter := func(context.Context, string) (*restores.Restore, error) {
			return nil, wantErr
		}

		_, err := waitForRestoredComputeInstanceID(context.Background(), "restore-2", time.Second, time.Millisecond, getter)
		if !errors.Is(err, wantErr) {
			t.Fatalf("waitForRestoredComputeInstanceID() error = %v, want wrapped %v", err, wantErr)
		}
	})

	t.Run("times out while the restore is pending", func(t *testing.T) {
		getter := func(_ context.Context, restoreID string) (*restores.Restore, error) {
			return &restores.Restore{
				ID:              restoreID,
				ResourcesStatus: map[string]string{"OS::Nova::Server#server-3": "restoring"},
			}, nil
		}

		_, err := waitForRestoredComputeInstanceID(context.Background(), "restore-3", time.Millisecond, time.Hour, getter)
		if err == nil || !strings.Contains(err.Error(), "restore-3") {
			t.Fatalf("waitForRestoredComputeInstanceID() error = %v, want timeout mentioning restore ID", err)
		}
	})
}
