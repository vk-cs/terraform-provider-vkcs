package cdn

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/cdn/datasource_ssl_certificate"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/cdn/v1/ssldata"
)

var (
	_ datasource.DataSource              = (*sslCertificateDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*sslCertificateDataSource)(nil)
)

func NewSslCertificateDataSource() datasource.DataSource {
	return &sslCertificateDataSource{}
}

type sslCertificateDataSource struct {
	config clients.Config
}

func (d *sslCertificateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cdn_ssl_certificate"
}

func (d *sslCertificateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_ssl_certificate.SslCertificateDataSourceSchema(ctx)
}

func (d *sslCertificateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *sslCertificateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_ssl_certificate.SslCertificateModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = d.config.GetRegion()
	}
	data.Region = types.StringValue(region)

	client, err := d.config.CDNV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating CDN API client", err.Error())
		return
	}

	tflog.Trace(ctx, "Calling CDN API to list SSL certificates")

	sslCerts, err := ssldata.List(client, d.config.GetTenantID()).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling CDN API to list SSL certificates", err.Error())
		return
	}

	tflog.Trace(ctx, "Called CDN API to list SSL certificates", map[string]interface{}{"ssl_certs": fmt.Sprintf("%#v", sslCerts)})

	name := data.Name.ValueString()
	i := slices.IndexFunc(sslCerts, func(cert ssldata.SSLCertificate) bool {
		return cert.Name == name
	})
	if i == -1 {
		resp.Diagnostics.AddError("Error finding a SSL certificate", fmt.Sprintf("SSL certificate with name %q not found", name))
		return
	}

	sslCert := sslCerts[i]

	data.Id = types.Int64Value(int64(sslCert.ID))
	data.Issuer = types.StringValue(sslCert.CertIssuer)
	data.Name = types.StringValue(sslCert.Name)
	data.SubjectCn = types.StringValue(sslCert.CertSubjectCN)
	data.ValidityNotAfter = types.StringValue(sslCert.ValidityNotAfter)
	data.ValidityNotBefore = types.StringValue(sslCert.ValidityNotBefore)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
