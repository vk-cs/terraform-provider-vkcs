package backup

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/checkpoints"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/plans"
)

func TestExpandObjectLock(t *testing.T) {
	config := &PlanResourceObjectLockModel{
		Type:        types.StringValue("governance"),
		Days:        types.Int64Null(),
		RetainUntil: types.Int64Value(1900000000),
	}

	got := expandObjectLock(config)
	if got == nil {
		t.Fatal("expandObjectLock() returned nil")
	}
	want := &plans.ObjectLock{
		Type:        "governance",
		RetainUntil: int64Pointer(1900000000),
	}
	if got.Type != want.Type || got.Days != want.Days || got.RetainUntil == nil || *got.RetainUntil != *want.RetainUntil {
		t.Fatalf("expandObjectLock() = %#v, want %#v", got, want)
	}
}

func TestFlattenObjectLock(t *testing.T) {
	days := int64(30)
	got := flattenObjectLock(&plans.ObjectLock{Type: "compliance", Days: &days})
	if got == nil || got.Type.ValueString() != "compliance" || got.Days.ValueInt64() != days || !got.RetainUntil.IsNull() {
		t.Fatalf("flattenObjectLock() = %#v, want compliance lock for 30 days", got)
	}
	if got := flattenObjectLock(nil); got != nil {
		t.Fatalf("flattenObjectLock(nil) = %#v, want nil", got)
	}
}

func TestFlattenCheckpointResourcesKeepsLegalHold(t *testing.T) {
	resources := []checkpoints.CheckpointResource{
		{
			ID:   "server-id",
			Name: "server",
			Backups: []checkpoints.CheckpointBackup{
				{
					CheckpointID: "checkpoint-id",
					ObjectLock:   &plans.ObjectLock{Type: "legalhold"},
				},
			},
		},
	}

	flattened := flattenCheckpointResources(resources)
	if len(flattened) != 1 || len(flattened[0].Backups) != 1 {
		t.Fatalf("unexpected flattened resources: %#v", flattened)
	}
	if got := flattened[0].Backups[0].ObjectLock; got == nil || got.Type.ValueString() != "legalhold" {
		t.Fatalf("unexpected flattened object_lock: %#v", got)
	}
}

func int64Pointer(value int64) *int64 {
	return &value
}
