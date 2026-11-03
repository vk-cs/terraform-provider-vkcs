package backup

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/checkpoints"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/restores"
)

func TestSelectRestoreResourceID(t *testing.T) {
	resources := []checkpoints.CheckpointDetailResource{
		{OriginalID: "server-one", Type: NovaInstance},
		{OriginalID: "volume-one", Type: CinderVolume},
	}

	got, err := selectRestoreResourceID(resources, "")
	if err != nil {
		t.Fatalf("select sole Nova resource: %v", err)
	}
	if got != "server-one" {
		t.Fatalf("unexpected selected resource ID: %q", got)
	}
}

func TestSelectRestoreResourceIDRequiresChoiceWhenCheckpointHasMultipleVMs(t *testing.T) {
	resources := []checkpoints.CheckpointDetailResource{
		{OriginalID: "server-one", Type: NovaInstance},
		{OriginalID: "server-two", Type: NovaInstance},
	}

	if _, err := selectRestoreResourceID(resources, ""); err == nil {
		t.Fatal("expected an error when multiple VM resources are present")
	}
	got, err := selectRestoreResourceID(resources, "server-two")
	if err != nil {
		t.Fatalf("select explicit Nova resource: %v", err)
	}
	if got != "server-two" {
		t.Fatalf("unexpected selected resource ID: %q", got)
	}
}

func TestRestoreStatusesCompleteOnlyWhenAllResourcesAreTerminal(t *testing.T) {
	tests := []struct {
		name     string
		statuses map[string]string
		want     bool
	}{
		{name: "empty response", statuses: map[string]string{}, want: false},
		{name: "in progress", statuses: map[string]string{"server-id": "restoring"}, want: false},
		{name: "success", statuses: map[string]string{"server-id": "available"}, want: true},
		{name: "partial failure terminal", statuses: map[string]string{"server-one": "available", "server-two": "error"}, want: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := restoreStatusesComplete(test.statuses); got != test.want {
				t.Fatalf("restoreStatusesComplete() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestRestoreFailureDetailsAreSortedAndIncludeReasons(t *testing.T) {
	statuses := map[string]string{"server-b": "error", "server-a": "error", "server-c": "available"}
	reasons := map[string]string{"server-a": "not found", "server-b": "quota exceeded"}
	want := []string{"server-a: not found", "server-b: quota exceeded"}

	if got := restoreFailureDetails(statuses, reasons); !reflect.DeepEqual(got, want) {
		t.Fatalf("restoreFailureDetails() = %#v, want %#v", got, want)
	}
}

func TestRestoreOutcomeErrorFailsApplyWhenAnyResourceFailed(t *testing.T) {
	err := restoreOutcomeError(
		"restore-id",
		map[string]string{"server-a": "available", "server-b": "error"},
		map[string]string{"server-b": "quota exceeded"},
	)
	if err == nil {
		t.Fatal("expected a partial restore failure to return an error")
	}
	if got := err.Error(); got != "restore restore-id completed with errors: server-b: quota exceeded" {
		t.Fatalf("unexpected restore error: %q", got)
	}
}

func TestRestoreOutcomeErrorReturnsNilWhenAllResourcesSucceeded(t *testing.T) {
	err := restoreOutcomeError("restore-id", map[string]string{"server-id": "available"}, nil)
	if err != nil {
		t.Fatalf("restoreOutcomeError() = %v, want nil", err)
	}
}

func TestFindMatchingRestoreAdoptsEquivalentRequest(t *testing.T) {
	desired := restores.RestoreCreateOpts{
		ProviderID:   "nova-provider",
		CheckpointID: "checkpoint-id",
		Parameters: map[string]restores.RestoreParameters{
			"OS::Nova::Server#server-id": {
				RestoreName:    "restored-vm",
				Destination:    "new",
				EnableFIP:      boolPointer(false),
				VolumeSize:     restoreTestInt64Pointer(0),
				SecurityGroups: []string{"sg-b", "sg-a"},
			},
		},
	}
	existing := []restores.Restore{
		{
			ID:           "existing-restore-id",
			ProviderID:   "nova-provider",
			CheckpointID: "checkpoint-id",
			Parameters: map[string]restores.RestoreParameters{
				"OS::Nova::Server#server-id": {
					RestoreName:    "restored-vm",
					Destination:    "new",
					SecurityGroups: []string{"sg-a", "sg-b"},
				},
			},
		},
	}

	got := findMatchingRestore(existing, desired)
	if got == nil || got.ID != "existing-restore-id" {
		t.Fatalf("findMatchingRestore() = %#v, want existing restore", got)
	}
}

func TestFindMatchingRestoreDoesNotAdoptDifferentRestoreName(t *testing.T) {
	desired := restores.RestoreCreateOpts{
		ProviderID:   "nova-provider",
		CheckpointID: "checkpoint-id",
		Parameters: map[string]restores.RestoreParameters{
			"OS::Nova::Server#server-id": {RestoreName: "new-vm", Destination: "new"},
		},
	}
	existing := []restores.Restore{
		{
			ID:           "old-restore-id",
			ProviderID:   "nova-provider",
			CheckpointID: "checkpoint-id",
			Parameters: map[string]restores.RestoreParameters{
				"OS::Nova::Server#server-id": {RestoreName: "old-vm", Destination: "new"},
			},
		},
	}

	if got := findMatchingRestore(existing, desired); got != nil {
		t.Fatalf("findMatchingRestore() unexpectedly matched %#v", got)
	}
}

func TestGetOrCreateRestoreDoesNotCreateWhenMatchingRestoreExists(t *testing.T) {
	desired := restores.RestoreCreateOpts{
		ProviderID:   "nova-provider",
		CheckpointID: "checkpoint-id",
		Parameters: map[string]restores.RestoreParameters{
			"OS::Nova::Server#server-id": {RestoreName: "restored-vm", Destination: "new"},
		},
	}
	existing := []restores.Restore{
		{
			ID:           "existing-id",
			ProviderID:   desired.ProviderID,
			CheckpointID: desired.CheckpointID,
			Parameters:   desired.Parameters,
		},
	}
	createCalls := 0

	got, err := getOrCreateRestore(existing, desired, func() (*restores.Restore, error) {
		createCalls++
		return &restores.Restore{ID: "new-id"}, nil
	})
	if err != nil {
		t.Fatalf("getOrCreateRestore() error = %v", err)
	}
	if got.ID != "existing-id" || createCalls != 0 {
		t.Fatalf("got restore %#v and %d create calls, want existing restore and no create", got, createCalls)
	}
}

func TestGetOrCreateRestoreCreatesWhenNoMatchingRestoreExists(t *testing.T) {
	desired := restores.RestoreCreateOpts{
		ProviderID:   "nova-provider",
		CheckpointID: "checkpoint-id",
		Parameters: map[string]restores.RestoreParameters{
			"OS::Nova::Server#server-id": {RestoreName: "restored-vm", Destination: "new"},
		},
	}
	createCalls := 0

	got, err := getOrCreateRestore(nil, desired, func() (*restores.Restore, error) {
		createCalls++
		return &restores.Restore{ID: "new-id"}, nil
	})
	if err != nil {
		t.Fatalf("getOrCreateRestore() error = %v", err)
	}
	if got.ID != "new-id" || createCalls != 1 {
		t.Fatalf("got restore %#v and %d create calls, want newly created restore", got, createCalls)
	}
}

func TestExpandRestoreParametersDefaultsToNewVM(t *testing.T) {
	model := RestoreResourceModel{
		RestoreName:             types.StringValue("restored-server"),
		RestoreKeypair:          types.StringValue("keypair-name"),
		AvailabilityZone:        types.StringNull(),
		RestoreNetID:            types.StringValue("network-id"),
		Flavor:                  types.StringValue("flavor-id"),
		VolumeType:              types.StringNull(),
		VolumeSize:              types.Int64Null(),
		EnableFIP:               types.BoolNull(),
		SecurityGroups:          types.SetNull(types.StringType),
		VolumeTargetVMID:        types.StringNull(),
		VolumeAutoresizeEnabled: types.Int64Null(),
		VolumeAutoresizeMaxSize: types.Int64Null(),
		CloudMonitoringEnabled:  types.BoolNull(),
	}

	parameters, diags := expandRestoreParameters(context.Background(), model)
	if diags.HasError() {
		t.Fatalf("expandRestoreParameters() diagnostics: %v", diags)
	}
	if parameters.Destination != "new" || parameters.RestoreName != "restored-server" || parameters.RestoreNetID != "network-id" || parameters.RestoreKeypair != "keypair-name" || parameters.Flavor != "flavor-id" {
		t.Fatalf("unexpected expanded restore parameters: %#v", parameters)
	}
}

func restoreTestInt64Pointer(value int64) *int64 {
	return &value
}

func boolPointer(value bool) *bool {
	return &value
}
