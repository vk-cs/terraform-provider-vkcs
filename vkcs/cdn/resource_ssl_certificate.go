package cdn

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/cdn/resource_ssl_certificate"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	fwutils "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/framework/utils"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/cdn/v1/ssldata"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

var (
	_ resource.Resource                = (*sslCertificateResource)(nil)
	_ resource.ResourceWithConfigure   = (*sslCertificateResource)(nil)
	_ resource.ResourceWithImportState = (*sslCertificateResource)(nil)
)

func NewSslCertificateResource() resource.Resource {
	return &sslCertificateResource{}
}

type sslCertificateResource struct {
	config clients.Config
}

func (r *sslCertificateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cdn_ssl_certificate"
}

func (r *sslCertificateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_ssl_certificate.SslCertificateResourceSchema(ctx)
}

func (r *sslCertificateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *sslCertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_ssl_certificate.SslCertificateModel

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

	addOpts := ssldata.AddOpts{
		Name:           data.Name.ValueString(),
		SSLCertificate: data.Certificate.ValueString(),
		SSLPrivateKey:  data.PrivateKey.ValueString(),
	}

	tflog.Trace(ctx, "Calling CDN API to add SSL certificate", map[string]interface{}{"opts": fmt.Sprintf("%#v", &addOpts)})

	sslCert, err := ssldata.Add(client, r.config.GetProjectID(), &addOpts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling CDN API to add SSL certificate", err.Error())
		return
	}

	tflog.Trace(ctx, "Called CDN API to add SSL certificate", map[string]interface{}{"ssl_cert": fmt.Sprintf("%#v", sslCert)})

	id := types.Int64Value(int64(sslCert.ID))
	resp.State.SetAttribute(ctx, path.Root("id"), id)

	resp.Diagnostics.Append(data.UpdateFromSslCertificate(ctx, sslCert)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sslCertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_ssl_certificate.SslCertificateModel

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
	ctx = tflog.SetField(ctx, "ssl_certificate_id", id)

	tflog.Trace(ctx, "Calling CDN API to list SSL certificates")

	sslCerts, err := ssldata.List(client, r.config.GetProjectID()).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling CDN API to list SSL certificates", err.Error())
		return
	}

	tflog.Trace(ctx, "Called CDN API to list SSL certificates", map[string]interface{}{"ssl_certs": fmt.Sprintf("%#v", sslCerts)})

	i := slices.IndexFunc(sslCerts, func(cert ssldata.SSLCertificate) bool {
		return cert.ID == id
	})
	if i == -1 {
		tflog.Debug(ctx, "Removing the resource from the state due to missing corresponding API object")
		resp.State.RemoveResource(ctx)
		return
	}

	sslCert := &sslCerts[i]
	resp.Diagnostics.Append(data.UpdateFromSslCertificate(ctx, sslCert)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sslCertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan resource_ssl_certificate.SslCertificateModel
	var data resource_ssl_certificate.SslCertificateModel

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

	id := int(data.Id.ValueInt64())
	ctx = tflog.SetField(ctx, "ssl_certificate_id", id)

	updateOpts := ssldata.UpdateOpts{
		Name: plan.Name.ValueString(),
	}

	tflog.Trace(ctx, "Calling CDN API to update SSL certificate", map[string]interface{}{"opts": fmt.Sprintf("%#v", &updateOpts)})

	sslCert, err := ssldata.Update(client, r.config.GetProjectID(), id, &updateOpts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling CDN API to add SSL certificate", err.Error())
		return
	}

	tflog.Trace(ctx, "Called CDN API to add SSL certificate", map[string]interface{}{"ssl_cert": fmt.Sprintf("%#v", sslCert)})

	resp.Diagnostics.Append(data.UpdateFromSslCertificate(ctx, sslCert)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sslCertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_ssl_certificate.SslCertificateModel

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

	tflog.Trace(ctx, "Calling CDN API to delete SSL certificate")

	err = ssldata.Delete(client, r.config.GetProjectID(), id).ExtractErr()
	if errutil.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error calling CDN API to delete SSL certificate", err.Error())
		return
	}

	tflog.Trace(ctx, "Called CDN API to delete SSL certificate")
}

func (r *sslCertificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	fwutils.ImportStatePassthroughInt64ID(ctx, req, resp)
}
