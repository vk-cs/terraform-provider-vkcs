// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package datasource_ssl_certificate

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func SslCertificateDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				Description:         "ID of the SSL certificate.",
				MarkdownDescription: "ID of the SSL certificate.",
			},
			"issuer": schema.StringAttribute{
				Computed:            true,
				Description:         "Name of the certification center issued the SSL certificate.",
				MarkdownDescription: "Name of the certification center issued the SSL certificate.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "SSL certificate name.",
				MarkdownDescription: "SSL certificate name.",
			},
			"region": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The region in which to obtain the CDN client. If omitted, the `region` argument of the provider is used.",
				MarkdownDescription: "The region in which to obtain the CDN client. If omitted, the `region` argument of the provider is used.",
			},
			"subject_cn": schema.StringAttribute{
				Computed:            true,
				Description:         "Domain name that the SSL certificate secures.",
				MarkdownDescription: "Domain name that the SSL certificate secures.",
			},
			"validity_not_after": schema.StringAttribute{
				Computed:            true,
				Description:         "Date when certificate become untrusted (ISO 8601/RFC 3339 format, UTC.).",
				MarkdownDescription: "Date when certificate become untrusted (ISO 8601/RFC 3339 format, UTC.).",
			},
			"validity_not_before": schema.StringAttribute{
				Computed:            true,
				Description:         "Date when certificate become valid (ISO 8601/RFC 3339 format, UTC.).",
				MarkdownDescription: "Date when certificate become valid (ISO 8601/RFC 3339 format, UTC.).",
			},
		},
	}
}

type SslCertificateModel struct {
	Id                types.Int64  `tfsdk:"id"`
	Issuer            types.String `tfsdk:"issuer"`
	Name              types.String `tfsdk:"name"`
	Region            types.String `tfsdk:"region"`
	SubjectCn         types.String `tfsdk:"subject_cn"`
	ValidityNotAfter  types.String `tfsdk:"validity_not_after"`
	ValidityNotBefore types.String `tfsdk:"validity_not_before"`
}
