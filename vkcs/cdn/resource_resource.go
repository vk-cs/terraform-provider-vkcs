package cdn

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/cdn/resource_resource"
	resource_validators "github.com/vk-cs/terraform-provider-vkcs/vkcs/cdn/resource_resource/validators"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	fwutils "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/framework/utils"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/cdn/v1/resources"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

var (
	_ resource.Resource                     = (*resourceResource)(nil)
	_ resource.ResourceWithConfigure        = (*resourceResource)(nil)
	_ resource.ResourceWithConfigValidators = (*resourceResource)(nil)
	_ resource.ResourceWithImportState      = (*resourceResource)(nil)
	_ resource.ResourceWithValidateConfig   = (*resourceResource)(nil)
)

func NewResourceResource() resource.Resource {
	return &resourceResource{}
}

type resourceResource struct {
	config clients.Config
}

func (r *resourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cdn_resource"
}

func (r *resourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_resource.ResourceResourceSchema(ctx)
}

func (r *resourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *resourceResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resource_validators.ConflictingEnabled(
			path.MatchRoot("options").AtName("fetch_compressed"),
			path.MatchRoot("options").AtName("gzip_on"),
		),
		resource_validators.ConflictingEnabled(
			path.MatchRoot("options").AtName("forward_host_header"),
			path.MatchRoot("options").AtName("host_header"),
		),
		resource_validators.ConflictingEnabled(
			path.MatchRoot("options").AtName("ignore_query_string"),
			path.MatchRoot("options").AtName("query_params_blacklist"),
			path.MatchRoot("options").AtName("query_params_whitelist"),
		),
	}
}

func (r *resourceResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config resource_resource.ResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sslCertType := config.SslCertificate.SslCertificateType.ValueString()
	if sslCertType == string(resource_resource.SslCertificateProviderTypeOwn) && config.SslCertificate.Id.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("ssl_certificate"),
			"Missing Attribute Value",
			"`ssl_certificate.id` must be configured when `ssl_certificate.type` is \"own\"",
		)
	}

	shieldingEnabled := config.Shielding.Enabled.ValueBool()
	if shieldingEnabled && config.Shielding.PopId.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("shielding"),
			"Missing Attribute Value",
			"`shielding.pop_id` must be configured when `shielding.enabled` is \"true\"",
		)
	}
}

func (r *resourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_resource.ResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, resource_resource.ResourceReadyTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

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

	secondaryHostnames := make([]string, 0, len(data.SecondaryHostnames.Elements()))
	if !data.SecondaryHostnames.IsNull() && !data.SecondaryHostnames.IsUnknown() {
		resp.Diagnostics.Append(data.SecondaryHostnames.ElementsAs(ctx, &secondaryHostnames, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	opts, diags := data.Options.ToResourceOptions(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpts := resources.CreateOpts{
		Active:             data.Active.ValueBool(),
		CNAME:              data.Cname.ValueString(),
		Options:            opts,
		OriginGroup:        int(data.OriginGroup.ValueInt64()),
		OriginProtocol:     resources.ResourceOriginProtocol(data.OriginProtocol.ValueString()),
		SecondaryHostnames: secondaryHostnames,
	}

	if sslOpts := data.SslCertificate.ToSslOpts(); sslOpts != nil {
		createOpts.SSLEnabled = sslOpts.Enabled
		createOpts.SSLData = sslOpts.Data
	}

	projectID := r.config.GetTenantID()

	tflog.Trace(ctx, "Calling CDN API to create CDN resource", map[string]interface{}{"opts": fmt.Sprintf("%#v", createOpts)})

	resource, err := resources.Create(client, projectID, &createOpts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling CDN API to create CDN resource", err.Error())
		return
	}

	tflog.Trace(ctx, "Called CDN API to create CDN resource", map[string]interface{}{"resource": fmt.Sprintf("%#v", resource)})

	id := resource.ID
	resp.State.SetAttribute(ctx, path.Root("id"), id)

	resp.Diagnostics.Append(resource_resource.WaitForResourceReady(ctx, client, projectID, id, timeout)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sslCertType := data.SslCertificate.SslCertificateType.ValueString()
	if sslCertType == string(resource_resource.SslCertificateProviderTypeLetsEncrypt) {
		resp.Diagnostics.Append(issueLetsEncrypt(ctx, client, projectID, id)...)
		if resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append(resource_resource.WaitForResourceReady(ctx, client, projectID, id, timeout)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	shieldingOpts := data.Shielding.ToUpdateShieldingOpts()
	if shieldingOpts != nil {
		resourceShielding, diags := updateShielding(ctx, client, projectID, id, shieldingOpts)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		data.Shielding = resource_resource.ShieldingValue{}.FromResourceShielding(resourceShielding)

		resp.Diagnostics.Append(resource_resource.WaitForResourceReady(ctx, client, projectID, id, timeout)...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		data.Shielding = resource_resource.NewShieldingValueNull()
	}

	resource, diags = retrieveCDNResource(ctx, client, projectID, id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(data.UpdateFromResource(ctx, resource)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(updateFromLetsEncryptStatus(ctx, client, projectID, id, &data.SslCertificate)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_resource.ResourceModel

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
	ctx = tflog.SetField(ctx, "resource_id", id)
	projectID := r.config.GetTenantID()

	tflog.Trace(ctx, "Calling CDN API to retrieve CDN resource")

	resource, err := resources.Get(client, projectID, id).Extract()
	if errutil.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error calling CDN API to retrieve CDN resource", err.Error())
		return
	}

	tflog.Trace(ctx, "Called CDN API to retrieve CDN resource", map[string]interface{}{"resource": fmt.Sprintf("%#v", resource)})

	resp.Diagnostics.Append(data.UpdateFromResource(ctx, resource)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(updateFromLetsEncryptStatus(ctx, client, projectID, id, &data.SslCertificate)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Calling CDN API to retrieve origin shielding settings")

	resourceShielding, err := resources.GetShielding(client, projectID, id).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling CDN API to retrieve origin shielding settings", err.Error())
		return
	}

	tflog.Trace(ctx, "Called CDN API to retrieve origin shielding settings", map[string]interface{}{"resource_shielding": fmt.Sprintf("%#v", resourceShielding)})

	data.Shielding = resource_resource.ShieldingValue{}.FromResourceShielding(resourceShielding)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan resource_resource.ResourceModel
	var data resource_resource.ResourceModel

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

	timeout, diags := data.Timeouts.Update(ctx, resource_resource.ResourceReadyTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	client, err := r.config.CDNV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating CDN API client", err.Error())
		return
	}

	id := int(data.Id.ValueInt64())
	ctx = tflog.SetField(ctx, "resource_id", id)

	secondaryHostnames := make([]string, 0, len(plan.SecondaryHostnames.Elements()))
	if !plan.SecondaryHostnames.IsNull() && !plan.SecondaryHostnames.IsUnknown() {
		resp.Diagnostics.Append(plan.SecondaryHostnames.ElementsAs(ctx, &secondaryHostnames, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	opts, diags := plan.Options.ToResourceOptions(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateOpts := resources.UpdateOpts{
		Active:             plan.Active.ValueBoolPointer(),
		Options:            opts,
		OriginGroup:        int(plan.OriginGroup.ValueInt64()),
		OriginProtocol:     resources.ResourceOriginProtocol(plan.OriginProtocol.ValueString()),
		SecondaryHostnames: secondaryHostnames,
	}

	if sslOpts := plan.SslCertificate.ToSslOpts(); sslOpts != nil {
		updateOpts.SSLEnabled = &sslOpts.Enabled
		updateOpts.SSLData = sslOpts.Data
	}

	projectID := r.config.GetTenantID()

	tflog.Trace(ctx, "Calling CDN API to update CDN resource", map[string]interface{}{"opts": fmt.Sprintf("%#v", updateOpts)})

	resource, err := resources.Update(client, projectID, id, &updateOpts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling CDN API to update CDN resource", err.Error())
		return
	}

	tflog.Trace(ctx, "Called CDN API to update CDN resource", map[string]interface{}{"resource": fmt.Sprintf("%#v", resource)})

	resp.Diagnostics.Append(resource_resource.WaitForResourceReady(ctx, client, projectID, id, timeout)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sslCertType := plan.SslCertificate.SslCertificateType.ValueString()
	if !plan.SslCertificate.SslCertificateType.Equal(data.SslCertificate.SslCertificateType) && sslCertType == string(resource_resource.SslCertificateProviderTypeLetsEncrypt) {
		resp.Diagnostics.Append(issueLetsEncrypt(ctx, client, projectID, id)...)
		if resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append(resource_resource.WaitForResourceReady(ctx, client, projectID, id, timeout)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !plan.Shielding.Equal(data.Shielding) {
		shieldingOpts := plan.Shielding.ToUpdateShieldingOpts()
		if shieldingOpts != nil {
			resourceShielding, diags := updateShielding(ctx, client, projectID, id, shieldingOpts)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}

			data.Shielding = resource_resource.ShieldingValue{}.FromResourceShielding(resourceShielding)

			resp.Diagnostics.Append(resource_resource.WaitForResourceReady(ctx, client, projectID, id, timeout)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}
	}

	resource, diags = retrieveCDNResource(ctx, client, projectID, id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(data.UpdateFromResource(ctx, resource)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(updateFromLetsEncryptStatus(ctx, client, projectID, id, &data.SslCertificate)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_resource.ResourceModel

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

	id := data.Id.ValueInt64()
	ctx = tflog.SetField(ctx, "resource_id", id)

	tflog.Trace(ctx, "Calling CDN API to delete CDN resource")

	err = resources.Delete(client, r.config.GetTenantID(), int(id)).ExtractErr()
	if errutil.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error calling CDN API to delete CDN resource", err.Error())
		return
	}

	tflog.Trace(ctx, "Called CDN API to delete CDN resource")
}

func (r *resourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	fwutils.ImportStatePassthroughInt64ID(ctx, req, resp)
}

func retrieveCDNResource(ctx context.Context, client *gophercloud.ServiceClient, projectID string, resourceID int) (*resources.Resource, diag.Diagnostics) {
	var diags diag.Diagnostics
	tflog.Trace(ctx, "Calling CDN API to retrieve CDN resource")

	resource, err := resources.Get(client, projectID, resourceID).Extract()
	if err != nil {
		diags.AddError("Error calling CDN API to retrieve CDN resource", err.Error())
		return nil, diags
	}

	tflog.Trace(ctx, "Called CDN API to retrieve CDN resource", map[string]interface{}{"resource": fmt.Sprintf("%#v", resource)})
	return resource, diags
}

func issueLetsEncrypt(ctx context.Context, client *gophercloud.ServiceClient, projectID string, resourceID int) diag.Diagnostics {
	var diags diag.Diagnostics
	tflog.Trace(ctx, "Calling CDN API to issue a Let's Encrypt certificate")

	err := resources.IssueLetsEncrypt(client, projectID, resourceID).ExtractErr()
	if err != nil {
		diags.AddError("Error calling CDN API to issue a Let's Encrypt certificate", err.Error())
		return diags
	}

	tflog.Trace(ctx, "Called CDN API to issue a Let's Encrypt certificate")
	return diags
}

func updateFromLetsEncryptStatus(ctx context.Context, client *gophercloud.ServiceClient, projectID string, resourceID int, sslCertificateData *resource_resource.SslCertificateValue) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Trace(ctx, "Calling CDN API to get Let's Encrypt certificate issuing details")

	leStatus, err := resources.GetLetsEncryptStatus(client, projectID, resourceID).Extract()
	if err != nil {
		if !errutil.IsNotFound(err) {
			diags.AddError("Error calling CDN API to get Let's Encrypt certificate issuing details", err.Error())
		}
		return diags
	}

	if leStatus.Active {
		sslCertificateData.SslCertificateType = types.StringValue(string(resource_resource.SslCertificateProviderTypeLetsEncrypt))
		sslCertificateData.Status = types.StringValue(string(resource_resource.SslCertificateStatusBeingIssued))
	}

	tflog.Trace(ctx, "Called CDN API to get Let's Encrypt certificate issuing details", map[string]any{"status": fmt.Sprintf("%#v", leStatus)})
	return diags
}

func updateShielding(ctx context.Context, client *gophercloud.ServiceClient, projectID string, resourceID int, opts *resources.UpdateShieldingOpts) (*resources.ResourceShielding, diag.Diagnostics) {
	var diags diag.Diagnostics
	tflog.Trace(ctx, "Calling CDN API to update origin shielding settings")

	resourceShielding, err := resources.UpdateShielding(client, projectID, resourceID, opts).Extract()
	if err != nil {
		diags.AddError("Error calling CDN API to update origin shielding settings", err.Error())
		return nil, diags
	}

	tflog.Trace(ctx, "Called CDN API to update origin shielding settings")
	return resourceShielding, diags
}
