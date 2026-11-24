package compute

import (
	"context"
	"strings"
	"testing"

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
			name: "rejects adoption of the source VM",
			restore: &restores.Restore{
				ID:              "restore-source-id",
				ResourcesStatus: map[string]string{"OS::Nova::Server#source-id": "available"},
			},
			wantErrSubstr: "returned the source Nova server ID",
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
			gotID, gotReady, err := restoredComputeInstanceID(test.restore, "source-id")
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

				response := &schema.ValidateResourceConfigFuncResponse{}
				validateCheckpointRestoreRawConfig(context.Background(), schema.ValidateResourceConfigFuncRequest{RawConfig: config}, response)
				if !response.Diagnostics.HasError() || !strings.Contains(response.Diagnostics[0].Detail, "network."+attribute) {
					t.Fatalf("validateCheckpointRestoreRawConfig() diagnostics = %v, want unsupported network.%s error", response.Diagnostics, attribute)
				}
			})
		}
	}
}

func TestValidateCheckpointRestoreConfigDefersUnknownCheckpointIDToCreateGuard(t *testing.T) {
	network := cty.ObjectVal(map[string]cty.Value{
		"uuid":           cty.StringVal("network-id"),
		"name":           cty.NullVal(cty.String),
		"port":           cty.StringVal("unsupported-port"),
		"fixed_ip_v4":    cty.NullVal(cty.String),
		"mac":            cty.NullVal(cty.String),
		"access_network": cty.BoolVal(false),
	})
	deferredConfig := cty.ObjectVal(map[string]cty.Value{
		"checkpoint_id": cty.UnknownVal(cty.String),
		"network":       cty.ListVal([]cty.Value{network}),
	})

	if err := validateCheckpointRestoreConfig(deferredConfig); err != nil {
		t.Fatalf("schema-level raw-config validation should defer an unknown checkpoint_id, got %v", err)
	}
	resolvedConfig := cty.ObjectVal(map[string]cty.Value{
		"checkpoint_id": cty.StringVal("checkpoint-id"),
		"network":       cty.ListVal([]cty.Value{network}),
	})
	if err := validateCheckpointRestoreNetworkConfig(resolvedConfig); err == nil || !strings.Contains(err.Error(), "network.port") {
		t.Fatalf("create-time network validation error = %v, want unsupported network.port error", err)
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

	got, err := getOrCreateRestore(context.Background(), nil, nil, existing, desired, "source-id")
	if err != nil {
		t.Fatalf("getOrCreateRestore() error = %v", err)
	}
	if got.ID != "restore-id" {
		t.Fatalf("getOrCreateRestore() = %#v, want the matching in-progress restore", got)
	}
}

func TestCanReuseExistingRestoreAdoptsFailedRestore(t *testing.T) {
	restore := &restores.Restore{
		ID:              "restore-id",
		ResourcesStatus: map[string]string{"OS::Nova::Server#source-id": "error"},
	}

	got, err := canReuseExistingRestore(context.Background(), nil, restore, "source-id")
	if err != nil {
		t.Fatalf("canReuseExistingRestore() error = %v", err)
	}
	if !got {
		t.Fatal("canReuseExistingRestore() = false, want failed restore reused to report its failure")
	}
}
