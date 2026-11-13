package backup

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/plans"
)

func expandObjectLock(config *PlanResourceObjectLockModel) *plans.ObjectLock {
	if config == nil || config.Type.IsNull() || config.Type.IsUnknown() {
		return nil
	}

	objectLock := &plans.ObjectLock{Type: config.Type.ValueString()}
	if !config.Days.IsNull() && !config.Days.IsUnknown() {
		days := config.Days.ValueInt64()
		objectLock.Days = &days
	}
	if !config.RetainUntil.IsNull() && !config.RetainUntil.IsUnknown() {
		retainUntil := config.RetainUntil.ValueInt64()
		objectLock.RetainUntil = &retainUntil
	}

	return objectLock
}

func flattenObjectLock(objectLock *plans.ObjectLock) *PlanResourceObjectLockModel {
	if objectLock == nil {
		return nil
	}

	config := &PlanResourceObjectLockModel{
		Type:        types.StringValue(objectLock.Type),
		Days:        types.Int64Null(),
		RetainUntil: types.Int64Null(),
	}
	if objectLock.Days != nil {
		config.Days = types.Int64Value(*objectLock.Days)
	}
	if objectLock.RetainUntil != nil {
		config.RetainUntil = types.Int64Value(*objectLock.RetainUntil)
	}

	return config
}
