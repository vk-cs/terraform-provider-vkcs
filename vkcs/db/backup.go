package db

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/backups"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/datastores"
)

type BackupDataStoreModel struct {
	Type    types.String `tfsdk:"type"`
	Version types.String `tfsdk:"version"`
}

func (m BackupDataStoreModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":    types.StringType,
		"version": types.StringType,
	}
}

func flattenBackupDatastore(ctx context.Context, d datastores.DatastoreShort, respDiags *diag.Diagnostics) types.List {
	datastores, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: BackupDataStoreModel{}.AttrTypes()}, []BackupDataStoreModel{
		{
			Type:    types.StringValue(d.Type),
			Version: types.StringValue(d.Version),
		},
	})
	respDiags.Append(diags...)
	return datastores
}

func backupStateRefreshFunc(client *gophercloud.ServiceClient, backupID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		b, err := backups.Get(client, backupID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return b, dbBackupStatusDeleted, nil
			}
			return nil, "", err
		}

		if b.Status == dbBackupStatusError {
			return b, b.Status, fmt.Errorf("there was an error creating the database backup")
		}

		return b, b.Status, nil
	}
}
