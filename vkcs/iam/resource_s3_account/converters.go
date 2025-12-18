package resource_s3_account

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/iam/s3accounts"
)

func (m *S3AccountModel) UpdateFromCreateS3AccountResponse(ctx context.Context, createResp *s3accounts.CreateS3AccountResponse) diag.Diagnostics {
	if createResp == nil {
		return nil
	}

	var diags diag.Diagnostics

	diags.Append(m.UpdateFromS3Account(ctx, &createResp.S3Account)...)
	if diags.HasError() {
		return diags
	}

	m.SecretKey = types.StringValue(createResp.SecretKey)

	return nil
}

func (m *S3AccountModel) UpdateFromS3Account(ctx context.Context, s3Account *s3accounts.S3Account) diag.Diagnostics {
	if s3Account == nil {
		return nil
	}

	m.Id = types.StringValue(s3Account.ID)
	m.AccessKey = types.StringValue(s3Account.AccessKey)
	m.AccountId = types.StringValue(s3Account.AccountID)
	m.AccountName = types.StringValue(s3Account.AccountName)
	m.Name = types.StringValue(s3Account.Name)
	m.CreatedAt = types.StringValue(s3Account.CreatedAt)
	m.Description = types.StringValue(s3Account.Description)

	return nil
}
