package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/iam/resource_s3_account"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/iam/s3accounts"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

var (
	_ resource.Resource                = (*s3AccountResource)(nil)
	_ resource.ResourceWithConfigure   = (*s3AccountResource)(nil)
	_ resource.ResourceWithImportState = (*s3AccountResource)(nil)
)

func NewS3AccountResource() resource.Resource {
	return &s3AccountResource{}
}

type s3AccountResource struct {
	config clients.Config
}

func (r *s3AccountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_s3_account"
}

func (r *s3AccountResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_s3_account.S3AccountResourceSchema(ctx)
}

func (r *s3AccountResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *s3AccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_s3_account.S3AccountModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}
	data.Region = types.StringValue(region)

	client, err := r.config.IAMServiceUsersV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating IAM API client", err.Error())
		return
	}

	createOpts := s3accounts.CreateOpts{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
	}

	tflog.Trace(ctx, "Calling IAM Service Users API to create S3 account", map[string]any{"opts": fmt.Sprintf("%#v", createOpts)})

	createResp, err := s3accounts.Create(client, &createOpts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling IAM Service Users API to create S3 account", err.Error())
		return
	}

	tflog.Trace(ctx, "Called IAM Service Users API to create S3 account", map[string]any{"create_response": fmt.Sprintf("%#v", createResp)})

	resp.Diagnostics.Append(data.UpdateFromCreateS3AccountResponse(ctx, createResp)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *s3AccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_s3_account.S3AccountModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}
	data.Region = types.StringValue(region)

	client, err := r.config.IAMServiceUsersV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating IAM Service Users API client", err.Error())
		return
	}

	id := data.Id.ValueString()
	ctx = tflog.SetField(ctx, "s3_account_id", id)

	tflog.Trace(ctx, "Calling IAM Service Users API to retrieve S3 account")

	s3Account, err := s3accounts.Get(client, id).Extract()
	if errutil.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error calling IAM Service Users API to retrieve S3 account", err.Error())
		return
	}

	tflog.Trace(ctx, "Retrieved S3 account", map[string]any{"s3_account": fmt.Sprintf("%#v", s3Account)})

	resp.Diagnostics.Append(data.UpdateFromS3Account(ctx, s3Account)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *s3AccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_s3_account.S3AccountModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddError("Unable to update the S3 account",
		"Not implemented. Please report this issue to the provider developers.")
}

func (r *s3AccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_s3_account.S3AccountModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.config.IAMServiceUsersV1Client(data.Region.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error creating IAM Service Users API client", err.Error())
		return
	}

	id := data.Id.ValueString()
	ctx = tflog.SetField(ctx, "s3_account_id", id)

	tflog.Trace(ctx, "Calling IAM Service Users API to delete S3 account")

	err = s3accounts.Delete(client, id).ExtractErr()
	if errutil.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error calling IAM Service Users API to delete S3 account", err.Error())
		return
	}

	tflog.Trace(ctx, "Called IAM Service Users API to delete S3 account")
}

func (r *s3AccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
