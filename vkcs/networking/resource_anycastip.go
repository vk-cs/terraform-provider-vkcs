package networking

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	inetworking "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking/v2/anycastips"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/networking/resource_anycastip"
)

var (
	_ resource.Resource                   = (*anycastIPResource)(nil)
	_ resource.ResourceWithConfigure      = (*anycastIPResource)(nil)
	_ resource.ResourceWithImportState    = (*anycastIPResource)(nil)
	_ resource.ResourceWithValidateConfig = (*anycastIPResource)(nil)
)

func NewAnycastIPResource() resource.Resource {
	return &anycastIPResource{}
}

type anycastIPResource struct {
	config clients.Config
}

func (r *anycastIPResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networking_anycastip"

}
func (r *anycastIPResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = resource_anycastip.AnycastipResourceSchema(ctx)
}

func (r *anycastIPResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *anycastIPResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config resource_anycastip.AnycastipModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	associations, diags := resource_anycastip.ExpandAssociations(ctx, config.Associations)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	healthCheck := resource_anycastip.ExpandHealthCheck(ctx, config.HealthCheck)

	// Editing HealthCheck is forbidden for Octavia associations
	if resource_anycastip.HasOctavia(associations) && healthCheck != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("health_check"),
			"Conflicting Attribute Value",
			"setting health_check is forbidden for octavia associations")
	}
}

func (r *anycastIPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_anycastip.AnycastipModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}
	data.Region = types.StringValue(region)

	client, err := r.config.NetworkingV2Client(region, inetworking.SprutSDN)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS network client", err.Error())
		return
	}

	associations, diags := resource_anycastip.ExpandAssociations(ctx, data.Associations)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	healthCheck := resource_anycastip.ExpandHealthCheck(ctx, data.HealthCheck)

	createOpts := anycastips.CreateOpts{
		Name:         data.Name.ValueString(),
		Description:  data.Description.ValueString(),
		NetworkID:    data.NetworkId.ValueString(),
		Associations: associations,
		HealthCheck:  healthCheck,
	}

	tflog.Trace(ctx, "Calling Networking API to create anycast IP", map[string]interface{}{"opts": fmt.Sprintf("%#v", createOpts)})

	anycastIP, err := anycastips.Create(client, &createOpts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling Networking API to create anycast IP", err.Error())
		return
	}

	tflog.Trace(ctx, "Called Networking API to create anycast IP", map[string]interface{}{"anycast_ip": fmt.Sprintf("%#v", anycastIP)})

	id := types.StringValue(anycastIP.ID)
	resp.State.SetAttribute(ctx, path.Root("id"), id)

	resp.Diagnostics.Append(data.UpdateFromAnycastIP(ctx, anycastIP)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *anycastIPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_anycastip.AnycastipModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}
	data.Region = types.StringValue(region)

	client, err := r.config.NetworkingV2Client(region, inetworking.SprutSDN)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Networking API client", err.Error())
		return
	}

	id := data.Id.ValueString()
	ctx = tflog.SetField(ctx, "anycast_ip_id", id)

	tflog.Trace(ctx, "Calling Networking API to retrieve anycast IP")

	anycastIP, err := anycastips.Get(client, id).Extract()
	if errutil.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error calling Networking API to retrieve anycast IP", err.Error())
		return
	}

	tflog.Trace(ctx, "Called Networking API to retrieve anycast IP", map[string]interface{}{"anycast_ip": fmt.Sprintf("%#v", anycastIP)})

	resp.Diagnostics.Append(data.UpdateFromAnycastIP(ctx, anycastIP)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *anycastIPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan resource_anycastip.AnycastipModel
	var data resource_anycastip.AnycastipModel

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

	client, err := r.config.NetworkingV2Client(region, inetworking.SprutSDN)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Networking API client", err.Error())
		return
	}

	associations, diags := resource_anycastip.ExpandAssociations(ctx, plan.Associations)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	healthCheck := resource_anycastip.ExpandHealthCheck(ctx, plan.HealthCheck)

	id := data.Id.ValueString()
	ctx = tflog.SetField(ctx, "anycast_ip_id", id)

	updateOpts := anycastips.UpdateOpts{
		Name:         plan.Name.ValueString(),
		Description:  plan.Description.ValueString(),
		Associations: associations,
		HealthCheck:  healthCheck,
	}

	tflog.Trace(ctx, "Calling Networking API to update anycast IP", map[string]interface{}{"opts": fmt.Sprintf("%#v", updateOpts)})

	anycastIP, err := anycastips.Update(client, id, &updateOpts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling Networking API to update anycast IP", err.Error())
		return
	}

	tflog.Trace(ctx, "Called Networking API to update anycast IP", map[string]interface{}{"anycast_ip": fmt.Sprintf("%#v", anycastIP)})

	resp.Diagnostics.Append(data.UpdateFromAnycastIP(ctx, anycastIP)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *anycastIPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_anycastip.AnycastipModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	client, err := r.config.NetworkingV2Client(region, inetworking.SprutSDN)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Networking API client", err.Error())
		return
	}

	id := data.Id.ValueString()
	ctx = tflog.SetField(ctx, "anycast_ip_id", id)

	tflog.Trace(ctx, "Calling Networking API to delete anycast IP")

	err = anycastips.Delete(client, id).ExtractErr()
	if errutil.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error calling Networking API to delete anycast IP", err.Error())
		return
	}

	tflog.Trace(ctx, "Called Networking API to delete anycast IP")
}

func (r *anycastIPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
