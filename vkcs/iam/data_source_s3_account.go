package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/iam/datasource_s3_account"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/iam/s3accounts"
)

var (
	_ datasource.DataSource                     = (*s3AccountDataSource)(nil)
	_ datasource.DataSourceWithConfigure        = (*s3AccountDataSource)(nil)
	_ datasource.DataSourceWithConfigValidators = (*s3AccountDataSource)(nil)
)

func NewS3AccountDataSource() datasource.DataSource {
	return &s3AccountDataSource{}
}

type s3AccountDataSource struct {
	config clients.Config
}

func (d *s3AccountDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_s3_account"
}

func (d *s3AccountDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_s3_account.S3AccountDataSourceSchema(ctx)
}

func (d *s3AccountDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *s3AccountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_s3_account.S3AccountModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = d.config.GetRegion()
	}
	data.Region = types.StringValue(region)

	client, err := d.config.IAMServiceUsersV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating IAM Service Users API client", err.Error())
		return
	}

	var s3Account *s3accounts.S3Account

	if id := data.Id.ValueString(); id != "" {
		ctx = tflog.SetField(ctx, "s3_account_id", id)
		tflog.Trace(ctx, "Calling IAM Service Users API to get S3 account")

		var err error
		s3Account, err = s3accounts.Get(client, id).Extract()
		if err != nil {
			resp.Diagnostics.AddError("Error calling IAM Service Users API to get S3 account", err.Error())
			return
		}

		tflog.Trace(ctx, "Called IAM Service Users API to get S3 account", map[string]any{"s3_account": fmt.Sprintf("%#v", s3Account)})
	} else {
		listOpts := s3accounts.ListOpts{
			Name: data.Name.ValueString(),
		}

		tflog.Trace(ctx, "Calling IAM Service Users API to list S3 accounts", map[string]any{"opts": fmt.Sprintf("%#v", listOpts)})

		allPages, err := s3accounts.List(client, listOpts).AllPages()
		if err != nil {
			resp.Diagnostics.AddError("Error calling IAM Service Users API to list S3 accounts", err.Error())
			return
		}

		tflog.Trace(ctx, "Called IAM Service Users API to list S3 accounts", map[string]any{"all_pages": fmt.Sprintf("%#v", allPages)})

		s3Accounts, err := s3accounts.ExtractS3Accounts(allPages)
		if err != nil {
			resp.Diagnostics.AddError("Error extracting S3 accounts", err.Error())
			return
		}

		tflog.Trace(ctx, "Extracted list of S3 accounts", map[string]any{"s3_accounts": fmt.Sprintf("%#v", s3Accounts)})

		if len(s3Accounts) < 1 {
			resp.Diagnostics.AddError("Your query returned no results", "Please change your search criteria and try again")
			return
		}

		if len(s3Accounts) > 1 {
			resp.Diagnostics.AddError("Your query returned more than one result", "Please try a more specific search criteria")
			return
		}

		s3Account = &s3Accounts[0]
	}

	resp.Diagnostics.Append(data.UpdateFromS3Account(ctx, s3Account)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *s3AccountDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("name"),
		),
	}
}
