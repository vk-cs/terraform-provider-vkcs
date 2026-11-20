package backup

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	if got := flattenObjectLock(&plans.ObjectLock{}); got != nil {
		t.Fatalf("flattenObjectLock(empty object) = %#v, want nil", got)
	}
	retainUntil := int64(1900000000)
	got = flattenObjectLock(&plans.ObjectLock{Type: "governance", RetainUntil: &retainUntil})
	if got == nil {
		t.Fatal("flattenObjectLock(retain_until) returned nil")
	}

	hasGovernanceType := got.Type.ValueString() == "governance"
	hasExpectedRetainUntil := !got.RetainUntil.IsNull() && got.RetainUntil.ValueInt64() == retainUntil
	hasNoRetentionDays := got.Days.IsNull()
	if !hasGovernanceType || !hasExpectedRetainUntil || !hasNoRetentionDays {
		t.Fatalf("flattenObjectLock(retain_until) = %#v, want original API value", got)
	}
}

func TestObjectLockFromResponsePreservesStateWhenAPIOmitsLock(t *testing.T) {
	days := int64(30)
	previous := &PlanResourceObjectLockModel{
		Type:        types.StringValue("compliance"),
		Days:        types.Int64Value(days),
		RetainUntil: types.Int64Null(),
	}

	for _, response := range []*plans.ObjectLock{nil, {}} {
		if got := objectLockFromResponse(previous, response); got != previous {
			t.Fatalf("objectLockFromResponse() = %#v, want previous value %#v", got, previous)
		}
	}
}

func TestObjectLockFromResponseUsesAPIValuesWhenReturned(t *testing.T) {
	previous := &PlanResourceObjectLockModel{
		Type:        types.StringValue("compliance"),
		Days:        types.Int64Value(30),
		RetainUntil: types.Int64Null(),
	}
	retainUntil := int64(1900000000)
	response := &plans.ObjectLock{Type: "compliance", RetainUntil: &retainUntil}

	got := objectLockFromResponse(previous, response)
	if got == previous || got == nil || got.RetainUntil.ValueInt64() != retainUntil || !got.Days.IsNull() {
		t.Fatalf("objectLockFromResponse() = %#v, want API-returned retain_until value", got)
	}
}

func TestExpandObjectLockUpdateDistinguishesOmittedAndClearedLock(t *testing.T) {
	previous := &PlanResourceObjectLockModel{
		Type:        types.StringValue("compliance"),
		Days:        types.Int64Value(30),
		RetainUntil: types.Int64Null(),
	}

	if got := expandObjectLockUpdate(nil, nil); got != nil {
		t.Fatalf("expandObjectLockUpdate(nil, nil) = %#v, want field omitted", got)
	}
	if got := expandObjectLockUpdate(nil, previous); got == nil || *got != nil {
		t.Fatalf("expandObjectLockUpdate(nil, previous) = %#v, want explicit null to clear", got)
	}
	updated := &PlanResourceObjectLockModel{
		Type:        types.StringValue("governance"),
		Days:        types.Int64Null(),
		RetainUntil: types.Int64Value(1900000000),
	}
	if got := expandObjectLockUpdate(updated, previous); got == nil || *got == nil || (*got).Type != "governance" {
		t.Fatalf("expandObjectLockUpdate(updated, previous) = %#v, want governance lock", got)
	}
}

func TestValidateObjectLockSettings(t *testing.T) {
	tests := []struct {
		name       string
		objectLock PlanResourceObjectLockModel
		wantErrors int
	}{
		{
			name: "compliance with days",
			objectLock: PlanResourceObjectLockModel{
				Type: types.StringValue("compliance"), Days: types.Int64Value(30), RetainUntil: types.Int64Null(),
			},
		},
		{
			name: "governance with retain_until",
			objectLock: PlanResourceObjectLockModel{
				Type: types.StringValue("governance"), Days: types.Int64Null(), RetainUntil: types.Int64Value(1900000000),
			},
		},
		{
			name: "legalhold without retention",
			objectLock: PlanResourceObjectLockModel{
				Type: types.StringValue("legalhold"), Days: types.Int64Null(), RetainUntil: types.Int64Null(),
			},
		},
		{
			name: "compliance requires one retention field",
			objectLock: PlanResourceObjectLockModel{
				Type: types.StringValue("compliance"), Days: types.Int64Null(), RetainUntil: types.Int64Null(),
			},
			wantErrors: 1,
		},
		{
			name: "governance rejects both retention fields",
			objectLock: PlanResourceObjectLockModel{
				Type: types.StringValue("governance"), Days: types.Int64Value(30), RetainUntil: types.Int64Value(1900000000),
			},
			wantErrors: 1,
		},
		{
			name: "legalhold rejects retention days",
			objectLock: PlanResourceObjectLockModel{
				Type: types.StringValue("legalhold"), Days: types.Int64Value(30), RetainUntil: types.Int64Null(),
			},
			wantErrors: 1,
		},
		{
			name: "legalhold rejects retain_until",
			objectLock: PlanResourceObjectLockModel{
				Type: types.StringValue("legalhold"), Days: types.Int64Null(), RetainUntil: types.Int64Value(1900000000),
			},
			wantErrors: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			diagnostics := validateObjectLockSettings(path.Root("object_lock"), test.objectLock)
			if diagnostics.ErrorsCount() != test.wantErrors {
				t.Fatalf("validateObjectLockSettings() returned %d errors, want %d: %v", diagnostics.ErrorsCount(), test.wantErrors, diagnostics)
			}
		})
	}
}

func TestObjectLockConfigValidator(t *testing.T) {
	ctx := context.Background()
	objectType := map[string]attr.Type{
		"type":         types.StringType,
		"days":         types.Int64Type,
		"retain_until": types.Int64Type,
	}

	for _, test := range []struct {
		name       string
		attributes map[string]attr.Value
		wantError  bool
	}{
		{
			name: "valid compliance retention",
			attributes: map[string]attr.Value{
				"type":         types.StringValue("compliance"),
				"days":         types.Int64Value(30),
				"retain_until": types.Int64Null(),
			},
		},
		{
			name: "compliance missing retention",
			attributes: map[string]attr.Value{
				"type":         types.StringValue("compliance"),
				"days":         types.Int64Null(),
				"retain_until": types.Int64Null(),
			},
			wantError: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			configValue, diags := types.ObjectValue(objectType, test.attributes)
			if diags.HasError() {
				t.Fatalf("build object config: %v", diags)
			}

			response := &validator.ObjectResponse{}
			objectLockConfigValidator{}.ValidateObject(ctx, validator.ObjectRequest{
				Path:        path.Root("object_lock"),
				ConfigValue: configValue,
			}, response)
			if response.Diagnostics.HasError() != test.wantError {
				t.Fatalf("Object Lock validator errors = %v, want error=%t", response.Diagnostics, test.wantError)
			}
		})
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
