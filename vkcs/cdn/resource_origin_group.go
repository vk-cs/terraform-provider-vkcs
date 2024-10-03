package cdn

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/cdn/resource_origin_group"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	fwutils "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/framework/utils"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/cdn/v1/origingroups"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

var (
	_ resource.Resource                = (*originGroupResource)(nil)
	_ resource.ResourceWithConfigure   = (*originGroupResource)(nil)
	_ resource.ResourceWithImportState = (*originGroupResource)(nil)
)

func NewOriginGroupResource() resource.Resource {
	return &originGroupResource{}
}

type originGroupResource struct {
	config clients.Config
}

func (r *originGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cdn_origin_group"
}

func (r *originGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_origin_group.OriginGroupResourceSchema(ctx)
}

func (r *originGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *originGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_origin_group.OriginGroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}
	data.Region = types.StringValue(region)

	client, err := r.config.CDNV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating CDN API client", err.Error())
		return
	}

	origins, diags := resource_origin_group.ExpandOrigins(ctx, data.Origins)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpts := origingroups.CreateOpts{
		Name:    data.Name.ValueString(),
		Origins: origins,
		UseNext: data.UseNext.ValueBool(),
	}

	tflog.Trace(ctx, "Calling CDN API to create origin group", map[string]interface{}{"opts": fmt.Sprintf("%#v", createOpts)})

	originGroup, err := origingroups.Create(client, r.config.GetTenantID(), &createOpts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling CDN API to create origin group", err.Error())
		return
	}

	tflog.Trace(ctx, "Called CDN API to create origin group", map[string]interface{}{"origin_group": fmt.Sprintf("%#v", originGroup)})

	id := types.Int64Value(int64(originGroup.ID))
	resp.State.SetAttribute(ctx, path.Root("id"), id)

	resp.Diagnostics.Append(data.UpdateFromOriginGroup(ctx, originGroup)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *originGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_origin_group.OriginGroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}
	data.Region = types.StringValue(region)

	client, err := r.config.CDNV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating CDN API client", err.Error())
		return
	}

	id := int(data.Id.ValueInt64())
	ctx = tflog.SetField(ctx, "origin_group_id", id)

	tflog.Trace(ctx, "Calling CDN API to retrieve origin group")

	originGroup, err := origingroups.Get(client, r.config.GetTenantID(), id).Extract()
	if errutil.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error calling CDN API to retrieve origin group", err.Error())
		return
	}

	tflog.Trace(ctx, "Called CDN API to retrieve origin group", map[string]interface{}{"origin_group": fmt.Sprintf("%#v", originGroup)})

	resp.Diagnostics.Append(data.UpdateFromOriginGroup(ctx, originGroup)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *originGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan resource_origin_group.OriginGroupModel
	var data resource_origin_group.OriginGroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := plan.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}
	data.Region = types.StringValue(region)

	client, err := r.config.CDNV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating CDN API client", err.Error())
		return
	}

	origins, diags := resource_origin_group.ExpandOrigins(ctx, plan.Origins)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := int(data.Id.ValueInt64())
	ctx = tflog.SetField(ctx, "origin_group_id", id)

	updateOpts := origingroups.UpdateOpts{
		Name:    plan.Name.ValueString(),
		Origins: origins,
		UseNext: plan.UseNext.ValueBool(),
	}

	tflog.Trace(ctx, "Calling CDN API to update origin group", map[string]interface{}{"opts": fmt.Sprintf("%#v", updateOpts)})

	originGroup, err := origingroups.Update(client, r.config.GetTenantID(), id, &updateOpts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling CDN API to update origin group", err.Error())
		return
	}

	tflog.Trace(ctx, "Called CDN API to update origin group", map[string]interface{}{"origin_group": fmt.Sprintf("%#v", originGroup)})

	resp.Diagnostics.Append(data.UpdateFromOriginGroup(ctx, originGroup)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *originGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_origin_group.OriginGroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	client, err := r.config.CDNV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating CDN API client", err.Error())
		return
	}

	id := int(data.Id.ValueInt64())
	ctx = tflog.SetField(ctx, "origin_group_id", id)

	tflog.Trace(ctx, "Calling CDN API to delete origin group")

	err = origingroups.Delete(client, r.config.GetTenantID(), id).ExtractErr()
	if errutil.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error calling CDN API to delete origin group", err.Error())
		return
	}

	tflog.Trace(ctx, "Called CDN API to delete origin group")
}

func (r *originGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	fwutils.ImportStatePassthroughInt64ID(ctx, req, resp)
}
