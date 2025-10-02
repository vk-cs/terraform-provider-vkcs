package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/iam/resource_service_user"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/iam/serviceusers"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

var (
	_ resource.Resource                = (*serviceUserResource)(nil)
	_ resource.ResourceWithConfigure   = (*serviceUserResource)(nil)
	_ resource.ResourceWithImportState = (*serviceUserResource)(nil)
)

func NewServiceUserResource() resource.Resource {
	return &serviceUserResource{}
}

type serviceUserResource struct {
	config clients.Config
}

func (r *serviceUserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_service_user"
}

func (r *serviceUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_service_user.ServiceUserResourceSchema(ctx)
}

func (r *serviceUserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *serviceUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_service_user.ServiceUserModel

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

	var roleNames []string
	resp.Diagnostics.Append(data.RoleNames.ElementsAs(ctx, &roleNames, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpts := serviceusers.CreateOpts{
		Name:        data.Name.ValueString(),
		RoleNames:   roleNames,
		Description: data.Description.ValueString(),
	}

	tflog.Trace(ctx, "Calling IAM Service Users API to create service user", map[string]any{"opts": fmt.Sprintf("%#v", createOpts)})

	createResp, err := serviceusers.Create(client, &createOpts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling IAM Service Users API to create service user", err.Error())
		return
	}

	tflog.Trace(ctx, "Called IAM Service Users API to create service user", map[string]any{"create_response": fmt.Sprintf("%#v", createResp)})

	resp.Diagnostics.Append(data.UpdateFromCreateServiceUserResponse(ctx, createResp)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *serviceUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_service_user.ServiceUserModel

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
	ctx = tflog.SetField(ctx, "service_user_id", id)

	tflog.Trace(ctx, "Calling IAM Service Users API to retrieve service user")

	serviceUser, err := serviceusers.Get(client, id).Extract()
	if errutil.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error calling IAM Service Users API to retrieve service user", err.Error())
		return
	}

	tflog.Trace(ctx, "Retrieved service user", map[string]any{"service_user": fmt.Sprintf("%#v", serviceUser)})

	resp.Diagnostics.Append(data.UpdateFromServiceUser(ctx, serviceUser)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *serviceUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_service_user.ServiceUserModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddError("Unable to update the service user",
		"Not implemented. Please report this issue to the provider developers.")
}

func (r *serviceUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_service_user.ServiceUserModel

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
	ctx = tflog.SetField(ctx, "service_user_id", id)

	tflog.Trace(ctx, "Calling IAM Service Users API to delete service user")

	err = serviceusers.Delete(client, id).ExtractErr()
	if errutil.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error calling IAM Service Users API to delete service user", err.Error())
		return
	}

	tflog.Trace(ctx, "Called IAM Service Users API to delete service user")
}

func (r *serviceUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
